package store

import (
	"context"
	"database/sql"

	"github.com/bignyap/kafka-go/pkg/models"
	"github.com/bignyap/kafka-go/pkg/producer"
	"github.com/gorilla/websocket"
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
	MessageProducer interface {
		ProduceChatMessages(context.Context, string, string) error
	}
	MessageBroadcaster interface {
		CreateConnection(string, *websocket.Conn) error
		SendMessage(string, string) error
	}
}

func NewStore(db *sql.DB, producer producer.KafkaProducer) Store {
	return Store{
		ChatRoom:           &ChatRoomStore{db},
		Message:            &MessageStore{db},
		MessageProducer:    KafkaProducerStore{producer},
		MessageBroadcaster: &WebSocketMessageSender{},
	}
}
