package kafka

import (
	"github.com/soukengo/gopkg/component/queue/core"
	"github.com/soukengo/gopkg/infra/storage/kafka"
	"github.com/soukengo/gopkg/log"
)

type kafkaConsumer struct {
	consumer *kafka.Consumer
}

func NewConsumer(cfg *Config, logger log.Logger) (core.Consumer, error) {
	var topics []string
	for _, topic := range cfg.Consumer.Topics {
		topics = append(topics, string(topic))
	}
	consumer, err := kafka.NewKafkaConsumer(cfg.Config,
		&kafka.ConsumerConfig{Topics: topics, GroupId: cfg.Consumer.GroupId, Workers: cfg.Consumer.Workers},
		logger)
	if err != nil {
		return nil, err
	}
	return &kafkaConsumer{consumer: consumer}, nil
}

func (c *kafkaConsumer) Start() (err error) {
	return c.consumer.Start()
}

func (c *kafkaConsumer) Stop() error {
	return c.consumer.Stop()
}

func (c *kafkaConsumer) Subscribe(topic core.Topic, handler core.HandlerFunc) {
	c.consumer.Subscribe(string(topic), func(topic string, value kafka.Value) (err error) {
		return handler(core.Topic(topic), core.Value(value))
	})
}
