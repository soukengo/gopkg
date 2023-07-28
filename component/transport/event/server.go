package event

import (
	"context"
)

type Server interface {
	Start(context.Context) error
	Stop(context.Context) error
	Consumer
}

type server struct {
	consumer Consumer
}

func NewServer(consumer Consumer) Server {
	return &server{consumer: consumer}
}

func (s *server) Start(ctx context.Context) error {
	return nil
}

func (s *server) Stop(ctx context.Context) error {
	return nil
}

func (s *server) Subscribe(key Key, handler Listener) {
	s.consumer.Subscribe(key, handler)
}
