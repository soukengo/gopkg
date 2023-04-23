package redis

import (
	"time"
)

const (
	Max = "+inf"
	Min = "-inf"
)

type Config struct {
	Addr         []string
	Password     string
	Prefix       string
	Active       int
	Idle         int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type Configs map[string]*Config

func (m Configs) Use(key string) *Config {
	if item, ok := m[key]; ok {
		return item
	}
	return nil
}

type Reference struct {
	Key    string
	Config *Config
}

func (r *Reference) Parse(configs Configs) {
	r.Config = configs.Use(r.Key)
}
