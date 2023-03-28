package json

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
)

var json = jsoniter.Config{
	EscapeHTML:             true,
	ValidateJsonRawMessage: true,
}.Froze()

func init() {
	extra.RegisterFuzzyDecoders()
}

func Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}
func Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func ToString(v any) (result string, err error) {
	if v == nil {
		return
	}
	result, err = json.MarshalToString(v)
	return
}
func TryToString(v any) (result string) {
	result, _ = ToString(v)
	return
}
func FromString(str string, v any) error {
	return json.UnmarshalFromString(str, v)
}
