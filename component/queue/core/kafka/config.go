package kafka

import (
	"github.com/soukengo/gopkg/component/queue/core"
	"github.com/soukengo/gopkg/infra/storage/kafka"
)

type Config struct {
	kafka.Reference
	Consumer *ConsumerConfig
}

type ConsumerConfig struct {
	Topics  []core.Topic
	Workers int
	GroupId string
}
