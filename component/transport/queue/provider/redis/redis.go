package redis

import (
	"github.com/soukengo/gopkg/component/transport/queue/iface"
	"github.com/soukengo/gopkg/log"
	"time"
)

const (
	delayIdleTime      = time.Millisecond * 100
	delayConsumerBatch = 20
	defaultWorkers     = 1
	keyPrefix          = "iface."
	keyDelayPrefix     = keyPrefix + "delay.wait."
	keyDelayLockPrefix = keyPrefix + "delay.lock."
)

type redisQueue struct {
	iface.Consumer
	iface.Producer
}

func NewRedisQueue(cfg *Config, logger log.Logger) iface.Server {
	return &redisQueue{
		Consumer: NewConsumer(cfg, logger),
		Producer: NewProducer(cfg, logger),
	}
}

func readyKey(topic iface.Topic) string {
	return topic.GenKey(keyPrefix, string(topic))
}
func delayKey(topic iface.Topic) string {
	return topic.GenKey(keyDelayPrefix, string(topic))
}
func delayLockKey(topic iface.Topic) string {
	return topic.GenKey(keyDelayLockPrefix, string(topic))
}
