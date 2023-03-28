package discovery

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/registry"
)

type Discovery interface {
	registry.Discovery
	registry.Registrar
}

const (
	Etcd      = "etcd"
	Zookeeper = "zookeeper"
	Nacos     = "nacos"
)

type handler func(config *Config) (Discovery, error)

var (
	processor = map[string]handler{
		Zookeeper: NewZookeeper,
		Etcd:      NewEtcd,
		Nacos:     NewNacos,
	}
)

func NewDiscovery(cfg *Config) (registry.Discovery, error) {
	engine := cfg.Engine
	f, ok := processor[cfg.Engine]
	if !ok {
		return nil, fmt.Errorf("engine [%s] not supported", engine)
	}
	return f(cfg)
}

func NewRegistrar(cfg *Config) (registry.Registrar, error) {
	engine := cfg.Engine
	f, ok := processor[cfg.Engine]
	if !ok {
		return nil, fmt.Errorf("engine [%s] not supported", engine)
	}
	return f(cfg)
}
