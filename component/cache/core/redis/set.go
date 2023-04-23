package redis

import (
	"context"
	"github.com/soukengo/gopkg/component/cache/core"
	"github.com/soukengo/gopkg/infra/storage/redis"
)

type setCache[T any] struct {
	redisCache[T]
}

func NewSetCache[T any](client *redis.Client, category *core.Category, opts ...core.Option) core.SetCache[T] {
	return newSetCache[T](client, category, core.Apply(category, opts...))
}

func newSetCache[T any](cli *redis.Client, category *core.Category, opts *core.Options) core.SetCache[T] {
	return &setCache[T]{redisCache: newRedisCache[T](cli, category, opts)}
}

func (c *setCache[T]) Add(ctx context.Context, parts core.KeyParts, v *T) (err error) {
	data, err := c.encode(v)
	if err != nil {
		return
	}
	_, err = c.cli.SAdd(ctx, c.key(parts), data)
	return
}

func (c *setCache[T]) AddSlice(ctx context.Context, parts core.KeyParts, values []*T) (err error) {
	var records = make([]any, 0)
	for i, v := range values {
		var data []byte
		data, err = c.encode(v)
		if err != nil {
			return
		}
		records[i] = data
	}
	_, err = c.cli.SAdd(ctx, c.key(parts), records...)
	return
}

func (c *setCache[T]) Scan(ctx context.Context, parts core.KeyParts, cursor uint64, count int64) (ret []*T, retCursor uint64, err error) {
	values, retCursor, err := c.cli.SScan(ctx, c.key(parts), cursor, "*", count)
	if err != nil {
		return
	}
	if retCursor == 0 {
		var exists bool
		exists, err = c.Exists(ctx, parts)
		if err != nil {
			return
		}
		if !exists {
			err = core.ErrNoCache
			return
		}
	}
	ret = make([]*T, len(values))
	for i, v := range values {
		var item *T
		item, err = c.decodeStr(v)
		if err != nil {
			return
		}
		ret[i] = item
	}
	return
}

func (c *setCache[T]) DeleteMember(ctx context.Context, parts core.KeyParts, v *T) (err error) {
	data, err := c.encode(v)
	if err != nil {
		return
	}
	_, err = c.cli.SRem(ctx, c.key(parts), data)
	return
}

func (c *setCache[T]) ExistsMember(ctx context.Context, parts core.KeyParts, v *T) (isMember bool, err error) {
	data, err := c.encode(v)
	if err != nil {
		return
	}
	isMember, err = c.cli.SIsMember(ctx, c.key(parts), data)
	if err != nil {
		return
	}
	if !isMember {
		var exists bool
		exists, err = c.Exists(ctx, parts)
		if err != nil {
			return
		}
		if !exists {
			err = core.ErrNoCache
			return
		}
	}
	return
}

func (c *setCache[T]) Count(ctx context.Context, parts core.KeyParts) (count int64, err error) {
	count, err = c.cli.SCard(ctx, c.key(parts))
	if err != nil {
		return
	}
	if count == 0 {
		var exists bool
		exists, err = c.Exists(ctx, parts)
		if err != nil {
			return
		}
		if !exists {
			err = core.ErrNoCache
			return
		}
	}
	return
}
