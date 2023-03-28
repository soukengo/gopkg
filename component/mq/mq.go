package mq

import (
	"context"
)

type Producer interface {
	Send(ctx context.Context, topic string, data []byte) error
}

type Consumer interface {
	Start() error
	Close() error
}

type Handler interface {
	OnReceived(topic string, data []byte) error
}
