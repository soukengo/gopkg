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

type redisDelayed struct {
	Queue
	cli       *redis.Client
	topics    map[string]*delayItem
	ctx       context.Context
	cancel    context.CancelFunc
	startOnce *sync.Once
	stopOnce  *sync.Once
}

func NewRedisDelayed(cli *redis.Client, topics []Topic, workers int) Delayed {
	queue := NewRedisQueue(cli, topics, workers)
	ctx, cancel := context.WithCancel(context.Background())
	ins := &redisDelayed{Queue: queue, cli: cli, ctx: ctx, cancel: cancel, startOnce: new(sync.Once), stopOnce: new(sync.Once)}
	items := make(map[string]*delayItem)
	for _, topic := range topics {
		items[string(topic)] = newDelayItem(ins, topic)
	}
	ins.topics = items
	return ins
}

func (q *redisDelayed) Start() (err error) {
	q.startOnce.Do(func() {
		err = q.Queue.Start()
		if err != nil {
			return
		}
		for _, t := range q.topics {
			t.runQueue()
		}
	})
	return
}
func (q *redisDelayed) Stop() (err error) {
	q.stopOnce.Do(func() {
		q.cancel()
	})
	return
}
func (q *redisDelayed) SubmitDelay(topic Topic, values []string, delay int64, overwritten bool) (err error) {
	key := topic.waitKey()
	score := time.Now().UnixMilli() + delay
	members := make([]*redisdriver.Z, len(values))
	for i, value := range values {
		members[i] = &redisdriver.Z{Score: float64(score), Member: value}
	}
	if overwritten {
		_, err = q.cli.ZAdd(q.ctx, key, members...)
		return
	}
	_, err = q.cli.ZAddNX(q.ctx, key, members...)
	return
}

func (q *redisDelayed) RemoveDelay(topic Topic, values []string) (deleted bool, err error) {
	waitKey := topic.waitKey()
	members := make([]any, len(values))
	for i, value := range values {
		members[i] = value
	}
	count, err := q.cli.ZRem(q.ctx, waitKey, members...)
	if err != nil {
		return
	}
	deleted = count > 0
	return
}

type delayItem struct {
	queue *redisDelayed
	topic Topic
	once  *sync.Once
}

func newDelayItem(queue *redisDelayed, topic Topic) *delayItem {
	return &delayItem{queue: queue, topic: topic, once: new(sync.Once)}
}

func (t *delayItem) runQueue() {
	t.once.Do(func() {
		go runtimes.TryCatch(func() {
			t.producer()
		})
	})
}

// 从 wait队列移动到ready队列
func (t *delayItem) moveToReady() int {
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
	members := make([]any, 0)
	memberValues := make([]string, 0)
	for _, value := range values {
		v := value
		members = append(members, v.Member)
		memberValues = append(memberValues, v.Member.(string))
	}
	err := t.queue.Submit(t.topic, memberValues)
	if err != nil {
		log.Errorf("Submit: %v error: %v", memberValues, err)
		return 0
	}
	_, err = t.queue.cli.ZRem(ctx, waitKey, members...)
	if err != nil {
		log.Errorf("t.redis.ZRem waitKey: %v error: %v", waitKey, err)
		return 0
	}
	return 1
}

func (t *delayItem) producer() {
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
