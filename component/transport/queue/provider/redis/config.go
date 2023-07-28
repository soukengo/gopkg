package redis

import (
	"github.com/soukengo/gopkg/infra/storage/redis"
)

type Config struct {
	redis.Reference `mapstructure:",squash"`
	Consumer        *ConsumerConfig
}

type ConsumerConfig struct {
	Workers int
}
