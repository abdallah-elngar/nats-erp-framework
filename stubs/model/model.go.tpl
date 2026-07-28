package models

import (
    "time"

    "gorm.io/gorm"
)

// {{.Model.Name}} يمثل {{.Model.Name}} في النظام
type {{.Model.Name}} struct {
    ID        uint           `gorm:"primaryKey"`
    {{range .Model.Fields}}
    {{.Name}} {{.Type}} `gorm:"{{.Tags}}"`
    {{end}}
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName يعيد اسم الجدول
func ({{.Model.Name}}) TableName() string {
    return "{{.Model.Name | lower}}s"
}

// BeforeCreate يقوم بمعالجة قبل الإنشاء
func (m *{{.Model.Name}}) BeforeCreate(tx *gorm.DB) error {
    // أي معالجة إضافية
    return nil
}

// BeforeUpdate يقوم بمعالجة قبل التحديث
func (m *{{.Model.Name}}) BeforeUpdate(tx *gorm.DB) error {
    // أي معالجة إضافية
    return nil
}