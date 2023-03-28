package host

import (
	"github.com/soukengo/gopkg/util/ip"
	"net/url"
	"strings"
)

const (
	LocalHost = "localhost"
	AnyHost   = "0.0.0.0"
)

const (
	SchemeGrpc = "grpc"
	SchemeHttp = "http"
	SchemeTcp  = "tcp"
	SchemeWS   = "ws"
	SchemeWSS  = "wss"
)

func FixAddr(addr string) string {
	if strings.HasPrefix(addr, ":") {
		innerIP := ip.InternalIP()
		addr = innerIP + addr
	}
	return addr
}

func EndPoint(scheme string, addr string) *url.URL {
	return EndPointWithPath(scheme, addr, "")
}
func EndPointWithPath(scheme string, addr string, path string) *url.URL {
	return &url.URL{Scheme: scheme, Host: FixAddr(addr), Path: path}
}

func FindEndpoint(endpoints []string, scheme string) (string, error) {
	for _, e := range endpoints {
		u, err := url.Parse(e)
		if err != nil {
			return "", err
		}
		if u.Scheme == scheme {
			return u.Host, nil
		}
	}
	return "", nil
}
