package queue

import (
	"context"
	"github.com/soukengo/gopkg/component/transport/queue/iface"
	"github.com/soukengo/gopkg/component/transport/queue/options"
)

type Producer = iface.Producer
type Consumer = iface.Consumer
type Delayed = iface.Delayed

type DelayedProducer = iface.DelayedProducer

type Topic = iface.Topic
type Message = iface.Message

func WrapHandler[T any](fn func(context.Context, *T) error) *iface.Handler {
	return WrapHandlerWithOptions(fn, options.Consumer())
}

func WrapHandlerWithOptions[T any](fn func(context.Context, *T) error, opts *options.ConsumerOptions) *iface.Handler {
	return iface.WrapHandlerWithOptions(fn, opts)
}
