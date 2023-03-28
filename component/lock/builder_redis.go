package lock

import (
	"github.com/soukengo/gopkg/component/database/redis"
	"github.com/soukengo/gopkg/log"
)

type redisBuilder struct {
	logger log.Logger
	cli    *redis.Client
}

func NewRedisBuilder(logger log.Logger, cli *redis.Client) Builder {
	return &redisBuilder{cli: cli, logger: logger}
}

func (b *redisBuilder) Build(key Key, parts ...string) DistributedLock {
	return NewRedisDistributedLock(b.logger, b.cli, key.GenKey(parts...))
}
