package socket

import (
	"context"
	"time"
)

type Context struct {
	ctx context.Context
	ch  Channel
}

func NewContext(ctx context.Context, ch Channel) *Context {
	return &Context{ctx: ctx, ch: ch}
}

func (c *Context) Channel() Channel {
	return c.ch
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

func (c *Context) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *Context) Err() error {
	return c.ctx.Err()
}

func (c *Context) Value(key any) any {
	return c.ctx.Value(key)
}

func FromContext(ctx context.Context) (Channel, bool) {
	c, ok := ctx.(*Context)
	if !ok {
		return nil, false
	}
	return c.ch, true
}
