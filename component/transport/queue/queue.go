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

type Handler = iface.Handler

type Action = iface.Action

func Handle[T any](fn func(context.Context, *T) iface.Action) *iface.Handler {
	return HandleWithOptions(fn, options.Consumer())
}

func HandleWithOptions[T any](fn func(context.Context, *T) iface.Action, opts *options.ConsumerOptions) *iface.Handler {
	return iface.HandleWithOptions(fn, opts)
}

func NewRawMessage(topic Topic, data any) *iface.RawMessage {
	return iface.NewRawMessage(topic, data)
}

func NewBytesMessage(topic Topic, data []byte) *iface.BytesMessage {
	return iface.NewBytesMessage(topic, data)
}

func ConsumerOptions() *options.ConsumerOptions {
	return options.Consumer()
}
func ProducerOptions() *options.ProducerOptions {
	return options.Producer()
}
func DelayedOptions() *options.DelayedProducerOptions {
	return options.Delayed()
}
func DelayedRemoveOptions() *options.DelayedProducerRemoveOptions {
	return options.DelayedRemove()
}
