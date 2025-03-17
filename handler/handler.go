package handler

import (
	"log"
	"net/http"

	"github.com/bignyap/kafka-go/pkg/middleware"
	"github.com/bignyap/kafka-go/pkg/producer"
	"github.com/bignyap/kafka-go/pkg/ws"
)

func StartWebServer(
	kafkaProducer producer.KafkaProducer,
	wsms *ws.WebSocketMessageSender,
) {
	mux := http.NewServeMux()
	mux.HandleFunc("/send-message", SendMessageHandler(kafkaProducer))
	mux.HandleFunc("/ws", WebSocketHandler(wsms))

	middlewareMux := middleware.ChainMiddleware(
		mux,
		middleware.CorsMiddleware,
		middleware.AuthMiddleware,
		middleware.LoggingMiddleware,
	)
	log.Fatal(http.ListenAndServe(":8080", middlewareMux))
}
