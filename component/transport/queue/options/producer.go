package options

import "github.com/soukengo/gopkg/util/codec"

type ProducerOptions struct {
	encoder codec.Encoder
}

func Producer() *ProducerOptions {
	return &ProducerOptions{encoder: codec.JSON}
}

func (o *ProducerOptions) Encoder() codec.Encoder {
	return o.encoder
}

func (o *ProducerOptions) SetEncoder(encoder codec.Encoder) *ProducerOptions {
	o.encoder = encoder
	return o
}

type DelayedProducerOptions struct {
	ProducerOptions
	delay       int64
	overwritten bool
}

func (o *DelayedProducerOptions) Delay() int64 {
	return o.delay
}

func (o *DelayedProducerOptions) Overwritten() bool {
	return o.overwritten
}

func (o *DelayedProducerOptions) SetDelay(delay int64) *DelayedProducerOptions {
	o.delay = delay
	return o
}

func (o *DelayedProducerOptions) SetOverwritten(overwritten bool) *DelayedProducerOptions {
	o.overwritten = overwritten
	return o
}

type DelayedProducerRemoveOptions struct {
	ProducerOptions
}
