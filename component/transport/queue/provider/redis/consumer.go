package redis

import (
	"context"
	"github.com/soukengo/gopkg/component/transport/queue/iface"
	"github.com/soukengo/gopkg/infra/storage/redis"
	"github.com/soukengo/gopkg/log"
	"github.com/soukengo/gopkg/util/dispatcher"
	"github.com/soukengo/gopkg/util/runtimes"
	"sync"
	"time"
)

type consumer struct {
	cfg       *Config
	cli       *redis.Client
	topics    map[iface.Topic]*redisQueueItem
	ctx       context.Context
	cancel    context.CancelFunc
	startOnce *sync.Once
	stopOnce  *sync.Once
	mutex     sync.Mutex
}

func NewConsumer(cfg *Config, logger log.Logger) iface.Consumer {
	cli := redis.NewClient(cfg.Config, logger)
	ctx, cancel := context.WithCancel(context.Background())
	ins := &consumer{cfg: cfg, cli: cli, ctx: ctx, cancel: cancel, startOnce: new(sync.Once), stopOnce: new(sync.Once)}
	items := make(map[iface.Topic]*redisQueueItem)
	ins.topics = items
	return ins
}

type redisQueueItem struct {
	queue    *consumer
	topic    iface.Topic
	messages chan *iface.BytesMessage
	once     *sync.Once
	handlers []*iface.Handler
}

func newQueueTopic(c *consumer, topic iface.Topic, workers int) *redisQueueItem {
	return &redisQueueItem{
		queue:    c,
		topic:    topic,
		messages: make(chan *iface.BytesMessage, workers),
		once:     new(sync.Once),
	}
}

func (c *consumer) Start(context.Context) (err error) {
	c.startOnce.Do(func() {
		for _, t := range c.topics {
			t.runQueue()
		}
	})
	return
}
func (c *consumer) Stop(context.Context) (err error) {
	c.stopOnce.Do(func() {
		c.cancel()
		for _, t := range c.topics {
			close(t.messages)
		}
	})
	return
}

func (c *consumer) Subscribe(topic iface.Topic, handler *iface.Handler) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	item, ok := c.topics[topic]
	if !ok {
		item = newQueueTopic(c, topic, c.cfg.Consumer.Workers)
	}
	item.handlers = append(item.handlers, handler)
	c.topics[topic] = item
}

// runQueue 启动队列
func (t *redisQueueItem) runQueue() {
	t.once.Do(func() {
		runtimes.Async(t.consumer)
		runtimes.Async(t.dispatch)
	})
}

// consumer 消费数据
func (t *redisQueueItem) consumer() {
	key := readyKey(t.topic)
	for {
		select {
		case <-t.queue.ctx.Done():
			return
		default:
			values, err := t.queue.cli.ZPopMin(context.TODO(), key, delayConsumerBatch)
			if err != nil {
				log.Errorf("redis.SPopN key: %s,error: %v", key, err)
				time.Sleep(delayIdleTime)
				break
			}
			if len(values) == 0 {
				time.Sleep(delayIdleTime)
				break
			}
			for _, value := range values {
				t.messages <- iface.NewBytesMessage(t.topic, []byte(value.Member.(string)))
			}
		}
	}
}

// dispatch 分发
func (t *redisQueueItem) dispatch() {
	dispatcher.Dispatch(t.queue.ctx, t.messages, func(message *iface.BytesMessage) {
		for _, item := range t.handlers {
			var handler = item
			handler.Process(message, func(action iface.Action) {
				// to do
			})
		}
	})
}
