package redis

import (
	"context"
	redisdriver "github.com/go-redis/redis/v8"
	"github.com/soukengo/gopkg/component/transport/queue/iface"
	"github.com/soukengo/gopkg/component/transport/queue/options"
	"github.com/soukengo/gopkg/infra/storage/redis"
	"github.com/soukengo/gopkg/log"
	"time"
)

type producer struct {
	cli *redis.Client
}

func NewProducer(cfg *Config, logger log.Logger) iface.Producer {
	cli := redis.NewClient(cfg.Config, logger)
	return &producer{cli: cli}
}

func (q *producer) Publish(ctx context.Context, message iface.Message, opts *options.ProducerOptions) (err error) {
	value, err := iface.EncodeMessage(message, opts)
	if err != nil {
		return
	}
	key := readyKey(message.Topic())
	score := time.Now().UnixMilli()
	members := []*redisdriver.Z{{Score: float64(score), Member: string(value)}}
	_, err = q.cli.ZAddNX(ctx, key, members...)
	return
}
