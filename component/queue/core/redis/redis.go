package redis

import (
	"github.com/soukengo/gopkg/component/queue/core"
	"github.com/soukengo/gopkg/log"
	"time"
)

const (
	delayIdleTime      = time.Millisecond * 100
	delayConsumerBatch = 20
	defaultWorkers     = 1
	keyPrefix          = "core."
	keyDelayPrefix     = keyPrefix + "delay.wait."
	keyDelayLockPrefix = keyPrefix + "delay.lock."
)

type redisQueue struct {
	core.Consumer
	core.Producer
}

func NewRedisQueue(cfg *Config, logger log.Logger) core.Queue {
	return &redisQueue{
		Consumer: NewConsumer(cfg, logger),
		Producer: NewProducer(cfg, logger),
	}
}

func readyKey(topic core.Topic) string {
	return topic.GenKey(keyPrefix, string(topic))
}
func delayKey(topic core.Topic) string {
	return topic.GenKey(keyDelayPrefix, string(topic))
}
func delayLockKey(topic core.Topic) string {
	return topic.GenKey(keyDelayLockPrefix, string(topic))
}
