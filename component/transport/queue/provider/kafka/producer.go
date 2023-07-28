package kafka

import (
	"context"
	"github.com/soukengo/gopkg/component/transport/queue/iface"
	"github.com/soukengo/gopkg/component/transport/queue/options"
	"github.com/soukengo/gopkg/infra/storage/kafka"
	"github.com/soukengo/gopkg/log"
)

type kafkaProducer struct {
	producer *kafka.Producer
}

func NewProducer(cfg *Config, logger log.Logger) (iface.Producer, error) {
	producer, err := kafka.NewKafkaProducer(cfg.Config)
	if err != nil {
		return nil, err
	}
	return &kafkaProducer{producer: producer}, nil
}

func (c *kafkaProducer) Close() error {
	return c.producer.Close()
}

func (c *kafkaProducer) Publish(ctx context.Context, message iface.Message, opts *options.ProducerOptions) (err error) {
	value, err := iface.EncodeMessage(message, opts)
	if err != nil {
		return
	}
	return c.producer.Send(ctx, string(message.Topic()), value)
}
