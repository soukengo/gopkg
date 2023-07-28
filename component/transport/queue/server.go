package queue

import (
	"errors"
	"github.com/soukengo/gopkg/component/transport/queue/iface"
	"github.com/soukengo/gopkg/component/transport/queue/provider/kafka"
	"github.com/soukengo/gopkg/component/transport/queue/provider/redis"
	"github.com/soukengo/gopkg/log"
)

type Server interface {
	iface.Server
}

var (
	ErrInvalidConfig = errors.New("invalid queue configuration")
)

func NewServer(cfg *Config, logger log.Logger) (Server, error) {
	if cfg.General != nil {
		return NewGeneralServer(cfg.General, logger)
	}
	return NewDelayServer(cfg.Delayed, logger)
}

func NewGeneralServer(cfg *GeneralConfig, logger log.Logger) (Server, error) {
	if cfg.Redis != nil {
		return redis.NewServer(cfg.Redis, logger), nil
	}
	if cfg.Kafka != nil {
		return kafka.NewServer(cfg.Kafka, logger)
	}
	return nil, ErrInvalidConfig
}
func NewDelayServer(cfg *DelayedConfig, logger log.Logger) (Server, error) {
	if cfg.Redis != nil {
		return redis.NewRedisDelayed(cfg.Redis, logger), nil
	}
	return nil, ErrInvalidConfig
}
