package migrations

import (
    "gorm.io/gorm"
)

// Up_20260724232353_create_sale_table يقوم بتطبيق الهجرة
func Up_20260724232353_create_sale_table(db *gorm.DB) error {
    // TODO: Add your migration logic here
    return nil
}

// Down_20260724232353_create_sale_table يقوم بالتراجع عن الهجرة
func Down_20260724232353_create_sale_table(db *gorm.DB) error {
    // TODO: Add your rollback logic here
    return nil
}
