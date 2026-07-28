package {{.AppName}}

import (
    "github.com/nats-framework/nats/pkg/engine"
)

// Register يسجل التطبيق في المحرك
func Register(app *engine.Engine) error {
    // إنشاء التطبيق
    appInstance := NewApp(app)

    // تسجيل التطبيق
    if err := appInstance.Register(); err != nil {
        return err
    }

    // تهيئة التطبيق
    if err := appInstance.Boot(); err != nil {
        return err
    }

    return nil
}

func init() {
    // تسجيل التطبيق في المحرك
    // engine.RegisterApp("{{.AppName}}", Register)
}