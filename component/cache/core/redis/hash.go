package redis

import (
	"context"
	"github.com/soukengo/gopkg/component/cache/core"
	"github.com/soukengo/gopkg/errors"
	"github.com/soukengo/gopkg/infra/storage/redis"
)

type hashCache[T any] struct {
	redisCache[T]
}

func NewHashCache[T any](client *redis.Client, category *core.Category, opts ...core.Option) core.HashCache[T] {
	return newHashCache[T](client, category, core.Apply(category, opts...))
}

func newHashCache[T any](cli *redis.Client, category *core.Category, opts *core.Options) core.HashCache[T] {
	return &hashCache[T]{redisCache: newRedisCache[T](cli, category, opts)}
}

func (c *hashCache[T]) Put(ctx context.Context, parts core.KeyParts, field string, v *T) (err error) {
	data, err := c.encode(v)
	if err != nil {
		return
	}
	_, err = c.cli.HSet(ctx, c.key(parts), field, data)
	if err != nil {
		return
	}
	c.setTTL(ctx, parts)
	return
}

func (c *hashCache[T]) PutMap(ctx context.Context, parts core.KeyParts, hash map[string]*T) (err error) {
	var records = make(map[string]any)
	for k, v := range hash {
		var data []byte
		data, err = c.encode(v)
		if err != nil {
			return
		}
		records[k] = data
	}
	_, err = c.cli.HMSet(ctx, c.key(parts), records)
	if err != nil {
		return
	}
	c.setTTL(ctx, parts)
	return
}

func (c *hashCache[T]) Get(ctx context.Context, parts core.KeyParts, field string) (ret *T, err error) {
	v, err := c.cli.HGet(ctx, c.key(parts), field)
	if err != nil {
		if errors.IsNotFound(err) {
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
	return c.decodeStr(v)
}

func (c *hashCache[T]) MGet(ctx context.Context, parts core.KeyParts, fields []string) (ret []*T, err error) {
	list, err := c.cli.HMGet(ctx, c.key(parts), fields...)
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
func (c *hashCache[T]) GetAll(ctx context.Context, parts core.KeyParts) (ret map[string]*T, err error) {
	hash, err := c.cli.HGetAll(ctx, c.key(parts))
	if err != nil {
		if errors.IsNotFound(err) {
			err = core.ErrNoCache
		}
		return
	}
	ret = make(map[string]*T)
	for k, v := range hash {
		var item *T
		item, err = c.decodeStr(v)
		if err != nil {
			return
		}
		ret[k] = item
	}
	return
}

func (c *hashCache[T]) ExistsField(ctx context.Context, parts core.KeyParts, field string) (existsField bool, err error) {
	existsField, err = c.cli.HExists(ctx, c.key(parts), field)
	if err != nil {
		return
	}
	if !existsField {
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

func (c *hashCache[T]) DeleteField(ctx context.Context, parts core.KeyParts, field string) (err error) {
	_, err = c.cli.HDel(ctx, c.key(parts), field)
	return
}
