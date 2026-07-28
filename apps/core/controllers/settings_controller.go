package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/nats-framework/nats/pkg/template"
)

// SettingsController متحكم الإعدادات
type SettingsController struct {
	BaseController
}

// NewSettingsController ينشئ متحكم الإعدادات
func NewSettingsController(engine *template.Engine) *SettingsController {
	return &SettingsController{
		BaseController: BaseController{engine: engine},
	}
}

// Index يعرض صفحة الإعدادات (HTML)
func (c *SettingsController) Index(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":  "الإعدادات",
		"Locale": "ar",
	}

	if err := c.Render(w, "settings/index", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Update يقوم بتحديث الإعدادات (API)
func (c *SettingsController) Update(w http.ResponseWriter, r *http.Request) {
	var settings map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		c.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	c.Success(w, map[string]interface{}{
		"message":  "Settings updated successfully",
		"settings": settings,
	})
}
