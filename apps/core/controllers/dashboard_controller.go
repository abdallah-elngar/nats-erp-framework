package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/nats-framework/nats/pkg/response"
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

// ============================================
// ✅ Developer Dashboard (HTML)
// ============================================

// DeveloperDashboard يعرض لوحة المطور الرئيسية
func (c *DashboardController) DeveloperDashboard(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":  "Developer Dashboard",
		"Locale": "ar",
		"User": map[string]interface{}{
			"Username": "admin",
			"FullName": "Administrator",
		},
	}

	// ✅ استخدام المسار المطابق لمكان الملف الفعلي: dashboard/developer
	if err := c.engine.Render(w, "dashboard/developer", data); err != nil {
		// ✅ طباعة الخطأ بالتفصيل
		http.Error(w, fmt.Sprintf("Template Error: %v", err), http.StatusInternalServerError)
		return
	}
}

// ProductionDashboard يعرض لوحة المستخدم النهائي
func (c *DashboardController) ProductionDashboard(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":   "Dashboard",
		"Locale":  "ar",
		"Version": "1.0.0",
		"User": map[string]interface{}{
			"Username": "admin",
			"FullName": "Administrator",
		},
	}

	c.engine.Render(w, "production/dashboard", data)
}

// ============================================
// ✅ صفحات التطبيقات (HTML)
// ============================================

// AppsView يعرض قائمة التطبيقات
func (c *DashboardController) AppsView(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":  "Applications",
		"Locale": "ar",
		"View":   "apps",
	}
	c.engine.Render(w, "developer/apps", data)
}

// CreateAppView يعرض نموذج إنشاء تطبيق
func (c *DashboardController) CreateAppView(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":  "Create Application",
		"Locale": "ar",
		"View":   "create_app",
	}
	c.engine.Render(w, "developer/create_app", data)
}

// EditAppView يعرض نموذج تعديل تطبيق
func (c *DashboardController) EditAppView(w http.ResponseWriter, r *http.Request) {
	app := r.URL.Path[len("/admin/apps/"):]
	app = strings.TrimSuffix(app, "/edit")

	data := map[string]interface{}{
		"Title":  "Edit Application",
		"Locale": "ar",
		"View":   "edit_app",
		"App":    app,
	}
	c.engine.Render(w, "developer/edit_app", data)
}

// ModelsView يعرض نماذج التطبيق
func (c *DashboardController) ModelsView(w http.ResponseWriter, r *http.Request) {
	app := r.URL.Path[len("/admin/apps/"):]
	app = strings.TrimSuffix(app, "/models")

	data := map[string]interface{}{
		"Title":  "Models",
		"Locale": "ar",
		"View":   "models",
		"App":    app,
	}
	c.engine.Render(w, "developer/models", data)
}

// FieldsView يعرض حقول النموذج
func (c *DashboardController) FieldsView(w http.ResponseWriter, r *http.Request) {
	// استخراج app و model من المسار
	path := r.URL.Path
	parts := strings.Split(path, "/")
	app := parts[3]
	model := parts[5]

	data := map[string]interface{}{
		"Title":  "Fields",
		"Locale": "ar",
		"View":   "fields",
		"App":    app,
		"Model":  model,
	}
	c.engine.Render(w, "developer/fields", data)
}

// ============================================
// ✅ صفحات العلاقات والهجرات والمستخدمين
// ============================================

// RelationsView يعرض العلاقات
func (c *DashboardController) RelationsView(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":  "Relations",
		"Locale": "ar",
		"View":   "relations",
	}
	c.engine.Render(w, "developer/relations", data)
}

// MigrationsView يعرض الهجرات
func (c *DashboardController) MigrationsView(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":  "Migrations",
		"Locale": "ar",
		"View":   "migrations",
	}
	c.engine.Render(w, "developer/migrations", data)
}

// UsersView يعرض المستخدمين
func (c *DashboardController) UsersView(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":  "Users",
		"Locale": "ar",
		"View":   "users",
	}
	c.engine.Render(w, "developer/users", data)
}

// ============================================
// ✅ Stats API
// ============================================

// Stats يعيد إحصائيات النظام (JSON)
func (c *DashboardController) Stats(w http.ResponseWriter, r *http.Request) {
	stats := map[string]interface{}{
		"apps":      3,
		"models":    6,
		"users":     1,
		"relations": 2,
	}
	response.Success(w, stats)
}
