package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bignyap/kafka-go/pkg/middleware"
	"github.com/bignyap/kafka-go/pkg/producer"
	"github.com/bignyap/kafka-go/pkg/utils"
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
	port := utils.GetEnvString("APPLICATION_PORT", "8080")
	port = fmt.Sprintf(":%s", port)
	log.Fatal(http.ListenAndServe(port, middlewareMux))
}
