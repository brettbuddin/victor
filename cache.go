package victor

import (
	"github.com/brettbuddin/victor/adapter"
	"sync"
)

type Cache struct {
	*sync.RWMutex
	cache map[string]adapter.Cacheable
}

func NewCache() *Cache {
	return &Cache{&sync.RWMutex{}, map[string]adapter.Cacheable{}}
}

func (c *Cache) Add(o adapter.Cacheable) {
	c.Lock()
	defer c.Unlock()
	c.cache[o.CacheKey()] = o
}

func (c *Cache) Get(key string) adapter.Cacheable {
	c.RLock()
	defer c.RUnlock()
	if result, ok := c.cache[key]; ok {
		return result
	}
	return nil
}

func (c *Cache) Exists(key string) bool {
	if result := c.Get(key); result == nil {
		return false
	}

	return true
}

func (c *Cache) Delete(key string) {
	o := c.Get(key)

	if o == nil {
		return
	}

	c.Lock()
	defer c.Unlock()
	delete(c.cache, o.CacheKey())
}
