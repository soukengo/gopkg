package queue

import (
	"github.com/soukengo/gopkg/component/transport/queue/provider/kafka"
	"github.com/soukengo/gopkg/component/transport/queue/provider/redis"
	"github.com/soukengo/gopkg/infra/storage"
)

type Config struct {
	General *GeneralConfig
	Delayed *DelayedConfig
}

// GeneralConfig 普通队列
type GeneralConfig struct {
	Redis *redis.Config
	Kafka *kafka.Config
}

// DelayedConfig 延时队列
type DelayedConfig struct {
	Redis *redis.Config
}

func (c *Config) Parse(configs *storage.Config) {
	if c.General != nil {
		if c.General.Redis != nil {
			c.General.Redis.Parse(configs.Redis)
		}
		if c.General.Kafka != nil {
			c.General.Kafka.Parse(configs.Kafka)
		}
	}
	if c.Delayed != nil {
		if c.Delayed.Redis != nil {
			c.Delayed.Redis.Parse(configs.Redis)
		}
	}
}
