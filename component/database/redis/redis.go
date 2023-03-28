package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/soukengo/gopkg/log"
	"time"
)

type Client struct {
	config            *Config
	redis             redis.UniversalClient
	keyChangeHandlers map[string]OnKeyChanged
}

type OnKeyChanged func(key string)

func NewClient(conf *Config, logger log.Logger) *Client {
	c := new(Client)
	c.config = conf
	c.keyChangeHandlers = make(map[string]OnKeyChanged)
	redisOpts := &redis.Options{
		Addr:         c.config.Addr[0],
		MinIdleConns: c.config.Idle,
		PoolSize:     c.config.Active,
		IdleTimeout:  time.Duration(c.config.IdleTimeout),
	}
	if len(c.config.Password) > 0 {
		redisOpts.Password = c.config.Password
	}
	r := redis.NewClient(redisOpts)
	r.AddHook(ErrorHook{})
	r.AddHook(TracingHook{})
	pong, err := r.Ping(context.Background()).Result()
	logger.Helper().Infof("redis ping result val=%v err=%v", pong, err)
	if err != nil {
		panic(err)
	}
	c.redis = r
	return c
}

func (c *Client) Conn() redis.UniversalClient {
	return c.redis
}

func (c *Client) genKey(key string) string {
	return c.config.Prefix + key
}

func (c *Client) Close() error {
	return c.redis.Close()
}

func (c *Client) Exists(ctx context.Context, keys ...string) (int64, error) {
	var nkeys = make([]string, len(keys))
	for i, value := range keys {
		nkeys[i] = c.genKey(value)
	}
	return c.redis.Exists(ctx, nkeys...).Result()
}
func (c *Client) Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return c.redis.Expire(ctx, c.genKey(key), expiration).Result()
}
func (c *Client) ExpireAt(ctx context.Context, key string, tm time.Time) (bool, error) {
	return c.redis.ExpireAt(ctx, c.genKey(key), tm).Result()
}

func (c *Client) TTL(ctx context.Context, key string) (time.Duration, error) {
	return c.redis.TTL(ctx, c.genKey(key)).Result()
}
func (c *Client) Type(ctx context.Context, key string) (string, error) {
	return c.redis.Type(ctx, c.genKey(key)).Result()
}
func (c *Client) Append(ctx context.Context, key, value string) (int64, error) {
	return c.redis.Append(ctx, c.genKey(key), value).Result()
}
func (c *Client) Decr(ctx context.Context, key string) (int64, error) {
	return c.redis.Decr(ctx, c.genKey(key)).Result()
}
func (c *Client) DecrBy(ctx context.Context, key string, decrement int64) (int64, error) {
	return c.redis.DecrBy(ctx, c.genKey(key), decrement).Result()
}
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.redis.Get(ctx, c.genKey(key)).Result()
}
func (c *Client) GetRange(ctx context.Context, key string, start, end int64) (string, error) {
	return c.redis.GetRange(ctx, c.genKey(key), start, end).Result()
}
func (c *Client) GetSet(ctx context.Context, key string, value any) (string, error) {
	return c.redis.GetSet(ctx, c.genKey(key), value).Result()
}
func (c *Client) Incr(ctx context.Context, key string) (int64, error) {
	return c.redis.Incr(ctx, c.genKey(key)).Result()
}
func (c *Client) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	return c.redis.IncrBy(ctx, c.genKey(key), value).Result()
}
func (c *Client) IncrByFloat(ctx context.Context, key string, value float64) (float64, error) {
	return c.redis.IncrByFloat(ctx, c.genKey(key), value).Result()
}
func (c *Client) MGet(ctx context.Context, keys ...string) ([]any, error) {
	var nkeys = make([]string, len(keys))
	for i, value := range keys {
		nkeys[i] = c.genKey(value)
	}
	return c.redis.MGet(ctx, nkeys...).Result()
}

func (c *Client) HExists(ctx context.Context, key, field string) (bool, error) {
	return c.redis.HExists(ctx, c.genKey(key), field).Result()
}
func (c *Client) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return c.redis.HGetAll(ctx, c.genKey(key)).Result()
}
func (c *Client) HGet(ctx context.Context, key, field string) (string, error) {
	return c.redis.HGet(ctx, c.genKey(key), field).Result()
}

