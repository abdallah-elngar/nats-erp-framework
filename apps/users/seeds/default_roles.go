package seeds

import (
	"gorm.io/gorm"

	"github.com/nats-framework/nats/apps/users/models"
)

// SeedDefaultRoles يضيف الأدوار الافتراضية
func SeedDefaultRoles(db *gorm.DB) error {
	roles := []models.Role{
		{
			Name:        "admin",
			DisplayName: "Administrator",
			Description: "Full system access",
			IsDefault:   false,
			IsSystem:    true,
		},
		{
			Name:        "manager",
			DisplayName: "Manager",
			Description: "Manage users and content",
			IsDefault:   false,
			IsSystem:    true,
		},
		{
			Name:        "user",
			DisplayName: "User",
			Description: "Basic user access",
			IsDefault:   true,
			IsSystem:    true,
		},
	}

	for _, role := range roles {
		var existing models.Role
		if err := db.Where("name = ?", role.Name).First(&existing).Error; err == gorm.ErrRecordNotFound {
			if err := db.Create(&role).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
