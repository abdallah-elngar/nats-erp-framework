package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache تخزين مؤقت باستخدام Redis
type RedisCache struct {
	client *redis.Client
	prefix string
	ctx    context.Context
}

// NewRedisCache ينشئ تخزيناً مؤقتاً باستخدام Redis
func NewRedisCache(config *Config) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + string(rune(config.Port)),
		Password: config.Password,
		DB:       config.Database,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisCache{
		client: client,
		prefix: config.Prefix,
		ctx:    ctx,
	}, nil
}

// Get يحصل على قيمة من التخزين المؤقت
func (r *RedisCache) Get(ctx context.Context, key string) (interface{}, error) {
	fullKey := r.prefix + key
	val, err := r.client.Get(ctx, fullKey).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		return val, nil
	}

	return result, nil
}

// Set يضع قيمة في التخزين المؤقت
func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	fullKey := r.prefix + key

	var val string
	switch v := value.(type) {
	case string:
		val = v
	default:
		data, err := json.Marshal(value)
		if err != nil {
			return err
		}
		val = string(data)
	}

	return r.client.Set(ctx, fullKey, val, ttl).Err()
}

// Delete يحذف قيمة من التخزين المؤقت
func (r *RedisCache) Delete(ctx context.Context, key string) error {
	fullKey := r.prefix + key
	return r.client.Del(ctx, fullKey).Err()
}

// Clear يمسح التخزين المؤقت
func (r *RedisCache) Clear(ctx context.Context) error {
	pattern := r.prefix + "*"
	iter := r.client.Scan(ctx, 0, pattern, 0).Iterator()

	for iter.Next(ctx) {
		if err := r.client.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}

	return iter.Err()
}

// Has يتحقق من وجود مفتاح
func (r *RedisCache) Has(ctx context.Context, key string) bool {
	fullKey := r.prefix + key
	exists, err := r.client.Exists(ctx, fullKey).Result()
	if err != nil {
		return false
	}
	return exists > 0
}

// GetMultiple يحصل على قيم متعددة
func (r *RedisCache) GetMultiple(ctx context.Context, keys []string) (map[string]interface{}, error) {
	fullKeys := make([]string, len(keys))
	for i, key := range keys {
		fullKeys[i] = r.prefix + key
	}

	vals, err := r.client.MGet(ctx, fullKeys...).Result()
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	for i, val := range vals {
		if val != nil {
			result[keys[i]] = val
		}
	}

	return result, nil
}

// SetMultiple يضع قيماً متعددة
func (r *RedisCache) SetMultiple(ctx context.Context, items map[string]interface{}, ttl time.Duration) error {
	pipe := r.client.Pipeline()

	for key, value := range items {
		fullKey := r.prefix + key
		var val string
		switch v := value.(type) {
		case string:
			val = v
		default:
			data, err := json.Marshal(value)
			if err != nil {
				return err
			}
			val = string(data)
		}
		pipe.Set(ctx, fullKey, val, ttl)
	}

	_, err := pipe.Exec(ctx)
	return err
}

// DeleteMultiple يحذف قيماً متعددة
func (r *RedisCache) DeleteMultiple(ctx context.Context, keys []string) error {
	fullKeys := make([]string, len(keys))
	for i, key := range keys {
		fullKeys[i] = r.prefix + key
	}

	return r.client.Del(ctx, fullKeys...).Err()
}

// Increment يزيد قيمة
func (r *RedisCache) Increment(ctx context.Context, key string, value int64) (int64, error) {
	fullKey := r.prefix + key
	return r.client.IncrBy(ctx, fullKey, value).Result()
}

// Decrement ينقص قيمة
func (r *RedisCache) Decrement(ctx context.Context, key string, value int64) (int64, error) {
	fullKey := r.prefix + key
	return r.client.DecrBy(ctx, fullKey, value).Result()
}

// Close يغلق التخزين المؤقت
func (r *RedisCache) Close() error {
	return r.client.Close()
}
