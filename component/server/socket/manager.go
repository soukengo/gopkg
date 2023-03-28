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

func (s *server) JoinGroup(groupId string, channel Channel) error {
	return s.Bucket(groupId).JoinGroup(groupId, channel)
}
func (s *server) QuitGroup(groupId string, channel Channel) error {
	s.Bucket(groupId).QuitGroup(groupId, channel)
	return nil
}
func (s *server) PushGroup(groupId string, p packet.IPacket) {
	for _, bucket := range s.buckets {
		bucket.PushGroup(groupId, p)
	}
}
