package http

import (
	"time"
)

type Option func(opts *options)

type options struct {
	timeout   time.Duration
	baseUrl   string
	basicAuth *basicAuth
	proxyURL  string
}

type basicAuth struct {
	username string
	password string
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
func WithBaseUrl(baseUrl string) Option {
	return func(o *options) { o.baseUrl = baseUrl }
}
func WithBasicAuth(username, password string) Option {
	return func(o *options) { o.basicAuth = &basicAuth{username: username, password: password} }
}

func WithProxy(proxyURL string) Option {
	return func(o *options) { o.proxyURL = proxyURL }
}
