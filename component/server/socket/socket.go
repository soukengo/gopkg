package socket

import (
	"context"
	"github.com/soukengo/gopkg/component/server/socket/network/tcp"
	"github.com/soukengo/gopkg/component/server/socket/network/ws"
	"github.com/soukengo/gopkg/component/server/socket/packet"
)

type Server interface {
	Start(context.Context) error
	Stop(context.Context) error
	Manager
	Registry
}

type Manager interface {
	Channel(channelId string) (Channel, bool)
	JoinGroup(groupId string, channel Channel) error
	QuitGroup(groupId string, channel Channel) error
	PushGroup(groupId string, p packet.IPacket)
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
	AddGroup(groupId string)
	DelGroup(groupId string)
	Groups() []string
	MarkAuthenticated()
	Authenticated() bool
	Attributes
}

type Group interface {
	ID() string
	Put(ch Channel) (err error)
	Del(ch Channel) bool
	Push(p packet.IPacket)
	Close()
}

type Attributes interface {
	// SetAttr sets attributes
	SetAttr(key string, value any)
	// DelAttr delete an attribute
	DelAttr(key string)
	// Attr get an attribute
	Attr(key string) (value any, ok bool)
	// StringAttr get an string attribute
	StringAttr(key string) (value string)
	// Int64Attr get an int64 attribute
	Int64Attr(key string) (value int64)
}
