package {{.AppName}}

import (
    "github.com/nats-framework/nats/pkg/engine"
)

// App يمثل تطبيق {{.AppTitle}}
type App struct {
    engine *engine.Engine
    name   string
}

// NewApp ينشئ تطبيقاً جديداً
func NewApp(engine *engine.Engine) *App {
    return &App{
        engine: engine,
        name:   "{{.AppName}}",
    }
}

// Name يعيد اسم التطبيق
func (a *App) Name() string {
    return a.name
}

// Register يسجل التطبيق
func (a *App) Register() error {
    // تسجيل المسارات
    // تسجيل النماذج
    // تسجيل الخدمات
    return nil
}

// Boot يقوم بتهيئة التطبيق
func (a *App) Boot() error {
    // تهيئة التطبيق
    return nil
}