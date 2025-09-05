package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt    time.Time
	val        []byte
}

type Cache struct {
	data 	map[string]cacheEntry
	ttl  	time.Duration
	mutex 	sync.RWMutex
}

func NewCache(ttl time.Duration) *Cache {
	cache := &Cache{
		data: make(map[string]cacheEntry),
		ttl:  ttl,
	}
	
	// Start the reap loop in a goroutine
	go cache.reapLoop(time.Second * 5)
	
	return cache
}

func (c *Cache) Add(key string, value []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data[key] = cacheEntry{
		createdAt: time.Now(),
		val:       value,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	entry, exists := c.data[key]
	if !exists {
		return nil, false
	}
	if time.Since(entry.createdAt) > c.ttl {
		delete(c.data, key)
		return nil, false
	}
	return entry.val, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		c.mutex.Lock()
		now := time.Now()
		for key, entry := range c.data {
			if now.Sub(entry.createdAt) > c.ttl {
				delete(c.data, key)
			}
		}
		c.mutex.Unlock()
	}
}