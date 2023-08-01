package event

type Key string

type Event interface {
	Key() Key
}

//func Wrap[T Event](fn func(context.Context, *T)) Handler {
//	return func(ctx context.Context, data any) {
//		fn(ctx, data.(*T))
//		return
//	}
//}
