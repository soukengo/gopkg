package codec

import "github.com/soukengo/gopkg/util/encoding/json"

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
