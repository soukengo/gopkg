package socket

import "sync"

type attributes struct {
	session any
	lock    sync.RWMutex
}

func newAttributes() *attributes {
	return &attributes{}
}

func (a *attributes) SetSession(session any) {
	a.session = session
}
func (a *attributes) Session() any {
	return a.session
}
