package core

import "strings"

type Topic string

const (
	keySeparator = ":"
)

func (t Topic) GenKey(parts ...string) string {
	return strings.Join(parts, keySeparator)
}
