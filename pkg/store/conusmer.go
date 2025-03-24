package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/IBM/sarama"
	"github.com/bignyap/kafka-go/pkg/consumer"
	"github.com/bignyap/kafka-go/pkg/models"
)

type KafkaConsumerStore struct {
	cmm consumer.KafkaConsumer
	crm *sql.DB
	ms  *WebSocketMessageSender
}

func (kafkaConsumer KafkaConsumerStore) ConsumeChatMessages(
	ctx context.Context,
	topic string, message string,
) error {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	fmt.Printf("The message is %s\n", string(messageBytes))
	consumerHandler := NewConsumerHandler(
		kafkaConsumer.crm, *kafkaConsumer.ms,
	)
	kafkaConsumer.cmm.Consume(
		ctx, []string{topic}, consumerHandler,
	)
	return nil
}

type consumerHandler struct {
	ds DataStore
	ms WebSocketMessageSender
}

func NewConsumerHandler(
	dbConn *sql.DB,
	ms WebSocketMessageSender,
) *consumerHandler {
	return &consumerHandler{
		ds: NewDataStore(dbConn),
		ms: ms,
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
		err = h.ds.Message.SendMessageToRoom(
			sess.Context(),
			chatMessage.RoomID,
			chatMessage.Message,
		)
		if err != nil {
			log.Printf("Error while storing message: %v", err)
			continue
		}

		// Get the members of the chat room
		members, err := h.ds.ChatRoom.GetChatRoomsForMember(
			sess.Context(), chatMessage.RoomID,
		)
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
