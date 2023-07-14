package server

import (
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/soukengo/gopkg/component/server/job"
)

type HttpServerRegistry interface {
	RegisterHttp(*http.Server)
}

type GrpcServerRegistry interface {
	RegisterGrpc(*grpc.Server)
}
type JobServerRegistry interface {
	RegisterJob(job.Server)
}
