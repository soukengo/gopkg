package socket

import (
	"context"
	"errors"
	"github.com/soukengo/gopkg/component/transport/socket/network"
	"github.com/soukengo/gopkg/component/transport/socket/options"
)

type server struct {
	buckets []*Bucket
	handler Handler
	servers map[string]network.Server
	opt     *options.Options
}

func NewServer(handler Handler, opts ...options.Option) Server {
	m := &server{handler: handler}
	// set options
	opt := options.Default()
	opt.ParseOptions(opts...)
	// channel buckets
	m.opt = opt
	m.buckets = make([]*Bucket, opt.BucketSize)
	for i := uint32(0); i < opt.BucketSize; i++ {
		m.buckets[i] = newBucket(opt.ChannelSize, opt.RoomSize, opt.RoutineAmount, opt.RoutineSize)
	}
	m.servers = make(map[string]network.Server)
	return m
}

func (s *server) Start(ctx context.Context) (err error) {
	if len(s.servers) == 0 {
		return errors.New("running server manager without servers")
	}
	for _, srv := range s.servers {
		err = srv.Start()
		if err != nil {
			return
		}
	}
	return
}

func (s *server) Stop(ctx context.Context) (err error) {
	for _, srv := range s.servers {
		err = srv.Close()
		if err != nil {
			return
		}
	}
	return
}
