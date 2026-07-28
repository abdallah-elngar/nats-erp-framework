package core

import (
	"github.com/nats-framework/nats/apps/core/routes"
	"github.com/nats-framework/nats/pkg/router"
	"github.com/nats-framework/nats/pkg/template"
)

// App يمثل التطبيق الأساسي
type App struct {
	router   *router.Router
	template *template.Engine
	name     string
}

// NewApp ينشئ تطبيقاً أساسياً جديداً
func NewApp(router *router.Router, template *template.Engine) *App {
	return &App{
		router:   router,
		template: template,
		name:     "core",
	}
}

// Name يعيد اسم التطبيق
func (a *App) Name() string {
	return a.name
}

// Register يسجل التطبيق
func (a *App) Register() error {
	// ✅ تسجيل المسارات
	routes.RegisterRoutes(a.router, a.template)
	return nil
}

// Boot يقوم بتهيئة التطبيق
func (a *App) Boot() error {
	return nil
}
