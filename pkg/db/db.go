package db

import "database/sql"

type ChatMessageManager interface {
	SendMessage(roomID string, message string) error
}

type ChatMessageManagerImpl struct {
	DB *sql.DB
}

// Implement the methods for ChatMessageManagerImpl

type ChatRoomManager interface {
	GetMembers(roomID string) ([]Member, error)
}

type ChatRoomManagerImpl struct {
	DB *sql.DB
}

// Implement the methods for ChatRoomManagerImpl
