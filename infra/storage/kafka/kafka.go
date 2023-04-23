package kafka

const (
	defaultWorkers = 1
)

type Value []byte

type HandlerFunc func(topic string, value Value) (err error)
