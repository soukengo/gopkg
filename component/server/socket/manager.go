package socket

import (
	"github.com/soukengo/gopkg/component/server/socket/packet"
	"github.com/zhenjl/cityhash"
)

func (s *server) Bucket(channelId string) *Bucket {
	idx := cityhash.CityHash32([]byte(channelId), uint32(len(channelId))) % s.opt.BucketSize
	return s.buckets[idx]
}

func (s *server) Channel(channelId string) (Channel, bool) {
	return s.Bucket(channelId).Channel(channelId)
}

func (s *server) JoinRoom(roomId string, channel Channel) error {
	return s.Bucket(roomId).JoinRoom(roomId, channel)
}
func (s *server) QuitRoom(roomId string, channel Channel) error {
	s.Bucket(roomId).QuitRoom(roomId, channel)
	return nil
}
func (s *server) PushRoom(roomId string, p packet.IPacket) {
	for _, bucket := range s.buckets {
		bucket.PushRoom(roomId, p)
	}
}
