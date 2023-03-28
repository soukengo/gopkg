package discovery

import (
	"errors"
	impl "github.com/go-kratos/kratos/contrib/registry/zookeeper/v2"
	"github.com/go-zookeeper/zk"
	"time"
)

func NewZookeeper(cfg *Config) (res Discovery, err error) {
	if cfg.Zookeeper == nil || cfg.Zookeeper.Nodes == nil || len(cfg.Zookeeper.Nodes) == 0 {
		err = errors.New("discovery zookeeper配置错误")
		return
	}
	conn, _, err := zk.Connect(cfg.Zookeeper.Nodes, time.Minute*60)
	if err != nil {
		return
	}
	var opts []impl.Option
	if len(cfg.Zookeeper.RootPath) > 0 {
		opts = append(opts, impl.WithRootPath(cfg.Zookeeper.RootPath))
	}
	res = impl.New(conn, opts...)
	return
}
