package event

import (
	"context"
)

type Key string

type Event interface {
	Key() Key
}

type Handler func(ctx context.Context, data any)

type Listener struct {
	Mode    Mode
	Handler Handler
}

type Consumer interface {
	Subscribe(key Key, handler Listener)
}

type Producer interface {
	Publish(ctx context.Context, ev Event) error
}

type Dispatcher interface {
	Consumer
	Producer
}

func Wrap[T Event](fn func(context.Context, *T)) Handler {
	return func(ctx context.Context, data any) {
		fn(ctx, data.(*T))
		return
	}
}
