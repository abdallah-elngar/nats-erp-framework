// apps/core/controllers/dashboard_controller.go - النسخة النهائية

package controllers

import (
	"net/http"

	"github.com/nats-framework/nats/pkg/template"
)

// DashboardController متحكم لوحة التحكم
type DashboardController struct {
	engine *template.Engine
}

// NewDashboardController ينشئ متحكم لوحة التحكم
func NewDashboardController(engine *template.Engine) *DashboardController {
	return &DashboardController{
		engine: engine,
	}
}

// DeveloperDashboard يعرض لوحة تحكم المطور
func (c *DashboardController) DeveloperDashboard(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":  "Developer Dashboard",
		"Locale": "ar",
	}

	// ✅ محاولة عرض القالب مع التعامل مع الأخطاء
	if err := c.engine.RenderWriter(w, "dashboard/developer", data); err != nil {
		// ✅ إذا فشل، جرب بدون المجلد dashboard
		if err2 := c.engine.RenderWriter(w, "developer", data); err2 != nil {
			// ✅ إذا فشل، عرض خطأ واضح
			http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		}
	}
}

// ProductionDashboard يعرض لوحة تحكم المستخدم النهائي
func (c *DashboardController) ProductionDashboard(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":  "Dashboard",
		"Locale": "ar",
	}

	if err := c.engine.RenderWriter(w, "dashboard/production", data); err != nil {
		if err2 := c.engine.RenderWriter(w, "production", data); err2 != nil {
			http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		}
	}
}
