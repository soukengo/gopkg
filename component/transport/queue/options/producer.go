package options

import (
	"github.com/soukengo/gopkg/util/codec"
	"time"
)

type ProducerOptions struct {
	encoder codec.Encoder
}

func Producer() *ProducerOptions {
	return &ProducerOptions{encoder: codec.JSON}
}

func (o *ProducerOptions) Encoder() codec.Encoder {
	if o.encoder == nil {
		o.encoder = codec.JSON
	}
	return o.encoder
}

func (o *ProducerOptions) SetEncoder(encoder codec.Encoder) *ProducerOptions {
	o.encoder = encoder
	return o
}

type DelayedProducerOptions struct {
	ProducerOptions
	delay       time.Duration
	overwritten bool
}

func Delayed() *DelayedProducerOptions {
	producer := Producer()
	return &DelayedProducerOptions{ProducerOptions: *producer}
}

func (o *DelayedProducerOptions) Delay() time.Duration {
	return o.delay
}

func (o *DelayedProducerOptions) Overwritten() bool {
	return o.overwritten
}

func (o *DelayedProducerOptions) SetDelay(delay time.Duration) *DelayedProducerOptions {
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

func DelayedRemove() *DelayedProducerRemoveOptions {
	producer := Producer()
	return &DelayedProducerRemoveOptions{ProducerOptions: *producer}
}
