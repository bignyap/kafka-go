package main

import (
	"context"
	"expvar"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"

	"github.com/bignyap/kafka-go/handler"
	"github.com/bignyap/kafka-go/pkg/consumer"
	"github.com/bignyap/kafka-go/pkg/db"
	"github.com/bignyap/kafka-go/pkg/store"
	"github.com/bignyap/kafka-go/pkg/utils"
	"github.com/bignyap/kafka-go/pkg/ws"
	"go.uber.org/zap"
)

func init() {
	utils.LoadEnv()
}

func main() {

	// Define the logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// Create the application configuration
	appConfig := handler.AppConfig{
		Address: utils.GetEnvString("APPLICATION_PORT", "8080"),
		ApiURL:  utils.GetEnvString("API_URL", "localhost:8080/api"),
		Env:     utils.GetEnvString("ENV", "development"),
		Version: utils.GetEnvString("VERSION", "1.0"),
		DBConfig: db.DBConfig{
			SQLDriver: utils.GetEnvString("SQL_DRIVER", "postgres"),
			URLConfig: db.DBURLConfig{
				Username: utils.GetEnvString("DB_USER", "root"),
				Password: utils.GetEnvString("DB_PASSWORD", "root"),
				Database: utils.GetEnvString("DB_NAME", "chatapp"),
			},
			PoolConfig: db.DBPoolConfig{
				MaxOpenConnections: 10,
				MaxIdleConnections: 10,
				MaxIdleTime:        10,
			},
		},
		KafkaConfig: handler.KafkaConfig{
			Address: utils.GetEnvString("APPLICATION_PORT", "8080"),
		},
	}

	// Create the Database Connection
	db, err := appConfig.DBConfig.Connect()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	logger.Info("database connection pool established")

	// Define the Kafka consumer
	consumer, err := consumer.NewKafkaConsumer(
		utils.GetEnvString("KAFKA_URL", "localhost:9092"),
		utils.GetEnvString("CHAT_GROUP", "test"),
	)
	if err != nil {
		logger.Fatal(err)
	}
	defer consumer.Close()

	// Define the new repository store with db and producer
	store := store.NewConsumerStore(db, consumer)
	app := &handler.Application{
		Config: appConfig,
		Store:  store,
		Logger: logger,
	}

	expvar.NewString("version").Set(appConfig.Version)
	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	messageSender := ws.NewWebSocketMessageSender()
	// Here you need to implement the logic to manage WebSocket connections
	// This includes opening connections when a member joins, and closing connections when a member leaves

	dbConn := db.NewDBConn()
	defer dbConn.Close()

	consumerHandler := consumer.NewConsumerHandler(dbConn, messageSender)

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
