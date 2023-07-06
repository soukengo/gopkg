package runtimes

import (
	"github.com/panjf2000/ants/v2"
	"github.com/soukengo/gopkg/log"
)

const (
	maxPoolSize = ants.DefaultAntsPoolSize
)

var (
	runPool *ants.Pool
)

func init() {
	runPool, _ = ants.NewPool(maxPoolSize)
}

func Async(task func()) {
	var err error
	if runPool == nil {
		err = ants.Submit(task)
	} else {
		err = runPool.Submit(task)
	}
	if err != nil {
		log.Errorf("Async err: %v", err)
	}
}
