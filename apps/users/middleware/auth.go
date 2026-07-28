package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/nats-framework/nats/pkg/auth"
	"github.com/nats-framework/nats/pkg/response"
)

// AuthMiddleware يتحقق من المصادقة
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// الحصول على التوكن من الـ Header
		token := r.Header.Get("Authorization")
		if token == "" {
			// التحقق من وجود كوكي التوكن
			cookie, err := r.Cookie("token")
			if err == nil && cookie != nil {
				token = cookie.Value
			}
		}

		if token == "" {
			response.Unauthorized(w, "Authorization required")
			return
		}

		// إزالة "Bearer " من التوكن
		token = strings.TrimPrefix(token, "Bearer ")
		token = strings.TrimPrefix(token, "bearer ")

		if token == "" {
			response.Unauthorized(w, "Invalid token format")
			return
		}

		// ✅ إنشاء خدمة JWT والتحقق من التوكن
		jwtService := auth.NewJWTService(auth.JWTConfig{
			Secret:     "your-secret-key", // يجب أن يكون من الإعدادات
			Expiration: 24 * 3600,         // 24 ساعة
			Issuer:     "nats",
		})

		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			response.Unauthorized(w, "Invalid or expired token")
			return
		}

		// وضع المستخدم في السياق
		ctx := context.WithValue(r.Context(), "user", claims)
		ctx = context.WithValue(ctx, "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "token", token)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserFromContext يحصل على المستخدم من السياق
func GetUserFromContext(ctx context.Context) (*auth.Claims, bool) {
	user, ok := ctx.Value("user").(*auth.Claims)
	return user, ok
}

// GetUserIDFromContext يحصل على معرف المستخدم من السياق
func GetUserIDFromContext(ctx context.Context) (uint, bool) {
	userID, ok := ctx.Value("user_id").(uint)
	return userID, ok
}

// GetTokenFromContext يحصل على التوكن من السياق
func GetTokenFromContext(ctx context.Context) (string, bool) {
	token, ok := ctx.Value("token").(string)
	return token, ok
}
