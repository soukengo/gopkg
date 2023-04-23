package lock

import (
	"context"
	"github.com/google/uuid"
	"github.com/soukengo/gopkg/infra/storage/redis"
	"github.com/soukengo/gopkg/log"
	"time"
)

const (
	redisKeyPrefix = "LOCK:"
)

type redisDistributedLock struct {
	logger  log.Logger
	cli     *redis.Client
	lockKey string
	lockVal string
}

func NewRedisDistributedLock(logger log.Logger, cli *redis.Client, key string) DistributedLock {
	return &redisDistributedLock{logger: logger, cli: cli, lockVal: genLockKey(key)}
}

func (lock *redisDistributedLock) TryLock(ctx context.Context) (locked bool) {
	locked, err := lock.addLock(ctx, defaultExpire)
	if err != nil {
		lock.logger.Helper().WithContext(ctx).Errorf("addLock error(%v)", err)
		return
	}
	return
}

func (lock *redisDistributedLock) Lock(ctx context.Context) bool {
	return lock.LockUntil(ctx, int(defaultWait.Seconds()))
}

func (lock *redisDistributedLock) LockUntil(ctx context.Context, waitSeconds int) (locked bool) {
	var countDown = time.Second * time.Duration(waitSeconds)
	for countDown.Seconds() > 0 {
		var err error
		locked, err = lock.addLock(ctx, defaultExpire)
		if err != nil {
			lock.logger.Helper().WithContext(ctx).Errorf("addLock error(%v)", err)
			return
		}
		if locked {
			break
		}
		time.Sleep(delay)
		countDown = countDown - delay
	}
	return
}

func (lock *redisDistributedLock) UnLock() {
	if len(lock.lockVal) == 0 {
		return
	}
	ctx := context.Background()
	err := lock.delLock(ctx)
	if err != nil {
		lock.logger.Helper().Errorf("delLock lockKey: %s,  error(%v)", lock.lockKey, err)
	}
}

func (lock *redisDistributedLock) addLock(ctx context.Context, expire time.Duration) (locked bool, err error) {
	lockVal := uuid.New().String()
	locked, err = lock.cli.SetNX(ctx, lock.lockKey, lockVal, expire)
	if err != nil || !locked {
		return
	}
	lock.lockVal = lockVal
	return
}

func (lock *redisDistributedLock) delLock(ctx context.Context) (err error) {
	val, err := lock.cli.Get(ctx, lock.lockKey)
	if err != nil {
		return
	}
	if val == lock.lockVal {
		_, err = lock.cli.Del(ctx, lock.lockKey)
		return err
	}
	return
}

func genLockKey(key string) string {
	return redisKeyPrefix + key
}
