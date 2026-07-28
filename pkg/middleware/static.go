package middleware

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Static يخدم الملفات الثابتة
func Static(prefix, dir string) Middleware {
	// تحويل المسار إلى مسار مطلق
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// التحقق من أن المسار يبدأ بالبادئة
			if !strings.HasPrefix(r.URL.Path, prefix) {
				next.ServeHTTP(w, r)
				return
			}

			// إزالة البادئة للحصول على المسار النسبي
			path := strings.TrimPrefix(r.URL.Path, prefix)
			filePath := filepath.Join(absDir, path)

			// التحقق من وجود الملف
			info, err := os.Stat(filePath)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			// إذا كان مجلداً، عرض index.html
			if info.IsDir() {
				indexPath := filepath.Join(filePath, "index.html")
				if _, err := os.Stat(indexPath); err == nil {
					http.ServeFile(w, r, indexPath)
					return
				}
				next.ServeHTTP(w, r)
				return
			}

			// خدمة الملف
			http.ServeFile(w, r, filePath)
		})
	}
}
