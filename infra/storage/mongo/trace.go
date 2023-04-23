package mongo

import (
	"go.mongodb.org/mongo-driver/event"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
	"golang.org/x/net/context"
)

type tracer struct {
	impl *event.CommandMonitor
}

func newTracer() Monitor {
	return &tracer{impl: otelmongo.NewMonitor()}
}

func (t tracer) Started(ctx context.Context, evt *event.CommandStartedEvent) {
	t.impl.Started(ctx, evt)
}

func (t tracer) Succeeded(ctx context.Context, evt *event.CommandSucceededEvent) {
	t.impl.Succeeded(ctx, evt)
}

func (t tracer) Failed(ctx context.Context, evt *event.CommandFailedEvent) {
	t.impl.Failed(ctx, evt)
}