func (c *Client) HMGet(ctx context.Context, key string, fields ...string) ([]any, error) {
	return c.redis.HMGet(ctx, c.genKey(key), fields...).Result()
}
func (c *Client) HSet(ctx context.Context, key string, values ...any) (int64, error) {
	return c.redis.HSet(ctx, c.genKey(key), values...).Result()
}
func (c *Client) HMSet(ctx context.Context, key string, values ...any) (bool, error) {
	return c.redis.HMSet(ctx, c.genKey(key), values...).Result()
}
func (c *Client) HSetNX(ctx context.Context, key, field string, value any) (bool, error) {
	return c.redis.HSetNX(ctx, c.genKey(key), field, value).Result()
}

func (c *Client) HDel(ctx context.Context, key string, fields ...string) (int64, error) {
	return c.redis.HDel(ctx, c.genKey(key), fields...).Result()
}

func (c *Client) Set(ctx context.Context, key string, value any, expiration time.Duration) (string, error) {
	return c.redis.Set(ctx, c.genKey(key), value, expiration).Result()
}
func (c *Client) SetWithoutResult(ctx context.Context, key string, value any, expiration time.Duration) error {
	return c.redis.Set(ctx, c.genKey(key), value, expiration).Err()
}

func (c *Client) SetNX(ctx context.Context, key string, value any, expiration time.Duration) (bool, error) {
	return c.redis.SetNX(ctx, c.genKey(key), value, expiration).Result()
}
func (c *Client) SetXX(ctx context.Context, key string, value any, expiration time.Duration) (bool, error) {
	return c.redis.SetXX(ctx, c.genKey(key), value, expiration).Result()
}

func (c *Client) Del(ctx context.Context, keys ...string) (int64, error) {
	var nkeys = make([]string, len(keys))
	for i, value := range keys {
		nkeys[i] = c.genKey(value)
	}
	return c.redis.Del(ctx, nkeys...).Result()
}

func (c *Client) LLen(ctx context.Context, key string) (int64, error) {
	return c.redis.LLen(ctx, c.genKey(key)).Result()
}
func (c *Client) LPop(ctx context.Context, key string) (string, error) {
	return c.redis.LPop(ctx, c.genKey(key)).Result()
}
func (c *Client) LPush(ctx context.Context, key string, values ...any) (int64, error) {
	return c.redis.LPush(ctx, c.genKey(key), values...).Result()
}

func (c *Client) RPop(ctx context.Context, key string) (string, error) {
	return c.redis.RPop(ctx, c.genKey(key)).Result()
}
func (c *Client) BRPop(ctx context.Context, timeout time.Duration, key string) ([]string, error) {
	return c.redis.BRPop(ctx, timeout, c.genKey(key)).Result()
}
func (c *Client) RPopCount(ctx context.Context, key string, count int) ([]string, error) {
	return c.redis.RPopCount(ctx, c.genKey(key), count).Result()
}
func (c *Client) SAdd(ctx context.Context, key string, members ...any) (int64, error) {
	return c.redis.SAdd(ctx, c.genKey(key), members...).Result()
}
func (c *Client) SCard(ctx context.Context, key string) (int64, error) {
	return c.redis.SCard(ctx, c.genKey(key)).Result()
}
func (c *Client) SDiff(ctx context.Context, keys ...string) ([]string, error) {
	var nkeys []string
	for _, value := range keys {
		nkeys = append(nkeys, c.genKey(value))
	}
	return c.redis.SDiff(ctx, nkeys...).Result()
}
func (c *Client) SDiffStore(ctx context.Context, destination string, keys ...string) (int64, error) {
	var nkeys = make([]string, len(keys))
	for i, value := range keys {
		nkeys[i] = c.genKey(value)
	}
	return c.redis.SDiffStore(ctx, c.genKey(destination), nkeys...).Result()
}
func (c *Client) SInter(ctx context.Context, keys ...string) ([]string, error) {
	var nkeys = make([]string, len(keys))
	for i, value := range keys {
		nkeys[i] = c.genKey(value)
	}
	return c.redis.SInter(ctx, nkeys...).Result()
}
func (c *Client) SInterStore(ctx context.Context, destination string, keys ...string) (int64, error) {
	var nkeys = make([]string, len(keys))
	for i, value := range keys {
		nkeys[i] = c.genKey(value)
	}
	return c.redis.SInterStore(ctx, c.genKey(destination), nkeys...).Result()
}
func (c *Client) SIsMember(ctx context.Context, key string, member any) (bool, error) {
	return c.redis.SIsMember(ctx, c.genKey(key), member).Result()
}
func (c *Client) SMembers(ctx context.Context, key string) ([]string, error) {
	return c.redis.SMembers(ctx, c.genKey(key)).Result()
}
func (c *Client) SMembersMap(ctx context.Context, key string) (map[string]struct{}, error) {
	return c.redis.SMembersMap(ctx, c.genKey(key)).Result()
}
func (c *Client) SMove(ctx context.Context, source, destination string, member any) (bool, error) {
	return c.redis.SMove(ctx, c.genKey(source), c.genKey(destination), member).Result()
}
func (c *Client) SPop(ctx context.Context, key string) (string, error) {
	return c.redis.SPop(ctx, c.genKey(key)).Result()
}
func (c *Client) SPopN(ctx context.Context, key string, count int64) ([]string, error) {
	return c.redis.SPopN(ctx, c.genKey(key), count).Result()
}
func (c *Client) SRandMember(ctx context.Context, key string) (string, error) {
	return c.redis.SRandMember(ctx, c.genKey(key)).Result()
}
func (c *Client) SRandMemberN(ctx context.Context, key string, count int64) ([]string, error) {
	return c.redis.SRandMemberN(ctx, c.genKey(key), count).Result()
}
func (c *Client) SRem(ctx context.Context, key string, members ...any) (int64, error) {
	return c.redis.SRem(ctx, c.genKey(key), members...).Result()
}
func (c *Client) SUnion(ctx context.Context, keys ...string) ([]string, error) {
	var nkeys = make([]string, len(keys))
	for i, value := range keys {
		nkeys[i] = c.genKey(value)
	}
	return c.redis.SUnion(ctx, nkeys...).Result()
}
func (c *Client) SUnionStore(ctx context.Context, destination string, keys ...string) (int64, error) {
	var nkeys []string
	for _, value := range keys {
		nkeys = append(nkeys, c.genKey(value))
	}
	return c.redis.SUnionStore(ctx, c.genKey(destination), nkeys...).Result()
}

