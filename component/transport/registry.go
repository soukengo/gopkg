package transport

import (
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/soukengo/gopkg/component/transport/queue"
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
