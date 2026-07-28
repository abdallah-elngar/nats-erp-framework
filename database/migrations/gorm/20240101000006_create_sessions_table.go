package migrations

import (
	"gorm.io/gorm"

	"github.com/nats-framework/nats/apps/users/models"
)

// Up_20240101000000 ينشئ جدول الجلسات
func Up_20240101000000(db *gorm.DB) error {
	return db.AutoMigrate(&models.Session{})
}

// Down_20240101000000 يحذف جدول الجلسات
func Down_20240101000000(db *gorm.DB) error {
	return db.Migrator().DropTable("sessions")
}
