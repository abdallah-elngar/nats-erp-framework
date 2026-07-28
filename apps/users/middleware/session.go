package middleware

import (
	"context"
	"net/http"

	"github.com/gorilla/sessions"

	"github.com/nats-framework/nats/pkg/response"
)

// SessionKey مفتاح الجلسة في السياق
type contextKey string

const sessionContextKey contextKey = "session"

// SessionMiddleware يدير الجلسات
func SessionMiddleware(store sessions.Store, cookieName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := store.Get(r, cookieName)
			if err != nil {
				response.InternalError(w, "Session error")
				return
			}

			// وضع الجلسة في السياق
			ctx := context.WithValue(r.Context(), sessionContextKey, session)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

// GetSession يحصل على الجلسة من السياق
func GetSession(r *http.Request) *sessions.Session {
	if session, ok := r.Context().Value(sessionContextKey).(*sessions.Session); ok {
		return session
	}
	return nil
}
