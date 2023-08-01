package transport

import (
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/soukengo/gopkg/component/transport/queue"
	"github.com/soukengo/gopkg/component/transport/socket"
)

type HttpServerRegistry interface {
	RegisterHttp(*http.Server)
}

type GrpcServerRegistry interface {
	RegisterGrpc(*grpc.Server)
}
type QueueServerRegistry interface {
	RegisterQueue(queue.Server)
}
type SocketServerRegistry interface {
	RegisterSocket(server socket.Server)
}
