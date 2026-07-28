package users

import (
	"github.com/nats-framework/nats/apps/users/models"
	"github.com/nats-framework/nats/apps/users/routes"
	"github.com/nats-framework/nats/pkg/engine"
	"gorm.io/gorm"
)

// App يمثل تطبيق المستخدمين
type App struct {
	engine *engine.Engine
	name   string
}

// NewApp ينشئ تطبيق مستخدمين جديداً
func NewApp(engine *engine.Engine) *App {
	return &App{
		engine: engine,
		name:   "users",
	}
}

// Name يعيد اسم التطبيق
func (a *App) Name() string {
	return a.name
}

// Register يسجل التطبيق
func (a *App) Register() error {
	var db *gorm.DB

	// تسجيل النماذج في قاعدة البيانات
	if a.engine.DB != nil {
		db = a.engine.DB.DB()
		if db != nil {
			// تشغيل الهجرات التلقائية
			if err := db.AutoMigrate(&models.User{}, &models.Role{}, &models.Permission{}, &models.Session{}); err != nil {
				return err
			}
		}
	}

	// ✅ تسجيل المسارات - نمرر الـ Router و قاعدة البيانات
	routes.RegisterRoutes(a.engine.GetRouter(), db)

	return nil
}

// Boot يقوم بتهيئة التطبيق
func (a *App) Boot() error {
	// تهيئة المستخدمين (إضافة بيانات افتراضية إذا لزم الأمر)
	if a.engine.DB != nil {
		db := a.engine.DB.DB()
		if db != nil {
			// إضافة الأدوار والصلاحيات الافتراضية
			// سيتم تنفيذها لاحقاً
		}
	}
	return nil
}
