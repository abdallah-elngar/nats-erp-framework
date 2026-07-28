package sale

import (
	"github.com/nats-framework/nats/apps/core/routes"
	"github.com/nats-framework/nats/pkg/engine"
)

// App يمثل التطبيق الأساسي
type App struct {
	engine *engine.Engine
	name   string
}

// NewApp ينشئ تطبيقاً أساسياً جديداً
func NewApp(engine *engine.Engine) *App {
	return &App{
		engine: engine,
		name:   "core",
	}
}

// Name يعيد اسم التطبيق
func (a *App) Name() string {
	return a.name
}

// Register يسجل التطبيق
func (a *App) Register() error {
	// ✅ استخدام GetRouter() الذي يعيد *router.Router
	routes.RegisterRoutes(a.engine.GetRouter(), a.engine.Template)
	return nil
}

// Boot يقوم بتهيئة التطبيق
func (a *App) Boot() error {
	return nil
}
