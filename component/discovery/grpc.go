package discovery

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	ggrpc "google.golang.org/grpc"
	"time"
)

const serverPrefix = "discovery://default/"

func appServerAddr(app string) string {
	return serverPrefix + app
}

const (
	timeout = time.Second * 10
)

func NewGrpcClient(app string, r registry.Discovery) (*ggrpc.ClientConn, error) {
	return grpc.DialInsecure(
		context.Background(),
		grpc.WithTimeout(timeout),
		grpc.WithEndpoint(appServerAddr(app)),
		grpc.WithDiscovery(r),
		grpc.WithMiddleware(
			tracing.Client(),
			recovery.Recovery(),
		),
	)
}
