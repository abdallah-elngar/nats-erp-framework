package controllers

import (
	"context"
	"net/http"
	"strings"

	"github.com/nats-framework/nats/pkg/response"
)

// AuthController متحكم المصادقة
type AuthController struct{}

// NewAuthController ينشئ متحكم مصادقة جديداً
func NewAuthController() *AuthController {
	return &AuthController{}
}

// Login يسجل دخول المستخدم
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	if err := response.BindJSON(r, &req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	// TODO: تنفيذ منطق تسجيل الدخول
	response.Success(w, map[string]interface{}{
		"message": "Login successful",
		"token":   "mock-jwt-token",
		"user": map[string]interface{}{
			"id":       1,
			"username": req.Username,
			"email":    "admin@example.com",
		},
	})
}

// Register يسجل مستخدم جديداً
func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username" validate:"required,min=3,max=50"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
		FullName string `json:"full_name"`
	}

	if err := response.BindJSON(r, &req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	// TODO: تنفيذ منطق التسجيل
	response.Created(w, map[string]interface{}{
		"message": "User registered successfully",
		"user":    req,
	})
}

// Logout يسجل خروج المستخدم
func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	response.SuccessMessage(w, "Logged out successfully", nil)
}

// Check يتحقق من حالة المصادقة
func (c *AuthController) Check(w http.ResponseWriter, r *http.Request) {
	// التحقق من وجود توكن
	token := r.Header.Get("Authorization")
	if token == "" {
		response.Success(w, map[string]interface{}{
			"authenticated": false,
		})
		return
	}

	// TODO: التحقق من صحة التوكن
	response.Success(w, map[string]interface{}{
		"authenticated": true,
		"user": map[string]interface{}{
			"id":       1,
			"username": "admin",
			"email":    "admin@example.com",
		},
	})
}

// ForgotPassword يرسل رابط استعادة كلمة المرور
func (c *AuthController) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email" validate:"required,email"`
	}

	if err := response.BindJSON(r, &req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.SuccessMessage(w, "Password reset link sent to your email", nil)
}

// ResetPassword يعيد تعيين كلمة المرور
func (c *AuthController) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token           string `json:"token" validate:"required"`
		Password        string `json:"password" validate:"required,min=8"`
		ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
	}

	if err := response.BindJSON(r, &req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.SuccessMessage(w, "Password reset successfully", nil)
}

// AuthMiddleware يتحقق من المصادقة (ميدلوير)
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// الحصول على التوكن من الـ Header
		token := r.Header.Get("Authorization")
		if token == "" {
			response.Unauthorized(w, "Authorization required")
			return
		}

		// إزالة "Bearer " من التوكن
		token = strings.TrimPrefix(token, "Bearer ")
		if token == "" {
			response.Unauthorized(w, "Invalid token format")
			return
		}

		// TODO: التحقق من صحة التوكن
		// وضع المستخدم في السياق
		ctx := context.WithValue(r.Context(), "user", map[string]interface{}{
			"id":       1,
			"username": "admin",
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
