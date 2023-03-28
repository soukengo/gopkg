package options

import "github.com/soukengo/gopkg/component/server/socket/packet"

type Option func(o *Options)

type Options struct {
	BucketSize    uint32
	ChannelSize   uint32
	RoutineAmount uint64
	RoutineSize   int
	Factory       packet.IFactory
	RecvQueueSize uint32
	SendQueueSize uint32
	Parser        *packet.Parser
}

func WithBucketSize(bucketSize uint32) func(*Options) {
	return func(o *Options) { o.BucketSize = bucketSize }
}
func WithPacketFactory(factory packet.IFactory) func(*Options) {
	return func(o *Options) { o.Factory = factory }
}
func WithRecvQueueSize(recvQueueSize uint32) func(*Options) {
	return func(o *Options) { o.RecvQueueSize = recvQueueSize }
}
func WithSendQueueSize(sendQueueSize uint32) func(*Options) {
	return func(o *Options) { o.SendQueueSize = sendQueueSize }
}

func Default() *Options {
	opt := new(Options)
	opt.BucketSize = 32
	opt.ChannelSize = 1024
	opt.RoutineAmount = 32
	opt.RoutineSize = 1024

	opt.Factory = packet.DefaultPacketFactory()
	opt.RecvQueueSize = 10
	opt.SendQueueSize = 10
	return opt
}

func (opt *Options) ParseOptions(opts ...Option) {
	for _, option := range opts {
		option(opt)
	}
	opt.Parser = packet.NewPacketParser(opt.Factory)
}
