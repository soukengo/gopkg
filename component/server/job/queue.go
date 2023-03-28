package job

import (
	"context"
	"github.com/soukengo/gopkg/component/queue"
)

type delayQueueServer struct {
	queue queue.Delayed
}

func NewDelayQueueServer(queue queue.Delayed) Server {
	return &delayQueueServer{queue: queue}
}

func (s *delayQueueServer) Start(ctx context.Context) error {
	return s.queue.Start()
}

func (s *delayQueueServer) Stop(ctx context.Context) error {
	return s.queue.Stop()
}

func (s *delayQueueServer) Register(topic string, h Handler) {
	s.queue.Subscribe(queue.Topic(topic), func(topic queue.Topic, value string) {
		ctx := context.Background()
		_ = h(ctx, string(topic), []byte(value))
	})
}
