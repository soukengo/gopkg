package json

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var json = jsoniter.Config{
	EscapeHTML:             true,
	ValidateJsonRawMessage: true,
}.Froze()

func init() {
	extra.RegisterFuzzyDecoders()
}

func Marshal(v any) ([]byte, error) {
	if pv, ok := v.(proto.Message); ok {
		return protojson.Marshal(pv)
	}
	return json.Marshal(v)
}
func Unmarshal(data []byte, v any) error {
	if pv, ok := v.(proto.Message); ok {
		return protojson.Unmarshal(data, pv)
	}
	return json.Unmarshal(data, v)
}

func ToString(v any) (result string, err error) {
	if v == nil {
		return
	}
	bytes, err := Marshal(v)
	if err != nil {
		return
	}
	result = string(bytes)
	return
}
func TryToString(v any) (result string) {
	result, _ = ToString(v)
	return
}
func FromString(str string, v any) error {
	return json.UnmarshalFromString(str, v)
}
