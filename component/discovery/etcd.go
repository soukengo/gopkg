package discovery

import (
	"errors"
	impl "github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"time"
)

func NewEtcd(cfg *Config) (res Discovery, err error) {
	if cfg.Etcd == nil || cfg.Etcd.Nodes == nil || len(cfg.Etcd.Nodes) == 0 {
		err = errors.New("discovery etcd配置错误")
		return
	}
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   cfg.Etcd.Nodes,
		Username:    cfg.Etcd.Username,
		Password:    cfg.Etcd.Password,
		DialTimeout: time.Second * 10, DialOptions: []grpc.DialOption{grpc.WithBlock()},
	})
	if err != nil {
		return
	}
	res = impl.New(client)
	return
}
