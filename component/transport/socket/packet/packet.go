package packet

import (
	"errors"
	"io"
)

var (
	ErrInvalidPacket = errors.New("invalid packet")
)

type IPacket interface {
	UnPackFrom(r io.Reader) (err error)
	PackTo(w io.Writer) (err error)
}

type IParser interface {
	Parse(connId string, reader io.Reader) (IPacket, error)
}
