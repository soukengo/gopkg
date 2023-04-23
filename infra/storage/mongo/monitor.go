package mongo

import (
	"go.mongodb.org/mongo-driver/event"
	"golang.org/x/net/context"
)

type Monitor interface {
	Started(context.Context, *event.CommandStartedEvent)
	Succeeded(context.Context, *event.CommandSucceededEvent)
	Failed(context.Context, *event.CommandFailedEvent)
}

type monitor struct {
	opts     *options
	monitors []Monitor
	entry    *event.CommandMonitor
}

func newMonitor(opts *options) *monitor {
	m := &monitor{
		opts:     opts,
		monitors: make([]Monitor, 0),
	}
	m.entry = &event.CommandMonitor{
		Started:   m.Started,
		Succeeded: m.Succeeded,
		Failed:    m.Failed,
	}
	return m
}

func (m *monitor) Entry() *event.CommandMonitor {
	return m.entry
}

func (m *monitor) WithTrace() {
	m.monitors = append(m.monitors, newTracer())
	return
}
func (m *monitor) WithLog() {
	m.monitors = append(m.monitors, newLogger(m.opts))
	return
}

func (m *monitor) Started(ctx context.Context, evt *event.CommandStartedEvent) {
	for _, cm := range m.monitors {
		cm.Started(ctx, evt)
	}
}
func (m *monitor) Succeeded(ctx context.Context, evt *event.CommandSucceededEvent) {
	for _, cm := range m.monitors {
		cm.Succeeded(ctx, evt)
	}

}
func (m *monitor) Failed(ctx context.Context, evt *event.CommandFailedEvent) {
	for _, cm := range m.monitors {
		cm.Failed(ctx, evt)
	}
}
