package controllers

import (
	"net/http"

	"github.com/nats-framework/nats/pkg/response"
	"github.com/nats-framework/nats/pkg/template"
)

// BaseController هو المتحكم الأساسي
type BaseController struct {
	engine *template.Engine
}

// NewBaseController ينشئ متحكم أساسي جديد
func NewBaseController(engine *template.Engine) *BaseController {
	return &BaseController{
		engine: engine,
	}
}

// Render يعرض قالباً (استخدام مباشر)
func (c *BaseController) Render(w http.ResponseWriter, name string, data interface{}) error {
	if data == nil {
		data = make(map[string]interface{})
	}
	return c.engine.RenderWriter(w, name, data)
}

// JSON يرسل استجابة JSON
func (c *BaseController) JSON(w http.ResponseWriter, status int, data interface{}) {
	response.JSON(w, status, data)
}

// Success يرسل استجابة نجاح
func (c *BaseController) Success(w http.ResponseWriter, data interface{}) {
	response.Success(w, data)
}

// Error يرسل استجابة خطأ
func (c *BaseController) Error(w http.ResponseWriter, status int, message string) {
	response.Error(w, status, message)
}
