package producer

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/IBM/sarama"
)

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

func NewProducer(addr string) (KafkaProducer, error) {

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Return.Errors = true

	brokers := strings.Split(addr, ",")

	producer, err := NewSaramaProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	// Start listening for errors for kafka producer client
	producer.StartErrorListener()

	return producer, nil
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
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return producer.SendMessage(topic, messageBytes)
}
