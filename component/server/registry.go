package server

import (
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/soukengo/gopkg/component/server/job"
)

type HttpServerRegistry interface {
	Register(*http.Server)
}

type GrpcServerRegistry interface {
	Register(*grpc.Server)
}
type JobServerRegistry interface {
	Register(job.Server)
}
