package middleware

import (
    "net/http"
    "strings"

    "github.com/nats-framework/nats/pkg/response"
)

// AuthMiddleware يتحقق من صلاحيات المستخدم لتطبيق sale
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            response.Unauthorized(w, "Authorization required")
            return
        }

        token = strings.TrimPrefix(token, "Bearer ")
        if token == "" {
            response.Unauthorized(w, "Invalid token")
            return
        }

        // TODO: التحقق من صلاحية المستخدم لتطبيق sale
        next.ServeHTTP(w, r)
    })
}

// RateLimitMiddleware يحد من معدل الطلبات
func RateLimitMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // TODO: تطبيق تحديد معدل الطلبات
        next.ServeHTTP(w, r)
    })
}
