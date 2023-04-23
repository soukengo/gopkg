package cache

import (
	"github.com/soukengo/gopkg/component/cache/core"
	rediscache "github.com/soukengo/gopkg/component/cache/core/redis"
	"github.com/soukengo/gopkg/infra/storage/redis"
	"github.com/soukengo/gopkg/log"
)

type Manager struct {
	cfg      *Config
	redisCli *redis.Client
}

func NewManager(cfg *Config, logger log.Logger) *Manager {
	cli := redis.NewClient(cfg.Redis.Config, logger)
	return &Manager{cfg: cfg, redisCli: cli}
}

func NewValueCache[T any](mgr *Manager, category *Category, opts ...Option) ValueCache[T] {
	return rediscache.NewValueCache[T](mgr.redisCli, (*core.Category)(category), convertOptions(opts)...)
}
func NewHashCache[T any](mgr *Manager, category *Category, opts ...Option) HashCache[T] {
	return rediscache.NewHashCache[T](mgr.redisCli, (*core.Category)(category), convertOptions(opts)...)
}
func NewSetCache[T any](mgr *Manager, category *Category, opts ...Option) SetCache[T] {
	return rediscache.NewSetCache[T](mgr.redisCli, (*core.Category)(category), convertOptions(opts)...)
}
func NewSortedSetCache[T any](mgr *Manager, category *Category, opts ...Option) SortedSetCache[T] {
	return rediscache.NewSortedSetCache[T](mgr.redisCli, (*core.Category)(category), convertOptions(opts)...)
}

func convertOptions(opts []Option) []core.Option {
	out := make([]core.Option, len(opts))
	for i, opt := range opts {
		out[i] = core.Option(opt)
	}
	return out
}
