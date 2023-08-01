package socket

import (
	"context"
	"github.com/soukengo/gopkg/component/transport/socket/network"
	"github.com/soukengo/gopkg/component/transport/socket/packet"
	"github.com/soukengo/gopkg/log"
)

type defaultHandler struct {
}

func (h *defaultHandler) OnCreated(ctx *Context) {
	log.WithContext(ctx).Info("OnClosed", log.Pairs{"channelId": ctx.Channel().Id()})
}

func (h *defaultHandler) OnClosed(ctx *Context) {
	log.WithContext(ctx).Info("OnClosed", log.Pairs{"channelId": ctx.Channel().Id()})
}

func (h *defaultHandler) OnReceived(ctx *Context, packet packet.IPacket) {
	log.WithContext(ctx).Info("OnReceived", log.Pairs{"channelId": ctx.Channel().Id(), "packet": packet})
}

func (s *server) SetHandler(handler Handler) {
	s.handler = handler
}
func (s *server) OnConnect(conn network.Connection) {
	ch := newChannel(conn, s.opt.SendQueueSize)
	s.Bucket(conn.Id()).PutChannel(ch)
	ctx := NewContext(context.Background(), s, ch)
	s.handler.OnCreated(ctx)
}

func (s *server) OnDisConnect(conn network.Connection) {
	channelId := conn.Id()
	ch, ok := s.Channel(conn.Id())
	if !ok {
		return
	}
	s.Bucket(channelId).DelChannel(ch)
	ctx := NewContext(context.Background(), s, ch)
	s.handler.OnClosed(ctx)
}

func (s *server) OnReceived(conn network.Connection, p packet.IPacket) {
	ch, ok := s.Channel(conn.Id())
	if !ok {
		_ = conn.Close()
		return
	}
	ctx := NewContext(context.Background(), s, ch)
	s.handler.OnReceived(ctx, p)
}
