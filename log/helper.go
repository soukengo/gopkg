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

func (h Helper) Debug(args ...interface{}) {
	h.logger.Log(LevelDebug, Pairs{MessageKey: fmt.Sprint(args...)})
}

func (h Helper) Debugf(format string, args ...interface{}) {
	h.logger.Log(LevelDebug, Pairs{MessageKey: fmt.Sprintf(format, args...)})
}

func (h Helper) Info(args ...interface{}) {
	h.logger.Log(LevelInfo, Pairs{MessageKey: fmt.Sprint(args...)})
}

func (h Helper) Infof(format string, args ...interface{}) {
	h.logger.Log(LevelInfo, Pairs{MessageKey: fmt.Sprintf(format, args...)})
}

func (h Helper) Warn(args ...interface{}) {
	h.logger.Log(LevelWarn, Pairs{MessageKey: fmt.Sprint(args...)})
}

func (h Helper) Warnf(format string, args ...interface{}) {
	h.logger.Log(LevelWarn, Pairs{MessageKey: fmt.Sprintf(format, args...)})
}

func (h Helper) Error(args ...interface{}) {
	h.logger.Log(LevelError, Pairs{MessageKey: fmt.Sprint(args...)})
}

func (h Helper) Errorf(format string, args ...interface{}) {
	h.logger.Log(LevelError, Pairs{MessageKey: fmt.Sprintf(format, args...)})
}

func (h Helper) WithContext(ctx context.Context) *Helper {
	return &Helper{logger: h.logger.WithContext(ctx)}
}
