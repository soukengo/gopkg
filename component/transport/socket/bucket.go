package socket

import (
	"github.com/google/uuid"
	"github.com/soukengo/gopkg/component/transport/socket/packet"
	"github.com/soukengo/gopkg/util/runtimes"
	"sync"
	"sync/atomic"
)

type Bucket struct {
	channelSize   uint32
	id            string
	chs           map[string]Channel
	cLock         sync.RWMutex
	gLock         sync.RWMutex
	rooms         map[string]Room
	routines      []chan *PushRoomPacket
	routineAmount uint64
	routineSize   int

	routinesNum uint64
}

func newBucket(channelSize uint32, roomSize uint32, routineAmount uint64, routineSize int) *Bucket {
	b := &Bucket{
		id:            uuid.New().String(),
		chs:           make(map[string]Channel, channelSize),
		rooms:         make(map[string]Room, roomSize),
		routineAmount: routineAmount,
		routineSize:   routineSize,
	}
	b.routines = make([]chan *PushRoomPacket, routineAmount)
	for i := uint64(0); i < routineAmount; i++ {
		c := make(chan *PushRoomPacket, routineSize)
		b.routines[i] = c
		runtimes.Async(func() {
			b.roomproc(c)
		})
	}
	return b
}

func (b *Bucket) PutChannel(ch Channel) {
	b.cLock.Lock()
	b.chs[ch.Id()] = ch
	b.cLock.Unlock()
}
func (b *Bucket) DelChannel(ch Channel) {
	b.cLock.Lock()
	delete(b.chs, ch.Id())
	rooms := ch.Rooms()
	for _, roomId := range rooms {
		b.QuitRoom(roomId, ch)
	}
	b.cLock.Unlock()
}

func (b *Bucket) Channel(channelId string) (ch Channel, ok bool) {
	b.cLock.RLock()
	ch, ok = b.chs[channelId]
	b.cLock.RUnlock()
	return
}

func (b *Bucket) Room(roomId string) (room Room, ok bool) {
	b.gLock.RLock()
	room, ok = b.rooms[roomId]
	b.gLock.RUnlock()
	return
}
func (b *Bucket) DelRoom(room Room) {
	b.gLock.Lock()
	delete(b.rooms, room.ID())
	b.gLock.Unlock()
	room.Close()
}

func (b *Bucket) JoinRoom(roomId string, ch Channel) (err error) {
	var (
		g  Room
		ok bool
	)
	b.gLock.Lock()
	if g, ok = b.rooms[roomId]; !ok {
		g = NewRoom(roomId)
		b.rooms[roomId] = g
	}
	ch.AddRoom(roomId)
	b.gLock.Unlock()
	if g != nil {
		err = g.Put(ch)
	}
	return
}

func (b *Bucket) QuitRoom(roomId string, ch Channel) {
	b.gLock.RLock()
	g, ok := b.rooms[roomId]
	b.gLock.RUnlock()
	if ok {
		if g.Del(ch) {
			b.DelRoom(g)
		}
	}
	ch.DelRoom(roomId)
	return
}

// PushRoom broadcast a message to specified room
func (b *Bucket) PushRoom(roomId string, p packet.IPacket) {
	num := atomic.AddUint64(&b.routinesNum, 1) % b.routineAmount
	b.routines[num] <- &PushRoomPacket{RoomId: roomId, Packet: p}
}

// roomproc
func (b *Bucket) roomproc(c chan *PushRoomPacket) {
	for {
		arg := <-c
		if item, ok := b.Room(arg.RoomId); ok {
			var g = item
			g.Push(arg.Packet)
		}
	}
}
