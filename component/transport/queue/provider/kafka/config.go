package kafka

import (
	"github.com/soukengo/gopkg/infra/storage/kafka"
)

type Config struct {
	kafka.Reference `mapstructure:",squash"`
}

type ConsumerConfig struct {
	Workers int
	GroupId string
}
