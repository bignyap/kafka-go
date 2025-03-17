package main

import (
	"log"
	"strings"

	"github.com/IBM/sarama"
	"github.com/bignyap/kafka-go/handler"
	"github.com/bignyap/kafka-go/pkg/producer"
	"github.com/bignyap/kafka-go/pkg/utils"
	"github.com/bignyap/kafka-go/pkg/ws"
)

func init() {
	utils.LoadEnv()
}

func main() {

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Return.Errors = true

	brokerEnv := utils.GetEnvString("KAFKA_URL", "localhost:9092")
	brokers := strings.Split(brokerEnv, ",")

	producer, err := producer.NewSaramaProducer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	// Start listening for errors for kafka producer client
	producer.StartErrorListener()

	wsms := ws.NewWebSocketMessageSender()

	handler.StartWebServer(producer, wsms)

}
