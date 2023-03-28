package queue

import "strings"

type Topic string

const (
	lockPrefix   = "delay.lock."
	waitPrefix   = "delay.wait."
	readyPrefix  = "delay.ready."
	keySeparator = ":"
)

func (t Topic) genKey(parts ...string) string {
	return strings.Join(parts, keySeparator)
}

func (t Topic) waitKey() string {
	parts := []string{waitPrefix, string(t)}
	return t.genKey(parts...)
}
func (t Topic) readyKey() string {
	parts := []string{readyPrefix, string(t)}
	return t.genKey(parts...)
}
func (t Topic) lockKey() string {
	parts := []string{lockPrefix, string(t)}
	return t.genKey(parts...)
}
