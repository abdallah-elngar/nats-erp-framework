package dto

import "time"

// PermissionDTO بيانات الصلاحية
type PermissionDTO struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Description string    `json:"description"`
	Resource    string    `json:"resource"`
	Action      string    `json:"action"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreatePermissionRequest طلب إنشاء صلاحية
type CreatePermissionRequest struct {
	Name        string `json:"name" validate:"required"`
	DisplayName string `json:"display_name" validate:"required"`
	Description string `json:"description"`
	Resource    string `json:"resource" validate:"required"`
	Action      string `json:"action" validate:"required"`
}
