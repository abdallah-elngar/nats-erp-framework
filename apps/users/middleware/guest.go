package middleware

import (
	"net/http"

	"github.com/nats-framework/nats/pkg/response"
)

// GuestMiddleware يتحقق من أن المستخدم ضيف (غير مسجل)
func GuestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// التحقق من وجود توكن في الـ Header
		token := r.Header.Get("Authorization")
		if token != "" {
			// المستخدم مسجل بالفعل
			response.BadRequest(w, "Already authenticated")
			return
		}

		// التحقق من وجود كوكي جلسة
		cookie, err := r.Cookie("session_token")
		if err == nil && cookie != nil && cookie.Value != "" {
			response.BadRequest(w, "Already authenticated")
			return
		}

		next.ServeHTTP(w, r)
	})
}
