package socket

import (
	"github.com/google/uuid"
	"github.com/soukengo/gopkg/component/server/socket/packet"
	"sync"
	"sync/atomic"
)

type Bucket struct {
	channelSize   uint32
	id            string
	chs           map[string]Channel
	cLock         sync.RWMutex
	gLock         sync.RWMutex
	groups        map[string]Group
	routines      []chan *PushGroupPacket
	routineAmount uint64
	routineSize   int

	routinesNum uint64
}

func newBucket(channelSize uint32, routineAmount uint64, routineSize int) *Bucket {
	b := &Bucket{
		id:            uuid.New().String(),
		chs:           make(map[string]Channel, channelSize),
		routineAmount: routineAmount,
		routineSize:   routineSize,
	}
	b.routines = make([]chan *PushGroupPacket, routineAmount)
	for i := uint64(0); i < routineAmount; i++ {
		c := make(chan *PushGroupPacket, routineSize)
		b.routines[i] = c
		go b.groupproc(c)
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
	groups := ch.Groups()
	for _, groupId := range groups {
		b.QuitGroup(groupId, ch)
	}
	b.cLock.Unlock()
}

func (b *Bucket) Channel(channelId string) (ch Channel, ok bool) {
	b.cLock.RLock()
	ch, ok = b.chs[channelId]
	b.cLock.RUnlock()
	return
}

func (b *Bucket) Group(groupId string) (group Group, ok bool) {
	b.gLock.RLock()
	group, ok = b.groups[groupId]
	b.gLock.RUnlock()
	return
}
func (b *Bucket) DelGroup(group Group) {
	b.gLock.Lock()
	delete(b.groups, group.ID())
	b.gLock.Unlock()
	group.Close()
}

func (b *Bucket) JoinGroup(groupId string, ch Channel) (err error) {
	var (
		g  Group
		ok bool
	)
	b.gLock.Lock()
	if g, ok = b.groups[groupId]; !ok {
		g = NewGroup(groupId)
		b.groups[groupId] = g
	}
	ch.AddGroup(groupId)
	b.gLock.Unlock()
	if g != nil {
		err = g.Put(ch)
	}
	return
}

func (b *Bucket) QuitGroup(groupId string, ch Channel) {
	b.gLock.RLock()
	g, ok := b.groups[groupId]
	b.gLock.RUnlock()
	if ok {
		if g.Del(ch) {
			b.DelGroup(g)
		}
	}
	ch.DelGroup(groupId)
	return
}

// PushGroup broadcast a message to specified group
func (b *Bucket) PushGroup(groupId string, p packet.IPacket) {
	num := atomic.AddUint64(&b.routinesNum, 1) % b.routineAmount
	b.routines[num] <- &PushGroupPacket{GroupId: groupId, Packet: p}
}

// groupproc
func (b *Bucket) groupproc(c chan *PushGroupPacket) {
	for {
		arg := <-c
		if item, ok := b.Group(arg.GroupId); ok {
			var g = item
			g.Push(arg.Packet)
		}
	}
}
