package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/bignyap/kafka-go/pkg/producer"
)

func SendMessageHandler(producer producer.KafkaProducer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, fmt.Sprintf("error reading request body: %v", err), http.StatusBadRequest)
			return
		}

		if err := producer.ProduceMsgToKafka("test", string(body)); err != nil {
			http.Error(w, fmt.Sprintf("error producing Kafka message: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Message sent")
	}
}

func StartWebServer(producer producer.KafkaProducer) {
	mux := http.NewServeMux()
	mux.HandleFunc("/send-message", SendMessageHandler(producer))

	log.Fatal(http.ListenAndServe(":8080", mux))
}
