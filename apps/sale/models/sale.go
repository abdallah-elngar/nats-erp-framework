package models

import (
    "time"

    "gorm.io/gorm"
)

// Sale يمثل Sale في النظام
type Sale struct {
    ID        uint           `gorm:"primaryKey"`
    SaName string `gorm:"not null;unique;"`
    Price float64 
    Quantity int 
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName يعيد اسم الجدول
func (Sale) TableName() string {
    return "sales"
}
