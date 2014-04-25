package store

import (
	"sync"
)

func init() {
	Register("memory", func() Adapter {
		return &MemoryStore{
			data: make(map[string]string),
		}
	})
}

type MemoryStore struct {
	sync.RWMutex
	data map[string]string
}

func (s *MemoryStore) Get(key string) (string, bool) {
	s.RLock()
	defer s.RUnlock()
	val, ok := s.data[key]
	return val, ok
}

func (s *MemoryStore) Set(key string, val string) {
	s.Lock()
	defer s.Unlock()
	s.data[key] = val
}

func (s *MemoryStore) Delete(key string) {
	s.Lock()
	defer s.Unlock()
	delete(s.data, key)
}

func (s *MemoryStore) All() map[string]string {
	s.RLock()
	defer s.RUnlock()
	return s.data
}
