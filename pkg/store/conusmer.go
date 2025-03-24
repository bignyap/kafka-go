package store

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bignyap/kafka-go/pkg/consumer"
)

type KafkaConsumerStore struct {
	producer consumer.KafkaConsumer
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
	return nil
}
