package core

import (
	"github.com/soukengo/gopkg/util/codec"
	"time"
)

const (
	// default cache expiration time
	defaultExpire = 0
	// empty cache expiration time
	emptyExpire = time.Minute * 5
)

type Option func(opts *Options)

type Options struct {
	expire      time.Duration
	emptyExpire time.Duration
	codec       codec.Codec
}

func Apply(category *Category, opts ...Option) *Options {
	return options(category).apply(opts...)
}

func WithExpire(expire time.Duration) Option {
	return func(o *Options) { o.expire = expire }
}

func WithEmptyExpire(emptyExpire time.Duration) Option {
	return func(o *Options) { o.emptyExpire = emptyExpire }
}

func WithCodec(codec codec.Codec) Option {
	return func(o *Options) { o.codec = codec }
}

func options(category *Category) *Options {
	opts := defaultOptions()
	if category != nil {
		opts = opts.apply(category.Options()...)
	}
	return opts
}

func defaultOptions() *Options {
	return &Options{expire: defaultExpire, emptyExpire: emptyExpire, codec: codec.JSON}
}

func (o *Options) Expire() time.Duration {
	return o.expire
}

func (o *Options) EmptyExpire() time.Duration {
	return o.emptyExpire
}

func (o *Options) Codec() codec.Codec {
	if o.codec == nil {
		o.codec = codec.NewJsonCodec()
	}
	return o.codec
}

func (o *Options) apply(opts ...Option) *Options {
	for _, opt := range opts {
		opt(o)
	}
	return o
}
