package redis

import (
	"context"
	rds "github.com/go-redis/redis/v8"
	"github.com/soukengo/gopkg/component/cache"
	"github.com/soukengo/gopkg/component/database/redis"
	"github.com/soukengo/gopkg/errors"
	"rockimserver/apis/rockim/shared"
	"strconv"
)

type sortedSetCache[T any] struct {
	redisCache[T]
}

func NewSortedSetCache[T any](client *redis.Client, category *cache.Category, opts ...cache.Option) cache.SortedSetCache[T] {
	return newSortedSetCache[T](client, category, cache.Apply(category, opts...))
}

func newSortedSetCache[T any](cli *redis.Client, category *cache.Category, opts *cache.Options) cache.SortedSetCache[T] {
	return &sortedSetCache[T]{redisCache: newRedisCache[T](cli, category, opts)}
}

func (c *sortedSetCache[T]) Add(ctx context.Context, parts cache.KeyParts, member *cache.SortedMember[T]) (err error) {
	data, err := c.encode(member.Value)
	if err != nil {
		return
	}
	_, err = c.cli.ZAdd(ctx, c.key(parts), &rds.Z{Score: float64(member.Score), Member: data})
	return
}

func (c *sortedSetCache[T]) AddSlice(ctx context.Context, parts cache.KeyParts, values []*cache.SortedMember[T]) (err error) {
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

func (c *sortedSetCache[T]) Paginate(ctx context.Context, parts cache.KeyParts, paginate *shared.Paginating) (ret []*T, paginated *shared.Paginated, err error) {
	key := c.key(parts)
	total, err := c.cli.ZCard(ctx, key)
	if err != nil {
		return
	}
	paginated = &shared.Paginated{Total: total}
	if total == 0 {
		var exists bool
		exists, err = c.Exists(ctx, parts)
		if err != nil {
			return
		}
		if !exists {
			err = cache.ErrNoCache
			return
		}
		return
	}
	values, err := c.cli.ZRangeByScore(ctx, key, &rds.ZRangeBy{Count: int64(paginate.PageSize), Offset: paginate.Offset()})
	if err != nil {
		return
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

func (c *sortedSetCache[T]) Tail(ctx context.Context, parts cache.KeyParts, size int64) (ret []*T, err error) {
	key := c.key(parts)
	values, err := c.cli.ZRangeByScore(ctx, key, &rds.ZRangeBy{Count: size, Offset: 0})
	if err != nil {
		return
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

func (c *sortedSetCache[T]) DeleteMember(ctx context.Context, parts cache.KeyParts, v *T) (err error) {
	data, err := c.encode(v)
	if err != nil {
		return
	}
	_, err = c.cli.ZRem(ctx, c.key(parts), data)
	return
}

func (c *sortedSetCache[T]) ExistsMember(ctx context.Context, parts cache.KeyParts, v *T) (isMember bool, err error) {
	data, err := c.encode(v)
	if err != nil {
		return
	}
	_, err = c.cli.ZRank(ctx, c.key(parts), string(data))
	if err != nil {
		if errors.IsNotFound(err) {
			var exists bool
			exists, err = c.Exists(ctx, parts)
			if err != nil {
				return
			}
			if !exists {
				err = cache.ErrNoCache
				return
			}
			return
		}
	}
	return
}

func (c *sortedSetCache[T]) Count(ctx context.Context, parts cache.KeyParts) (count int64, err error) {
	count, err = c.cli.ZCard(ctx, c.key(parts))
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
			err = cache.ErrNoCache
			return
		}
	}
	return
}

func (c *sortedSetCache[T]) Clear(ctx context.Context, parts cache.KeyParts, keep uint64, reverse bool) (err error) {
	// 默认从最小开始删
	start := int64(0)
	stop := int64(-keep) + 1
	// 删除最大的
	if reverse {
		start = int64(keep)
		stop = -1
	}
	if keep == 0 {
		start = 0
		stop = -1
	}
	_, err = c.cli.ZRemRangeByRank(ctx, c.key(parts), start, stop)
	return
}
func (c *sortedSetCache[T]) ClearByScore(ctx context.Context, parts cache.KeyParts, score int64, reverse bool) (err error) {
	scoreStr := strconv.Itoa(int(score))
	min := redis.Min
	max := scoreStr
	if reverse {
		min = scoreStr
		max = redis.Max
	}
	_, err = c.cli.ZRemRangeByScore(ctx, c.key(parts), min, max)
	return
}
