package config

import (
	"fmt"
	"io"
	"path/filepath"
)

type EnvSource struct {
	cfg     *Config
	envFile string
	source  Source
	err     error
}

func NewEnvSource(cfg *Config, configName string) Source {
	s := &EnvSource{cfg: cfg}
	switch s.cfg.Engine {
	case File:
		s.source = NewFileSource(filepath.Join(s.cfg.File.Dir, configName))
	case Nacos:
		s.cfg.Nacos.Client.Namespace = s.cfg.Env
		s.source = NewNacosSource(&s.cfg.Nacos, configName)
	default:
		s.err = fmt.Errorf("invalid engine: %s", s.cfg.Engine)
	}
	return s
}

func (s *EnvSource) From() (io.Reader, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.source.From()
}

func (s *EnvSource) ConfigType() string {
	if s.err != nil {
		return ""
	}
	return s.source.ConfigType()
}

func (s *EnvSource) Close() error {
	if s.err != nil {
		return nil
	}
	return s.source.Close()
}
