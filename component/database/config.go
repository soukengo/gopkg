package database

import (
	"github.com/soukengo/gopkg/component/database/mongo"
	"github.com/soukengo/gopkg/component/database/redis"
)

type Config struct {
	Mongodb *mongo.Config
	Redis   *redis.Config
}
