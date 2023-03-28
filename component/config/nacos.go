package config

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"io"
	"path"
	"strings"
)

type NacosConfig struct {
	Server NacosServerConfig
	Client NacosClientConfig
}

type NacosServerConfig struct {
	Addr string
	Port uint64
}
type NacosClientConfig struct {
	GroupId   string
	Namespace string
}

type NacosSource struct {
	cfg        *NacosConfig
	configType string
	content    string
	err        error
	dataId     string
}

func NewNacosSource(cfg *NacosConfig, dataId string) Source {
	cli, err := clients.NewConfigClient(vo.NacosClientParam{
		ServerConfigs: []constant.ServerConfig{
			*constant.NewServerConfig(cfg.Server.Addr, cfg.Server.Port),
		},
		ClientConfig: &constant.ClientConfig{
			AppName:     cfg.Client.GroupId,
			NamespaceId: cfg.Client.Namespace,
		},
	})
	configType := strings.TrimPrefix(path.Ext(dataId), ".")
	ins := &NacosSource{cfg: cfg, configType: configType, err: err, dataId: dataId}
	content, err := cli.GetConfig(vo.ConfigParam{Group: cfg.Client.GroupId, DataId: dataId})
	ins.content = content
	return ins
}

func (s *NacosSource) From() (io.Reader, error) {
	if s.err != nil {
		return nil, s.err
	}
	return strings.NewReader(s.content), nil
}

func (s *NacosSource) ConfigType() string {
	return s.configType
}

func (s *NacosSource) Close() error {
	return nil
}
