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

	"github.com/bignyap/kafka-go/pkg/middleware"
	"github.com/bignyap/kafka-go/pkg/store"
	"github.com/bignyap/kafka-go/pkg/utils"
	"go.uber.org/zap"
)

type application struct {
	config config
	store  store.Store
	logger *zap.SugaredLogger
}

type config struct {
	addr     string
	apiURL   string
	env      string // Whether production or development
	dbConfig dbConfig
	kafkaCfg kafkaConfig
}

type kafkaConfig struct {
	addr string
}

type dbConfig struct {
	username     string
	password     string
	database     string
	maxOpenConns int
	maxIdleConns int
	// maxIdleTime  string
}

func (app *application) StartWebServer(
// kafkaProducer producer.KafkaProducer,
// wsms *ws.WebSocketMessageSender,
) {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(
			fmt.Sprintf("Server running on port %s", utils.GetEnvString("APPLICATION_PORT", "8080"))),
		)
	})
	mux.HandleFunc("/send-message", app.SendMessageHandler())
	mux.HandleFunc("/ws", app.WebSocketHandler(wsms))

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
		s := <-quit

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
