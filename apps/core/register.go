package core

import (
	"github.com/nats-framework/nats/pkg/router"
	"github.com/nats-framework/nats/pkg/template"
)

// Register يسجل التطبيق في المحرك
func Register(router *router.Router, template *template.Engine) error {
	appInstance := NewApp(router, template)

	if err := appInstance.Register(); err != nil {
		return err
	}

	if err := appInstance.Boot(); err != nil {
		return err
	}

	return nil
}
