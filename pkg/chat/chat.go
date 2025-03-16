package models

import (
	"database/sql"

	"github.com/bignyap/kafka-go/pkg/db"
	"github.com/bignyap/kafka-go/pkg/models"
)

type ChatRoomManager interface {
	AddMember(roomID int, memberID int) error
	RemoveMember(roomID int, memberID int) error
	GetMembers(roomID int) ([]models.Member, error)
	GetChatRooms(memberID int) ([]models.ChatRoom, error)
}

type ChatMessageManager interface {
	SendMessage(roomID int, message string) error
	GetMessages(roomID int) ([]models.ChatMessage, error)
}

type ChatRoomManagerImpl struct {
	DB *sql.DB
}

func NewChatRoomManagerImpl(db *sql.DB) ChatRoomManager {
	return &ChatRoomManagerImpl{DB: db}
}

func (crm *ChatRoomManagerImpl) AddMember(roomID int, memberID int) error {
	return db.AddMemberToRoom(crm.DB, roomID, memberID)
}

func (crm *ChatRoomManagerImpl) RemoveMember(roomID int, memberID int) error {
	return db.RemoveMemberFromRoom(crm.DB, roomID, memberID)
}

func (crm *ChatRoomManagerImpl) GetMembers(roomID int) ([]models.Member, error) {
	return db.GetMembersFromRoom(crm.DB, roomID)
}

func (crm *ChatRoomManagerImpl) GetChatRooms(memberID int) ([]models.ChatRoom, error) {
	return db.GetChatRoomsForMember(crm.DB, memberID)
}

type ChatMessageManagerImpl struct {
	DB *sql.DB
}

func NewChatMessageManagerImpl(db *sql.DB) ChatMessageManager {
	return &ChatMessageManagerImpl{DB: db}
}

func (cmm *ChatMessageManagerImpl) SendMessage(roomID int, message string) error {
	return db.SendMessageToRoom(cmm.DB, roomID, message)
}

func (cmm *ChatMessageManagerImpl) GetMessages(roomID int) ([]models.ChatMessage, error) {
	return db.GetMessagesFromRoom(cmm.DB, roomID)
}
