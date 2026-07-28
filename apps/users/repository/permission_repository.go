package repository

import (
	"gorm.io/gorm"

	"github.com/nats-framework/nats/apps/users/models"
	"github.com/nats-framework/nats/pkg/database"
)

// PermissionRepository مستودع الصلاحيات
type PermissionRepository struct {
	db *gorm.DB
}

// NewPermissionRepository ينشئ مستودع صلاحيات جديد
func NewPermissionRepository(db ...*gorm.DB) *PermissionRepository {
	var gormDB *gorm.DB
	if len(db) > 0 && db[0] != nil {
		gormDB = db[0]
	} else {
		gormDB = database.DB()
	}
	return &PermissionRepository{
		db: gormDB,
	}
}

// NewPermissionRepositoryWithDB ينشئ مستودع صلاحيات جديد مع قاعدة بيانات محددة
func NewPermissionRepositoryWithDB(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{
		db: db,
	}
}

// Create ينشئ صلاحية جديدة
func (r *PermissionRepository) Create(permission *models.Permission) error {
	return r.db.Create(permission).Error
}

// FindByID يبحث عن صلاحية بالمعرف
func (r *PermissionRepository) FindByID(id uint) (*models.Permission, error) {
	var permission models.Permission
	err := r.db.First(&permission, id).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// FindByName يبحث عن صلاحية بالاسم
func (r *PermissionRepository) FindByName(name string) (*models.Permission, error) {
	var permission models.Permission
	err := r.db.Where("name = ?", name).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// FindAll يعيد جميع الصلاحيات
func (r *PermissionRepository) FindAll() ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Find(&permissions).Error
	return permissions, err
}

// FindByResource يعيد الصلاحيات حسب المورد
func (r *PermissionRepository) FindByResource(resource string) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Where("resource = ?", resource).Find(&permissions).Error
	return permissions, err
}

// FindByAction يعيد الصلاحيات حسب الإجراء
func (r *PermissionRepository) FindByAction(action string) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Where("action = ?", action).Find(&permissions).Error
	return permissions, err
}

// Update يحدث صلاحية
func (r *PermissionRepository) Update(permission *models.Permission) error {
	return r.db.Save(permission).Error
}

// Delete يحذف صلاحية
func (r *PermissionRepository) Delete(id uint) error {
	return r.db.Delete(&models.Permission{}, id).Error
}

// Exists يتحقق من وجود صلاحية
func (r *PermissionRepository) Exists(name string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Permission{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}

// GetByNames يعيد الصلاحيات حسب الأسماء
func (r *PermissionRepository) GetByNames(names []string) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Where("name IN ?", names).Find(&permissions).Error
	return permissions, err
}

// CountByResource يعيد عدد الصلاحيات حسب المورد
func (r *PermissionRepository) CountByResource(resource string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Permission{}).Where("resource = ?", resource).Count(&count).Error
	return count, err
}

// FindByResourceAndAction يعيد صلاحية حسب المورد والإجراء
func (r *PermissionRepository) FindByResourceAndAction(resource, action string) (*models.Permission, error) {
	var permission models.Permission
	err := r.db.Where("resource = ? AND action = ?", resource, action).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}
