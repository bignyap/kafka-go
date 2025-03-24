package store

import (
	"context"
	"encoding/json"

	"github.com/bignyap/kafka-go/pkg/producer"
)

type KafkaProducerStore struct {
	producer producer.KafkaProducer
}

func (kafkaProducer KafkaProducerStore) ProduceChatMessages(
	ctx context.Context,
	topic string, message string,
) error {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return kafkaProducer.producer.SendMessage(topic, messageBytes)
}
