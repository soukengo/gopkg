package lock

import (
	"github.com/soukengo/gopkg/infra/storage/redis"
	"github.com/soukengo/gopkg/log"
)

type redisBuilder struct {
	logger log.Logger
	cli    *redis.Client
}

func NewRedisBuilder(cfg *redis.Config, logger log.Logger) Builder {
	return &redisBuilder{cli: redis.NewClient(cfg, logger), logger: logger}
}

func (b *redisBuilder) Build(key Key, parts ...string) DistributedLock {
	return NewRedisDistributedLock(b.logger, b.cli, key.GenKey(parts...))
}
