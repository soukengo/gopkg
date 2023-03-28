package discovery

import (
	"errors"
	impl "github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func NewNacos(c *Config) (res Discovery, err error) {
	if c.Nacos == nil || len(c.Nacos.Server.Addr) == 0 || c.Nacos.Server.Port <= 0 {
		err = errors.New("discovery nacos配置错误")
		return
	}
	cfg := c.Nacos
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(cfg.Server.Addr, cfg.Server.Port),
	}
	cc := constant.ClientConfig{
		AppName:             c.AppId,
		BeatInterval:        10000,
		NamespaceId:         cfg.Client.Namespace,
		RegionId:            cfg.Client.RegionId,
		AccessKey:           cfg.Client.AccessKey,
		SecretKey:           cfg.Client.SecretKey,
		Username:            cfg.Client.Username,
		Password:            cfg.Client.Password,
		TimeoutMs:           cfg.Client.TimeoutMs,
		NotLoadCacheAtStart: true,
		LogDir:              cfg.Client.LogDir,
		CacheDir:            cfg.Client.CacheDir,
		LogLevel:            cfg.Client.LogLevel,
	}
	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		return
	}
	return impl.New(client, impl.WithGroup(cfg.Client.GroupId)), nil
}
