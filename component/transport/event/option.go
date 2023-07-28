package event

import "github.com/soukengo/gopkg/util/codec"

type options struct {
	codec codec.Codec
}

type Option func(opts *options)

func defaultOptions() *options {
	return &options{codec: codec.NewJsonCodec()}
}

func (o *options) apply(opts ...Option) *options {
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func Codec(codec codec.Codec) Option {
	return func(o *options) { o.codec = codec }
}
