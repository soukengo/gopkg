package kafka

import (
	"context"
	"github.com/soukengo/gopkg/component/queue/core"
	"github.com/soukengo/gopkg/infra/storage/kafka"
	"github.com/soukengo/gopkg/log"
)

type kafkaProducer struct {
	producer *kafka.Producer
}

func NewProducer(cfg *Config, logger log.Logger) (core.Producer, error) {
	producer, err := kafka.NewKafkaProducer(cfg.Config)
	if err != nil {
		return nil, err
	}
	return &kafkaProducer{producer: producer}, nil
}

func (c *kafkaProducer) Close() error {
	return c.producer.Close()
}

func (c *kafkaProducer) Submit(ctx context.Context, topic core.Topic, values []core.Value) (err error) {
	var kafkaValues = make([]kafka.Value, len(values))
	for i, value := range values {
		kafkaValues[i] = kafka.Value(value)
	}
	return c.producer.Submit(ctx, string(topic), kafkaValues)
}
