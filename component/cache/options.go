package cache

import "time"

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
	codec       Codec
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

func WithCodec(codec Codec) Option {
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
	return &Options{expire: defaultExpire, emptyExpire: emptyExpire, codec: NewJsonCodec()}
}

func (o *Options) Expire() time.Duration {
	return o.expire
}

func (o *Options) EmptyExpire() time.Duration {
	return o.emptyExpire
}

func (o *Options) Codec() Codec {
	if o.codec == nil {
		o.codec = NewJsonCodec()
	}
	return o.codec
}

func (o *Options) apply(opts ...Option) *Options {
	for _, opt := range opts {
		opt(o)
	}
	return o
}
