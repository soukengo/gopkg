package iface

import (
	"context"
	"github.com/soukengo/gopkg/component/transport/queue/options"
)

type Producer interface {
	Publish(ctx context.Context, message Message, opts *options.ProducerOptions) error
}

type Consumer interface {
	// Start 启动队列
	Start(context.Context) error
	// Stop 停止队列
	Stop(context.Context) error
	// Subscribe 订阅TOPIC
	Subscribe(topic Topic, handler *Handler)
}

type Server interface {
	Consumer
	Producer
}

func HandleWithOptions[T any](fn func(context.Context, *T) Action, opts *options.ConsumerOptions) *Handler {
	if opts == nil {
		opts = options.Consumer()
	}
	return &Handler{
		Func: func(ctx context.Context, m *BytesMessage) Action {
			var req = new(T)
			err := opts.Decoder().Decode(m.Bytes(), req)
			if err != nil {
				return ReconsumeLater
			}
			return fn(ctx, req)
		},
		Options: opts,
	}
}
