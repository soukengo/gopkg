package job

import "context"

type Server interface {
	Start(context.Context) error
	Stop(context.Context) error
	Register(topic string, h Handler)
}

type Handler func(ctx context.Context, topic string, data []byte) error
