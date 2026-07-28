package services

import (
	"errors"

	"gorm.io/gorm"

	"github.com/nats-framework/nats/apps/users/dto"
	"github.com/nats-framework/nats/apps/users/models"
	"github.com/nats-framework/nats/apps/users/repository"
)

// PermissionService خدمة الصلاحيات
type PermissionService struct {
	permissionRepo *repository.PermissionRepository
}

// NewPermissionService ينشئ خدمة صلاحيات جديدة (باستخدام قاعدة البيانات الافتراضية)
func NewPermissionService() *PermissionService {
	return &PermissionService{
		permissionRepo: repository.NewPermissionRepository(),
	}
}

// NewPermissionServiceWithDB ينشئ خدمة صلاحيات جديدة مع قاعدة بيانات محددة
func NewPermissionServiceWithDB(db *gorm.DB) *PermissionService {
	return &PermissionService{
		permissionRepo: repository.NewPermissionRepository(db),
	}
}

// GetAll يعيد جميع الصلاحيات
func (s *PermissionService) GetAll() ([]models.Permission, error) {
	return s.permissionRepo.FindAll()
}

// GetByID يعيد صلاحية بالمعرف
func (s *PermissionService) GetByID(id uint) (*models.Permission, error) {
	return s.permissionRepo.FindByID(id)
}

// GetByResource يعيد الصلاحيات حسب المورد
func (s *PermissionService) GetByResource(resource string) ([]models.Permission, error) {
	return s.permissionRepo.FindByResource(resource)
}

// Create ينشئ صلاحية جديدة
func (s *PermissionService) Create(req dto.CreatePermissionRequest) (*models.Permission, error) {
	// التحقق من وجود الصلاحية
	exists, err := s.permissionRepo.Exists(req.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("permission already exists")
	}

	permission := &models.Permission{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		Resource:    req.Resource,
		Action:      req.Action,
	}

	if err := s.permissionRepo.Create(permission); err != nil {
		return nil, err
	}

	return permission, nil
}

// Delete يحذف صلاحية
func (s *PermissionService) Delete(id uint) error {
	// التحقق من وجود الصلاحية
	_, err := s.permissionRepo.FindByID(id)
	if err != nil {
		return errors.New("permission not found")
	}

	return s.permissionRepo.Delete(id)
}

// GetByNames يعيد الصلاحيات حسب الأسماء
func (s *PermissionService) GetByNames(names []string) ([]models.Permission, error) {
	return s.permissionRepo.GetByNames(names)
}

// GetByResourceAndAction يعيد صلاحية حسب المورد والإجراء
func (s *PermissionService) GetByResourceAndAction(resource, action string) (*models.Permission, error) {
	return s.permissionRepo.FindByResourceAndAction(resource, action)
}
