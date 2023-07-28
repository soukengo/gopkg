package codec

var (
	JSON     = NewJsonCodec()
	Protobuf = NewProtobufCodec()
)

type Codec interface {
	Encoder
	Decoder
}

type Encoder interface {
	Encode(any) ([]byte, error)
}

type Decoder interface {
	Decode([]byte, any) error
}
