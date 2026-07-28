package migrations

import (
	"gorm.io/gorm"

	"github.com/nats-framework/nats/apps/users/models"
)

// Up_0001 ينشئ جميع جداول المستخدمين
func Up_0001(db *gorm.DB) error {
	// ترتيب الإنشاء حسب الاعتماديات
	if err := db.AutoMigrate(&models.User{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&models.Role{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&models.Permission{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&models.Session{}); err != nil {
		return err
	}
	return nil
}

// Up_0001_relations ينشئ جداول العلاقات
func Up_0001_relations(db *gorm.DB) error {
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS user_roles (
			user_id BIGINT NOT NULL,
			role_id BIGINT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (user_id, role_id),
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
		)
	`).Error; err != nil {
		return err
	}

	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS role_permissions (
			role_id BIGINT NOT NULL,
			permission_id BIGINT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (role_id, permission_id),
			FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
			FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
		)
	`).Error; err != nil {
		return err
	}
	return nil
}

// Up_0001_add_is_super يضيف عمود is_super
func Up_0001_add_is_super(db *gorm.DB) error {
	return db.Exec(`
		ALTER TABLE users ADD COLUMN IF NOT EXISTS is_super BOOLEAN DEFAULT false
	`).Error
}

// Down_0001 يحذف جميع الجداول
func Down_0001(db *gorm.DB) error {
	if err := db.Migrator().DropTable("user_roles"); err != nil {
		return err
	}
	if err := db.Migrator().DropTable("role_permissions"); err != nil {
		return err
	}
	if err := db.Migrator().DropTable("sessions"); err != nil {
		return err
	}
	if err := db.Migrator().DropTable("permissions"); err != nil {
		return err
	}
	if err := db.Migrator().DropTable("roles"); err != nil {
		return err
	}
	if err := db.Migrator().DropTable("users"); err != nil {
		return err
	}
	return nil
}