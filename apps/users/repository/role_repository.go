package repository

import (
	"gorm.io/gorm"

	"github.com/nats-framework/nats/apps/users/models"
	"github.com/nats-framework/nats/pkg/database"
)

// RoleRepository مستودع الأدوار
type RoleRepository struct {
	db *gorm.DB
}

// NewRoleRepository ينشئ مستودع أدوار جديد
func NewRoleRepository(db ...*gorm.DB) *RoleRepository {
	var gormDB *gorm.DB
	if len(db) > 0 && db[0] != nil {
		gormDB = db[0]
	} else {
		gormDB = database.DB()
	}
	return &RoleRepository{
		db: gormDB,
	}
}

// NewRoleRepositoryWithDB ينشئ مستودع أدوار جديد مع قاعدة بيانات محددة
func NewRoleRepositoryWithDB(db *gorm.DB) *RoleRepository {
	return &RoleRepository{
		db: db,
	}
}

// Create ينشئ دوراً جديداً
func (r *RoleRepository) Create(role *models.Role) error {
	return r.db.Create(role).Error
}

// FindByID يبحث عن دور بالمعرف
func (r *RoleRepository) FindByID(id uint) (*models.Role, error) {
	var role models.Role
	err := r.db.Preload("Permissions").First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// FindByName يبحث عن دور بالاسم
func (r *RoleRepository) FindByName(name string) (*models.Role, error) {
	var role models.Role
	err := r.db.Preload("Permissions").Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// FindAll يعيد جميع الأدوار
func (r *RoleRepository) FindAll() ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Preload("Permissions").Find(&roles).Error
	return roles, err
}

// FindDefault يعيد الأدوار الافتراضية
func (r *RoleRepository) FindDefault() ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Where("is_default = ?", true).Find(&roles).Error
	return roles, err
}

// Update يحدث دوراً
func (r *RoleRepository) Update(role *models.Role) error {
	return r.db.Save(role).Error
}

// Delete يحذف دوراً
func (r *RoleRepository) Delete(id uint) error {
	return r.db.Delete(&models.Role{}, id).Error
}

// AssignPermission يعين صلاحية لدور
func (r *RoleRepository) AssignPermission(roleID, permissionID uint) error {
	var role models.Role
	if err := r.db.First(&role, roleID).Error; err != nil {
		return err
	}

	var permission models.Permission
	if err := r.db.First(&permission, permissionID).Error; err != nil {
		return err
	}

	return r.db.Model(&role).Association("Permissions").Append(&permission)
}

// RemovePermission يزيل صلاحية من دور
func (r *RoleRepository) RemovePermission(roleID, permissionID uint) error {
	var role models.Role
	if err := r.db.First(&role, roleID).Error; err != nil {
		return err
	}

	var permission models.Permission
	if err := r.db.First(&permission, permissionID).Error; err != nil {
		return err
	}

	return r.db.Model(&role).Association("Permissions").Delete(&permission)
}

// GetPermissions يعيد صلاحيات الدور
func (r *RoleRepository) GetPermissions(roleID uint) ([]models.Permission, error) {
	var role models.Role
	if err := r.db.First(&role, roleID).Error; err != nil {
		return nil, err
	}

	var permissions []models.Permission
	err := r.db.Model(&role).Association("Permissions").Find(&permissions)
	return permissions, err
}

// SyncPermissions يزامن صلاحيات الدور
func (r *RoleRepository) SyncPermissions(roleID uint, permissionIDs []uint) error {
	var role models.Role
	if err := r.db.First(&role, roleID).Error; err != nil {
		return err
	}

	var permissions []models.Permission
	if err := r.db.Find(&permissions, permissionIDs).Error; err != nil {
		return err
	}

	return r.db.Model(&role).Association("Permissions").Replace(permissions)
}

// Exists يتحقق من وجود دور
func (r *RoleRepository) Exists(name string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Role{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}
