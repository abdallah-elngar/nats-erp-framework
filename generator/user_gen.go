package generator

import (
	"fmt"
)

// UserGenerator يولد ملفات المستخدمين
type UserGenerator struct {
	generator *Generator
}

// NewUserGenerator ينشئ مولد مستخدمين جديد
func NewUserGenerator(generator *Generator) *UserGenerator {
	return &UserGenerator{
		generator: generator,
	}
}

// GenerateAuth يولد نظام المصادقة
func (ug *UserGenerator) GenerateAuth(appName string) error {
	templates := map[string]string{
		"auth_controller.go.tpl": fmt.Sprintf("apps/%s/controllers/auth_controller.go", appName),
		"auth_middleware.go.tpl": fmt.Sprintf("apps/%s/middleware/auth.go", appName),
		"auth_service.go.tpl":    fmt.Sprintf("apps/%s/services/auth_service.go", appName),
		"auth_routes.go.tpl":     fmt.Sprintf("apps/%s/routes/auth.go", appName),
		"login.html.tpl":         fmt.Sprintf("apps/%s/templates/pages/auth/login.html", appName),
		"register.html.tpl":      fmt.Sprintf("apps/%s/templates/pages/auth/register.html", appName),
	}

	for tpl, output := range templates {
		if err := ug.generator.Generate(tpl, output, map[string]string{
			"AppName": appName,
		}); err != nil {
			return err
		}
	}

	return nil
}
