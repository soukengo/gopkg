package queue

import (
	"context"
	redisdriver "github.com/go-redis/redis/v8"
	"github.com/soukengo/gopkg/component/database/redis"
	"github.com/soukengo/gopkg/log"
	"github.com/soukengo/gopkg/util/runtimes"
	"strconv"
	"sync"
	"time"
)

const (
	delayIdleTime      = time.Millisecond * 100
	delayConsumerBatch = 20
	defaultWorkers     = 1
)

type HandlerFunc func(topic Topic, value string)

type Queue interface {
	// Start 启动队列
	Start() error
	// Stop 停止队列
	Stop() error
	// Subscribe 订阅TOPIC
	Subscribe(topic Topic, handler HandlerFunc)
	// Submit 往队列提交数据
	Submit(topic Topic, value string) error
}

type Delayed interface {
	Queue
	// SubmitDelay 提交延时数据
	SubmitDelay(topic Topic, value string, delayMills int64, overwritten bool) error
	// RemoveDelay 从延时队列移除任务
	RemoveDelay(topic Topic, value string) (deleted bool, err error)
}

type redisDelayed struct {
	cli       *redis.Client
	topics    map[string]*queueItem
	ctx       context.Context
	cancel    context.CancelFunc
	startOnce *sync.Once
	stopOnce  *sync.Once
}

func NewRedisDelayed(cli *redis.Client, topics []Topic, workers int) Delayed {
	if workers <= 0 {
		workers = defaultWorkers
	}
	ctx, cancel := context.WithCancel(context.Background())
	queue := &redisDelayed{cli: cli, ctx: ctx, cancel: cancel, startOnce: new(sync.Once), stopOnce: new(sync.Once)}
	items := make(map[string]*queueItem)
	for _, topic := range topics {
		items[string(topic)] = newQueueTopic(queue, topic, workers)
	}
	queue.topics = items
	return queue
}

type queueItem struct {
	queue    *redisDelayed
	topic    Topic
	messages chan string
	once     *sync.Once
	handlers []HandlerFunc
}

func newQueueTopic(queue *redisDelayed, topic Topic, workers int) *queueItem {
	return &queueItem{
		queue:    queue,
		topic:    topic,
		messages: make(chan string, workers),
		once:     new(sync.Once),
		handlers: []HandlerFunc{},
	}
}

func (q *redisDelayed) Start() (err error) {
	q.startOnce.Do(func() {
		for _, t := range q.topics {
			t.runQueue()
		}
	})
	return
}
func (q *redisDelayed) Stop() (err error) {
	q.stopOnce.Do(func() {
		q.cancel()
		for _, t := range q.topics {
			close(t.messages)
		}
	})
	return
}

func (q *redisDelayed) Subscribe(topic Topic, handler HandlerFunc) {
	item := q.topics[string(topic)]
	item.handlers = append(item.handlers, handler)
}

func (q *redisDelayed) Submit(topic Topic, value string) (err error) {
	key := topic.readyKey()
	_, err = q.cli.SAdd(q.ctx, key, value)
	return
}

func (q *redisDelayed) SubmitDelay(topic Topic, value string, delay int64, overwritten bool) (err error) {
	key := topic.waitKey()
	score := time.Now().UnixMilli() + delay
	if overwritten {
		_, err = q.cli.ZAdd(q.ctx, key, &redisdriver.Z{Score: float64(score), Member: value})
		return
	}
	_, err = q.cli.ZAddNX(q.ctx, key, &redisdriver.Z{Score: float64(score), Member: value})
	return
}

func (q *redisDelayed) RemoveDelay(topic Topic, value string) (deleted bool, err error) {
	waitKey := topic.waitKey()
	count, err := q.cli.ZRem(q.ctx, waitKey, value)
	if err != nil {
		return
	}
	deleted = count > 0
	return
}

// runQueue 启动队列
func (t *queueItem) runQueue() {
	t.once.Do(func() {
		go runtimes.TryCatch(func() {
			t.producer()
		})
		go runtimes.TryCatch(func() {
			t.consumer()
		})
		go runtimes.TryCatch(func() {
			t.dispatch()
		})
	})
}

func (t *queueItem) producer() {
	for {
		select {
		case <-t.queue.ctx.Done():
			return
		default:
			runtimes.TryCatch(func() {
				res := t.moveToReady()
				// 未读取到数据
				if res == 0 {
					time.Sleep(delayIdleTime)
				}
			})
		}
	}
}

// 从 wait队列移动到ready队列
func (t *queueItem) moveToReady() int {
	ctx := t.queue.ctx
	lockKey := t.topic.lockKey()
	locked, _ := t.queue.cli.SetNX(ctx, lockKey, "1", 10*time.Second)
	defer t.queue.cli.Del(ctx, lockKey)
	if !locked {
		return 0
	}
	waitKey := t.topic.waitKey()
	max := time.Now().UnixMilli()
	opt := redisdriver.ZRangeBy{Min: "0", Max: strconv.Itoa(int(max)), Offset: 0, Count: 100}
	values, _ := t.queue.cli.ZRangeByScoreWithScores(ctx, waitKey, &opt)
	if len(values) == 0 {
		return 0
	}
	readyKey := t.topic.readyKey()
	members := make([]interface{}, 0)
	for _, value := range values {
		v := value
		members = append(members, v.Member)
	}
	_, err := t.queue.cli.SAdd(ctx, readyKey, members...)
	if err != nil {
		log.Errorf("t.redis.SAdd readyKey: %v error: %v", readyKey, err)
		return 0
	}
	_, err = t.queue.cli.ZRem(ctx, waitKey, members...)
	if err != nil {
		log.Errorf("t.redis.ZRem waitKey: %v error: %v", waitKey, err)
		return 0
	}
	return 1
}

// consumer 消费数据
func (t *queueItem) consumer() {
	key := t.topic.readyKey()
	for {
		select {
		case <-t.queue.ctx.Done():
			return
		default:
			values, err := t.queue.cli.SPopN(context.TODO(), key, delayConsumerBatch)
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
				t.messages <- value
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
