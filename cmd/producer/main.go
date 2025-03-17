package main

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/bignyap/kafka-go/handler"
	"github.com/bignyap/kafka-go/pkg/producer"
	"github.com/bignyap/kafka-go/pkg/ws"
)

func main() {

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Return.Errors = true

	producer, err := producer.NewSaramaProducer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	// Start listening for errors for kafka producer client
	producer.StartErrorListener()

	wsms := ws.NewWebSocketMessageSender()

	handler.StartWebServer(producer, wsms)

}
