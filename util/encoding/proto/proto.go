// Package proto defines the protobuf codec. Importing this package will
// register the codec.
package proto

import (
	"errors"
	"github.com/go-kratos/kratos/v2/encoding"

	"google.golang.org/protobuf/proto"
)

var (
	ErrInvalidMessage = errors.New("invalid message")
)

// Name is the name registered for the proto compressor.
const Name = "x-protobuf"

func init() {
	encoding.RegisterCodec(codec{})
}

// codec is a Codec implementation with protobuf. It is the default codec for Transport.
type codec struct{}

func (codec) Marshal(v any) ([]byte, error) {
	return Marshal(v)
}

func (codec) Unmarshal(data []byte, v any) error {
	return Unmarshal(data, v)
}

func (codec) Name() string {
	return Name
}

func Marshal(v any) ([]byte, error) {
	m, ok := v.(proto.Message)
	if !ok {
		return nil, ErrInvalidMessage
	}
	return proto.Marshal(m)
}

func Unmarshal(data []byte, v any) error {
	m, ok := v.(proto.Message)
	if !ok {
		return ErrInvalidMessage
	}
	return proto.Unmarshal(data, m)
}
