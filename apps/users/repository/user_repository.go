package repository

import (
	"time"

	"gorm.io/gorm"

	"github.com/nats-framework/nats/apps/users/models"
	"github.com/nats-framework/nats/pkg/database"
)

// UserRepository مستودع المستخدمين
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository ينشئ مستودع مستخدمين جديد
func NewUserRepository(db ...*gorm.DB) *UserRepository {
	var gormDB *gorm.DB
	if len(db) > 0 && db[0] != nil {
		gormDB = db[0]
	} else {
		gormDB = database.DB()
	}
	return &UserRepository{
		db: gormDB,
	}
}

// NewUserRepositoryWithDB ينشئ مستودع مستخدمين جديد مع قاعدة بيانات محددة
// (مستخدم للتوافق مع الكود القديم)
func NewUserRepositoryWithDB(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create ينشئ مستخدماً جديداً
func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// FindByID يبحث عن مستخدم بالمعرف
func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Roles").Preload("Roles.Permissions").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername يبحث عن مستخدم باسم المستخدم
func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Roles").Preload("Roles.Permissions").Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail يبحث عن مستخدم بالبريد الإلكتروني
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Roles").Preload("Roles.Permissions").Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindAll يعيد جميع المستخدمين
func (r *UserRepository) FindAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Preload("Roles").Find(&users).Error
	return users, err
}

// Update يحدث مستخدماً
func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete يحذف مستخدماً (ناعم)
func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

// DeletePermanently يحذف مستخدماً بشكل دائم
func (r *UserRepository) DeletePermanently(id uint) error {
	return r.db.Unscoped().Delete(&models.User{}, id).Error
}

// Exists يتحقق من وجود مستخدم
func (r *UserRepository) Exists(query string, args ...interface{}) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where(query, args...).Count(&count).Error
	return count > 0, err
}

// FindByRole يعيد المستخدمين حسب الدور
func (r *UserRepository) FindByRole(roleName string) ([]models.User, error) {
	var users []models.User
	err := r.db.Joins("JOIN user_roles ON user_roles.user_id = users.id").
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Where("roles.name = ?", roleName).
		Find(&users).Error
	return users, err
}

// UpdateLastLogin يحدث وقت آخر تسجيل دخول
func (r *UserRepository) UpdateLastLogin(id uint) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("last_login", time.Now()).Error
}

// AssignRole يعين دوراً لمستخدم
func (r *UserRepository) AssignRole(userID, roleID uint) error {
	var user models.User
	if err := r.db.First(&user, userID).Error; err != nil {
		return err
	}

	var role models.Role
	if err := r.db.First(&role, roleID).Error; err != nil {
		return err
	}

	return r.db.Model(&user).Association("Roles").Append(&role)
}

// RemoveRole يزيل دوراً من مستخدم
func (r *UserRepository) RemoveRole(userID, roleID uint) error {
	var user models.User
	if err := r.db.First(&user, userID).Error; err != nil {
		return err
	}

	var role models.Role
	if err := r.db.First(&role, roleID).Error; err != nil {
		return err
	}

	return r.db.Model(&user).Association("Roles").Delete(&role)
}

// GetRoles يعيد أدوار المستخدم
func (r *UserRepository) GetRoles(userID uint) ([]models.Role, error) {
	var user models.User
	if err := r.db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	var roles []models.Role
	err := r.db.Model(&user).Association("Roles").Find(&roles)
	return roles, err
}
