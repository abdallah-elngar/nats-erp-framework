package main

import (
	"fmt"
	"log"

	"github.com/nats-framework/nats/pkg/engine"
)

func main() {
	// إنشاء محرك التطبيق
	app := engine.New()

	// تحميل الإعدادات
	if err := app.LoadConfig(); err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}

	// تهيئة قاعدة البيانات
	if err := app.InitDatabase(); err != nil {
		log.Fatalf("❌ Failed to init database: %v", err)
	}

	// تحميل التطبيقات
	if err := app.LoadApps(); err != nil {
		log.Fatalf("❌ Failed to load apps: %v", err)
	}

	// تشغيل الخادم
	if err := app.Run(); err != nil {
		log.Fatalf("❌ Failed to run server: %v", err)
	}

	fmt.Println("🚀 NATS Framework is running!")
	fmt.Printf("🌐 Server: http://%s:%d\n", app.Config.Server.Host, app.Config.Server.Port)
	fmt.Printf("📊 Admin: http://%s:%d/admin\n", app.Config.Server.Host, app.Config.Server.Port)

	// انتظار إشارة الإيقاف
	app.WaitForShutdown()
}
