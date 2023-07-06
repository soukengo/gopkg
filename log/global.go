package log

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
)

func Global() Logger {
	return global
}

func SetLogger(l Logger) {
	p, ok := l.(*proxy)
	global = l
	if ok {
		global = p.logger
	}
	if c, ok := global.(*logger); ok {
		log.SetLogger(&kratosLogger{logger: c.provider})
	}
}

func Debug(args ...interface{}) {
	global.Log(LevelDebug, Pairs{MessageKey: fmt.Sprint(args...)})
}

func Debugf(format string, args ...interface{}) {
	global.Log(LevelDebug, Pairs{MessageKey: fmt.Sprintf(format, args...)})
}

func Info(args ...interface{}) {
	global.Log(LevelInfo, Pairs{MessageKey: fmt.Sprint(args...)})
}

func Infof(format string, args ...interface{}) {
	global.Log(LevelInfo, Pairs{MessageKey: fmt.Sprintf(format, args...)})
}

func Warn(args ...interface{}) {
	global.Log(LevelWarn, Pairs{MessageKey: fmt.Sprint(args...)})
}

func Warnf(format string, args ...interface{}) {
	global.Log(LevelWarn, Pairs{MessageKey: fmt.Sprintf(format, args...)})
}

func Error(args ...interface{}) {
	global.Log(LevelError, Pairs{MessageKey: fmt.Sprint(args...)})
}

func Errorf(format string, args ...interface{}) {
	global.Log(LevelError, Pairs{MessageKey: fmt.Sprintf(format, args...)})
}

func WithContext(ctx context.Context) *Helper {
	return global.WithContext(ctx).Helper()
}
func With(pairs Pairs) *Helper {
	return global.With(pairs).Helper()
}
