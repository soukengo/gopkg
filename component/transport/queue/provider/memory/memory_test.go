package memory

import (
	"context"
	"github.com/soukengo/gopkg/component/transport/queue/iface"
	"github.com/soukengo/gopkg/component/transport/queue/options"
	"testing"
	"time"
)

const (
	eventGroupJoined = "GroupJoinedEvent"
)

type GroupJoinedEvent struct {
	GroupId string
}

func TestMemory(t *testing.T) {

	q := NewQueue(10)

	q.Subscribe(eventGroupJoined, iface.HandleWithOptions(func(ctx context.Context, g *GroupJoinedEvent) iface.Action {
		t.Logf("OnGroupJoinedEvent: %v", g.GroupId)
		return iface.None
	}, options.Consumer().SetMode(options.Async)))
	ctx := context.TODO()
	q.Start(ctx)
	q.Publish(ctx, iface.NewRawMessage(eventGroupJoined, &GroupJoinedEvent{
		GroupId: "10001",
	}), nil)
	q.Publish(ctx, iface.NewRawMessage(eventGroupJoined, &GroupJoinedEvent{
		GroupId: "10002",
	}), nil)
	q.Publish(ctx, iface.NewRawMessage(eventGroupJoined, &GroupJoinedEvent{
		GroupId: "10003",
	}), nil)

	time.Sleep(time.Minute)

}
