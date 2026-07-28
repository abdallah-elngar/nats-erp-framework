// pkg/template/static_middleware.go
package template

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// StaticMiddleware وسيط للملفات الثابتة
type StaticMiddleware struct {
	manager *StaticManager
	config  StaticConfig
}

// NewStaticMiddleware ينشئ وسيطاً جديداً
func NewStaticMiddleware(config StaticConfig) *StaticMiddleware {
	return &StaticMiddleware{
		manager: NewStaticManager(config),
		config:  config,
	}
}

// Handler يعيد معالج HTTP
func (m *StaticMiddleware) Handler() http.Handler {
	return m.manager.Handler()
}

// Middleware يعيد وسيطاً
func (m *StaticMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// التحقق مما إذا كان الطلب لملف ثابت
		if strings.HasPrefix(r.URL.Path, m.config.Prefix) {
			m.manager.Handler().ServeHTTP(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ServeFile يخدم ملفاً محدداً
func (m *StaticMiddleware) ServeFile(w http.ResponseWriter, r *http.Request, path string) {
	m.manager.ServeFile(w, r, path)
}

// Exists يتحقق من وجود ملف
func (m *StaticMiddleware) Exists(path string) bool {
	fullPath := filepath.Join(m.config.Dir, path)
	if m.config.FS != nil {
		_, err := m.config.FS.ReadFile(path)
		return err == nil
	}
	_, err := os.Stat(fullPath)
	return err == nil
}
