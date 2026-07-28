package models

import (
	"time"

	"gorm.io/gorm"
)

// Session يمثل جلسة مستخدم في قاعدة البيانات
type Session struct {
	ID        string    `gorm:"primaryKey;size:64"`
	UserID    uint      `gorm:"index;not null"`
	Data      string    `gorm:"type:text"`
	IP        string    `gorm:"size:45"`
	UserAgent string    `gorm:"size:255"`
	ExpiresAt time.Time `gorm:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// العلاقات
	User User `gorm:"foreignKey:UserID"`
}

// TableName يعيد اسم الجدول
func (Session) TableName() string {
	return "sessions"
}

// IsExpired يتحقق من انتهاء الجلسة
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// GetTTL يعيد الوقت المتبقي للجلسة
func (s *Session) GetTTL() time.Duration {
	if s.IsExpired() {
		return 0
	}
	return time.Until(s.ExpiresAt)
}

// BeforeCreate يقوم بمعالجة قبل الإنشاء
func (s *Session) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = generateSessionID()
	}
	if s.ExpiresAt.IsZero() {
		s.ExpiresAt = time.Now().Add(24 * time.Hour) // الجلسة لمدة 24 ساعة
	}
	return nil
}

// generateSessionID يولد معرف جلسة فريد
func generateSessionID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(32)
}

// randomString يولد نصاً عشوائياً
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}
