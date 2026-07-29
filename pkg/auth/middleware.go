package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/nats-framework/nats/pkg/response"
)

// MiddlewareAuth يتحقق من المصادقة
func MiddlewareAuth(jwtService *JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
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

			// التحقق من التوكن
			claims, err := jwtService.ValidateToken(token)
			if err != nil {
				response.Unauthorized(w, "Invalid or expired token")
				return
			}

			// إنشاء مستخدم من المطالبات
			user := &User{
				ID:       claims.UserID,
				Username: claims.Username,
				Email:    claims.Email,
				Roles:    claims.Roles,
			}

			// وضع المستخدم في السياق (استخدام دوال auth)
			ctx := r.Context()
			ctx = SetUserInContext(ctx, user)
			ctx = SetUserIDInContext(ctx, user.ID)
			ctx = SetTokenInContext(ctx, token)
			ctx = SetPermissionsInContext(ctx, user.Roles)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// GetUserFromRequestContext يحصل على المستخدم من سياق الطلب (دالة مساعدة للـ Controllers)
func GetUserFromRequest(r *http.Request) interface{} {
	if r == nil {
		return nil
	}
	user, ok := GetUserFromContext(r.Context())
	if !ok {
		return nil
	}
	return user
}

// GetUserIDFromRequestContext يحصل على معرف المستخدم من سياق الطلب
func GetUserIDFromRequestContext(ctx context.Context) (uint, bool) {
	return GetUserIDFromContext(ctx)
}

// GetUsernameFromRequestContext يحصل على اسم المستخدم من سياق الطلب
func GetUsernameFromRequestContext(ctx context.Context) string {
	return GetUsernameFromContext(ctx)
}

// IsAuthenticatedRequest يتحقق من وجود مستخدم في سياق الطلب
func IsAuthenticatedRequest(ctx context.Context) bool {
	return IsAuthenticated(ctx)
}
