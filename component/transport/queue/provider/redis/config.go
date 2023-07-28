package redis

import (
	"github.com/soukengo/gopkg/component/transport/queue/iface"
	"github.com/soukengo/gopkg/infra/storage/redis"
)

type Config struct {
	redis.Reference `mapstructure:",squash"`
	Consumer        *ConsumerConfig
}

type ConsumerConfig struct {
	Topics  []iface.Topic
	Workers int
}
