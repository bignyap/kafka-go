package consumer

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"strconv"

	"github.com/IBM/sarama"
	"github.com/bignyap/kafka-go/pkg/db"
	"github.com/bignyap/kafka-go/pkg/models"
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
	Client sarama.ConsumerGroup
}

func NewKafkaConsumer(brokers []string, group string, config *sarama.Config) (*KafkaConsumer, error) {
	client, err := sarama.NewConsumerGroup(brokers, group, config)
	if err != nil {
		return nil, err
	}
	return &KafkaConsumer{Client: client}, nil
}

func (kc *KafkaConsumer) Consume(ctx context.Context, handler ConsumerHandler) error {
	return kc.Client.Consume(ctx, []string{"test"}, handler)
}

type consumerHandler struct {
	cmm db.ChatMessageManager
	crm db.ChatRoomManager
	ms  websocket.MessageSender
}

func NewConsumerHandler(dbConn *sql.DB, ms websocket.MessageSender) *consumerHandler {
	return &consumerHandler{
		cmm: db.NewChatMessageManagerImpl(dbConn),
		crm: db.NewChatRoomManagerImpl(dbConn),
		ms:  ms,
	}
}

func (h *consumerHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var chatMessage models.ChatMessage
		err := json.Unmarshal(msg.Value, &chatMessage)
		if err != nil {
			log.Printf("Error while decoding message: %v", err)
			continue
		}

		// Store the message in the database
		err = h.cmm.SendMessage(strconv.Itoa(chatMessage.RoomID), chatMessage.Message)
		if err != nil {
			log.Printf("Error while storing message: %v", err)
			continue
		}

		// Get the members of the chat room
		members, err := h.crm.GetMembers(strconv.Itoa(chatMessage.RoomID))
		if err != nil {
			log.Printf("Error while getting members: %v", err)
			continue
		}

		// Send the message to the members
		for _, member := range members {
			if err := h.ms.SendMessage(strconv.Itoa(member.ID), chatMessage.Message); err != nil {
				log.Printf("Error while sending message: %v", err)
				continue
			}
		}

		sess.MarkMessage(msg, "")
		sess.Commit()
	}
	return nil
}
