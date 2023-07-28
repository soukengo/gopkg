package dispatcher

import "context"

func Dispatch[T any](ctx context.Context, ch chan T, fn func(T)) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-ch:
			fn(msg)
		}
	}
}
