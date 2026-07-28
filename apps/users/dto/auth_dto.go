package dto

import "time"

// LoginRequest طلب تسجيل الدخول
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse استجابة تسجيل الدخول
type LoginResponse struct {
	Token string  `json:"token"`
	User  UserDTO `json:"user"`
}

// RegisterRequest طلب التسجيل
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"full_name"`
}

// UserDTO بيانات المستخدم
type UserDTO struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Avatar    string    `json:"avatar"`
	Status    string    `json:"status"`
	Roles     []RoleDTO `json:"roles"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
