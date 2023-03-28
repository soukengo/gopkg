package lock

type Builder interface {
	Build(key Key, parts ...string) DistributedLock
}
