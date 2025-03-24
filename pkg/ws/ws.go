package ws

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func UpgradeToWebSocket(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to WebSocket: %v", err)
		return nil, fmt.Errorf("failed to upgrade connection to WebSocket: %w", err)
	}

	return conn, nil
}
