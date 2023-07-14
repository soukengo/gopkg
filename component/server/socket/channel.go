package socket

import (
	"context"
	log "github.com/golang/glog"
	"github.com/soukengo/gopkg/component/server/socket/network"
	"github.com/soukengo/gopkg/component/server/socket/packet"
	"github.com/soukengo/gopkg/util/runtimes"
	"net"
	"sync"
)

type channel struct {
	clientIP  string
	conn      network.Connection
	done      *sync.Once
	sendQueue chan packet.IPacket
	ctx       context.Context
	cancel    context.CancelFunc
	rooms     map[string]struct{}
	rlock     sync.RWMutex
	Attributes
}

func newChannel(conn network.Connection, sendQueueSize uint32) Channel {
	ctx, cancel := context.WithCancel(context.Background())
	ins := &channel{
		ctx:        ctx,
		cancel:     cancel,
		done:       new(sync.Once),
		sendQueue:  make(chan packet.IPacket, sendQueueSize),
		conn:       conn,
		rooms:      map[string]struct{}{},
		Attributes: newAttributes(),
	}
	ins.clientIP, _, _ = net.SplitHostPort(conn.RemoteAddr().String())
	runtimes.Async(ins.dispatch)
	return ins
}

func (c *channel) Id() string {
	return c.conn.Id()
}

func (c *channel) ClientIP() string {
	return c.clientIP
}

func (c *channel) Send(data packet.IPacket) error {
	select {
	case <-c.ctx.Done():
		return c.ctx.Err()
	default:
		c.sendQueue <- data
	}
	return nil
}

func (c *channel) AddRoom(roomId string) {
	c.rlock.Lock()
	c.rooms[roomId] = struct{}{}
	c.rlock.Unlock()
}

func (c *channel) DelRoom(roomId string) {
	c.rlock.Lock()
	delete(c.rooms, roomId)
	c.rlock.Unlock()
}
func (c *channel) Rooms() (ret []string) {
	c.rlock.RLock()
	for id := range c.rooms {
		ret = append(ret, id)
	}
	c.rlock.RUnlock()
	return
}

func (c *channel) Close() (err error) {
	c.done.Do(func() {
		c.cancel()
		err = c.conn.Close()
	})
	return
}

// dispatch send packet to connection
func (c *channel) dispatch() {
	defer func() {
		_ = c.Close()
	}()
	for {
		select {
		case <-c.ctx.Done():
			return
		case p := <-c.sendQueue:
			err := c.conn.Send(p)
			if err != nil {
				log.Errorf("conn.Send err: %v", err)
				return
			}
		}
	}
}
