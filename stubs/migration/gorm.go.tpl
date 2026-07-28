package migrations

import (
    "gorm.io/gorm"
)

// {{.ID}}_create_{{.Model.Name | lower}}_table ينشئ جدول {{.Model.Name | lower}}
func Up_{{.ID}}(db *gorm.DB) error {
    return db.AutoMigrate(&{{.Model.Name}}{})
}

// Down_{{.ID}} يحذف جدول {{.Model.Name | lower}}
func Down_{{.ID}}(db *gorm.DB) error {
    return db.Migrator().DropTable("{{.Model.Name | lower}}s")
}