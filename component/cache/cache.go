package cache

import (
	"github.com/soukengo/gopkg/component/cache/core"
)

type Option = core.Option
type Category = core.Category
type Key = core.Key
type KeyParts = core.KeyParts

type ValueCache[T any] core.ValueCache[T]
type HashCache[T any] core.HashCache[T]
type SetCache[T any] core.SetCache[T]
type SortedSetCache[T any] core.SortedSetCache[T]

func Parts(parts ...string) KeyParts {
	return core.Parts(parts...)
}

func Member[T any](score int64, value *T) *core.SortedMember[T] {
	return core.Member[T](score, value)
}

func IsErrNoCache(err error) bool {
	return core.IsErrNoCache(err)
}
