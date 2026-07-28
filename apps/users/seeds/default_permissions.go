package seeds

import (
	"gorm.io/gorm"

	"github.com/nats-framework/nats/apps/users/models"
)

// SeedDefaultPermissions يضيف الصلاحيات الافتراضية
func SeedDefaultPermissions(db *gorm.DB) error {
	permissions := []models.Permission{
		// صلاحيات المستخدمين
		{Name: "users.view", DisplayName: "View Users", Resource: "users", Action: "view"},
		{Name: "users.create", DisplayName: "Create Users", Resource: "users", Action: "create"},
		{Name: "users.edit", DisplayName: "Edit Users", Resource: "users", Action: "edit"},
		{Name: "users.delete", DisplayName: "Delete Users", Resource: "users", Action: "delete"},

		// صلاحيات الأدوار
		{Name: "roles.view", DisplayName: "View Roles", Resource: "roles", Action: "view"},
		{Name: "roles.create", DisplayName: "Create Roles", Resource: "roles", Action: "create"},
		{Name: "roles.edit", DisplayName: "Edit Roles", Resource: "roles", Action: "edit"},
		{Name: "roles.delete", DisplayName: "Delete Roles", Resource: "roles", Action: "delete"},

		// صلاحيات النظام
		{Name: "system.settings", DisplayName: "System Settings", Resource: "system", Action: "settings"},
		{Name: "system.apps", DisplayName: "Manage Apps", Resource: "system", Action: "apps"},
	}

	for _, perm := range permissions {
		var existing models.Permission
		if err := db.Where("name = ?", perm.Name).First(&existing).Error; err == gorm.ErrRecordNotFound {
			if err := db.Create(&perm).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
