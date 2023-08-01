package event

import (
	"context"
	"github.com/soukengo/gopkg/component/transport/queue"
	"github.com/soukengo/gopkg/component/transport/queue/iface"
	"github.com/soukengo/gopkg/component/transport/queue/provider/memory"
	"github.com/soukengo/gopkg/log"
)

type Server interface {
	// Start 启动队列
	Start(context.Context) error
	// Stop 停止队列
	Stop(context.Context) error
	// Subscribe 订阅事件
	Subscribe(key Key, handler *Handler)
	// Publish 发布事件
	Publish(ctx context.Context, ev *Event, opts *ProducerOptions)
}

type DistributedServer interface {
	Server
}

type server struct {
	queue  queue.Server
	logger log.Logger
}

func NewMemoryServer(cfg *Memory, logger log.Logger) (s Server) {
	return &server{queue: memory.NewServer(cfg, logger), logger: logger}
}

func NewDistributedServer(cfg *Distributed, logger log.Logger) (s DistributedServer, err error) {
	q, err := queue.NewGeneralServer(
		&queue.GeneralConfig{
			Redis: cfg.Redis,
			Kafka: cfg.Kafka,
		}, logger)
	if err != nil {
		return
	}
	return &server{queue: q}, nil
}

func (s *server) Start(ctx context.Context) error {
	return s.queue.Start(ctx)
}

func (s *server) Stop(ctx context.Context) error {
	return s.queue.Stop(ctx)
}

func (s *server) Subscribe(key Key, handler *Handler) {
	s.queue.Subscribe(key.Topic(), &queue.Handler{
		Options: handler.Options,
		Func: func(ctx context.Context, message *iface.BytesMessage) iface.Action {
			handler.Func(ctx, key, message.Bytes())
			return iface.CommitMessage
		},
	})
}

func (s *server) Publish(ctx context.Context, ev *Event, opts *ProducerOptions) {
	err := s.queue.Publish(ctx, iface.NewRawMessage(ev.Key.Topic(), ev.Value), opts)
	if err != nil {
		s.logger.Helper().Error("Event Publish error", log.Pairs{"ev": ev, "err": err})
	}
}
