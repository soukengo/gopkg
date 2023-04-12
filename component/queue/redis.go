package queue

import (
	"context"
	redisdriver "github.com/go-redis/redis/v8"
	"github.com/soukengo/gopkg/component/database/redis"
	"github.com/soukengo/gopkg/log"
	"github.com/soukengo/gopkg/util/runtimes"
	"sync"
	"time"
)

var _ Queue = &redisQueue{}

type redisQueue struct {
	cli       *redis.Client
	topics    map[string]*queueItem
	ctx       context.Context
	cancel    context.CancelFunc
	startOnce *sync.Once
	stopOnce  *sync.Once
}

func NewRedisQueue(cli *redis.Client, topics []Topic, workers int) Queue {
	if workers <= 0 {
		workers = defaultWorkers
	}
	queue := newRedisQueue(cli)
	items := make(map[string]*queueItem)
	for _, topic := range topics {
		items[string(topic)] = newQueueTopic(queue, topic, workers)
	}
	queue.topics = items
	return queue
}

func newRedisQueue(cli *redis.Client) *redisQueue {
	ctx, cancel := context.WithCancel(context.Background())
	return &redisQueue{cli: cli, ctx: ctx, cancel: cancel, startOnce: new(sync.Once), stopOnce: new(sync.Once)}
}

type queueItem struct {
	queue    *redisQueue
	topic    Topic
	messages chan string
	once     *sync.Once
	handlers []HandlerFunc
}

func newQueueTopic(queue *redisQueue, topic Topic, workers int) *queueItem {
	return &queueItem{
		queue:    queue,
		topic:    topic,
		messages: make(chan string, workers),
		once:     new(sync.Once),
		handlers: []HandlerFunc{},
	}
}

func (q *redisQueue) Start() (err error) {
	q.startOnce.Do(func() {
		for _, t := range q.topics {
			t.runQueue()
		}
	})
	return
}
func (q *redisQueue) Stop() (err error) {
	q.stopOnce.Do(func() {
		q.cancel()
		for _, t := range q.topics {
			close(t.messages)
		}
	})
	return
}

func (q *redisQueue) Subscribe(topic Topic, handler HandlerFunc) {
	item := q.topics[string(topic)]
	item.handlers = append(item.handlers, handler)
}

func (q *redisQueue) Submit(topic Topic, values []string) (err error) {
	key := topic.readyKey()
	score := time.Now().UnixMilli()
	members := make([]*redisdriver.Z, len(values))
	for i, value := range values {
		members[i] = &redisdriver.Z{Score: float64(score), Member: value}
	}
	_, err = q.cli.ZAddNX(q.ctx, key, members...)
	return
}

// runQueue 启动队列
func (t *queueItem) runQueue() {
	t.once.Do(func() {
		go runtimes.TryCatch(func() {
			t.consumer()
		})
		go runtimes.TryCatch(func() {
			t.dispatch()
		})
	})
}

// consumer 消费数据
func (t *queueItem) consumer() {
	key := t.topic.readyKey()
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
				t.messages <- value.Member.(string)
			}
		}
	}
}

// dispatch 分发
func (t *queueItem) dispatch() {
	for {
		select {
		case <-t.queue.ctx.Done():
			return
		case msg := <-t.messages:
			var record = msg
			for _, h := range t.handlers {
				var h1 = h
				h1(t.topic, record)
			}
		}
	}
}
