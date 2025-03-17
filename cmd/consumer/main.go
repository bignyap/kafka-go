package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"

	"github.com/IBM/sarama"
	"github.com/bignyap/kafka-go/pkg/consumer"
	"github.com/bignyap/kafka-go/pkg/db"
	"github.com/bignyap/kafka-go/pkg/utils"
	"github.com/bignyap/kafka-go/pkg/ws"
)

func init() {
	utils.LoadEnv()
}

func main() {

	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Offsets.AutoCommit.Enable = false

	brokerEnv := utils.GetEnvString("KAFKA_URL", "localhost:9092")
	brokers := strings.Split(brokerEnv, ",")
	consumerClient, err := consumer.NewKafkaConsumer(brokers, "test-group", config)
	if err != nil {
		log.Fatalf("unable to create kafka consumer: %v", err)
	}
	defer consumerClient.Client.Close()

	messageSender := ws.NewWebSocketMessageSender()
	// Here you need to implement the logic to manage WebSocket connections
	// This includes opening connections when a member joins, and closing connections when a member leaves

	dbUrl := fmt.Sprintf(
		"%s:%s@/%s",
		utils.GetEnvString("DB_USER", ""),
		utils.GetEnvString("DB_PASSWORD", ""),
		utils.GetEnvString("DB_NAME", ""),
	)
	dbConfig := &db.DBConfig{
		SQLDriver:     utils.GetEnvString("DB_DRIVER", "mysql"),
		ConnectionURL: dbUrl,
	}
	db, err := dbConfig.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	consumerHandler := consumer.NewConsumerHandler(db, messageSender)

	ctx, cancel := context.WithCancel(context.Background())
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			if err := consumerClient.Consume(ctx, consumerHandler); err != nil {
				log.Printf("consume error: %v", err)
			}

			select {
			case <-signals:
				cancel()
				return
			default:
			}
		}
	}()

	wg.Wait()
}
