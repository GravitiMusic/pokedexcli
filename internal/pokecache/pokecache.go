package pokecache

import (
	"time"
	"sync"
)

type cacheEntry struct {
	createdAt time.Time
	val []byte
}

type Cache struct {
	store map[string]cacheEntry
	mu sync.RWMutex
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		store:    make(map[string]cacheEntry),
	}

	go cache.reapLoop(interval)

	return cache
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		for key, entry := range c.store {
			if time.Since(entry.createdAt) > interval {
				delete(c.store, key)
			}
		}
		c.mu.Unlock()
	}
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	c.store[key] = cacheEntry{
		createdAt: time.Now(),
		val: val,
	}
	c.mu.Unlock()
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	entry, exists := c.store[key]
	c.mu.RUnlock()

	if !exists {
		return nil, false
	} else {
		return entry.val, true
	}
}