package security

import (
	"sync"
)

// Permission يمثل صلاحية
type Permission struct {
	ID       uint
	Name     string
	Resource string
	Action   string
}

// Role يمثل دوراً
type Role struct {
	ID          uint
	Name        string
	DisplayName string
	Permissions []Permission
}

// User يمثل مستخدم
type User struct {
	ID    uint
	Roles []Role
}

// RBAC يدير نظام الصلاحيات
type RBAC struct {
	permissions map[string]Permission
	roles       map[string]Role
	mu          sync.RWMutex
}

// NewRBAC ينشئ نظام RBAC جديد
func NewRBAC() *RBAC {
	return &RBAC{
		permissions: make(map[string]Permission),
		roles:       make(map[string]Role),
	}
}

// RegisterPermission يسجل صلاحية
func (r *RBAC) RegisterPermission(name, resource, action string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.permissions[name] = Permission{
		Name:     name,
		Resource: resource,
		Action:   action,
	}
}

// RegisterRole يسجل دوراً
func (r *RBAC) RegisterRole(name, displayName string, permissions []string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var perms []Permission
	for _, permName := range permissions {
		if perm, ok := r.permissions[permName]; ok {
			perms = append(perms, perm)
		}
	}

	r.roles[name] = Role{
		Name:        name,
		DisplayName: displayName,
		Permissions: perms,
	}
}

// HasPermission يتحقق من وجود صلاحية
func (r *RBAC) HasPermission(user *User, resource, action string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, role := range user.Roles {
		if role.Name == "admin" {
			return true
		}

		if r.hasRolePermission(role, resource, action) {
			return true
		}
	}

	return false
}

// hasRolePermission يتحقق من صلاحية دور
func (r *RBAC) hasRolePermission(role Role, resource, action string) bool {
	for _, perm := range role.Permissions {
		if perm.Resource == resource && perm.Action == action {
			return true
		}
	}
	return false
}

// GetRole يعيد دوراً
func (r *RBAC) GetRole(name string) (Role, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	role, ok := r.roles[name]
	return role, ok
}

// GetPermission يعيد صلاحية
func (r *RBAC) GetPermission(name string) (Permission, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	perm, ok := r.permissions[name]
	return perm, ok
}
