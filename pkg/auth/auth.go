package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

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

// ============================================
// تعريفات الأنواع
// ============================================

// ContextKey هو نوع مخصص لمفاتيح السياق
type ContextKey string

const (
	// UserKey هو مفتاح المستخدم في السياق
	UserKey ContextKey = "user"
	// UserIDKey هو مفتاح معرف المستخدم في السياق
	UserIDKey ContextKey = "user_id"
	// TokenKey هو مفتاح التوكن في السياق
	TokenKey ContextKey = "token"
	// PermissionsKey هو مفتاح الصلاحيات في السياق
	PermissionsKey ContextKey = "permissions"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrTokenExpired       = errors.New("token expired")
	ErrInvalidToken       = errors.New("invalid token")
)

// ============================================
// تعريفات الأنواع
// ============================================

// User يمثل مستخدم النظام
type User struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Password  string    `json:"-"`
	IsSuper   bool      `json:"is_super"`
	Status    string    `json:"status"`
	Roles     []string  `json:"roles"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Claims يمثل مطالبات JWT
type Claims struct {
	UserID   uint     `json:"user_id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

// JWTConfig يمثل إعدادات JWT
type JWTConfig struct {
	Secret     string
	Expiration time.Duration
	Issuer     string
}

// ============================================
// دوال المستخدم (User)
// ============================================

// IsActive يتحقق من نشاط المستخدم
func (u *User) IsActive() bool {
	return u.Status == "active"
}

// IsAdmin يتحقق من أن المستخدم مدير
func (u *User) IsAdmin() bool {
	return u.IsSuper || u.HasRole("admin")
}

// HasRole يتحقق من وجود دور معين
func (u *User) HasRole(role string) bool {
	for _, r := range u.Roles {
		if r == role {
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
	return false
}

// ============================================
// دوال السياق (Context Functions)
// ============================================

// GetUserFromContext يحصل على المستخدم من السياق
func GetUserFromContext(ctx context.Context) (*User, bool) {
	user, ok := ctx.Value(UserKey).(*User)
	return user, ok
}

// GetUserIDFromContext يحصل على معرف المستخدم من السياق
func GetUserIDFromContext(ctx context.Context) (uint, bool) {
	userID, ok := ctx.Value(UserIDKey).(uint)
	return userID, ok
}

// GetTokenFromContext يحصل على التوكن من السياق
func GetTokenFromContext(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(TokenKey).(string)
	return token, ok
}

// GetPermissionsFromContext يحصل على صلاحيات المستخدم من السياق
func GetPermissionsFromContext(ctx context.Context) ([]string, bool) {
	perms, ok := ctx.Value(PermissionsKey).([]string)
	return perms, ok
}

// SetUserInContext يضع المستخدم في السياق
func SetUserInContext(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, UserKey, user)
}

// SetUserIDInContext يضع معرف المستخدم في السياق
func SetUserIDInContext(ctx context.Context, userID uint) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// SetTokenInContext يضع التوكن في السياق
func SetTokenInContext(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, TokenKey, token)
}

// SetPermissionsInContext يضع الصلاحيات في السياق
func SetPermissionsInContext(ctx context.Context, permissions []string) context.Context {
	return context.WithValue(ctx, PermissionsKey, permissions)
}

// IsAuthenticated يتحقق من وجود مستخدم في السياق
func IsAuthenticated(ctx context.Context) bool {
	_, ok := GetUserFromContext(ctx)
	return ok
}

// GetUsernameFromContext يحصل على اسم المستخدم من السياق
func GetUsernameFromContext(ctx context.Context) string {
	if user, ok := GetUserFromContext(ctx); ok && user != nil {
		return user.Username
	}
	return "system"
}

// GetUserEmailFromContext يحصل على بريد المستخدم من السياق
func GetUserEmailFromContext(ctx context.Context) string {
	if user, ok := GetUserFromContext(ctx); ok && user != nil {
		return user.Email
	}
	return ""
}

// ============================================
// خدمة JWT
// ============================================

// JWTService خدمة JWT
type JWTService struct {
	secret     []byte
	expiration time.Duration
	issuer     string
}

// NewJWTService ينشئ خدمة JWT جديدة
func NewJWTService(config JWTConfig) *JWTService {
	return &JWTService{
		secret:     []byte(config.Secret),
		expiration: config.Expiration,
		issuer:     config.Issuer,
	}
}

// GenerateToken يولد توكن JWT
func (j *JWTService) GenerateToken(user *User) (string, error) {
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Roles:    user.Roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    j.issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

// ValidateToken يتحقق من صحة التوكن
func (j *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return j.secret, nil
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

// RefreshToken يقوم بتجديد التوكن
func (j *JWTService) RefreshToken(tokenString string) (string, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	user := &User{
		ID:       claims.UserID,
		Username: claims.Username,
		Email:    claims.Email,
		Roles:    claims.Roles,
	}

	return j.GenerateToken(user)
}

// ============================================
// دوال مساعدة عامة (Global Helpers)
// ============================================

var defaultJWTService *JWTService

// InitJWTService يقوم بتهيئة خدمة JWT الافتراضية
func InitJWTService(config JWTConfig) {
	defaultJWTService = NewJWTService(config)
}

// ValidateToken يقوم بالتحقق من صحة التوكن (دالة عامة)
func ValidateToken(tokenString string) (*Claims, error) {
	if defaultJWTService == nil {
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

// HashPassword يقوم بتشفير كلمة المرور
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyPassword يقوم بالتحقق من كلمة المرور
func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
