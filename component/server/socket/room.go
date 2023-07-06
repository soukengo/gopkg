package socket

import (
	"errors"
	"github.com/soukengo/gopkg/component/server/socket/packet"
	"sync"
)

var (
	ErrRoomDropped = errors.New("room dropped")
)

type room struct {
	id       string
	rLock    sync.RWMutex
	channels map[string]Channel
	drop     bool
	online   int32 // dirty read is ok
}

func NewRoom(id string) Room {
	r := new(room)
	r.id = id
	r.drop = false
	r.online = 0
	r.channels = make(map[string]Channel)
	return r
}

func (r *room) ID() string {
	return r.id
}
func (r *room) Put(ch Channel) (err error) {
	r.rLock.Lock()
	if !r.drop {
		if r.channels[ch.Id()] == nil {
			r.channels[ch.Id()] = ch
			r.online++
		}
	} else {
		err = ErrRoomDropped
	}
	r.rLock.Unlock()
	return
}

func (r *room) Del(ch Channel) bool {
	r.rLock.Lock()
	delete(r.channels, ch.Id())
	r.online--
	r.drop = r.online == 0
	r.rLock.Unlock()
	return r.drop
}

func (r *room) Push(p packet.IPacket) {
	r.rLock.RLock()
	for _, ch := range r.channels {
		_ = ch.Send(p)
	}
	r.rLock.RUnlock()
}

func (r *room) Close() {
	r.rLock.RLock()
	for _, ch := range r.channels {
		ch.DelRoom(r.id)
	}
	r.rLock.RUnlock()
}
