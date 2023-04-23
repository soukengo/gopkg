package redis

import (
	"context"
	redisdriver "github.com/go-redis/redis/v8"
	"github.com/soukengo/gopkg/component/queue/core"
	"github.com/soukengo/gopkg/infra/storage/redis"
	"github.com/soukengo/gopkg/log"
	"github.com/soukengo/gopkg/util/runtimes"
	"strconv"
	"sync"
	"time"
)

type redisDelayProducer struct {
	core.Producer
	cli       *redis.Client
	topics    map[string]*delayItem
	ctx       context.Context
	cancel    context.CancelFunc
	startOnce *sync.Once
	stopOnce  *sync.Once
}

func NewDelayProducer(cfg *Config, logger log.Logger) core.DelayedProducer {
	cli := redis.NewClient(cfg.Config, logger)
	ctx, cancel := context.WithCancel(context.Background())
	ins := &redisDelayProducer{
		cli:       cli,
		Producer:  NewProducer(cfg, logger),
		ctx:       ctx,
		cancel:    cancel,
		startOnce: new(sync.Once),
		stopOnce:  new(sync.Once),
	}
	items := make(map[string]*delayItem)
	for _, topic := range cfg.Consumer.Topics {
		items[string(topic)] = newDelayItem(ins, topic)
	}
	ins.topics = items
	return ins
}

func (q *redisDelayProducer) Start() (err error) {
	q.startOnce.Do(func() {
		for _, t := range q.topics {
			t.runQueue()
		}
	})
	return
}
func (q *redisDelayProducer) Stop() (err error) {
	q.stopOnce.Do(func() {
		q.cancel()
	})
	return
}

func (q *redisDelayProducer) SubmitDelay(ctx context.Context, topic core.Topic, values []core.Value, delay int64, overwritten bool) (err error) {
	waitKey := delayKey(topic)
	score := time.Now().UnixMilli() + delay
	members := make([]*redisdriver.Z, len(values))
	for i, value := range values {
		members[i] = &redisdriver.Z{Score: float64(score), Member: string(value)}
	}
	if overwritten {
		_, err = q.cli.ZAdd(ctx, waitKey, members...)
		return
	}
	_, err = q.cli.ZAddNX(ctx, waitKey, members...)
	return
}

func (q *redisDelayProducer) RemoveDelay(ctx context.Context, topic core.Topic, values []core.Value) (deleted bool, err error) {
	waitKey := delayKey(topic)
	members := make([]any, len(values))
	for i, value := range values {
		members[i] = string(value)
	}
	count, err := q.cli.ZRem(ctx, waitKey, members...)
	if err != nil {
		return
	}
	deleted = count > 0
	return
}

type delayItem struct {
	queue *redisDelayProducer
	topic core.Topic
	once  *sync.Once
}

func newDelayItem(queue *redisDelayProducer, topic core.Topic) *delayItem {
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
	lockKey := delayLockKey(t.topic)
	locked, _ := t.queue.cli.SetNX(ctx, lockKey, "1", 10*time.Second)
	defer t.queue.cli.Del(ctx, lockKey)
	if !locked {
		return 0
	}
	waitKey := delayKey(t.topic)
	max := time.Now().UnixMilli()
	opt := redisdriver.ZRangeBy{Min: "0", Max: strconv.Itoa(int(max)), Offset: 0, Count: 100}
	values, _ := t.queue.cli.ZRangeByScoreWithScores(ctx, waitKey, &opt)
	if len(values) == 0 {
		return 0
	}
	members := make([]any, 0)
	memberValues := make([]core.Value, 0)
	for _, value := range values {
		v := value
		members = append(members, v.Member)
		memberValues = append(memberValues, core.Value(v.Member.(string)))
	}
	err := t.queue.Submit(t.queue.ctx, t.topic, memberValues)
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
