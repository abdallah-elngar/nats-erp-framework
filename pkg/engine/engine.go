package engine

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-framework/nats/apps/core/routes" // ✅ استيراد المسارات
	"github.com/nats-framework/nats/pkg/config"
	"github.com/nats-framework/nats/pkg/database"
	"github.com/nats-framework/nats/pkg/logger"
	"github.com/nats-framework/nats/pkg/router"
	"github.com/nats-framework/nats/pkg/template"
)

// Engine هو محرك التطبيق الرئيسي
type Engine struct {
	Config     *config.Config
	Logger     *logger.Logger
	DB         *database.Manager
	Server     *http.Server
	Router     *router.Router
	Template   *template.Engine
	Apps       map[string]interface{}
	Context    context.Context
	CancelFunc context.CancelFunc
}

// New ينشئ محرك تطبيق جديد
func New() *Engine {
	ctx, cancel := context.WithCancel(context.Background())

	return &Engine{
		Context:    ctx,
		CancelFunc: cancel,
		Apps:       make(map[string]interface{}),
	}
}

// LoadConfig يقوم بتحميل الإعدادات
func (e *Engine) LoadConfig() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	e.Config = cfg
	return nil
}

// InitDatabase يقوم بتهيئة قاعدة البيانات
func (e *Engine) InitDatabase() error {
	db, err := database.NewManager(&e.Config.Database)
	if err != nil {
		return fmt.Errorf("failed to init database: %w", err)
	}
	e.DB = db
	return nil
}

// InitTemplate يقوم بتهيئة محرك القوالب
func (e *Engine) InitTemplate() error {
	config := template.Config{
		Dir:          "apps/core/templates",
		Layout:       "layouts/admin",
		Debug:        e.Config.App.Debug,
		CacheEnabled: true,
		CacheTTL:     5 * time.Minute,
		AutoReload:   e.Config.App.Debug,
		AuthEnabled:  true,
		AuthConfig: template.AuthConfig{
			LoginURL:    "/login",
			LogoutURL:   "/logout",
			RegisterURL: "/register",
			UserKey:     "user",
		},
		StaticConfig: template.StaticConfig{
			Dir:          "static",
			Prefix:       "/static/",
			CacheEnabled: true,
			CacheTTL:     24 * time.Hour,
			Debug:        e.Config.App.Debug,
			Minify:       true,
		},
	}

	e.Template = template.New(config)
	return nil
}

// InitRouter يقوم بتهيئة الموجه
func (e *Engine) InitRouter() error {
	routerConfig := router.RouterConfig{
		Debug:        e.Config.App.Debug,
		StaticPrefix: "/static/",
		StaticDir:    "static",
	}

	e.Router = router.NewRouter(e.Template, routerConfig)

	// ✅ تسجيل المسارات - استدعاء RegisterRoutes مباشرة
	routes.RegisterRoutes(e.Router, e.Template)

	return nil
}

// LoadApps يقوم بتحميل التطبيقات
func (e *Engine) LoadApps() error {
	if err := e.InitTemplate(); err != nil {
		return err
	}

	if err := e.InitRouter(); err != nil {
		return err
	}

	return nil
}

// Run يقوم بتشغيل الخادم
func (e *Engine) Run() error {
	addr := fmt.Sprintf("%s:%d", e.Config.Server.Host, e.Config.Server.Port)

	e.Server = &http.Server{
		Addr:         addr,
		Handler:      e.Router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	fmt.Printf("🚀 Server started on %s\n", addr)
	return e.Server.ListenAndServe()
}

// WaitForShutdown ينتظر إشارة الإيقاف
func (e *Engine) WaitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down gracefully...")

	e.CancelFunc()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := e.Server.Shutdown(ctx); err != nil {
		log.Printf("❌ Server shutdown error: %v", err)
	}

	if e.DB != nil {
		if err := e.DB.Close(); err != nil {
			log.Printf("❌ Database close error: %v", err)
		}
	}

	log.Println("✅ Shutdown complete")
}

// GetRouter يعيد الموجه
func (e *Engine) GetRouter() *router.Router {
	return e.Router
}
