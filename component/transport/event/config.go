package event

import (
	"github.com/soukengo/gopkg/component/transport/queue/provider/kafka"
	"github.com/soukengo/gopkg/component/transport/queue/provider/memory"
	"github.com/soukengo/gopkg/component/transport/queue/provider/redis"
	"github.com/soukengo/gopkg/infra/storage"
)

type Config struct {
	Memory      *Memory
	Distributed *Distributed
}

type Memory = memory.Config

type Distributed struct {
	Redis *redis.Config
	Kafka *kafka.Config
}

func (c *Config) Parse(configs *storage.Config) {
	if c.Distributed != nil {
		c.Distributed.Parse(configs)
	}
}
func (c *Distributed) Parse(configs *storage.Config) {
	if c.Kafka != nil {
		c.Kafka.Parse(configs.Kafka)
	}
	if c.Redis != nil {
		c.Redis.Parse(configs.Redis)
	}
}
