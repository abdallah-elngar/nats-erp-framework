package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User يمثل مستخدم النظام
type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"uniqueIndex;not null;size:50"`
	Email     string `gorm:"uniqueIndex;not null;size:100"`
	Password  string `gorm:"not null;size:255"`
	FullName  string `gorm:"size:100"`
	Avatar    string `gorm:"size:255"`
	Status    string `gorm:"default:active;size:20"` // active, inactive, suspended
	IsSuper   bool   `gorm:"default:false;index"`    // ✅ عمود السوبر يوزر
	LastLogin *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// العلاقات
	Roles    []Role    `gorm:"many2many:user_roles;"`
	Sessions []Session `gorm:"foreignKey:UserID"`
}

// TableName يعيد اسم الجدول
func (User) TableName() string {
	return "users"
}

// BeforeCreate يقوم بمعالجة قبل الإنشاء
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashed)
	}
	return nil
}

// BeforeUpdate يقوم بمعالجة قبل التحديث
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	// فقط إذا تغيرت كلمة المرور
	if tx.Statement.Changed("Password") && u.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashed)
	}
	return nil
}

// CheckPassword يتحقق من كلمة المرور
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// IsActive يتحقق من نشاط المستخدم
func (u *User) IsActive() bool {
	return u.Status == "active"
}

// IsSuperUser يتحقق من أن المستخدم هو Superuser
func (u *User) IsSuperUser() bool {
	return u.IsSuper
}

// HasRole يتحقق من وجود دور معين
func (u *User) HasRole(roleName string) bool {
	for _, role := range u.Roles {
		if role.Name == roleName {
			return true
		}
	}
	return false
}

// HasPermission يتحقق من وجود صلاحية معينة
func (u *User) HasPermission(permission string) bool {
	if u.IsSuper {
		return true
	}

	for _, role := range u.Roles {
		for _, perm := range role.Permissions {
			if perm.Name == permission {
				return true
			}
		}
	}
	return false
}
