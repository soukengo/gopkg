package iface

import (
	"context"
	"github.com/soukengo/gopkg/component/transport/queue/options"
	"github.com/soukengo/gopkg/util/runtimes"
)

type HandlerFunc func(context.Context, *BytesMessage) Action

type Callback func(Action)

type Handler struct {
	Func    HandlerFunc
	Options *options.ConsumerOptions
}

func (h *Handler) Process(message *BytesMessage, callback Callback) {
	ctx := context.TODO()
	opts := h.Options
	if opts.Mode() == options.Async {
		runtimes.Async(func() {
			callback(h.Func(ctx, message))
		})
	} else {
		callback(h.Func(ctx, message))
	}
}
