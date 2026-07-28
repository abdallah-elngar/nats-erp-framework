package migrations

import (
    "gorm.io/gorm"

    "github.com/nats-framework/nats/apps/sale/models"
)

// Up_20260724232352 ينشئ جدول sales
func Up_20260724232352(db *gorm.DB) error {
    return db.AutoMigrate(&models.Sale{})
}

// Down_20260724232352 يحذف جدول sales
func Down_20260724232352(db *gorm.DB) error {
    return db.Migrator().DropTable("sales")
}
