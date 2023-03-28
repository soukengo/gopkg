package socket

import "github.com/soukengo/gopkg/component/server/socket/packet"

type PushGroupPacket struct {
	GroupId string
	Packet  packet.IPacket
}
