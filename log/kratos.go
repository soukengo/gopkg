package log

import (
	"github.com/go-kratos/kratos/v2/log"
)

type kratosLogger struct {
	logger log.Logger
}

func (k kratosLogger) Log(level log.Level, keyvals ...interface{}) error {
	return k.logger.Log(level, keyvals...)
}
