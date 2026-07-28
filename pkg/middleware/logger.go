package middleware

import (
	"log"
	"net/http"
	"time"
)

// ✅ إزالة تعريف responseWriter من هنا واستخدام التعريف من middleware.go

// LoggerMiddleware يسجل الطلبات
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// استخدام responseWriter من middleware.go
		wrapped := &responseWriter{ResponseWriter: w}

		next.ServeHTTP(wrapped, r)

		log.Printf("[%s] %s %s %d %v",
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			wrapped.status,
			time.Since(start),
		)
	})
}

// AdvancedLoggerMiddleware يسجل الطلبات بتفاصيل أكثر
func AdvancedLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &responseWriter{ResponseWriter: w}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)

		log.Printf("[%s] %s %s %d %v | User-Agent: %s | Referer: %s",
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			wrapped.status,
			duration,
			r.UserAgent(),
			r.Referer(),
		)
	})
}
