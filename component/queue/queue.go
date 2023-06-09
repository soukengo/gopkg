package queue

import (
	"errors"
	"github.com/soukengo/gopkg/component/queue/core"
	"github.com/soukengo/gopkg/component/queue/core/kafka"
	"github.com/soukengo/gopkg/component/queue/core/redis"
	"github.com/soukengo/gopkg/log"
)

type Queue = core.Queue

type Producer = core.Producer
type Consumer = core.Consumer
type Delayed = core.Delayed

type DelayedProducer = core.DelayedProducer

type Topic = core.Topic
type Value = core.Value

var (
	ErrInvalidConfig = errors.New("invalid queue configuration")
)

func NewProducer(cfg *GeneralConfig, logger log.Logger) (Producer, error) {
	if cfg.Redis != nil {
		return redis.NewProducer(cfg.Redis, logger), nil
	}
	if cfg.Kafka != nil {
		return kafka.NewProducer(cfg.Kafka, logger)
	}

	return nil, ErrInvalidConfig
}

func NewConsumer(cfg *Config, logger log.Logger) (Consumer, error) {
	if cfg.General != nil {
		return NewGeneralConsumer(cfg.General, logger)
	}
	return NewDelayedConsumer(cfg.Delayed, logger)
}
func NewGeneralConsumer(cfg *GeneralConfig, logger log.Logger) (Consumer, error) {
	if cfg.Redis != nil {
		return redis.NewConsumer(cfg.Redis, logger), nil
	}
	if cfg.Kafka != nil {
		return kafka.NewConsumer(cfg.Kafka, logger)
	}
	return nil, ErrInvalidConfig
}
func NewDelayedConsumer(cfg *DelayedConfig, logger log.Logger) (Consumer, error) {
	if cfg.Redis != nil {
		return redis.NewConsumer(cfg.Redis, logger), nil
	}
	return nil, ErrInvalidConfig
}

func NewDelayProducer(cfg *DelayedConfig, logger log.Logger) (DelayedProducer, error) {
	if cfg.Redis != nil {
		return redis.NewDelayProducer(cfg.Redis, logger), nil
	}
	return nil, ErrInvalidConfig
}
