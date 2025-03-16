package consumer

import (
	"context"
	"database/sql"

	"github.com/IBM/sarama"
	"github.com/bignyap/kafka-go/pkg/db"
	"github.com/bignyap/kafka-go/pkg/websocket"
)

type MessageConsumer interface {
	Consume(ctx context.Context, handler ConsumerHandler) error
}

type ConsumerHandler interface {
	Setup(sarama.ConsumerGroupSession) error
	Cleanup(sarama.ConsumerGroupSession) error
	ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error
}

type KafkaConsumer struct {
	client sarama.ConsumerGroup
}

// Implement the methods for KafkaConsumer

type consumerHandler struct {
	cmm db.ChatMessageManager
	crm db.ChatRoomManager
	ms  websocket.MessageSender
}

func NewConsumerHandler(db *sql.DB, ms websocket.MessageSender) *consumerHandler {
	return &consumerHandler{
		cmm: &db.ChatMessageManagerImpl{DB: db},
		crm: &db.ChatRoomManagerImpl{DB: db},
		ms:  ms,
	}
}

// Implement the methods for consumerHandler
