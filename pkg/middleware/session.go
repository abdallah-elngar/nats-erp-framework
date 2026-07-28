package middleware

import (
	"context"
	"net/http"

	"github.com/gorilla/sessions"
)

// SessionMiddleware يدير الجلسات
func SessionMiddleware(store sessions.Store, cookieName string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// الحصول على الجلسة
			session, err := store.Get(r, cookieName)
			if err != nil {
				http.Error(w, "Session error", http.StatusInternalServerError)
				return
			}

			// وضع الجلسة في السياق
			ctx := contextWithSession(r.Context(), session)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

// SessionKey مفتاح الجلسة في السياق
type contextKey string

const sessionKey contextKey = "session"

// contextWithSession يضع الجلسة في السياق
func contextWithSession(ctx context.Context, session *sessions.Session) context.Context {
	return context.WithValue(ctx, sessionKey, session)
}

// GetSession يحصل على الجلسة من السياق
func GetSession(ctx context.Context) *sessions.Session {
	if session, ok := ctx.Value(sessionKey).(*sessions.Session); ok {
		return session
	}
	return nil
}
