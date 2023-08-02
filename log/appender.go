package log

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

func newAppender(c *AppenderConfig) io.Writer {
	if c == nil || c.Console {
		return os.Stdout
	}
	return &lumberjack.Logger{
		Filename:   c.FileName,
		MaxSize:    c.MaxSize,
		MaxAge:     c.MaxAge,
		MaxBackups: c.MaxBackups,
		Compress:   c.Compress,
		LocalTime:  true,
	}
}
