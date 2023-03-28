package idgen

type Generator interface {
	GenID() (string, error)
}

type GeneratorType string

const (
	Sonyflake GeneratorType = "sonyflake"
	Mongo     GeneratorType = "mongo"
	XID       GeneratorType = "xid"
)
