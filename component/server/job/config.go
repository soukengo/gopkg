package job

import (
	"github.com/soukengo/gopkg/component/queue"
	"github.com/soukengo/gopkg/infra/storage"
)

type Config struct {
	Queue *queue.Config
}

func (c Config) Parse(configs *storage.Config) {
	if c.Queue != nil {
		c.Queue.Parse(configs)
	}
}
