package routes

import (
	"gorm.io/gorm"

	"github.com/nats-framework/nats/apps/users/controllers"
	"github.com/nats-framework/nats/pkg/router"
)

// RegisterRoutes يسجل مسارات تطبيق المستخدمين
func RegisterRoutes(r *router.Router, db *gorm.DB) {
	authCtrl := controllers.NewAuthController()
	userCtrl := controllers.NewUserController(db)
	roleCtrl := controllers.NewRoleController(db)
	permCtrl := controllers.NewPermissionController(db)

	// ============================================
	// مسارات المصادقة (عامة)
	// ============================================
	r.Group(func(r *router.Router) {
		r.Post("/api/auth/login", authCtrl.Login)
		r.Post("/api/auth/register", authCtrl.Register)
		r.Post("/api/auth/logout", authCtrl.Logout)
		r.Get("/api/auth/check", authCtrl.Check)

		// مسارات استعادة كلمة المرور
		r.Post("/api/auth/forgot-password", authCtrl.ForgotPassword)
		r.Post("/api/auth/reset-password", authCtrl.ResetPassword)
	})

	// ============================================
	// مسارات المستخدمين (محمية)
	// ============================================
	r.Group(func(r *router.Router) {
		// إضافة ميدلوير المصادقة
		r.Use(controllers.AuthMiddleware)

		// المستخدمين
		r.Get("/api/users", userCtrl.Index)
		r.Get("/api/users/{id}", userCtrl.Show)
		r.Post("/api/users", userCtrl.Create)
		r.Put("/api/users/{id}", userCtrl.Update)
		r.Delete("/api/users/{id}", userCtrl.Delete)

		// الملف الشخصي
		r.Get("/api/users/profile", userCtrl.Profile)
		r.Put("/api/users/profile", userCtrl.UpdateProfile)

		// الأدوار
		r.Get("/api/roles", roleCtrl.Index)
		r.Get("/api/roles/{id}", roleCtrl.Show)
		r.Post("/api/roles", roleCtrl.Create)
		r.Put("/api/roles/{id}", roleCtrl.Update)
		r.Delete("/api/roles/{id}", roleCtrl.Delete)

		// الصلاحيات
		r.Get("/api/permissions", permCtrl.Index)
		r.Get("/api/permissions/{id}", permCtrl.Show)
		r.Post("/api/permissions", permCtrl.Create)
		r.Delete("/api/permissions/{id}", permCtrl.Delete)
	})
}
