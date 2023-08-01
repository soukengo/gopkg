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
	keyPrefix          = "queue"
	separator          = "."
	keyDelayPrefix     = keyPrefix + separator + "delay.wait"
	keyDelayLockPrefix = keyPrefix + separator + "delay.lock"
)

type Queue struct {
	iface.Consumer
	iface.Producer
}

func NewServer(cfg *Config, logger log.Logger) iface.Server {
	return &Queue{
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
