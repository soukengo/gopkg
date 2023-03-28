package socket

import (
	"context"
	log "github.com/golang/glog"
	"github.com/soukengo/gopkg/component/server/socket/network"
	"github.com/soukengo/gopkg/component/server/socket/packet"
)

func (s *server) OnConnect(conn network.Connection) {
	channelId := conn.Id()
	log.Infof("OnConnect: %v", channelId)
	ch := newChannel(conn, s.opt.SendQueueSize)
	s.Bucket(conn.Id()).PutChannel(ch)
	ctx := NewContext(context.Background(), ch)
	s.handler.OnCreated(ctx)
}

func (s *server) OnDisConnect(conn network.Connection) {
	channelId := conn.Id()
	log.Infof("OnDisConnect: %v", channelId)
	ch, ok := s.Channel(conn.Id())
	if !ok {
		return
	}
	s.Bucket(channelId).DelChannel(ch)
	ctx := NewContext(context.Background(), ch)
	s.handler.OnClosed(ctx)
}

func (s *server) OnReceived(conn network.Connection, p packet.IPacket) {
	log.Infof("OnReceived")
	ch, ok := s.Channel(conn.Id())
	if !ok {
		_ = conn.Close()
		return
	}
	ctx := NewContext(context.Background(), ch)
	s.handler.OnReceived(ctx, p)
}
