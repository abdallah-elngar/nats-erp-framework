package template

import (
	"sync"
	"time"
)

type Cache struct {
	items map[string]*CacheItem
	mu    sync.RWMutex
	ttl   time.Duration
}

type CacheItem struct {
	Value     interface{}
	CreatedAt time.Time
}

func NewCache(ttl time.Duration) *Cache {
	if ttl == 0 {
		ttl = 5 * time.Minute
	}
	return &Cache{
		items: make(map[string]*CacheItem),
		ttl:   ttl,
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	item, exists := c.items[key]
	if !exists {
		return nil, false
	}
	
	if time.Since(item.CreatedAt) > c.ttl {
		return nil, false
	}
	
	return item.Value, true
}

func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.items[key] = &CacheItem{
		Value:     value,
		CreatedAt: time.Now(),
	}
}

func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]*CacheItem)
}