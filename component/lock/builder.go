package lock

import (
	"github.com/soukengo/gopkg/log"
)

type Builder interface {
	Build(key Key, parts ...string) DistributedLock
}

func NewBuilder(cfg *Config, logger log.Logger) Builder {
	return NewRedisBuilder(cfg.Redis.Config, logger)
}
