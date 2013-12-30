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

// Expired reports whether a key in the cache is expired or not
func (c *Cache) Expired(key string) bool {
	c.RLock()
	defer c.RUnlock()
	if record, ok := c.cache[key]; ok {
		return record.CachedAt.Add(c.cacheFor).After(time.Now())
	}

	return false
}

// Store sets a key/value pair in the cache
func (c *Cache) Store(o adapter.CacheKeyer) {
	c.Lock()
	defer c.Unlock()
	record := CacheRecord{o, time.Now()}
	c.cache[o.CacheKey()] = record
}

// Get returns a particular key's value from the cache
func (c *Cache) Get(key string) adapter.CacheKeyer {
	c.RLock()
	defer c.RUnlock()
	if record, ok := c.cache[key]; ok {
		return record.Object
	}
	return nil
}

// Exists checks to see if a key is set in the cache
func (c *Cache) Exists(key string) bool {
	if record := c.Get(key); record == nil {
		return false
	}
	return true
}

// Delete removes a key/value pair from the cache
func (c *Cache) Delete(key string) {
	o := c.Get(key)

	if o == nil {
		return
	}

	c.Lock()
	defer c.Unlock()
	delete(c.cache, o.CacheKey())
}
