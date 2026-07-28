package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"sync"
)

// CSRFConfig يمثل إعدادات CSRF
type CSRFConfig struct {
	TokenLength int
	CookieName  string
	HeaderName  string
	Secure      bool
	HttpOnly    bool
}

// DefaultCSRFConfig يعيد إعدادات CSRF الافتراضية
func DefaultCSRFConfig() *CSRFConfig {
	return &CSRFConfig{
		TokenLength: 32,
		CookieName:  "csrf_token",
		HeaderName:  "X-CSRF-Token",
		Secure:      false,
		HttpOnly:    true,
	}
}

// CSRF ينشئ وسيط CSRF
func CSRF(config *CSRFConfig) Middleware {
	if config == nil {
		config = DefaultCSRFConfig()
	}

	var mu sync.Mutex
	tokens := make(map[string]bool)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// توليد توكن CSRF جديد
			token := generateCSRFToken(config.TokenLength)

			// حفظ التوكن
			mu.Lock()
			tokens[token] = true
			mu.Unlock()

			// تعيين الكوكي
			http.SetCookie(w, &http.Cookie{
				Name:     config.CookieName,
				Value:    token,
				Path:     "/",
				Secure:   config.Secure,
				HttpOnly: config.HttpOnly,
			})

			// التحقق من التوكن للطلبات غير الآمنة
			if r.Method != "GET" && r.Method != "HEAD" && r.Method != "OPTIONS" {
				// الحصول على التوكن من الرأس أو من النموذج
				tokenFromHeader := r.Header.Get(config.HeaderName)
				tokenFromCookie, _ := r.Cookie(config.CookieName)

				var tokenToVerify string
				if tokenFromHeader != "" {
					tokenToVerify = tokenFromHeader
				} else if tokenFromCookie != nil {
					tokenToVerify = tokenFromCookie.Value
				}

				// التحقق من التوكن
				mu.Lock()
				valid := tokens[tokenToVerify]
				delete(tokens, tokenToVerify)
				mu.Unlock()

				if !valid {
					http.Error(w, "CSRF token validation failed", http.StatusForbidden)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

// generateCSRFToken يولد توكن CSRF
func generateCSRFToken(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)
}
