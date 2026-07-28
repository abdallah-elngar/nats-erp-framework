package seeds

import (
	"gorm.io/gorm"

	"github.com/nats-framework/nats/apps/users/seeds"
)

// SeedAll يقوم بتشغيل جميع البذور
func SeedAll(db *gorm.DB) error {
	// إضافة الأدوار الافتراضية
	if err := seeds.SeedDefaultRoles(db); err != nil {
		return err
	}

	// إضافة الصلاحيات الافتراضية
	if err := seeds.SeedDefaultPermissions(db); err != nil {
		return err
	}

	return nil
}

// SeedApp يقوم بتشغيل بذور تطبيق معين
func SeedApp(db *gorm.DB, appName string) error {
	switch appName {
	case "users":
		if err := seeds.SeedDefaultRoles(db); err != nil {
			return err
		}
		if err := seeds.SeedDefaultPermissions(db); err != nil {
			return err
		}
	default:
		return nil
	}

	return nil
}
