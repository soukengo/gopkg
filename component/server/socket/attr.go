package socket

import "sync"

type attributes struct {
	attrs map[string]any
	lock  sync.RWMutex
}

func newAttributes() *attributes {
	return &attributes{
		attrs: map[string]any{},
	}
}

func (a *attributes) SetAttr(key string, value any) {
	a.lock.Lock()
	a.attrs[key] = value
	a.lock.Unlock()
}
func (a *attributes) DelAttr(key string) {
	a.lock.Lock()
	delete(a.attrs, key)
	a.lock.Unlock()
}

func (a *attributes) Attr(key string) (value any, ok bool) {
	value, ok = a.attrs[key]
	return
}

func (a *attributes) StringAttr(key string) (value string) {
	v, ok := a.Attr(key)
	if !ok {
		return
	}
	value, ok = v.(string)
	return
}

func (a *attributes) Int64Attr(key string) (value int64) {
	v, ok := a.Attr(key)
	if !ok {
		return
	}
	value, ok = v.(int64)
	return
}
