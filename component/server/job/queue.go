package job

import (
	"context"
	"github.com/soukengo/gopkg/component/queue"
	"github.com/soukengo/gopkg/log"
)

type queueConsumerServer struct {
	logger   log.Logger
	consumer queue.Consumer
}

func NewQueueServer(cfg *queue.Config, logger log.Logger) (Server, error) {
	consumer, err := queue.NewConsumer(cfg, logger)
	if err != nil {
		return nil, err
	}
	return &queueConsumerServer{consumer: consumer, logger: logger}, nil
}

func (s *queueConsumerServer) Start(ctx context.Context) (err error) {
	err = s.consumer.Start()
	if err != nil {
		return
	}
	log.WithContext(ctx).Infof("queue server started")
	return
}

func (s *queueConsumerServer) Stop(ctx context.Context) error {
	return s.consumer.Stop()
}

func (s *queueConsumerServer) Register(topic string, h Handler) {
	s.consumer.Subscribe(queue.Topic(topic), func(topic queue.Topic, value queue.Value) error {
		ctx := context.Background()
		return h(ctx, string(topic), value)
	})
}
