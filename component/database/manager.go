package database

import (
	"github.com/soukengo/gopkg/infra/storage/mongo"
)

type Manager struct {
	cfg      *Config
	mongoCli *mongo.Client
}

func NewManager(cfg *Config) *Manager {
	cli := mongo.NewClient(cfg.Mongodb.Config)
	return &Manager{cfg: cfg, mongoCli: cli}
}

func (m *Manager) Mongo() *mongo.Client {
	return m.mongoCli
}
