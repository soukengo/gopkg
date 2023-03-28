package nbio

import "github.com/soukengo/gopkg/log"

type logger struct {
	level int
}

func (l *logger) SetLevel(lvl int) {
	l.level = lvl
}

func (l *logger) Debug(format string, v ...interface{}) {
	log.Debugf(format, v...)
}

func (l *logger) Info(format string, v ...interface{}) {
	log.Infof(format, v...)
}

func (l *logger) Warn(format string, v ...interface{}) {
	log.Warnf(format, v...)
}

func (l *logger) Error(format string, v ...interface{}) {
	log.Errorf(format, v...)
}
