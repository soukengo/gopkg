package kafka

import (
	"context"
	"github.com/soukengo/gopkg/component/transport/queue/iface"
	"github.com/soukengo/gopkg/component/transport/queue/options"
	"github.com/soukengo/gopkg/infra/storage/kafka"
	"github.com/soukengo/gopkg/log"
	"github.com/soukengo/gopkg/util/runtimes"
)

type kafkaConsumer struct {
	cfg      *Config
	consumer kafka.Consumer
	logger   log.Logger
}

func NewConsumer(cfg *Config, logger log.Logger) (iface.Consumer, error) {
	consumer, err := kafka.NewConsumer(cfg.Config, logger)
	if err != nil {
		return nil, err
	}
	return &kafkaConsumer{cfg: cfg, consumer: consumer, logger: logger}, nil
}

func (c *kafkaConsumer) Start() (err error) {
	return c.consumer.Start()
}

func (c *kafkaConsumer) Stop() (err error) {
	return c.consumer.Stop()
}

func (c *kafkaConsumer) Subscribe(topic iface.Topic, handler *iface.Handler) {
	c.consumer.Subscribe(string(topic), c.messageHandler(handler))
}

func (c *kafkaConsumer) messageHandler(handler *iface.Handler) kafka.HandlerFunc {
	return func(m *kafka.Message) {
		ctx := context.TODO()
		msg := iface.NewBytesMessage(iface.Topic(m.Topic()), m.Value())
		if handler.Options.Mode() == options.Async {
			runtimes.Async(func() {
				err := handler.Func(ctx, msg)
				if err != nil {
					c.consumer.Ack(m)
				}
			})
		} else {
			err := handler.Func(ctx, msg)
			if err != nil {
				c.consumer.Ack(m)
			}
		}
	}
}
