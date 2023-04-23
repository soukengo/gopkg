package redis

import (
	"context"
	"github.com/soukengo/gopkg/component/cache/core"
	"github.com/soukengo/gopkg/errors"
	"github.com/soukengo/gopkg/infra/storage/redis"
)

type valueCache[T any] struct {
	redisCache[T]
}

func NewValueCache[T any](client *redis.Client, category *core.Category, opts ...core.Option) core.ValueCache[T] {
	return newValueCache[T](client, category, core.Apply(category, opts...))
}

func newValueCache[T any](cli *redis.Client, category *core.Category, opts *core.Options) core.ValueCache[T] {
	return &valueCache[T]{redisCache: newRedisCache[T](cli, category, opts)}
}

func (c *valueCache[T]) Set(ctx context.Context, parts core.KeyParts, v *T) (err error) {
	data, err := c.encode(v)
	if err != nil {
		return
	}
	expire := c.opts.Expire()
	if len(data) == 0 {
		expire = c.opts.EmptyExpire()
	}
	_, err = c.cli.Set(ctx, c.key(parts), data, expire)
	return
}

func (c *valueCache[T]) Get(ctx context.Context, parts core.KeyParts) (ret *T, err error) {
	v, err := c.cli.Get(ctx, c.key(parts))
	if err != nil {
		if errors.IsNotFound(err) {
			err = core.ErrNoCache
		}
		return
	}
	ret, err = c.decodeStr(v)
	return
}
func (c *valueCache[T]) List(ctx context.Context, parts core.KeyParts, ids []string) (ret []*T, err error) {
	var keys []string
	for _, id := range ids {
		values := parts.Parts()
		values = append(values, id)
		keys = append(keys, c.key(core.Parts(values...)))
	}
	list, err := c.cli.MGet(ctx, keys...)
	if err != nil {
		if errors.IsNotFound(err) {
			err = core.ErrNoCache
		}
		return
	}
	ret = make([]*T, 0)
	for _, v := range list {
		str, ok := v.(string)
		if !ok {
			continue
		}
		var item *T
		item, err = c.decodeStr(str)
		if err != nil {
			return
		}
		ret = append(ret, item)
	}
	return
}
