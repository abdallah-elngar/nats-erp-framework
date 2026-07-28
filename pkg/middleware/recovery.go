package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/nats-framework/nats/pkg/logger"
)

// Recovery يستعيد من الأخطاء
func Recovery(log *logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					// تسجيل الخطأ
					log.Error("Panic recovered",
						"error", err,
						"path", r.URL.Path,
						"method", r.Method,
						"stack", string(debug.Stack()),
					)

					// إرسال استجابة مناسبة
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)

					if log.Level() == slog.LevelDebug {
						fmt.Fprintf(w, `{"error": "%v", "stack": "%s"}`, err, string(debug.Stack()))
					} else {
						fmt.Fprintf(w, `{"error": "Internal Server Error"}`)
					}
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
