package socket

import "github.com/soukengo/gopkg/component/server/socket/packet"

type PushRoomPacket struct {
	RoomId string
	Packet packet.IPacket
}
