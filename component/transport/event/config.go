package event

import (
	"github.com/soukengo/gopkg/component/transport/queue"
	"github.com/soukengo/gopkg/component/transport/queue/provider/memory"
)

type Config struct {
	Memory      *memory.Config
	Distributed *queue.GeneralConfig
}
