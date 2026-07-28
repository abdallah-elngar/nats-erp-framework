package models

import (
	"time"

	"gorm.io/gorm"
)

// Role يمثل دوراً في النظام
type Role struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"uniqueIndex;not null;size:50"`
	DisplayName string `gorm:"size:100"`
	Description string `gorm:"size:255"`
	IsDefault   bool   `gorm:"default:false"`
	IsSystem    bool   `gorm:"default:false"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	// العلاقات
	Users       []User       `gorm:"many2many:user_roles;"`
	Permissions []Permission `gorm:"many2many:role_permissions;"`
}

// TableName يعيد اسم الجدول
func (Role) TableName() string {
	return "roles"
}

// BeforeCreate يقوم بمعالجة قبل الإنشاء
func (r *Role) BeforeCreate(tx *gorm.DB) error {
	// أي معالجة إضافية
	return nil
}
