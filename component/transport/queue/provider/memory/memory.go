package memory

import (
	"context"
	"github.com/soukengo/gopkg/component/transport/queue/iface"
	"github.com/soukengo/gopkg/component/transport/queue/options"
	"github.com/soukengo/gopkg/util/dispatcher"
	"github.com/soukengo/gopkg/util/runtimes"
	"sync"
)

type Queue struct {
	events  map[iface.Topic][]*iface.Handler
	lock    sync.RWMutex
	values  chan *iface.BytesMessage
	ctx     context.Context
	cancel  context.CancelFunc
	started sync.Once
}

func NewQueue(queueSize uint) *Queue {
	ctx, cancel := context.WithCancel(context.Background())
	return &Queue{
		events: make(map[iface.Topic][]*iface.Handler),
		values: make(chan *iface.BytesMessage, queueSize),
		ctx:    ctx,
		cancel: cancel,
	}
}

func (q *Queue) Subscribe(topic iface.Topic, handler *iface.Handler) {
	q.lock.Lock()
	listeners, ok := q.events[topic]
	if !ok {
		listeners = []*iface.Handler{}
	}
	listeners = append(listeners, handler)
	q.events[topic] = listeners
	q.lock.Unlock()
}

func (q *Queue) Publish(ctx context.Context, message iface.Message, opts *options.ProducerOptions) (err error) {
	if opts != nil {
		opts = options.Producer()
	}
	value, err := message.Encode(opts.Encoder())
	if err != nil {
		return
	}
	q.values <- iface.NewBytesMessage(message.Topic(), value)
	return
}

func (q *Queue) Start(context.Context) (err error) {
	q.started.Do(func() {
		runtimes.Async(q.dispatch)
	})
	return nil
}

func (q *Queue) Stop(context.Context) error {
	q.cancel()
	return nil
}

func (q *Queue) dispatch() {
	dispatcher.Dispatch(q.ctx, q.values, func(msg *iface.BytesMessage) {
		topic := msg.Topic()
		q.lock.RLock()
		handlers, ok := q.events[topic]
		if !ok {
			return
		}
		q.lock.RUnlock()
		q.trigger(msg, handlers)
	})
}

func (q *Queue) trigger(message *iface.BytesMessage, handlers []*iface.Handler) {
	for _, item := range handlers {
		var handler = item
		handler.Process(message, func(action iface.Action) {
			// do nothing
		})
	}
}
