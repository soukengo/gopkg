package redis

import (
	"context"
	"github.com/soukengo/gopkg/component/transport/queue/iface"
	"github.com/soukengo/gopkg/component/transport/queue/options"
	"github.com/soukengo/gopkg/log"
)

type redisDelayed struct {
	consumer iface.Consumer
	producer iface.DelayedProducer
}

func NewRedisDelayed(cfg *Config, logger log.Logger) iface.Delayed {
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

func (q *redisDelayed) Subscribe(topic iface.Topic, handler *iface.Handler) {
	q.consumer.Subscribe(topic, handler)
}

func (q *redisDelayed) Publish(ctx context.Context, message iface.Message, opts *options.ProducerOptions) error {
	return q.producer.Publish(ctx, message, opts)
}

func (q *redisDelayed) PublishDelay(ctx context.Context, message iface.Message, opts *options.DelayedProducerOptions) error {
	return q.producer.PublishDelay(ctx, message, opts)
}

func (q *redisDelayed) RemoveDelay(ctx context.Context, message iface.Message, opts *options.DelayedProducerRemoveOptions) (deleted bool, err error) {
	return q.producer.RemoveDelay(ctx, message, opts)
}
