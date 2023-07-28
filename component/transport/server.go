package transport

import (
	"github.com/go-kratos/kratos/v2/transport"
)

type Server = transport.Server

type ServerGroup interface {
	Register(Server)
	Servers() []Server
}

type group struct {
	list []Server
}

func NewServerGroup(list ...Server) ServerGroup {
	return &group{list: list}
}

func (g *group) Register(server Server) {
	g.list = append(g.list, server)
}

func (g *group) Servers() []Server {
	return g.list
}
