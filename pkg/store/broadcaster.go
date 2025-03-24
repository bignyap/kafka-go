package store

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type WebSocketMessageSender struct {
	connections map[string]*websocket.Conn // Map to hold active connections, userID -> connection
}

func NewWebSocketMessageSender() *WebSocketMessageSender {
	return &WebSocketMessageSender{
		connections: make(map[string]*websocket.Conn),
	}
}

func (ws *WebSocketMessageSender) CreateConnection(userID string, conn *websocket.Conn) error {
	if _, exists := ws.connections[userID]; exists {
		return fmt.Errorf("user %s already has an active connection", userID)
	}
	ws.connections[userID] = conn
	log.Printf("Connection created for user: %s", userID)
	return nil
}

func (wsms *WebSocketMessageSender) SendMessage(member string, message string) error {
	conn, ok := wsms.connections[member]
	if !ok {
		return fmt.Errorf("no connection found for member: %s", member)
	}
	return conn.WriteMessage(websocket.TextMessage, []byte(message))
}
