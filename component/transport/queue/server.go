package queue

import (
	"errors"
	"github.com/soukengo/gopkg/component/transport/queue/iface"
	"github.com/soukengo/gopkg/component/transport/queue/provider/kafka"
	"github.com/soukengo/gopkg/component/transport/queue/provider/memory"
	"github.com/soukengo/gopkg/component/transport/queue/provider/redis"
	"github.com/soukengo/gopkg/log"
)

type Server interface {
	iface.Server
}

var (
	ErrInvalidConfig = errors.New("invalid queue configuration")
)

func NewGeneralServer(cfg *GeneralConfig, logger log.Logger) (Server, error) {
	if cfg.Redis != nil {
		return redis.NewServer(cfg.Redis, logger), nil
	}
	if cfg.Kafka != nil {
		return kafka.NewServer(cfg.Kafka, logger)
	}
	if cfg.Memory != nil {
		return memory.NewServer(cfg.Memory, logger), nil
	}
	return nil, ErrInvalidConfig
}
func NewDelayServer(cfg *DelayedConfig, topics []iface.Topic, logger log.Logger) (Delayed, error) {
	if cfg.Redis != nil {
		return redis.NewRedisDelayed(cfg.Redis, topics, logger), nil
	}
	return nil, ErrInvalidConfig
}
