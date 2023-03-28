package lock

import (
	"context"
	"time"
)

const (
	defaultExpire = time.Second * 60
	delay         = time.Millisecond * 100
	defaultWait   = time.Second * 60
)

type DistributedLock interface {
	TryLock(ctx context.Context) bool
	Lock(ctx context.Context) bool
	LockUntil(ctx context.Context, waitSeconds int) bool
	UnLock()
}
