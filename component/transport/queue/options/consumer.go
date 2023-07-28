package options

import "github.com/soukengo/gopkg/util/codec"

type ConsumerOptions struct {
	decoder codec.Decoder
	mode    Mode
}

func Consumer() *ConsumerOptions {
	return &ConsumerOptions{decoder: codec.JSON}
}

func (o *ConsumerOptions) Decoder() codec.Decoder {
	return o.decoder
}

func (o *ConsumerOptions) SetDecoder(decoder codec.Decoder) *ConsumerOptions {
	o.decoder = decoder
	return o
}

func (o *ConsumerOptions) Mode() Mode {
	return o.mode
}

func (o *ConsumerOptions) SetMode(mode Mode) *ConsumerOptions {
	o.mode = mode
	return o
}
