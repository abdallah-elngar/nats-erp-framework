package services

import (
	"errors"

	"gorm.io/gorm"

	"github.com/nats-framework/nats/apps/users/dto"
	"github.com/nats-framework/nats/apps/users/models"
	"github.com/nats-framework/nats/apps/users/repository"
)

// RoleService خدمة الأدوار
type RoleService struct {
	roleRepo       *repository.RoleRepository
	permissionRepo *repository.PermissionRepository
}

// NewRoleService ينشئ خدمة أدوار جديدة (باستخدام قاعدة البيانات الافتراضية)
func NewRoleService() *RoleService {
	return &RoleService{
		roleRepo:       repository.NewRoleRepository(),
		permissionRepo: repository.NewPermissionRepository(),
	}
}

// NewRoleServiceWithDB ينشئ خدمة أدوار جديدة مع قاعدة بيانات محددة
func NewRoleServiceWithDB(db *gorm.DB) *RoleService {
	return &RoleService{
		roleRepo:       repository.NewRoleRepository(db),
		permissionRepo: repository.NewPermissionRepository(db),
	}
}

// GetAll يعيد جميع الأدوار
func (s *RoleService) GetAll() ([]models.Role, error) {
	return s.roleRepo.FindAll()
}

// GetByID يعيد دوراً بالمعرف
func (s *RoleService) GetByID(id uint) (*models.Role, error) {
	return s.roleRepo.FindByID(id)
}

// Create ينشئ دوراً جديداً
func (s *RoleService) Create(req dto.CreateRoleRequest) (*models.Role, error) {
	// التحقق من وجود الدور
	exists, err := s.roleRepo.Exists(req.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("role already exists")
	}

	role := &models.Role{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		IsDefault:   false,
		IsSystem:    false,
	}

	if err := s.roleRepo.Create(role); err != nil {
		return nil, err
	}

	// إضافة الصلاحيات
	if len(req.Permissions) > 0 {
		permissions, err := s.permissionRepo.GetByNames(req.Permissions)
		if err != nil {
			return nil, err
		}

		for _, perm := range permissions {
			if err := s.roleRepo.AssignPermission(role.ID, perm.ID); err != nil {
				return nil, err
			}
		}
	}

	return s.roleRepo.FindByID(role.ID)
}

// Update يحدث دوراً
func (s *RoleService) Update(id uint, req dto.UpdateRoleRequest) (*models.Role, error) {
	role, err := s.roleRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("role not found")
	}

	if req.DisplayName != "" {
		role.DisplayName = req.DisplayName
	}
	if req.Description != "" {
		role.Description = req.Description
	}

	if err := s.roleRepo.Update(role); err != nil {
		return nil, err
	}

	// تحديث الصلاحيات
	if req.Permissions != nil {
		permissions, err := s.permissionRepo.GetByNames(req.Permissions)
		if err != nil {
			return nil, err
		}

		var permissionIDs []uint
		for _, perm := range permissions {
			permissionIDs = append(permissionIDs, perm.ID)
		}

		if err := s.roleRepo.SyncPermissions(role.ID, permissionIDs); err != nil {
			return nil, err
		}
	}

	return s.roleRepo.FindByID(role.ID)
}

// Delete يحذف دوراً
func (s *RoleService) Delete(id uint) error {
	// التحقق من وجود الدور
	_, err := s.roleRepo.FindByID(id)
	if err != nil {
		return errors.New("role not found")
	}

	return s.roleRepo.Delete(id)
}

// AssignPermissions يعين صلاحيات لدور
func (s *RoleService) AssignPermissions(roleID uint, permissionIDs []uint) error {
	// التحقق من وجود الدور
	_, err := s.roleRepo.FindByID(roleID)
	if err != nil {
		return errors.New("role not found")
	}

	return s.roleRepo.SyncPermissions(roleID, permissionIDs)
}
