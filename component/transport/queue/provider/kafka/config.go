package kafka

import (
	"github.com/soukengo/gopkg/infra/storage/kafka"
)

type Config struct {
	kafka.Reference `mapstructure:",squash"`
	Consumer        *kafka.ConsumerConfig
}
