package job

import (
	"context"
	"github.com/soukengo/gopkg/component/queue"
)

type queueServer struct {
	queue queue.Queue
}

func NewQueueServer(queue queue.Queue) Server {
	return &queueServer{queue: queue}
}

func (s *queueServer) Start(ctx context.Context) error {
	return s.queue.Start()
}

func (s *queueServer) Stop(ctx context.Context) error {
	return s.queue.Stop()
}

func (s *queueServer) Register(topic string, h Handler) {
	s.queue.Subscribe(queue.Topic(topic), func(topic queue.Topic, value string) {
		ctx := context.Background()
		_ = h(ctx, string(topic), []byte(value))
	})
}
