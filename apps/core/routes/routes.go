package routes

import (
	"net/http"

	"github.com/nats-framework/nats/apps/core/controllers"
	"github.com/nats-framework/nats/pkg/router"
	"github.com/nats-framework/nats/pkg/template"
)

// RegisterRoutes يسجل مسارات التطبيق الأساسي
func RegisterRoutes(r *router.Router, engine *template.Engine) {
	// إنشاء المتحكمات
	dashboardCtrl := controllers.NewDashboardController(engine)
	appsCtrl := controllers.NewAppsController(engine)
	settingsCtrl := controllers.NewSettingsController(engine)

	// ✅ مسار الصحة (يجب أن يكون أولاً)
	r.Get("/health", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","message":"NATS Framework is running"}`))
	})

	// ✅ الصفحة الرئيسية - إعادة توجيه إلى Developer Dashboard
	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, "/admin/developer", http.StatusFound)
	})

	// ✅ Developer Dashboard
	r.Get("/admin/developer", dashboardCtrl.DeveloperDashboard)

	// ✅ Production Dashboard
	r.Get("/admin/production", dashboardCtrl.ProductionDashboard)

	// ✅ الإعدادات
	r.Get("/admin/settings", settingsCtrl.Index)
	r.Put("/api/admin/settings", settingsCtrl.Update)

	// ============================================
	// مسارات API
	// ============================================
	r.Route("/api/admin", func(r *router.Router) {
		// التطبيقات
		r.Get("/apps", appsCtrl.ListApps)
		r.Get("/apps/{app}/models", appsCtrl.GetAppModels)
		r.Post("/apps", appsCtrl.CreateApp)
		r.Delete("/apps/{app}", appsCtrl.DeleteApp)

		// العلاقات
		r.Get("/relations", appsCtrl.ListRelations)
		r.Post("/relations", appsCtrl.CreateRelation)
		r.Delete("/relations", appsCtrl.DeleteRelation)

		// الحقول
		r.Post("/apps/{app}/models/{model}/fields", appsCtrl.AddField)
		r.Delete("/apps/{app}/models/{model}/fields/{field}", appsCtrl.DeleteField)

		// المستخدمين
		r.Post("/users", appsCtrl.CreateUser)

		// الإحصائيات
		r.Get("/stats", dashboardCtrl.Stats)

		// الهجرات
		r.Post("/migrations/run", appsCtrl.RunMigrations)
		r.Post("/migrations/reset", appsCtrl.ResetMigrations)
	})
}
