package redis

import (
	"github.com/soukengo/gopkg/component/queue/core"
	"github.com/soukengo/gopkg/infra/storage/redis"
)

type Config struct {
	redis.Reference
	Consumer *ConsumerConfig
}

type ConsumerConfig struct {
	Topics  []core.Topic
	Workers int
}
