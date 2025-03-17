package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/bignyap/kafka-go/pkg/middleware"
	"github.com/bignyap/kafka-go/pkg/producer"
	"github.com/bignyap/kafka-go/pkg/ws"
	"github.com/gorilla/websocket"
)

func SendMessageHandler(kafkaProducer producer.KafkaProducer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, fmt.Sprintf("error reading request body: %v", err), http.StatusBadRequest)
			return
		}

		if err := producer.ProduceMsgToKafka(
			kafkaProducer, "test", string(body),
		); err != nil {
			http.Error(w, fmt.Sprintf("error producing Kafka message: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Message sent")
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func WebSocketHandler(wsms *ws.WebSocketMessageSender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		parsedToken, ok := r.Context().Value("parsedToken").(*middleware.ParsedToken)
		if !ok {
			http.Error(
				w, fmt.Sprintf("error reading token: %v", err),
				http.StatusBadRequest,
			)
			return
		}
		wsms.CreateConnection(parsedToken.Sub, conn)
	}
}

type Middleware func(http.Handler) http.Handler

func ChainMiddleware(mux http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		mux = middleware(mux)
	}
	return mux
}

func StartWebServer(
	kafkaProducer producer.KafkaProducer,
	wsms *ws.WebSocketMessageSender,
) {
	mux := http.NewServeMux()
	mux.HandleFunc("/send-message", SendMessageHandler(kafkaProducer))
	mux.HandleFunc("/ws", WebSocketHandler(wsms))

	middlewareMux := ChainMiddleware(
		mux,
		middleware.CorsMiddleware,
		middleware.AuthMiddleware,
		middleware.LoggingMiddleware,
	)
	log.Fatal(http.ListenAndServe(":8080", middlewareMux))
}
