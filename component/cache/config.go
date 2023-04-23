package cache

import (
	"github.com/soukengo/gopkg/infra/storage"
	"github.com/soukengo/gopkg/infra/storage/redis"
	"sync"
)

type Config struct {
	Default    Category
	Categories []*Category
	Redis      *redis.Reference

	cateMap map[Key]*Category
	lock    sync.RWMutex
}

func (c *Config) Parse(configs *storage.Config) {
	c.cateMap = make(map[Key]*Category)
	if c.Categories == nil {
		c.Categories = make([]*Category, 0)
	}
	c.lock.Lock()
	for _, category := range c.Categories {
		c.cateMap[category.Key] = c.Default.CopyFrom(category)
	}
	c.lock.Unlock()
	// 设置redis参数
	c.Redis.Parse(configs.Redis)
	return
}

func (c *Config) Category(key Key) *Category {
	c.lock.RLock()
	cate, ok := c.cateMap[key]
	c.lock.RUnlock()
	if ok {
		return cate
	}
	return c.new(key)
}

func (c *Config) new(key Key) *Category {
	c.lock.Lock()
	cate := c.Default.CopyFrom(&Category{Key: key})
	if c.cateMap == nil {
		c.cateMap = make(map[Key]*Category)
	}
	c.cateMap[key] = cate
	c.lock.Unlock()
	return cate
}
