package event

import "context"

type DistributedHandler func(ctx context.Context, data []byte) error

func WrapDistributed[T Event](fn func(context.Context, *T) error, opts ...Option) DistributedHandler {
	opt := defaultOptions().apply(opts...)
	return func(ctx context.Context, data []byte) (err error) {
		var req = new(T)
		err = opt.codec.Decode(data, req)
		if err != nil {
			return
		}
		err = fn(ctx, req)
		return
	}
}
