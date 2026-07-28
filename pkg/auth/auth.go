package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrTokenExpired       = errors.New("token expired")
	ErrInvalidToken       = errors.New("invalid token")
)

// User يمثل مستخدم النظام
type User struct {
	ID        uint
	Username  string
	Email     string
	Password  string
	Roles     []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Claims يمثل مطالبات JWT
type Claims struct {
	UserID   uint     `json:"user_id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

// AuthService خدمة المصادقة
type AuthService struct {
	secretKey  []byte
	expiration time.Duration
}

// NewAuthService ينشئ خدمة مصادقة جديدة
func NewAuthService(secretKey string, expiration time.Duration) *AuthService {
	return &AuthService{
		secretKey:  []byte(secretKey),
		expiration: expiration,
	}
}

// HashPassword يقوم بتشفير كلمة المرور
func (s *AuthService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyPassword يقوم بالتحقق من كلمة المرور
func (s *AuthService) VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateToken يقوم بإنشاء توكن JWT
func (s *AuthService) GenerateToken(user *User) (string, error) {
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Roles:    user.Roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "nats",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

// ValidateToken يقوم بالتحقق من صحة التوكن
func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return s.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// GetUserFromContext يحصل على المستخدم من السياق
func GetUserFromContext(ctx context.Context) (*User, bool) {
	user, ok := ctx.Value("user").(*User)
	return user, ok
}

// SetUserInContext يضع المستخدم في السياق
func SetUserInContext(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, "user", user)
}

// ============================================
// ✅ دوال عامة للاستخدام المباشر
// ============================================

var defaultJWTService *JWTService

// InitJWTService يقوم بتهيئة خدمة JWT الافتراضية
func InitJWTService(config JWTConfig) {
	defaultJWTService = NewJWTService(config)
}

// ValidateToken يقوم بالتحقق من صحة التوكن (دالة عامة)
func ValidateToken(tokenString string) (*Claims, error) {
	if defaultJWTService == nil {
		// إذا لم يتم تهيئة الخدمة، استخدم إعدادات افتراضية
		defaultJWTService = NewJWTService(JWTConfig{
			Secret:     "default-secret-key-change-in-production",
			Expiration: 24 * time.Hour,
			Issuer:     "nats",
		})
	}
	return defaultJWTService.ValidateToken(tokenString)
}

// GenerateToken يقوم بإنشاء توكن (دالة عامة)
func GenerateToken(user *User) (string, error) {
	if defaultJWTService == nil {
		defaultJWTService = NewJWTService(JWTConfig{
			Secret:     "default-secret-key-change-in-production",
			Expiration: 24 * time.Hour,
			Issuer:     "nats",
		})
	}
	return defaultJWTService.GenerateToken(user)
}
