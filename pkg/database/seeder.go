package database

import (
    "gorm.io/gorm"

    "github.com/nats-framework/nats/apps/users/seeds"
)

// SeedAll يقوم بتشغيل جميع البذور
func SeedAll(db *gorm.DB) error {
    // إضافة Superuser
    if err := seeds.SeedSuperUser(db); err != nil {
        return err
    }

    return nil
}