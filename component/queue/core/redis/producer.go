package redis

import (
	"context"
	redisdriver "github.com/go-redis/redis/v8"
	"github.com/soukengo/gopkg/component/queue/core"
	"github.com/soukengo/gopkg/infra/storage/redis"
	"github.com/soukengo/gopkg/log"
	"time"
)

type producer struct {
	cli *redis.Client
}

func NewProducer(cfg *Config, logger log.Logger) core.Producer {
	cli := redis.NewClient(cfg.Config, logger)
	return &producer{cli: cli}
}

func (q *producer) Submit(ctx context.Context, topic core.Topic, values []core.Value) (err error) {
	key := readyKey(topic)
	score := time.Now().UnixMilli()
	members := make([]*redisdriver.Z, len(values))
	for i, value := range values {
		members[i] = &redisdriver.Z{Score: float64(score), Member: string(value)}
	}
	_, err = q.cli.ZAddNX(ctx, key, members...)
	return
}
