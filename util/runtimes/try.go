package runtimes

import "github.com/soukengo/gopkg/log"

func Try(fun func(), handler func(e interface{})) {
	defer func() {
		if e := recover(); e != nil {
			handler(e)
		}
	}()
	fun()
}
func TryCatch(fun func()) {
	var handler = func(e interface{}) {
		log.Errorf(Stack(e))
	}
	defer func() {
		if e := recover(); e != nil {
			handler(e)
		}
	}()
	fun()
}
