package database

import (
	"github.com/soukengo/gopkg/infra/storage"
	"github.com/soukengo/gopkg/infra/storage/mongo"
)

type Config struct {
	Mongodb *mongo.Reference
}

func (c Config) Parse(configs *storage.Config) {
	if c.Mongodb != nil {
		c.Mongodb.Parse(configs.Mongodb)
	}
}
