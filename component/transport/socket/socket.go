package socket

import (
	"context"
	"github.com/soukengo/gopkg/component/transport/socket/network/tcp"
	"github.com/soukengo/gopkg/component/transport/socket/network/ws"
	"github.com/soukengo/gopkg/component/transport/socket/packet"
)

type Server interface {
	Start(context.Context) error
	Stop(context.Context) error
	SetHandler(handler Handler)
	Manager
	Registry
}

type Manager interface {
	Channel(channelId string) (Channel, bool)
	JoinRoom(roomId string, channel Channel) error
	QuitRoom(roomId string, channel Channel) error
	PushRoom(roomId string, p packet.IPacket)
}

type Registry interface {
	RegisterTCPServer(cfg *tcp.Config)
	RegisterWSServer(cfg *ws.Config)
}

type Handler interface {
	OnCreated(*Context)
	OnClosed(*Context)
	OnReceived(*Context, packet.IPacket)
}

type Channel interface {
	Id() string
	ClientIP() string
	Send(packet packet.IPacket) error
	Close() error
	AddRoom(roomId string)
	DelRoom(roomId string)
	Rooms() []string
	Attributes
}

type Room interface {
	ID() string
	Put(ch Channel) (err error)
	Del(ch Channel) bool
	Push(p packet.IPacket)
	Close()
}

type Attributes interface {
	SetSession(obj any)
	Session() any
}
