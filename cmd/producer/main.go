package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/IBM/sarama"
)

// https://www.tencentcloud.com/document/product/597/60360

var (
	initialRetryInterval = 2 * time.Second
	maxRetryInterval     = 30 * time.Second // Maximum wait time between retries
	maxRetries           = 3                // Maximum number of retries
	sendTimeout          = 5 * time.Second  // Timeout for sending messages
)

type KafkaProducer interface {
	SendMessage(topic string, message []byte) error
	Close() error
}

type SaramaProducer struct {
	producer sarama.AsyncProducer
}

func NewSaramaProducer(
	brokers []string,
	config *sarama.Config,
) (*SaramaProducer, error) {
	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}
	return &SaramaProducer{producer: producer}, nil
}

func (sp *SaramaProducer) SendMessage(topic string, message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}
	sp.producer.Input() <- msg
	return nil
}

func (sp *SaramaProducer) StartErrorListener() {
	go func() {
		for prodErr := range sp.producer.Errors() {
			log.Printf("Failed to send message to topic %s: %v", prodErr.Msg.Topic, prodErr.Err)

			retryInterval := initialRetryInterval
			for i := 0; i < maxRetries; i++ {
				time.Sleep(retryInterval) // Wait before retrying
				if retryErr := sp.retrySendMessage(prodErr.Msg); retryErr == nil {
					log.Printf("Message successfully resent to topic %s", prodErr.Msg.Topic)
					break // Exit the retry loop if successful
				} else if i == maxRetries-1 {
					log.Printf("Failed to resend message to topic %s after %d attempts", prodErr.Msg.Topic, maxRetries)
					// Implement additional error handling, such as alerting
				}
				retryInterval *= 2 // Exponential backoff
				if retryInterval > maxRetryInterval {
					retryInterval = maxRetryInterval
				}
			}
		}
	}()
}

func (sp *SaramaProducer) retrySendMessage(msg *sarama.ProducerMessage) error {
	select {
	case sp.producer.Input() <- msg:
		return nil // Message was successfully queued for sending
	case err := <-sp.producer.Errors():
		return err.Err // Return the error received while trying to send
	case <-time.After(sendTimeout): // Use the configurable timeout
		return fmt.Errorf("timeout while retrying message send to topic %s", msg.Topic)
	}
}

// Close closes the Sarama AsyncProducer.
func (sp *SaramaProducer) Close() error {
	return sp.producer.Close()
}

func ProduceMsgToKafka(producer KafkaProducer, topic string, message string) error {
	return producer.SendMessage(topic, []byte(message))
}

func SendMessageHandler(producer KafkaProducer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, fmt.Sprintf("error reading request body: %v", err), http.StatusBadRequest)
			return
		}

		if err := ProduceMsgToKafka(producer, "test", string(body)); err != nil {
			http.Error(w, fmt.Sprintf("error producing Kafka message: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Message sent")
	}
}

func StartWebServer(producer KafkaProducer) {
	mux := http.NewServeMux()
	mux.HandleFunc("/send-message", SendMessageHandler(producer))

	log.Fatal(http.ListenAndServe(":8080", mux))
}

func main() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Return.Errors = true

	producer, err := NewSaramaProducer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	// Start listening for errors
	producer.StartErrorListener()

	StartWebServer(producer)
}
