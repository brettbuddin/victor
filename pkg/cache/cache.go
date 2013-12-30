package cache

import (
	"github.com/brettbuddin/victor/pkg/adapter"
	"sync"
	"time"
)

type CacheRecord struct {
	Object   adapter.CacheKeyer
	CachedAt time.Time
}

type Cache struct {
	*sync.RWMutex
	cache    map[string]CacheRecord
	cacheFor time.Duration
}

func New(dur time.Duration) *Cache {
	return &Cache{&sync.RWMutex{}, map[string]CacheRecord{}, dur}
}

func (c *Cache) Expired(key string) bool {
	c.RLock()
	defer c.RUnlock()
	if record, ok := c.cache[key]; ok {
		return record.CachedAt.Add(c.cacheFor).After(time.Now())
	}

	return false
}

func (c *Cache) Store(o adapter.CacheKeyer) {
	c.Lock()
	defer c.Unlock()
	record := CacheRecord{o, time.Now()}
	c.cache[o.CacheKey()] = record
}

func (c *Cache) Get(key string) adapter.CacheKeyer {
	c.RLock()
	defer c.RUnlock()
	if record, ok := c.cache[key]; ok {
		return record.Object
	}
	return nil
}

func (c *Cache) Exists(key string) bool {
	if record := c.Get(key); record == nil {
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
