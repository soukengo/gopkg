package socket

import "github.com/soukengo/gopkg/component/transport/socket/packet"

type PushRoomPacket struct {
	RoomId string
	Packet packet.IPacket
}
