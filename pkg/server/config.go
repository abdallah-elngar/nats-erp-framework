package server

import (
	"time"
)

// Config يمثل إعدادات الخادم
type Config struct {
	Host            string
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	MaxHeaderBytes  int
	EnableCORS      bool
	EnableCSRF      bool
	EnableRateLimit bool
	RateLimit       int
}

// DefaultConfig يعيد الإعدادات الافتراضية
func DefaultConfig() *Config {
	return &Config{
		Host:            "0.0.0.0",
		Port:            8080,
		ReadTimeout:     30 * time.Second,
		WriteTimeout:    30 * time.Second,
		IdleTimeout:     120 * time.Second,
		MaxHeaderBytes:  1 << 20, // 1MB
		EnableCORS:      true,
		EnableCSRF:      true,
		EnableRateLimit: false,
		RateLimit:       100,
	}
}
