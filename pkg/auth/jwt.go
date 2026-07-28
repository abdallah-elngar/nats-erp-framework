package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTService خدمة JWT
type JWTService struct {
	secret     []byte
	expiration time.Duration
	issuer     string
}

// JWTConfig يمثل إعدادات JWT
type JWTConfig struct {
	Secret     string
	Expiration time.Duration
	Issuer     string
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
			return nil, errors.New("unexpected signing method")
		}
		return j.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshToken يقوم بتجديد التوكن
func (j *JWTService) RefreshToken(tokenString string) (string, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	// إنشاء توكن جديد
	user := &User{
		ID:       claims.UserID,
		Username: claims.Username,
		Email:    claims.Email,
		Roles:    claims.Roles,
	}

	return j.GenerateToken(user)
}