func (c *Client) SScan(ctx context.Context, key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return c.redis.SScan(ctx, c.genKey(key), cursor, match, count).Result()
}

func (c *Client) ZAdd(ctx context.Context, key string, members ...*redis.Z) (int64, error) {
	return c.redis.ZAdd(ctx, c.genKey(key), members...).Result()
}

func (c *Client) ZAddNX(ctx context.Context, key string, members ...*redis.Z) (int64, error) {
	return c.redis.ZAddNX(ctx, c.genKey(key), members...).Result()
}

func (c *Client) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return c.redis.ZRange(ctx, c.genKey(key), start, stop).Result()
}
func (c *Client) ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) ([]string, error) {
	return c.redis.ZRangeByScore(ctx, c.genKey(key), opt).Result()
}
func (c *Client) ZRangeByScoreWithScores(ctx context.Context, key string, opt *redis.ZRangeBy) ([]redis.Z, error) {
	return c.redis.ZRangeByScoreWithScores(ctx, c.genKey(key), opt).Result()
}
func (c *Client) ZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error) {
	return c.redis.ZRangeWithScores(ctx, c.genKey(key), start, stop).Result()
}
func (c *Client) ZRevRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) ([]string, error) {
	return c.redis.ZRevRangeByScore(ctx, c.genKey(key), opt).Result()
}

func (c *Client) ZRem(ctx context.Context, key string, members ...any) (int64, error) {
	return c.redis.ZRem(ctx, c.genKey(key), members...).Result()
}

func (c *Client) ZRemRangeByScore(ctx context.Context, key string, min, max string) (int64, error) {
	return c.redis.ZRemRangeByScore(ctx, c.genKey(key), min, max).Result()
}
func (c *Client) ZRemRangeByRank(ctx context.Context, key string, start, stop int64) (int64, error) {
	return c.redis.ZRemRangeByRank(ctx, c.genKey(key), start, stop).Result()
}
func (c *Client) ZCard(ctx context.Context, key string) (int64, error) {
	return c.redis.ZCard(ctx, c.genKey(key)).Result()
}

func (c *Client) ZPopMin(ctx context.Context, key string, count int64) ([]redis.Z, error) {
	return c.redis.ZPopMin(ctx, c.genKey(key), count).Result()
}

func (c *Client) ZPopMax(ctx context.Context, key string, count int64) ([]redis.Z, error) {
	return c.redis.ZPopMax(ctx, c.genKey(key), count).Result()
}
func (c *Client) ZRank(ctx context.Context, key string, member string) (int64, error) {
	return c.redis.ZRank(ctx, c.genKey(key), member).Result()
}

func (c *Client) ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return c.redis.ZScan(ctx, c.genKey(key), cursor, match, count).Result()
}
