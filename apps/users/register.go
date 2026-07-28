package users

import (
	"github.com/nats-framework/nats/pkg/engine"
)

// Register يسجل التطبيق في المحرك
func Register(app *engine.Engine) error {
	appInstance := NewApp(app)

	if err := appInstance.Register(); err != nil {
		return err
	}

	if err := appInstance.Boot(); err != nil {
		return err
	}

	return nil
}

// RegisterMigrations يسجل هجرات المستخدمين
func RegisterMigrations(migrator interface{}) error {
	// سيتم تنفيذها لاحقاً
	return nil
}
