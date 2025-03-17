package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"

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
		// Here you should identify the member. This is just a placeholder.
		member := "member1"
		wsms.CreateConnection(member, conn)
	}
}

func StartWebServer(
	kafkaProducer producer.KafkaProducer,
	wsms *ws.WebSocketMessageSender,
) {
	mux := http.NewServeMux()
	mux.HandleFunc("/send-message", SendMessageHandler(kafkaProducer))
	mux.HandleFunc("/ws", WebSocketHandler(wsms))

	log.Fatal(http.ListenAndServe(":8080", mux))
}
