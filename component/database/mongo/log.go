package mongo

import (
	"github.com/soukengo/gopkg/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"golang.org/x/net/context"
)

const (
	loggerName = "mongo"
)

type logger struct {
	opts   *options
	logger *log.Helper
}

func newLogger(opts *options) Monitor {
	return &logger{opts: opts, logger: opts.logger.Factory().Use(loggerName).Helper()}
}

func (l *logger) Started(ctx context.Context, event *event.CommandStartedEvent) {
	statement, _ := bson.MarshalExtJSON(event.Command, false, false)
	l.logger.WithContext(ctx).Infof("exec[%v] started command: %s, statement: %s", event.RequestID, event.CommandName, string(statement))
}

func (l *logger) Succeeded(ctx context.Context, event *event.CommandSucceededEvent) {
	l.logger.WithContext(ctx).Infof("exec[%v] succeeded command: %s,reply: %s", event.RequestID, event.CommandName, event.Reply)
}

func (l *logger) Failed(ctx context.Context, event *event.CommandFailedEvent) {
	l.logger.WithContext(ctx).Errorf("exec[%v] failed command: %s", event.RequestID, event.CommandName)
}
