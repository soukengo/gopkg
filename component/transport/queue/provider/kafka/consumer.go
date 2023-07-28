package kafka

import (
	"context"
	"github.com/soukengo/gopkg/component/transport/queue/iface"
	"github.com/soukengo/gopkg/infra/storage/kafka"
	"github.com/soukengo/gopkg/log"
)

type kafkaConsumer struct {
	cfg      *Config
	consumer kafka.Consumer
	logger   log.Logger
}

func NewConsumer(cfg *Config, logger log.Logger) (iface.Consumer, error) {
	consumer, err := kafka.NewConsumer(cfg.Config, cfg.Consumer, logger)
	if err != nil {
		return nil, err
	}
	return &kafkaConsumer{cfg: cfg, consumer: consumer, logger: logger}, nil
}

func (c *kafkaConsumer) Start(context.Context) (err error) {
	return c.consumer.Start()
}

func (c *kafkaConsumer) Stop(context.Context) (err error) {
	return c.consumer.Stop()
}

func (c *kafkaConsumer) Subscribe(topic iface.Topic, handler *iface.Handler) {
	c.consumer.Subscribe(string(topic), func(m *kafka.Message) {
		msg := iface.NewBytesMessage(iface.Topic(m.Topic()), m.Value())
		handler.Process(msg, func(action iface.Action) {
			if action == iface.CommitMessage {
				c.consumer.Ack(m)
			}
		})
	})
}
