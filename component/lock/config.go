package lock

import (
	"github.com/soukengo/gopkg/infra/storage"
	"github.com/soukengo/gopkg/infra/storage/redis"
)

type Config struct {
	Redis *redis.Reference
}

func (c Config) Parse(configs *storage.Config) {
	if c.Redis != nil {
		c.Redis.Parse(configs.Redis)
	}
}
