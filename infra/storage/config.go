package storage

import (
	"github.com/soukengo/gopkg/infra/storage/kafka"
	"github.com/soukengo/gopkg/infra/storage/mongo"
	"github.com/soukengo/gopkg/infra/storage/redis"
)

const (
	DefaultKey = "default"
)

type Config struct {
	Mongodb mongo.Configs
	Redis   redis.Configs
	Kafka   kafka.Configs
}
