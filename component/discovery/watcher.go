package discovery

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/soukengo/gopkg/log"
	"time"
)

type watcherHelper struct {
	serviceId string
	ctx       context.Context
	cancel    context.CancelFunc
	w         registry.Watcher
	onUpdate  OnInstanceUpdate
}

type OnInstanceUpdate func(ins []*registry.ServiceInstance)

func Watch(dis registry.Discovery, serviceId string, timeout time.Duration, onUpdate OnInstanceUpdate) (err error) {
	var w registry.Watcher
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{}, 1)
	go func() {
		w, err = dis.Watch(ctx, serviceId)
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(timeout):
		err = errors.New("discovery create watcher overtime")
	}
	if err != nil {
		cancel()
		return
	}
	h := &watcherHelper{
		serviceId: serviceId,
		ctx:       ctx,
		cancel:    cancel,
		w:         w,
		onUpdate:  onUpdate,
	}
	go h.watch()
	return
}

func (h *watcherHelper) watch() {
	for {
		select {
		case <-h.ctx.Done():
			return
		default:
		}
		ins, err := h.w.Next()
		if err != nil {
			if errors.Is(err, context.Canceled) {
				log.Errorf("[resolver] Failed to Watch discovery endpoint: %v", err)
				return
			}
			log.Errorf("[resolver] Failed to Watch discovery endpoint: %v", err)
			time.Sleep(time.Second)
			continue
		}
		h.onUpdate(ins)
	}
}
