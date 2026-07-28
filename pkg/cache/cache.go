package cache

import (
	"context"
	"time"
)

// Cache واجهة نظام التخزين المؤقت
type Cache interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Clear(ctx context.Context) error
	Has(ctx context.Context, key string) bool
	GetMultiple(ctx context.Context, keys []string) (map[string]interface{}, error)
	SetMultiple(ctx context.Context, items map[string]interface{}, ttl time.Duration) error
	DeleteMultiple(ctx context.Context, keys []string) error
	Increment(ctx context.Context, key string, value int64) (int64, error)
	Decrement(ctx context.Context, key string, value int64) (int64, error)
}

// Item يمثل عنصراً في التخزين المؤقت
type Item struct {
	Key        string
	Value      interface{}
	Expiration time.Time
}

// Config يمثل إعدادات التخزين المؤقت
type Config struct {
	Driver   string
	Host     string
	Port     int
	Password string
	Database int
	Prefix   string
}

// DefaultConfig يعيد الإعدادات الافتراضية
func DefaultConfig() *Config {
	return &Config{
		Driver: "memory",
		Prefix: "nats:",
	}
}
