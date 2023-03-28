package socket

import (
	"github.com/soukengo/gopkg/component/server/socket/network"
	"github.com/soukengo/gopkg/component/server/socket/network/tcp"
	"github.com/soukengo/gopkg/component/server/socket/network/tcp/nbio"
	"github.com/soukengo/gopkg/component/server/socket/network/ws"
)

func (s *server) RegisterTCPServer(cfg *tcp.Config) {
	ins := nbio.NewServer(cfg, s.opt.Parser)
	//s := gnet.NewServer(cfg, s.opt.Parser)
	s.register(ins)
	return
}
func (s *server) RegisterWSServer(cfg *ws.Config) {
	ins := ws.NewServer(cfg, s.opt.Parser)
	s.register(ins)
	return
}

func (s *server) register(srv network.Server) {
	srvId := srv.Id()
	s.servers[srvId] = srv
	srv.SetHandler(s)
}
