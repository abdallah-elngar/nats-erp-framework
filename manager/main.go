package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-framework/nats/apps/core/routes"
	"github.com/nats-framework/nats/pkg/router"
	"github.com/nats-framework/nats/pkg/template"
)

func main() {
	// ✅ 1. إنشاء محرك القوالب
	templateEngine := template.New(template.Config{
		Dir:          "apps/core/templates",
		Layout:       "admin",
		Debug:        true,
		CacheEnabled: false,
		AutoReload:   true,
		AuthEnabled:  true,
		AuthConfig: template.AuthConfig{
			LoginURL:  "/login",
			LogoutURL: "/logout",
			UserKey:   "user",
		},
		StaticConfig: template.StaticConfig{
			Dir:          "static",
			Prefix:       "/static/",
			CacheEnabled: false,
			Debug:        true,
		},
	})

	// ✅ 2. إنشاء الراوتر مع المعاملات الصحيحة
	routerConfig := router.RouterConfig{
		Debug:        true,
		StaticPrefix: "/static/",
		StaticDir:    "static",
	}

	r := router.NewRouter(templateEngine, routerConfig)

	// ✅ 3. تسجيل المسارات
	routes.RegisterRoutes(r, templateEngine)

	// ✅ 4. إنشاء الخادم
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	fmt.Println("🚀 NATS Framework Server")
	fmt.Println("========================================")
	fmt.Println("📊 Developer Dashboard: http://localhost:8080/admin/developer")
	fmt.Println("🌐 Production Dashboard: http://localhost:8080/admin/production")
	fmt.Println("📚 API: http://localhost:8080/api/admin")
	fmt.Println("❤️  Health: http://localhost:8080/health")
	fmt.Println("========================================")
	fmt.Println("Press Ctrl+C to stop")

	// ✅ 5. انتظار إشارة الإيقاف
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		fmt.Println("\n🛑 Shutting down gracefully...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("❌ Server shutdown error: %v", err)
		}
	}()

	// ✅ 6. تشغيل الخادم
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("❌ Server error: %v", err)
	}

	fmt.Println("✅ Server stopped")
}
