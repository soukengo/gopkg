package redis

import (
	"context"
	"github.com/soukengo/gopkg/component/queue/core"
	"github.com/soukengo/gopkg/log"
)

type redisDelayed struct {
	consumer core.Consumer
	producer core.DelayedProducer
}

func NewRedisDelayed(cfg *Config, logger log.Logger) core.Delayed {
	ins := &redisDelayed{
		consumer: NewConsumer(cfg, logger),
		producer: NewDelayProducer(cfg, logger),
	}
	return ins
}

func (q *redisDelayed) Start() (err error) {
	err = q.consumer.Start()
	if err != nil {
		return
	}
	err = q.producer.Start()
	return
}
func (q *redisDelayed) Stop() (err error) {
	err = q.consumer.Stop()
	if err != nil {
		return
	}
	err = q.producer.Stop()
	return
}

func (q *redisDelayed) Subscribe(topic core.Topic, handler core.HandlerFunc) {
	q.consumer.Subscribe(topic, handler)
}

func (q *redisDelayed) Submit(ctx context.Context, topic core.Topic, values []core.Value) error {
	return q.producer.Submit(ctx, topic, values)
}

func (q *redisDelayed) SubmitDelay(ctx context.Context, topic core.Topic, values []core.Value, delayMills int64, overwritten bool) error {
	return q.producer.SubmitDelay(ctx, topic, values, delayMills, overwritten)
}

func (q *redisDelayed) RemoveDelay(ctx context.Context, topic core.Topic, values []core.Value) (deleted bool, err error) {
	return q.producer.RemoveDelay(ctx, topic, values)
}
