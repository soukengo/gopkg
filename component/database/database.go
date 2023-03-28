package database

import (
	"github.com/soukengo/gopkg/component/database/mongo"
	"github.com/soukengo/gopkg/component/database/redis"
	"github.com/soukengo/gopkg/log"
)

func NewRedisClient(cfg *Config, logger log.Logger) *redis.Client {
	return redis.NewClient(cfg.Redis, logger)
}
func NewMongoClient(cfg *Config, logger log.Logger) *mongo.Client {
	return mongo.NewClientWithOptions(cfg.Mongodb, mongo.WithLogger(logger))
}
