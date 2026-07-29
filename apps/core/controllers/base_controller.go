package controllers

import (
	"net/http"

	"github.com/nats-framework/nats/pkg/response"
	"github.com/nats-framework/nats/pkg/template"
)

// BaseController هو الأساس لجميع المتحكمات
type BaseController struct {
	engine *template.Engine
}

// Success يرسل استجابة نجاح
func (c *BaseController) Success(w http.ResponseWriter, data interface{}) {
	response.Success(w, data)
}

// SuccessMessage يرسل استجابة نجاح مع رسالة
func (c *BaseController) SuccessMessage(w http.ResponseWriter, message string, data interface{}) {
	response.SuccessMessage(w, message, data)
}

// Error يرسل استجابة خطأ
func (c *BaseController) Error(w http.ResponseWriter, status int, message string) {
	response.Error(w, status, message)
}

// BadRequest يرسل خطأ 400
func (c *BaseController) BadRequest(w http.ResponseWriter, message string) {
	response.BadRequest(w, message)
}

// Unauthorized يرسل خطأ 401
func (c *BaseController) Unauthorized(w http.ResponseWriter, message string) {
	response.Unauthorized(w, message)
}

// Forbidden يرسل خطأ 403
func (c *BaseController) Forbidden(w http.ResponseWriter, message string) {
	response.Forbidden(w, message)
}

// NotFound يرسل خطأ 404
func (c *BaseController) NotFound(w http.ResponseWriter, message string) {
	response.NotFound(w, message)
}

// InternalError يرسل خطأ 500
func (c *BaseController) InternalError(w http.ResponseWriter, message string) {
	response.InternalError(w, message)
}

// Render يعرض قالباً
func (c *BaseController) Render(w http.ResponseWriter, name string, data interface{}) error {
	if c.engine == nil {
		return nil
	}
	return c.engine.Render(w, name, data)
}
