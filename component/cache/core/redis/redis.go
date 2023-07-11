package redis

import (
	"context"
	"github.com/soukengo/gopkg/component/cache/core"
	"github.com/soukengo/gopkg/infra/storage/redis"
	"github.com/soukengo/gopkg/log"
)

const (
	nullValue = "null"
)

type redisCache[T any] struct {
	cli      *redis.Client
	category *core.Category
	opts     *core.Options
}

func newRedisCache[T any](cli *redis.Client, category *core.Category, opts *core.Options) redisCache[T] {
	return redisCache[T]{cli: cli, category: category, opts: opts}
}

func (c *redisCache[T]) Exists(ctx context.Context, parts core.KeyParts) (bool, error) {
	exists, err := c.cli.Exists(ctx, c.key(parts))
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func (c *redisCache[T]) Delete(ctx context.Context, parts core.KeyParts) (err error) {
	_, err = c.cli.Del(ctx, c.key(parts))
	return
}

func (c *redisCache[T]) key(parts core.KeyParts) string {
	return c.category.GenKey(parts)
}

func (c *redisCache[T]) encode(v any) (data []byte, err error) {
	return c.opts.Codec().Encode(v)
}

func (c *redisCache[T]) decodeStr(data string) (ret *T, err error) {
	if data == nullValue {
		err = core.ErrNotFound
		return
	}
	b := []byte(data)
	return c.decode(b)
}
func (c *redisCache[T]) decode(data []byte) (ret *T, err error) {
	var value = new(T)
	err = c.opts.Codec().Decode(data, value)
	if err != nil {
		return
	}
	return value, nil
}

func (c *redisCache[T]) setTTL(ctx context.Context, parts core.KeyParts) {
	cacheKey := c.key(parts)
	if c.opts.Expire() <= 0 {
		return
	}
	_, err := c.cli.Expire(ctx, cacheKey, c.opts.Expire())
	if err != nil {
		log.WithContext(ctx).Errorf("expire cacheKey err: %v", err)
	}
}
