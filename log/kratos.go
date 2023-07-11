package log

import (
	"github.com/go-kratos/kratos/v2/log"
)

type kratosLogger struct {
	logger log.Logger
}

func Kratos(l Logger) log.Logger {
	if p, ok := l.(*proxy); ok {
		l = p.logger
	}
	if c, ok := l.(*logger); ok {
		return &kratosLogger{logger: c.provider}
	}
	return log.DefaultLogger
}

func (k kratosLogger) Log(level log.Level, keyvals ...interface{}) error {
	return k.logger.Log(level, keyvals...)
}
