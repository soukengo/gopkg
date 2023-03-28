package idgen

import "github.com/rs/xid"

type xidGenerator struct {
}

func NewXidGenerator() Generator {
	return &xidGenerator{}
}
func (x xidGenerator) GenID() (id string, err error) {
	id = xid.New().String()
	return
}
