package main

import (
	"expvar"
	"runtime"

	"github.com/bignyap/kafka-go/handler"
	"github.com/bignyap/kafka-go/pkg/producer"
	"github.com/bignyap/kafka-go/pkg/store"
	"github.com/bignyap/kafka-go/pkg/utils"
	"github.com/bignyap/kafka-go/pkg/ws"
	"go.uber.org/zap"
)

func init() {
	utils.LoadEnv()
}

func main() {

	config := &handler.config{
		addr:   utils.GetEnvString("APPLICATION_PORT", "8080"),
		apiURL: utils.GetEnvString("APPLICATION_PORT", "8080"),
		env:    utils.GetEnvString("APPLICATION_PORT", "8080"),
		dbConfig: handler.dbConfig{
			username:     utils.GetEnvString("APPLICATION_PORT", "8080"),
			password:     utils.GetEnvString("APPLICATION_PORT", "8080"),
			database:     utils.GetEnvString("APPLICATION_PORT", "8080"),
			maxOpenConns: 10,
			maxIdleConns: 10,
			// maxIdleTime:  utils.GetEnvString("APPLICATION_PORT", "8080"),
		},
		kafkaCfg: handler.kafkaConfig{
			addr: utils.GetEnvString("APPLICATION_PORT", "8080"),
		},
	}

	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	db, err := db.New(
		config.dbConfig.addr,
		config.dbConfig.maxOpenConns,
		config.dbConfig.maxIdleConns,
		config.dbConfig.maxIdleTime,
	)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Info("database connection pool established")

	producer, err := producer.NewProducer(
		utils.GetEnvString("KAFKA_URL", "localhost:9092"),
	)
	if err != nil {
		logger.Fatal(err)
	}
	defer producer.Close()

	store := store.NewStore(db)

	app := &application{
		config: config,
		store:  store,
		logger: logger,
	}

	expvar.NewString("version").Set(version)
	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	wsms := ws.NewWebSocketMessageSender()

	handler.StartWebServer(producer, wsms)

}
