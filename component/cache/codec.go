package cache

import "github.com/soukengo/gopkg/util/encoding/json"

type Codec interface {
	Encode(any) ([]byte, error)
	Decode([]byte, any) error
}

type jsonCodec struct {
}

func NewJsonCodec() Codec {
	return &jsonCodec{}
}

func (c *jsonCodec) Encode(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (c *jsonCodec) Decode(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
