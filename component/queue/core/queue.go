package core

import (
	"context"
)

type HandlerFunc func(topic Topic, value Value) (err error)

type Queue interface {
	Producer
	Consumer
}

type Value []byte

type Producer interface {
	// Submit 往队列提交数据
	Submit(ctx context.Context, topic Topic, values []Value) error
}

type Consumer interface {
	// Start 启动队列
	Start() error
	// Stop 停止队列
	Stop() error
	// Subscribe 订阅TOPIC
	Subscribe(topic Topic, handler HandlerFunc)
}

type Delayed interface {
	Consumer
	DelayedProducer
}

type DelayedProducer interface {
	Producer
	// SubmitDelay 提交延时数据
	SubmitDelay(ctx context.Context, topic Topic, values []Value, delayMills int64, overwritten bool) error
	// RemoveDelay 从延时队列移除任务
	RemoveDelay(ctx context.Context, topic Topic, values []Value) (deleted bool, err error)
	Start() error
	Stop() error
}
