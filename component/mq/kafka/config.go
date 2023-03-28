package kafka

import "github.com/soukengo/gopkg/component/mq"

type ConsumerConfig struct {
	Group   string
	Topics  []string
	Workers int32
	Kafka   *mq.Kafka
}
