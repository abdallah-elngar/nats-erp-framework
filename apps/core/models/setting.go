package models

import (
	"time"

	"gorm.io/gorm"
)

// Setting يمثل إعدادات النظام
type Setting struct {
	ID          uint   `gorm:"primaryKey"`
	Key         string `gorm:"uniqueIndex;not null;size:100"`
	Value       string `gorm:"type:text"`
	Group       string `gorm:"index;size:50"`
	Description string `gorm:"size:255"`
	IsPublic    bool   `gorm:"default:false"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// TableName يعيد اسم الجدول
func (Setting) TableName() string {
	return "settings"
}

// GetValue يعيد قيمة الإعداد
func (s *Setting) GetValue() string {
	return s.Value
}

// IsActive يتحقق من نشاط الإعداد
func (s *Setting) IsActive() bool {
	return s.DeletedAt.Time.IsZero()
}
