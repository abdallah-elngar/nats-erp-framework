package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/nats-framework/nats/pkg/config"
)

// Server يمثل خادم HTTP
type Server struct {
	config  *config.ServerConfig
	router  *chi.Mux
	http    *http.Server
	started bool
}

// New ينشئ خادم جديد
func New(cfg *config.ServerConfig) *Server {
	r := chi.NewRouter()

	// إضافة ميدلوير عامة
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.AllowContentType("application/json", "application/x-www-form-urlencoded", "multipart/form-data"))

	return &Server{
		config: cfg,
		router: r,
		http: &http.Server{
			Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			Handler:      r,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
	}
}

// Router يعيد الموجه
func (s *Server) Router() *chi.Mux {
	return s.router
}

// SetRouter يضع الموجه للخادم
func (s *Server) SetRouter(router *chi.Mux) {
	s.router = router
	s.http.Handler = router
}

// Start يقوم بتشغيل الخادم
func (s *Server) Start() error {
	if s.started {
		return fmt.Errorf("server already started")
	}

	s.started = true
	fmt.Printf("🚀 Server started on %s\n", s.http.Addr)
	return s.http.ListenAndServe()
}

// Shutdown يقوم بإيقاف الخادم
func (s *Server) Shutdown(ctx context.Context) error {
	if !s.started {
		return nil
	}

	s.started = false
	return s.http.Shutdown(ctx)
}

// GetAddr يعيد عنوان الخادم
func (s *Server) GetAddr() string {
	return s.http.Addr
}

// IsStarted يعيد حالة الخادم
func (s *Server) IsStarted() bool {
	return s.started
}

// GetActiveConnections يعيد عدد الاتصالات النشطة (محاكاة)
func (s *Server) GetActiveConnections() int {
	return 0
}

// ServeStatic يخدم الملفات الثابتة
func (s *Server) ServeStatic(prefix, dir string) {
	// ✅ إضافة خدمة الملفات الثابتة
	workDir, _ := os.Getwd()
	staticDir := filepath.Join(workDir, dir)

	// التحقق من وجود المجلد
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		fmt.Printf("⚠️ Static directory not found: %s\n", staticDir)
		return
	}

	// خدمة الملفات الثابتة
	s.router.Handle(prefix+"*", http.StripPrefix(prefix, http.FileServer(http.Dir(staticDir))))
	fmt.Printf("📁 Serving static files from: %s\n", staticDir)
}
