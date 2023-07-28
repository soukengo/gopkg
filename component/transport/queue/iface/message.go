package iface

type Message interface {
	Topic() Topic
	Data() any
}

type RawMessage struct {
	topic Topic
	data  any
}

func NewRawMessage(topic Topic, data any) *RawMessage {
	return &RawMessage{topic: topic, data: data}
}

func (m *RawMessage) Topic() Topic {
	return m.topic
}

func (m *RawMessage) Data() any {
	return m.data
}

type BytesMessage struct {
	topic Topic
	data  []byte
}

func NewBytesMessage(topic Topic, data []byte) *BytesMessage {
	return &BytesMessage{topic: topic, data: data}
}

func (m *BytesMessage) Topic() Topic {
	return m.topic
}

func (m *BytesMessage) Data() any {
	return m.data
}

func (m *BytesMessage) Bytes() []byte {
	return m.data
}
