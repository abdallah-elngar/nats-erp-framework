package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/nats-framework/nats/pkg/response"
)

// RateLimiter يحد من معدل الطلبات
type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
	limit    int
	window   time.Duration
}

// NewRateLimiter ينشئ محدد معدل جديد
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

// Allow يتحقق من السماح بالطلب
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	// تنظيف الطلبات القديمة
	requests := rl.requests[key]
	var valid []time.Time
	for _, t := range requests {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}

	// التحقق من الحد
	if len(valid) >= rl.limit {
		rl.requests[key] = valid
		return false
	}

	// إضافة الطلب الجديد
	valid = append(valid, now)
	rl.requests[key] = valid

	return true
}

// RateLimitMiddleware وسيط تحديد معدل الطلبات
func RateLimitMiddleware(limit int, window time.Duration) func(http.Handler) http.Handler {
	limiter := NewRateLimiter(limit, window)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// استخدام IP كمفتاح
			key := r.RemoteAddr
			if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
				key = forwarded
			}

			if !limiter.Allow(key) {
				response.Error(w, http.StatusTooManyRequests, "Rate limit exceeded. Please try again later.")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
