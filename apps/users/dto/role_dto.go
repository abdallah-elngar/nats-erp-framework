package dto

import "time"

// RoleDTO بيانات الدور
type RoleDTO struct {
	ID          uint            `json:"id"`
	Name        string          `json:"name"`
	DisplayName string          `json:"display_name"`
	Description string          `json:"description"`
	IsDefault   bool            `json:"is_default"`
	IsSystem    bool            `json:"is_system"`
	Permissions []PermissionDTO `json:"permissions,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// CreateRoleRequest طلب إنشاء دور
type CreateRoleRequest struct {
	Name        string   `json:"name" validate:"required"`
	DisplayName string   `json:"display_name" validate:"required"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}

// UpdateRoleRequest طلب تحديث دور
type UpdateRoleRequest struct {
	DisplayName string   `json:"display_name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}
