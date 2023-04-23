package redis

import (
	"context"
	xerrors "errors"
	"github.com/go-redis/redis/v8"
	"github.com/soukengo/gopkg/errors"
)

var (
	errNotFound = errors.NotFound("CACHE_NOT_FOUND", "cache not found")
)

type ErrorHook struct {
}

var _ redis.Hook = ErrorHook{}

func (ErrorHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (ErrorHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	if xerrors.Is(cmd.Err(), redis.Nil) {
		cmd.SetErr(errNotFound)
	}
	return nil
}

func (ErrorHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (ErrorHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	for _, cmd := range cmds {
		if xerrors.Is(cmd.Err(), redis.Nil) {
			cmd.SetErr(errNotFound)
		}
	}
	return nil
}
