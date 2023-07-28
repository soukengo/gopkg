package iface

import (
	"context"
	"github.com/soukengo/gopkg/component/transport/queue/options"
)

type Delayed interface {
	Consumer
	DelayedProducer
}

type DelayedProducer interface {
	Producer
	// PublishDelay 提交延时数据
	PublishDelay(ctx context.Context, message Message, opts *options.DelayedProducerOptions) error
	// RemoveDelay 从延时队列移除任务
	RemoveDelay(ctx context.Context, message Message, opts *options.DelayedProducerRemoveOptions) (deleted bool, err error)
	Start(context.Context) error
	Stop(context.Context) error
}
