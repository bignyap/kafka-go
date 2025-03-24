package consumer

import (
	"context"
	"log"
	"strings"

	"github.com/IBM/sarama"
)

type KafkaConsumer interface {
	Consume(context.Context, []string, ConsumerHandler) error
	Close() error
}

type ConsumerHandler interface {
	Setup(sarama.ConsumerGroupSession) error
	Cleanup(sarama.ConsumerGroupSession) error
	ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error
}

type SaramaConsumer struct {
	Client sarama.ConsumerGroup
}

func NewKafkaConsumer(addr string, group string) (KafkaConsumer, error) {

	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Offsets.AutoCommit.Enable = false

	brokers := strings.Split(addr, ",")
	consumerClient, err := NewSarmaConsumer(brokers, group, config)
	if err != nil {
		log.Fatalf("unable to create kafka consumer: %v", err)
	}
	defer consumerClient.Client.Close()

	return consumerClient, nil
}

func NewSarmaConsumer(
	brokers []string,
	group string,
	config *sarama.Config,
) (*SaramaConsumer, error) {
	client, err := sarama.NewConsumerGroup(brokers, group, config)
	if err != nil {
		return nil, err
	}
	return &SaramaConsumer{Client: client}, nil
}

func (kc *SaramaConsumer) Consume(
	ctx context.Context,
	topic []string,
	handler ConsumerHandler,
) error {
	return kc.Client.Consume(ctx, topic, handler)
}

func (kc *SaramaConsumer) Close() error {
	return kc.Client.Close()
}
