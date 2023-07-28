package kafka

import (
	"github.com/soukengo/gopkg/component/transport/queue/iface"
	"github.com/soukengo/gopkg/log"
)

type kafkaQueue struct {
	iface.Consumer
	iface.Producer
}

func NewQueue(cfg *Config, logger log.Logger) (iface.Server, error) {
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
