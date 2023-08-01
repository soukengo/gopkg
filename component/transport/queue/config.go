package queue

import (
	"github.com/soukengo/gopkg/component/transport/queue/provider/kafka"
	"github.com/soukengo/gopkg/component/transport/queue/provider/memory"
	"github.com/soukengo/gopkg/component/transport/queue/provider/redis"
	"github.com/soukengo/gopkg/infra/storage"
)

//type Config struct {
//	General *GeneralConfig
//	Delayed *DelayedConfig
//}

// GeneralConfig 普通队列
type GeneralConfig struct {
	Redis  *redis.Config
	Kafka  *kafka.Config
	Memory *memory.Config
}

// DelayedConfig 延时队列
type DelayedConfig struct {
	Redis *redis.Config
}

func (c *GeneralConfig) Parse(configs *storage.Config) {
	if c.Redis != nil {
		c.Redis.Parse(configs.Redis)
	}
	if c.Kafka != nil {
		c.Kafka.Parse(configs.Kafka)
	}
}
func (c *DelayedConfig) Parse(configs *storage.Config) {
	if c.Redis != nil {
		c.Redis.Parse(configs.Redis)
	}
}
