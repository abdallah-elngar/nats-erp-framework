package seeds

import (
    "gorm.io/gorm"

    "github.com/nats-framework/nats/apps/users/models"
)

// SeedSuperUser يضيف مستخدم Superuser افتراضي
func SeedSuperUser(db *gorm.DB) error {
    // التحقق من وجود Superuser
    var count int64
    db.Model(&models.User{}).Where("is_super = ?", true).Count(&count)

    if count > 0 {
        return nil // Superuser موجود بالفعل
    }

    // إنشاء Superuser
    superuser := &models.User{
        Username: "admin",
        Email:    "admin@example.com",
        Password: "admin123",
        FullName: "Super Administrator",
        Status:   "active",
        IsSuper:  true,
    }

    if err := db.Create(superuser).Error; err != nil {
        return err
    }

    // إنشاء دور Admin إذا لم يكن موجوداً
    var role models.Role
    if err := db.Where("name = ?", "admin").First(&role).Error; err == gorm.ErrRecordNotFound {
        role = models.Role{
            Name:        "admin",
            DisplayName: "Administrator",
            Description: "Full system access",
            IsSystem:    true,
            IsDefault:   false,
        }
        if err := db.Create(&role).Error; err != nil {
            return err
        }
    }

    // ربط Superuser بدور Admin
    if err := db.Model(superuser).Association("Roles").Append(&role); err != nil {
        return err
    }

    return nil
}