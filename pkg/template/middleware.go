package template

import (
	"context"
	"net/http"
)

// Middleware وسيط للمصادقة
type Middleware struct {
	engine *Engine
}

// NewMiddleware ينشئ وسيطاً جديداً
func NewMiddleware(engine *Engine) *Middleware {
	return &Middleware{
		engine: engine,
	}
}

// AuthMiddleware وسيط المصادقة
func (m *Middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !m.engine.config.AuthEnabled {
			next.ServeHTTP(w, r)
			return
		}

		// التحقق من الجلسة
		user := m.getUserFromSession(r)
		if user != nil {
			// إضافة المستخدم إلى السياق
			ctx := r.Context()
			ctx = context.WithValue(ctx, "user", user)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

// RequireAuth يتطلب مصادقة
func (m *Middleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := m.getUserFromSession(r)
		if user == nil {
			http.Redirect(w, r, m.engine.config.AuthConfig.LoginURL, http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RequirePermission يتطلب صلاحية
func (m *Middleware) RequirePermission(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := m.getUserFromSession(r)
			if user == nil {
				http.Redirect(w, r, m.engine.config.AuthConfig.LoginURL, http.StatusFound)
				return
			}

			if !m.hasPermission(user, permission) {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// getUserFromSession يحصل على المستخدم من الجلسة
func (m *Middleware) getUserFromSession(r *http.Request) interface{} {
	// تنفيذ منطق استخراج المستخدم من الجلسة
	// هذا مثال بسيط
	session, err := r.Cookie("session")
	if err != nil {
		return nil
	}

	// التحقق من صحة الجلسة
	if session.Value == "" {
		return nil
	}

	// إرجاع مستخدم افتراضي للاختبار
	return map[string]interface{}{
		"id":          1,
		"username":    "admin",
		"email":       "admin@example.com",
		"permissions": []string{"*"},
	}
}

// hasPermission يتحقق من الصلاحية
func (m *Middleware) hasPermission(user interface{}, permission string) bool {
	if user == nil {
		return false
	}

	// إذا كان المستخدم لديه صلاحية *
	if u, ok := user.(map[string]interface{}); ok {
		if perms, ok := u["permissions"].([]string); ok {
			for _, p := range perms {
				if p == "*" || p == permission {
					return true
				}
			}
		}
	}

	return false
}
