package transport

import "context"

type Server interface {
	Start(context.Context) error
	Stop(context.Context) error
}

type ServerGroup interface {
	Register(Server)
}

type group struct {
	list []Server
}

func NewGroup(list ...Server) ServerGroup {
	return &group{list: list}
}

func (g *group) Register(server Server) {
	g.list = append(g.list, server)
}
