package cache

import (
	"context"
	"sync"
	"time"
)

// MemoryCache تخزين مؤقت في الذاكرة
type MemoryCache struct {
	items  map[string]Item
	mu     sync.RWMutex
	prefix string
	stop   chan struct{}
}

// NewMemoryCache ينشئ تخزيناً مؤقتاً في الذاكرة
func NewMemoryCache(prefix string) *MemoryCache {
	c := &MemoryCache{
		items:  make(map[string]Item),
		prefix: prefix,
		stop:   make(chan struct{}),
	}

	// بدء تنظيف دوري
	go c.cleanup()

	return c
}

// Get يحصل على قيمة من التخزين المؤقت
func (c *MemoryCache) Get(ctx context.Context, key string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	fullKey := c.prefix + key
	item, ok := c.items[fullKey]
	if !ok {
		return nil, nil
	}

	// التحقق من انتهاء الصلاحية
	if !item.Expiration.IsZero() && time.Now().After(item.Expiration) {
		return nil, nil
	}

	return item.Value, nil
}

// Set يضع قيمة في التخزين المؤقت
func (c *MemoryCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	fullKey := c.prefix + key
	expiration := time.Time{}
	if ttl > 0 {
		expiration = time.Now().Add(ttl)
	}

	c.items[fullKey] = Item{
		Key:        fullKey,
		Value:      value,
		Expiration: expiration,
	}

	return nil
}

// Delete يحذف قيمة من التخزين المؤقت
func (c *MemoryCache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	fullKey := c.prefix + key
	delete(c.items, fullKey)

	return nil
}

// Clear يمسح التخزين المؤقت
func (c *MemoryCache) Clear(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]Item)
	return nil
}

// Has يتحقق من وجود مفتاح
func (c *MemoryCache) Has(ctx context.Context, key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	fullKey := c.prefix + key
	item, ok := c.items[fullKey]
	if !ok {
		return false
	}

	// التحقق من انتهاء الصلاحية
	if !item.Expiration.IsZero() && time.Now().After(item.Expiration) {
		return false
	}

	return true
}

// GetMultiple يحصل على قيم متعددة
func (c *MemoryCache) GetMultiple(ctx context.Context, keys []string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	for _, key := range keys {
		value, err := c.Get(ctx, key)
		if err != nil {
			return nil, err
		}
		if value != nil {
			result[key] = value
		}
	}

	return result, nil
}

// SetMultiple يضع قيماً متعددة
func (c *MemoryCache) SetMultiple(ctx context.Context, items map[string]interface{}, ttl time.Duration) error {
	for key, value := range items {
		if err := c.Set(ctx, key, value, ttl); err != nil {
			return err
		}
	}
	return nil
}

// DeleteMultiple يحذف قيماً متعددة
func (c *MemoryCache) DeleteMultiple(ctx context.Context, keys []string) error {
	for _, key := range keys {
		if err := c.Delete(ctx, key); err != nil {
			return err
		}
	}
	return nil
}

// Increment يزيد قيمة
func (c *MemoryCache) Increment(ctx context.Context, key string, value int64) (int64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	fullKey := c.prefix + key
	item, ok := c.items[fullKey]
	if !ok {
		c.items[fullKey] = Item{
			Key:   fullKey,
			Value: value,
		}
		return value, nil
	}

	if val, ok := item.Value.(int64); ok {
		newVal := val + value
		item.Value = newVal
		c.items[fullKey] = item
		return newVal, nil
	}

	return 0, nil
}

// Decrement ينقص قيمة
func (c *MemoryCache) Decrement(ctx context.Context, key string, value int64) (int64, error) {
	return c.Increment(ctx, key, -value)
}

// cleanup ينظف العناصر المنتهية
func (c *MemoryCache) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			now := time.Now()
			for key, item := range c.items {
				if !item.Expiration.IsZero() && now.After(item.Expiration) {
					delete(c.items, key)
				}
			}
			c.mu.Unlock()
		case <-c.stop:
			return
		}
	}
}

// Close يغلق التخزين المؤقت
func (c *MemoryCache) Close() {
	close(c.stop)
}
