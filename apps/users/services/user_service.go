package services

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/nats-framework/nats/apps/users/models"
	"github.com/nats-framework/nats/apps/users/repository"
)

// UserService خدمة المستخدمين
type UserService struct {
	userRepo *repository.UserRepository
	roleRepo *repository.RoleRepository
}

// NewUserService ينشئ خدمة مستخدمين جديدة
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		userRepo: repository.NewUserRepository(db),
		roleRepo: repository.NewRoleRepository(db),
	}
}

// GetAll يعيد جميع المستخدمين
func (s *UserService) GetAll() ([]models.User, error) {
	return s.userRepo.FindAll()
}

// GetByID يعيد مستخدم بالمعرف
func (s *UserService) GetByID(id uint) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

// Create ينشئ مستخدماً جديداً
func (s *UserService) Create(username, email, password, fullName, roleName string) (*models.User, error) {
	// التحقق من وجود المستخدم
	exists, err := s.userRepo.Exists("username = ? OR email = ?", username, email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username or email already exists")
	}

	// تشفير كلمة المرور
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// إنشاء المستخدم
	user := &models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
		FullName: fullName,
		Status:   "active",
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// تعيين الدور
	if roleName != "" {
		role, err := s.roleRepo.FindByName(roleName)
		if err == nil {
			user.Roles = []models.Role{*role}
			s.userRepo.Update(user)
		}
	}

	return user, nil
}

// Update يحدث مستخدماً
func (s *UserService) Update(id uint, email, fullName, roleName, status string) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if email != "" {
		user.Email = email
	}
	if fullName != "" {
		user.FullName = fullName
	}
	if status != "" {
		user.Status = status
	}

	if roleName != "" {
		role, err := s.roleRepo.FindByName(roleName)
		if err == nil {
			user.Roles = []models.Role{*role}
		}
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// ✅ UpdateProfile يحدث الملف الشخصي
func (s *UserService) UpdateProfile(id uint, email, fullName, avatar string) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if email != "" {
		user.Email = email
	}
	if fullName != "" {
		user.FullName = fullName
	}
	if avatar != "" {
		user.Avatar = avatar
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// ✅ ChangePassword يغير كلمة المرور
func (s *UserService) ChangePassword(id uint, currentPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}

	// التحقق من كلمة المرور الحالية
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword)); err != nil {
		return errors.New("current password is incorrect")
	}

	// تشفير كلمة المرور الجديدة
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return s.userRepo.Update(user)
}

// Delete يحذف مستخدماً
func (s *UserService) Delete(id uint) error {
	return s.userRepo.Delete(id)
}

// HasPermission يتحقق من صلاحية المستخدم
func (s *UserService) HasPermission(userID uint, permission string) (bool, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return false, err
	}

	// المدير لديه كل الصلاحيات
	for _, role := range user.Roles {
		if role.Name == "admin" {
			return true, nil
		}
	}

	// التحقق من الصلاحيات
	for _, role := range user.Roles {
		for _, perm := range role.Permissions {
			if perm.Name == permission {
				return true, nil
			}
		}
	}

	return false, nil
}
