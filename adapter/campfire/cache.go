package campfire

import "sync"

type Cacheable interface {
	CacheKey() string
}

type Cacher interface {
	Add(Cacheable)
	Get(string) Cacheable
	Exists(string) bool
}

type Cache struct {
	mutex *sync.RWMutex
	cache map[string]Cacheable
}

func NewCache() *Cache {
	return &Cache{&sync.RWMutex{}, map[string]Cacheable{}}
}

func (c *Cache) Add(o Cacheable) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cache[o.CacheKey()] = o
}

func (c *Cache) Get(key string) Cacheable {
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
