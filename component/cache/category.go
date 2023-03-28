package cache

import "time"

// Category cache category
type Category struct {
	Key         Key
	Expire      time.Duration
	EmptyExpire time.Duration
}

func (c *Category) Options() (opts []Option) {
	if c.Expire > 0 {
		opts = append(opts, WithExpire(c.Expire))
	}
	if c.EmptyExpire > 0 {
		opts = append(opts, WithEmptyExpire(c.EmptyExpire))
	}
	return
}

func (c *Category) GenKey(parts KeyParts) string {
	return c.Key.GenKey(parts)
}

func (c *Category) CopyFrom(source *Category) (ret *Category) {
	ret = &Category{Key: source.Key, Expire: c.Expire, EmptyExpire: c.EmptyExpire}
	if source.Expire > 0 {
		ret.Expire = source.Expire
	}
	if source.EmptyExpire > 0 {
		ret.EmptyExpire = source.EmptyExpire
	}
	return
}
