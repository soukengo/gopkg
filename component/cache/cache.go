package cache

import (
	"context"
	"rockimserver/apis/rockim/shared"
)

type Cache[T any] interface {
	Exists(ctx context.Context, parts KeyParts) (bool, error)
	Delete(ctx context.Context, parts KeyParts) error
}

type ValueCache[T any] interface {
	Cache[T]
	Get(ctx context.Context, parts KeyParts) (*T, error)
	Set(ctx context.Context, parts KeyParts, value *T) error
	List(ctx context.Context, parts KeyParts, ids []string) ([]*T, error)
}

type HashCache[T any] interface {
	Cache[T]
	Put(ctx context.Context, parts KeyParts, field string, value *T) error
	PutMap(ctx context.Context, parts KeyParts, hash map[string]*T) error
	Get(ctx context.Context, parts KeyParts, field string) (*T, error)
	MGet(ctx context.Context, parts KeyParts, fields []string) ([]*T, error)
	GetAll(ctx context.Context, parts KeyParts) (map[string]*T, error)
	ExistsField(ctx context.Context, parts KeyParts, field string) (bool, error)
	DeleteField(ctx context.Context, parts KeyParts, field string) error
}

type SetCache[T any] interface {
	Cache[T]
	Add(ctx context.Context, parts KeyParts, value *T) error
	AddSlice(ctx context.Context, parts KeyParts, values []*T) error
	Scan(ctx context.Context, parts KeyParts, cursor uint64, count int64) ([]*T, uint64, error)
	DeleteMember(ctx context.Context, parts KeyParts, value *T) error
	ExistsMember(ctx context.Context, parts KeyParts, value *T) (bool, error)
	Count(ctx context.Context, parts KeyParts) (int64, error)
}

type SortedSetCache[T any] interface {
	Cache[T]
	Add(ctx context.Context, parts KeyParts, member *SortedMember[T]) error
	AddSlice(ctx context.Context, parts KeyParts, values []*SortedMember[T]) error
	Paginate(ctx context.Context, parts KeyParts, paginate *shared.Paginating) ([]*T, *shared.Paginated, error)
	// Tail 按分数排序 查看最后的数据
	Tail(ctx context.Context, parts KeyParts, size int64) ([]*T, error)
	DeleteMember(ctx context.Context, parts KeyParts, value *T) error
	ExistsMember(ctx context.Context, parts KeyParts, value *T) (bool, error)
	Count(ctx context.Context, parts KeyParts) (int64, error)
	// Clear 清除数据
	// reverse: false -> 删除分数低的
	// reverse: false -> 删除分数高的
	Clear(ctx context.Context, parts KeyParts, keep uint64, reverse bool) error
	// ClearByScore 按照分数删除
	// reverse: false -> 删除分数低的
	// reverse: false -> 删除分数高的
	ClearByScore(ctx context.Context, parts KeyParts, score int64, reverse bool) error
}

type SortedMember[T any] struct {
	Score int64
	Value *T
}

func Member[T any](score int64, value *T) *SortedMember[T] {
	return &SortedMember[T]{Score: score, Value: value}
}
