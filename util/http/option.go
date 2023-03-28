package http

import (
	"time"
)

type Option func(opts *options)

type options struct {
	timeout time.Duration
}

func (o *options) apply(opts ...Option) *options {
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func WithTimeout(timeout time.Duration) Option {
	return func(o *options) { o.timeout = timeout }
}
