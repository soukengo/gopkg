package log

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"strings"
)

type Provider log.Logger

type Logger interface {
	Log(level Level, paris Pairs)
	With(paris Pairs) Logger
	WithContext(ctx context.Context) Logger
	Helper() *Helper
	HelperWith(paris Pairs) *Helper
	Factory() Factory
}

type logger struct {
	factory  Factory
	provider Provider
	helper   *Helper
}

func newLogger(factory Factory, provider Provider) *logger {
	l := &logger{factory: factory, provider: provider}
	l.helper = newHelper(l)
	return l
}

func newLoggerFromConfig(factory Factory, cfg *LoggerConfig, pairs Pairs) (l Logger, err error) {
	lv, ok := levelMap[strings.ToUpper(cfg.Level)]
	if !ok {
		lv = LevelInfo
	}
	var impl Provider
	if cfg.Appender == nil || cfg.Appender.Console || len(cfg.Appender.FileName) == 0 || !cfg.Split {
		impl = newSingleProvider(lv, cfg)
	} else {
		impl, err = newLeveledLogger(lv, cfg)
	}
	if err != nil {
		return
	}
	var kvs []any
	if pairs != nil {
		for k, v := range pairs {
			kvs = append(kvs, k, v)
		}
	}
	return newLogger(factory, log.With(impl, kvs...)), nil
}

func (l *logger) Log(level Level, paris Pairs) {
	var kv = make([]any, 0)
	for k, v := range paris {
		kv = append(kv, k, v)
	}
	_ = l.provider.Log(log.Level(level), kv...)
}

func (l *logger) With(paris Pairs) Logger {
	var kv = make([]any, 0)
	for k, v := range paris {
		kv = append(kv, k, v)
	}
	return newLogger(l.factory, log.With(l.provider, kv...))
}
func (l *logger) WithContext(ctx context.Context) Logger {
	return newLogger(l.factory, log.WithContext(ctx, l.provider))
}

func (l *logger) Helper() *Helper {
	return newHelper(l)
}
func (l *logger) HelperWith(paris Pairs) *Helper {
	return l.With(paris).Helper()
}

func (l *logger) Factory() Factory {
	return l.factory
}

var (
	empty = &emptyLogger{}
)

type emptyLogger struct {
}

func (l *emptyLogger) Log(level Level, paris Pairs) {
}

func (l *emptyLogger) With(paris Pairs) Logger {
	return &emptyLogger{}
}

func (l *emptyLogger) WithContext(ctx context.Context) Logger {
	return &emptyLogger{}
}

func (l *emptyLogger) Helper() *Helper {
	return &Helper{logger: l}
}

func (l *emptyLogger) HelperWith(paris Pairs) *Helper {
	return &Helper{logger: l.With(paris)}
}

func (l *emptyLogger) Factory() Factory {
	return defaultFactory
}
