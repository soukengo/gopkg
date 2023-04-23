package kafka

import (
	"github.com/soukengo/gopkg/component/queue/core"
	"github.com/soukengo/gopkg/log"
)

type kafkaQueue struct {
	core.Consumer
	core.Producer
}

func NewKafkaQueue(cfg *Config, logger log.Logger) (core.Queue, error) {
	consumer, err := NewConsumer(cfg, logger)
	if err != nil {
		return nil, err
	}
	producer, err := NewProducer(cfg, logger)
	if err != nil {
		return nil, err
	}
	return &kafkaQueue{
		Consumer: consumer,
		Producer: producer,
	}, nil
}
