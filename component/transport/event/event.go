package event

import (
	"context"
	"github.com/soukengo/gopkg/component/transport/queue"
	"github.com/soukengo/gopkg/component/transport/queue/options"
	"github.com/soukengo/gopkg/log"
)

var (
	TopicPrefix = "Event_"
)

type Key string

func (k Key) Topic() queue.Topic {
	return queue.Topic(TopicPrefix + string(k))
}

type Event struct {
	Key   Key
	Value any
}

func NewEvent(key Key, value any) *Event {
	return &Event{Key: key, Value: value}
}

type Handler struct {
	Func    HandlerFunc
	Options *ConsumerOptions
}

type HandlerFunc func(context.Context, Key, []byte)
type ConsumerOptions = options.ConsumerOptions
type ProducerOptions = options.ProducerOptions

func Consumer() *ConsumerOptions {
	return options.Consumer()
}
func Producer() *ProducerOptions {
	return options.Producer()
}

func Handle[T any](fn func(context.Context, *T), opts *ConsumerOptions) *Handler {
	if opts == nil {
		opts = options.Consumer()
	}
	return &Handler{
		Func: func(ctx context.Context, key Key, data []byte) {
			var req = new(T)
			err := opts.Decoder().Decode(data, req)
			if err != nil {
				log.Error("Event Decode error", log.Pairs{"err": err, "key": key})
			}
			fn(ctx, req)
		},
		Options: opts,
	}
}
