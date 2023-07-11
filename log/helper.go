package log

import (
	"context"
	"fmt"
)

const (
	MessageKey = "msg"
)

type Helper struct {
	logger Logger
}

func newHelper(logger Logger) *Helper {
	return &Helper{logger: logger}
}

func (h Helper) Debug(msg string, pairs Pairs) {
	h.logger.Log(LevelDebug, h.parse(msg, pairs))
}

func (h Helper) Debugf(format string, args ...interface{}) {
	h.logger.Log(LevelDebug, h.parseFormat(format, args...))
}

func (h Helper) Info(msg string, pairs Pairs) {
	h.logger.Log(LevelInfo, h.parse(msg, pairs))
}

func (h Helper) Infof(format string, args ...interface{}) {
	h.logger.Log(LevelInfo, h.parseFormat(format, args...))
}

func (h Helper) Warn(msg string, pairs Pairs) {
	h.logger.Log(LevelWarn, h.parse(msg, pairs))
}

func (h Helper) Warnf(format string, args ...interface{}) {
	h.logger.Log(LevelWarn, h.parseFormat(format, args...))
}

func (h Helper) Error(msg string, pairs Pairs) {
	h.logger.Log(LevelError, h.parse(msg, pairs))
}

func (h Helper) Errorf(format string, args ...interface{}) {
	h.logger.Log(LevelError, h.parseFormat(format, args...))
}

func (h Helper) WithContext(ctx context.Context) *Helper {
	return &Helper{logger: h.logger.WithContext(ctx)}
}

func (h Helper) parse(msg string, pairs Pairs) Pairs {
	var p = Pairs{MessageKey: msg}
	if pairs != nil {
		for k, v := range pairs {
			p[k] = v
		}
	}
	return p
}
func (h Helper) parseFormat(format string, args ...interface{}) Pairs {
	return Pairs{MessageKey: fmt.Sprintf(format, args...)}
}
