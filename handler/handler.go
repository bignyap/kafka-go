package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bignyap/kafka-go/pkg/db"
	"github.com/bignyap/kafka-go/pkg/middleware"
	"github.com/bignyap/kafka-go/pkg/store"
	"github.com/bignyap/kafka-go/pkg/utils"
	"go.uber.org/zap"
)

type Application struct {
	Config AppConfig
	Store  store.ProducerStore
	Logger *zap.SugaredLogger
}

type AppConfig struct {
	Address     string
	ApiURL      string
	Env         string // Whether production or development
	Version     string
	DBConfig    db.DBConfig
	KafkaConfig KafkaConfig
}

type KafkaConfig struct {
	Address string
}

func (app *Application) Run() {

	mux := http.NewServeMux()

	mux.HandleFunc("/health", app.HealthHandler)
	mux.HandleFunc("/send-message", app.SendMessageHandler())
	mux.HandleFunc("/ws", app.WebSocketHandler())

	middlewareMux := middleware.ChainMiddleware(
		mux,
		middleware.CorsMiddleware,
		middleware.AuthMiddleware,
		middleware.LoggingMiddleware,
	)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", utils.GetEnvString("APPLICATION_PORT", "8080")),
		Handler:      middlewareMux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		// s := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdown <- srv.Shutdown(ctx)
	}()

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}

	err = <-shutdown
	if err != nil {
		log.Fatal(err)
	}
}
