package dto

import "time"

// CreateUserRequest طلب إنشاء مستخدم
type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

// UpdateUserRequest طلب تحديث مستخدم
type UpdateUserRequest struct {
	Email    string `json:"email" validate:"omitempty,email"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
	Status   string `json:"status" validate:"omitempty,oneof=active inactive suspended"`
}

// UserResponse استجابة المستخدم
type UserResponse struct {
	ID        uint       `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	FullName  string     `json:"full_name"`
	Avatar    string     `json:"avatar"`
	Status    string     `json:"status"`
	LastLogin *time.Time `json:"last_login,omitempty"`
	Roles     []RoleDTO  `json:"roles"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// ChangePasswordRequest طلب تغيير كلمة المرور
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}
