package services

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/nats-framework/nats/apps/users/models"
	"github.com/nats-framework/nats/apps/users/repository"
	"github.com/nats-framework/nats/pkg/auth"
	"github.com/nats-framework/nats/pkg/database"
)

// AuthService خدمة المصادقة
type AuthService struct {
	userRepo   *repository.UserRepository
	jwtService *auth.JWTService
	db         *gorm.DB
}

// NewAuthService ينشئ خدمة مصادقة جديدة (باستخدام قاعدة البيانات الافتراضية)
func NewAuthService() *AuthService {
	db := database.DB()
	return &AuthService{
		userRepo: repository.NewUserRepository(db),
		db:       db,
		jwtService: auth.NewJWTService(auth.JWTConfig{
			Secret:     "your-secret-key",
			Expiration: 24 * time.Hour,
			Issuer:     "nats",
		}),
	}
}

// NewAuthServiceWithDB ينشئ خدمة مصادقة جديدة مع قاعدة بيانات محددة
func NewAuthServiceWithDB(db *gorm.DB) *AuthService {
	return &AuthService{
		userRepo: repository.NewUserRepository(db),
		db:       db,
		jwtService: auth.NewJWTService(auth.JWTConfig{
			Secret:     "your-secret-key",
			Expiration: 24 * time.Hour,
			Issuer:     "nats",
		}),
	}
}

// Authenticate يتحقق من بيانات المستخدم
func (s *AuthService) Authenticate(username, password string) (*models.User, error) {
	if s.db == nil {
		return nil, errors.New("database not initialized")
	}

	// البحث عن المستخدم
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// التحقق من كلمة المرور
	if !user.CheckPassword(password) {
		return nil, errors.New("invalid credentials")
	}

	// التحقق من النشاط
	if !user.IsActive() {
		return nil, errors.New("account is inactive")
	}

	return user, nil
}

// CreateSuperUser ينشئ Superuser جديد
func (s *AuthService) CreateSuperUser(username, email, password, fullName string) (*models.User, error) {
	if s.db == nil {
		return nil, errors.New("database not initialized")
	}

	// التحقق من وجود المستخدم
	exists, err := s.userRepo.Exists("username = ? OR email = ?", username, email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username or email already exists")
	}

	// بدء معاملة
	tx := s.db.Begin()

	// إنشاء المستخدم مع IsSuper = true
	user := &models.User{
		Username: username,
		Email:    email,
		Password: password,
		FullName: fullName,
		Status:   "active",
		IsSuper:  true, // ✅ تعيين كـ Superuser
	}

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// الحصول على دور Admin
	var role models.Role
	if err := tx.Where("name = ?", "admin").First(&role).Error; err != nil {
		// إنشاء دور Admin إذا لم يكن موجوداً
		role = models.Role{
			Name:        "admin",
			DisplayName: "Administrator",
			Description: "Full system access",
			IsSystem:    true,
			IsDefault:   false,
		}
		if err := tx.Create(&role).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// ربط المستخدم بدور Admin
	if err := tx.Model(user).Association("Roles").Append(&role); err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return user, nil
}

// GenerateToken يولد توكن JWT
func (s *AuthService) GenerateToken(user *models.User) (string, error) {
	// الحصول على أدوار المستخدم
	roles := make([]string, 0)
	for _, role := range user.Roles {
		roles = append(roles, role.Name)
	}

	return s.jwtService.GenerateToken(&auth.User{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Roles:    roles,
	})
}

// IsSuperUser يتحقق من أن المستخدم هو Superuser
func (s *AuthService) IsSuperUser(userID uint) bool {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return false
	}
	return user.IsSuper
}

// GetSuperUser يعيد Superuser
func (s *AuthService) GetSuperUser() (*models.User, error) {
	if s.db == nil {
		return nil, errors.New("database not initialized")
	}

	var user models.User
	err := s.db.Where("is_super = ?", true).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ListSuperUsers يعيد قائمة Superusers
func (s *AuthService) ListSuperUsers() ([]models.User, error) {
	if s.db == nil {
		return nil, errors.New("database not initialized")
	}

	var users []models.User
	err := s.db.Where("is_super = ?", true).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
