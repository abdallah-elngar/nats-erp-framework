package migrations

import (
	"gorm.io/gorm"

	"github.com/nats-framework/nats/apps/core/models"
)

// Up_20240101000000 ينشئ جدول الإعدادات
func Up_20240101000000(db *gorm.DB) error {
	return db.AutoMigrate(&models.Setting{})
}

// Down_20240101000000 يحذف جدول الإعدادات
func Down_20240101000000(db *gorm.DB) error {
	return db.Migrator().DropTable("settings")
}
