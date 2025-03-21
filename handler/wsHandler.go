package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bignyap/kafka-go/pkg/middleware"
	"github.com/bignyap/kafka-go/pkg/ws"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (app *application) WebSocketHandler(wsms *ws.WebSocketMessageSender) http.HandlerFunc {
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
