package kafka

const (
	defaultWorkers = 1
)

type Message struct {
	key, value  []byte
	originTopic string
	topic       string
	partition   int32
	offset      int64
}

func (m *Message) Topic() string {
	return m.topic
}

func (m *Message) Value() []byte {
	return m.value
}

type HandlerFunc func(m *Message)
