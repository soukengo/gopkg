package mongo

import "github.com/soukengo/gopkg/log"

type Option func(opts *options)

type options struct {
	logger log.Logger
}

func (o *options) apply(opts ...Option) *options {
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func WithLogger(logger log.Logger) Option {
	return func(o *options) { o.logger = logger }
}
