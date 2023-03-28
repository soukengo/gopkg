package log

import (
	"github.com/go-kratos/kratos/v2/log"
)

type Factory interface {
	Root() Logger
	Use(name string) Logger
}

type factory struct {
	root    Logger
	loggers map[string]Logger
}

var (
	defaultFactory = &factory{
		loggers: map[string]Logger{},
	}
	global Logger
)

func init() {
	defaultFactory.root = newLogger(defaultFactory, log.With(newSingleProvider(LevelInfo, &LoggerConfig{})))
	SetLogger(defaultFactory.Root())
}

func Configure(cfg *Config, pairs Pairs) (Logger, error) {
	f := &factory{
		loggers: map[string]Logger{},
	}
	root, err := newLoggerFromConfig(f, &cfg.LoggerConfig, pairs)
	if err != nil {
		return nil, err
	}
	f.root = root
	for _, item := range cfg.Loggers {
		var l Logger
		l, err = newLoggerFromConfig(f, &item, pairs)
		if err != nil {
			return nil, err
		}
		f.loggers[item.Name] = l
	}
	return newProxy(f.root), nil
}

func (m *factory) Root() Logger {
	if m.root == nil {
		return empty
	}
	return m.root
}

func (m *factory) Use(name string) Logger {
	l, ok := m.loggers[name]
	if !ok {
		return empty
	}
	return l
}
