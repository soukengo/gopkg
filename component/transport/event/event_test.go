package event

import (
	"context"
	"github.com/soukengo/gopkg/component/transport/queue/options"
	"github.com/soukengo/gopkg/log"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type GroupJoinedEvent struct {
	GroupId string
}

var key = Key("GroupJoined")

type MemoryEventSuite struct {
	suite.Suite
	ctx         context.Context
	eventServer Server
}

func (s *MemoryEventSuite) SetupTest() {
	s.ctx = context.TODO()
	s.eventServer = NewMemoryServer(&Memory{}, log.Global())
	err := s.eventServer.Start(s.ctx)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *MemoryEventSuite) subscribe() {
	s.eventServer.Subscribe(key, Handle[GroupJoinedEvent](func(ctx context.Context, e *GroupJoinedEvent) {
		s.T().Logf("On GroupJoinedEvent: %v", e)
	}, options.Consumer()))
}

func (s *MemoryEventSuite) publish() {
	s.eventServer.Publish(s.ctx, NewEvent(key, &GroupJoinedEvent{GroupId: "10001"}), options.Producer())
}

func (s *MemoryEventSuite) TestPublishAndSubscribe() {
	go s.subscribe()
	s.publish()
	time.Sleep(time.Second * 5)
}

func TestEvent(t *testing.T) {
	suite.Run(t, new(MemoryEventSuite))
}
