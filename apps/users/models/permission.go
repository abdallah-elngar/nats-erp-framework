package models

import (
	"time"

	"gorm.io/gorm"
)

// Permission يمثل صلاحية في النظام
type Permission struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"uniqueIndex;not null;size:100"`
	DisplayName string `gorm:"size:100"`
	Description string `gorm:"size:255"`
	Resource    string `gorm:"size:50"`
	Action      string `gorm:"size:50"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	// العلاقات
	Roles []Role `gorm:"many2many:role_permissions;"`
}

// TableName يعيد اسم الجدول
func (Permission) TableName() string {
	return "permissions"
}
