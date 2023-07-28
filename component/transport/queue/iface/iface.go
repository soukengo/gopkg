package iface

import (
	"context"
	"github.com/soukengo/gopkg/component/transport/queue/options"
)

type HandlerFunc func(context.Context, *BytesMessage) (err error)

type Handler struct {
	Func    HandlerFunc
	Options *options.ConsumerOptions
}

type Producer interface {
	Publish(ctx context.Context, message Message, opts *options.ProducerOptions) error
}

type Consumer interface {
	// Start 启动队列
	Start() error
	// Stop 停止队列
	Stop() error
	// Subscribe 订阅TOPIC
	Subscribe(topic Topic, handler *Handler)
}

type Server interface {
	Consumer
	Producer
}

func WrapHandlerWithOptions[T any](fn func(context.Context, *T) error, opts *options.ConsumerOptions) *Handler {
	if opts == nil {
		opts = options.Consumer()
	}
	return &Handler{
		Func: func(ctx context.Context, m *BytesMessage) (err error) {
			var req = new(T)
			err = opts.Decoder().Decode(m.Bytes(), req)
			if err != nil {
				return
			}
			err = fn(ctx, req)
			return
		}, Options: opts,
	}
}
func EncodeMessage(message Message, opts *options.ProducerOptions) (value []byte, err error) {
	if opts == nil {
		opts = options.Producer()
	}
	if bm, ok := message.(*BytesMessage); ok {
		value = bm.Bytes()
		return
	}
	return opts.Encoder().Encode(message.Data())
}
