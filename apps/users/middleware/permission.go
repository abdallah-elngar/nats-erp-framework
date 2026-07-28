package middleware

import (
	"net/http"

	"github.com/nats-framework/nats/apps/users/services"
	"github.com/nats-framework/nats/pkg/auth"
	"github.com/nats-framework/nats/pkg/response"
	"gorm.io/gorm"
)

// PermissionMiddleware يتحقق من الصلاحيات
func PermissionMiddleware(db *gorm.DB, permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// الحصول على المستخدم من السياق
			claims, ok := r.Context().Value("user").(*auth.Claims)
			if !ok {
				response.Unauthorized(w, "User not authenticated")
				return
			}

			// ✅ إنشاء خدمة المستخدم مع تمرير db
			userService := services.NewUserService(db)
			hasPermission, err := userService.HasPermission(claims.UserID, permission)
			if err != nil {
				response.InternalError(w, "Failed to check permission")
				return
			}

			if !hasPermission {
				response.Forbidden(w, "Insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// PermissionMiddlewareWithService يتحقق من الصلاحيات باستخدام خدمة موجودة
func PermissionMiddlewareWithService(userService *services.UserService, permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value("user").(*auth.Claims)
			if !ok {
				response.Unauthorized(w, "User not authenticated")
				return
			}

			hasPermission, err := userService.HasPermission(claims.UserID, permission)
			if err != nil {
				response.InternalError(w, "Failed to check permission")
				return
			}

			if !hasPermission {
				response.Forbidden(w, "Insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// HasPermissionMiddleware يتحقق من وجود صلاحية محددة
func HasPermissionMiddleware(db *gorm.DB, requiredPermission string) func(http.Handler) http.Handler {
	return PermissionMiddleware(db, requiredPermission)
}
