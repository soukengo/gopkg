package ws

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/soukengo/gopkg/component/server/socket/network"
	"github.com/soukengo/gopkg/component/server/socket/packet"
	"github.com/soukengo/gopkg/log"
	"github.com/soukengo/gopkg/util/host"
	"net/http"
	"net/url"
	"sync"
)

type wsServer struct {
	id        string
	cfg       *Config
	endpoint  *url.URL
	parser    *packet.Parser
	upgrader  *websocket.Upgrader
	handler   network.Handler
	startOnce *sync.Once
}

func NewServer(cfg *Config, parser *packet.Parser) network.Server {
	if len(cfg.Path) == 0 {
		cfg.Path = "/"
	}
	return &wsServer{
		id:       uuid.New().String(),
		cfg:      cfg,
		endpoint: host.EndPointWithPath(host.SchemeWS, cfg.Addr, cfg.Path),
		parser:   parser,
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  cfg.ReadBufSize,
			WriteBufferSize: cfg.WriteBufSize,
		},
		startOnce: new(sync.Once),
	}
}

func (s *wsServer) Id() string {
	return s.id
}

func (s *wsServer) Start() (err error) {
	s.startOnce.Do(func() {
		mux := http.NewServeMux()
		mux.HandleFunc(s.cfg.Path, s.handleRequest)
		go func() {
			if err = http.ListenAndServe(s.cfg.Addr, mux); err != nil {
				panic(err)
			}
		}()
		log.Infof("started ws server with endpoint %s", s.endpoint)
	})
	return
}

// handleRequest 处理请求
func (s *wsServer) handleRequest(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrade(w, r)
	if err != nil {
		log.Errorf("s.upgrade error: %v", err)
		return
	}
	s.handler.OnConnect(conn)
	// 循环读取数据
	conn.readLoop()
	// 执行到这里说明连接断开了
	s.handler.OnDisConnect(conn)
}

// upgrade 升级协议，返回socket连接
func (s *wsServer) upgrade(w http.ResponseWriter, r *http.Request) (conn *wsConn, err error) {
	ws, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	conn = newConn(s.cfg.MsgType, ws, s.parser, func(p packet.IPacket) {
		s.handler.OnReceived(conn, p)
	})
	return
}

func (s *wsServer) Close() error {
	return nil
}

func (s *wsServer) EndPoint() *url.URL {
	return s.endpoint
}

func (s *wsServer) SetHandler(listener network.Handler) {
	s.handler = listener
}
