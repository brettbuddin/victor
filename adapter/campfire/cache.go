package campfire

import (
	"github.com/brettbuddin/victor/adapter"
	"sync"
)

type Cache struct {
	mutex *sync.RWMutex
	cache map[string]adapter.Cacheable
}

func NewCache() *Cache {
	return &Cache{&sync.RWMutex{}, map[string]adapter.Cacheable{}}
}

func (c *Cache) Add(o adapter.Cacheable) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cache[o.CacheKey()] = o
}

func (c *Cache) Get(key string) adapter.Cacheable {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
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
