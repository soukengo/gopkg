package job

import (
	"context"
	"github.com/soukengo/gopkg/log"
)

type Server interface {
	Start(context.Context) error
	Stop(context.Context) error
	Register(topic string, h Handler)
}

type Handler func(ctx context.Context, topic string, data []byte) error

func NewServer(cfg *Config, logger log.Logger) (Server, error) {
	return NewQueueServer(cfg.Queue, logger)
}
