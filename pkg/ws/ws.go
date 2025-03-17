package ws

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type MessageSender interface {
	SendMessage(member string, message string) error
}

type WebSocketMessageSender struct {
	connections map[string]*websocket.Conn
}

func NewWebSocketMessageSender() *WebSocketMessageSender {
	return &WebSocketMessageSender{
		connections: make(map[string]*websocket.Conn),
	}
}

func (wsms *WebSocketMessageSender) CreateConnection(member string, conn *websocket.Conn) {
	wsms.connections[member] = conn
}

func (wsms *WebSocketMessageSender) SendMessage(member string, message string) error {
	conn, ok := wsms.connections[member]
	if !ok {
		return fmt.Errorf("no connection found for member: %s", member)
	}
	return conn.WriteMessage(websocket.TextMessage, []byte(message))
}
