package store

import (
	"context"
	"database/sql"

	"github.com/bignyap/kafka-go/pkg/consumer"
	"github.com/bignyap/kafka-go/pkg/models"
	"github.com/bignyap/kafka-go/pkg/producer"
	"github.com/gorilla/websocket"
)

type DataStore struct {
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

type WebSocketStore struct {
	MessageBroadcaster interface {
		CreateConnection(string, *websocket.Conn) error
		SendMessage(string, string) error
	}
}

type ProducerStore struct {
	DataStore
	WebSocketStore
	MessageProducer interface {
		ProduceChatMessages(context.Context, string, string) error
	}
}

type ConsumerStore struct {
	DataStore
	WebSocketStore
	MessageConsumer interface {
		ConsumeChatMessages(context.Context, string, string) error
	}
	MessageBroadcaster interface {
		CreateConnection(string, *websocket.Conn) error
		SendMessage(string, string) error
	}
}

func NewDataStore(db *sql.DB) DataStore {
	return DataStore{
		ChatRoom: &ChatRoomStore{db},
		Message:  &MessageStore{db},
	}
}

func NewProducerStore(db *sql.DB, producer producer.KafkaProducer) ProducerStore {
	return ProducerStore{
		DataStore: DataStore{
			ChatRoom: &ChatRoomStore{db},
			Message:  &MessageStore{db},
		},
		WebSocketStore: WebSocketStore{
			MessageBroadcaster: &WebSocketMessageSender{},
		},
		MessageProducer: KafkaProducerStore{producer},
	}
}

func NewConsumerStore(
	db *sql.DB, consumer consumer.KafkaConsumer,
) ConsumerStore {
	return ConsumerStore{
		DataStore: DataStore{
			ChatRoom: &ChatRoomStore{db},
			Message:  &MessageStore{db},
		},
		WebSocketStore: WebSocketStore{
			MessageBroadcaster: &WebSocketMessageSender{},
		},
		MessageConsumer: KafkaConsumerStore{consumer},
	}
}
