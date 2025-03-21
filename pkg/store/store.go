package store

import (
	"context"
	"database/sql"

	"github.com/bignyap/kafka-go/pkg/models"
)

type Store struct {
	ChatRoom interface {
		AddMemberToRoom(context.Context, int, int) error
		RemoveMemberFromRoom(context.Context, int, int) error
		GetMembersFromRoom(context.Context, int) ([]models.Member, error)
		GetChatRoomsForMember(context.Context, int) ([]models.ChatRoom, error)
	}
	Message interface {
		SendMessageToRoom(context.Context, int, string) error
		GetMessagesFromRoom(context.Context, int) ([]models.ChatMessage, error)
	}
}

func NewStore(db *sql.DB) Store {
	return Store{
		ChatRoom: &ChatRoomStore{db},
		Message:  &MessageStore{db},
	}
}
