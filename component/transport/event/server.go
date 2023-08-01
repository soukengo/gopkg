package event

import (
	"context"
	"github.com/soukengo/gopkg/component/transport/queue"
	"github.com/soukengo/gopkg/component/transport/queue/iface"
	"github.com/soukengo/gopkg/component/transport/queue/options"
)

type Server interface {
	queue.Server
}

type server struct {
	srv queue.Server
}

func NewServer() Server {
	return &server{}
}

func (s *server) Start(ctx context.Context) error {
	return nil
}

func (s *server) Stop(ctx context.Context) error {
	return nil
}

func (s *server) Subscribe(topic iface.Topic, handler *iface.Handler) {
	s.srv.Subscribe(topic, handler)
}

func (s *server) Publish(ctx context.Context, message iface.Message, opts *options.ProducerOptions) error {
	return s.srv.Publish(ctx, message, opts)
}
