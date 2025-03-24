package handler

import (
	"fmt"
	"net/http"

	"github.com/bignyap/kafka-go/pkg/middleware"
	"github.com/bignyap/kafka-go/pkg/ws"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (app *Application) WebSocketHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		conn, err := ws.UpgradeToWebSocket(w, r)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error upgrading to WebSocket: %v", err), http.StatusInternalServerError)
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

		app.Store.MessageBroadcaster.CreateConnection(parsedToken.Sub, conn)
	}
}
