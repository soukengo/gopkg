package lock

import "strings"

type Key string

const (
	lockKeySeparator = ":"
)

func (k Key) GenKey(parts ...string) string {
	return strings.Join(append([]string{string(k)}, parts...), lockKeySeparator)
}
