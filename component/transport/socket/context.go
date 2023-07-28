package socket

import (
	"context"
	"time"
)

type Context struct {
	ctx context.Context
	srv Server
	ch  Channel
}

func NewContext(ctx context.Context, srv Server, ch Channel) *Context {
	return &Context{ctx: ctx, srv: srv, ch: ch}
}

func (c *Context) Server() Server {
	return c.srv
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

func FromContext(ctx context.Context) (*Context, bool) {
	c, ok := ctx.(*Context)
	return c, ok
}
