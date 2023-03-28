package ws

import (
	"bytes"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/soukengo/gopkg/component/server/socket/packet"
	"net"
	"sync"
)

type OnReceived func(p packet.IPacket)
type wsConn struct {
	id       string
	core     *websocket.Conn
	parser   *packet.Parser
	msgType  int
	callback OnReceived
	once     *sync.Once
	done     chan struct{} // 关闭标识
}

func newConn(msgType int, core *websocket.Conn, parser *packet.Parser, callback OnReceived) *wsConn {
	c := &wsConn{
		id:       uuid.New().String(),
		core:     core,
		parser:   parser,
		msgType:  msgType,
		callback: callback,
		once:     new(sync.Once),
		done:     make(chan struct{}),
	}
	if c.msgType <= 0 {
		c.msgType = websocket.BinaryMessage
	}
	return c
}

func (c *wsConn) Id() string {
	return c.id
}

func (c *wsConn) Read() (p packet.IPacket, err error) {
	msgType, data, err := c.core.ReadMessage()
	if err != nil {
		return
	}
	if msgType != c.msgType {
		err = packet.ErrInvalidPacket
	}
	return c.parser.Parse(c.Id(), bytes.NewReader(data))
}

func (c *wsConn) Send(packet packet.IPacket) (err error) {
	var buf bytes.Buffer
	err = packet.PackTo(&buf)
	if err != nil {
		return
	}
	err = c.core.WriteMessage(c.msgType, buf.Bytes())
	return
}

func (c *wsConn) RemoteAddr() net.Addr {
	return c.core.RemoteAddr()
}
func (c *wsConn) Close() (err error) {
	c.once.Do(func() {
		close(c.done)
		err = c.core.Close()
	})
	return
}

func (c *wsConn) IsAsync() bool {
	return true
}

// readLoop 读取客户端数据
func (c *wsConn) readLoop() {
	defer c.Close()
	var (
		err error
		msg []byte
	)
	for {
		select {
		case <-c.done:
			return
		default:
			_, msg, err = c.core.ReadMessage()
			if err != nil {
				return
			}
			// 当客户端主动断开，可能导致数据为空，直接返回
			if len(msg) == 0 {
				return
			}
			var p packet.IPacket
			p, err = c.parser.Parse(c.id, bytes.NewReader(msg))
			if err != nil {
				return
			}
			c.callback(p)
		}

	}

}
