package mq

type Config struct {
	Kafka *Kafka
}

type Kafka struct {
	Brokers     []string
	TopicPrefix string
}
