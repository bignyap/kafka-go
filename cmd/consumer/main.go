package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"

	"github.com/IBM/sarama"
	"github.com/bignyap/kafka-go/pkg"
	"github.com/gorilla/websocket"
)

// MessageSender interface
type MessageSender interface {
	SendMessage(member string, message string) error
}

// WebSocketMessageSender struct
type WebSocketMessageSender struct {
	connections map[string]*websocket.Conn
}

// NewWebSocketMessageSender function
func NewWebSocketMessageSender() *WebSocketMessageSender {
	return &WebSocketMessageSender{
		connections: make(map[string]*websocket.Conn),
	}
}

// SendMessage method for WebSocketMessageSender
func (wsms *WebSocketMessageSender) SendMessage(member string, message string) error {
	conn, ok := wsms.connections[member]
	if !ok {
		return fmt.Errorf("no connection found for member: %s", member)
	}
	return conn.WriteMessage(websocket.TextMessage, []byte(message))
}

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

func NewKafkaConsumer(brokers []string, group string, config *sarama.Config) (*KafkaConsumer, error) {
	client, err := sarama.NewConsumerGroup(brokers, group, config)
	if err != nil {
		return nil, err
	}
	return &KafkaConsumer{client: client}, nil
}

func (kc *KafkaConsumer) Consume(ctx context.Context, handler ConsumerHandler) error {
	return kc.client.Consume(ctx, []string{"test"}, handler)
}

type consumerHandler struct {
	cmm pkg.ChatMessageManager
	crm pkg.ChatRoomManager
	ms  MessageSender
}

func NewConsumerHandler(db *sql.DB, ms MessageSender) *consumerHandler {
	return &consumerHandler{
		cmm: &pkg.ChatMessageManagerImpl{DB: db},
		crm: &pkg.ChatRoomManagerImpl{DB: db},
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
		var chatMessage pkg.ChatMessage
		err := json.Unmarshal(msg.Value, &chatMessage)
		if err != nil {
			log.Printf("Error while decoding message: %v", err)
			continue
		}

		// Store the message in the database
		err = h.cmm.SendMessage(chatMessage.RoomID, chatMessage.Message)
		if err != nil {
			log.Printf("Error while storing message: %v", err)
			continue
		}

		// Get the members of the chat room
		members, err := h.crm.GetMembers(chatMessage.RoomID)
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

func main() {

	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Offsets.AutoCommit.Enable = false

	brokers := []string{"localhost:9092"}
	consumer, err := NewKafkaConsumer(brokers, "test-group", config)
	if err != nil {
		log.Fatalf("unable to create kafka consumer: %v", err)
	}
	defer consumer.client.Close()

	messageSender := NewWebSocketMessageSender()
	// Here you need to implement the logic to manage WebSocket connections
	// This includes opening connections when a member joins, and closing connections when a member leaves

	db, err := sql.Open("mysql", "user:password@/dbname")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	consumerHandler := NewConsumerHandler(db, messageSender)

	ctx, cancel := context.WithCancel(context.Background())
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			if err := consumer.Consume(ctx, consumerHandler); err != nil {
				log.Printf("consume error: %v", err)
			}

			select {
			case <-signals:
				cancel()
				return
			default:
			}
		}
	}()

	wg.Wait()
}
