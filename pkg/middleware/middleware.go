package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/nats-framework/nats/pkg/auth"
	"github.com/nats-framework/nats/pkg/logger"
)

// Middleware هو دالة ميدلوير
type Middleware func(http.Handler) http.Handler

// ✅ تعريف responseWriter هنا (مرة واحدة فقط)
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

// Chain يربط عدة ميدلوير
func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

// AuthMiddleware يتحقق من المصادقة
func AuthMiddleware(authService *auth.AuthService) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			token = strings.TrimPrefix(token, "Bearer ")
			if token == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// التحقق من التوكن - استخدام الدالة العامة
			claims, err := auth.ValidateToken(token)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			user := &auth.User{
				ID:       claims.UserID,
				Username: claims.Username,
				Email:    claims.Email,
				Roles:    claims.Roles,
			}

			ctx := context.WithValue(r.Context(), "user", user)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

// GuestMiddleware يتحقق من أن المستخدم ضيف (غير مسجل)
func GuestMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, ok := auth.GetUserFromContext(r.Context()); ok {
				http.Redirect(w, r, "/dashboard", http.StatusFound)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// PermissionMiddleware يتحقق من الصلاحيات
func PermissionMiddleware(requiredPermission string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := auth.GetUserFromContext(r.Context())
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			hasPermission := false
			for _, role := range user.Roles {
				if role == "admin" {
					hasPermission = true
					break
				}
			}

			if !hasPermission {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// LoggerMiddleware يسجل الطلبات - ✅ إزالة هذا التعريف (موجود في logger.go)
// يتم استخدام التعريف من logger.go

// CORSMiddleware يضيف رؤوس CORS
func CORSMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RateLimitMiddleware يحد من معدل الطلبات
func RateLimitMiddleware(limit int) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// سيتم تنفيذ تحديد المعدل باستخدام Redis
			next.ServeHTTP(w, r)
		})
	}
}

// RecoveryMiddleware يستعيد من الأخطاء
func RecoveryMiddleware(log *logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Error("Panic recovered",
						"error", err,
						"path", r.URL.Path,
						"method", r.Method,
					)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// StaticMiddleware يخدم الملفات الثابتة
func StaticMiddleware(prefix string, dir string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, prefix) {
				http.StripPrefix(prefix, http.FileServer(http.Dir(dir))).ServeHTTP(w, r)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
