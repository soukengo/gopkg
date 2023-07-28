package codec

import "github.com/soukengo/gopkg/util/encoding/proto"

type protobufCodec struct {
}

func NewProtobufCodec() Codec {
	return &protobufCodec{}
}

func (c *protobufCodec) Encode(v any) ([]byte, error) {
	return proto.Marshal(v)
}

func (c *protobufCodec) Decode(data []byte, v any) error {
	return proto.Unmarshal(data, v)
}
