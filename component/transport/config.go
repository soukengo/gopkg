package transport

import (
	"github.com/soukengo/gopkg/component/transport/queue"
	"github.com/soukengo/gopkg/component/transport/socket/network/tcp"
	"github.com/soukengo/gopkg/component/transport/socket/network/ws"
	"time"
)

type Config struct {
	Http   *Http
	Grpc   *Grpc
	Socket *Socket
	Queue  *queue.Config
}

type Http struct {
	Network string
	Addr    string
	Timeout time.Duration
}

type Grpc struct {
	Network string
	Addr    string
	Timeout time.Duration
}

type Socket struct {
	TCP *tcp.Config
	WS  *ws.Config
}
