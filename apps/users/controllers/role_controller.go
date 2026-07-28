package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"github.com/nats-framework/nats/apps/users/dto"
	"github.com/nats-framework/nats/apps/users/services"
	"github.com/nats-framework/nats/pkg/response"
)

// RoleController متحكم الأدوار
type RoleController struct {
	roleService *services.RoleService
}

// NewRoleController ينشئ متحكم أدوار جديداً
func NewRoleController(db ...*gorm.DB) *RoleController {
	var roleService *services.RoleService
	if len(db) > 0 && db[0] != nil {
		roleService = services.NewRoleServiceWithDB(db[0])
	} else {
		roleService = services.NewRoleService()
	}
	return &RoleController{
		roleService: roleService,
	}
}

// Index يعيد قائمة الأدوار
func (c *RoleController) Index(w http.ResponseWriter, r *http.Request) {
	roles, err := c.roleService.GetAll()
	if err != nil {
		response.InternalError(w, err.Error())
		return
	}

	response.Success(w, roles)
}

// Show يعيد دوراً واحداً
func (c *RoleController) Show(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "Invalid ID")
		return
	}

	role, err := c.roleService.GetByID(uint(id))
	if err != nil {
		response.NotFound(w, "Role not found")
		return
	}

	response.Success(w, role)
}

// Create ينشئ دوراً جديداً
func (c *RoleController) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	role, err := c.roleService.Create(req)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Created(w, role)
}

// Update يحدث دوراً
func (c *RoleController) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "Invalid ID")
		return
	}

	var req dto.UpdateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	role, err := c.roleService.Update(uint(id), req)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Success(w, role)
}

// Delete يحذف دوراً
func (c *RoleController) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "Invalid ID")
		return
	}

	if err := c.roleService.Delete(uint(id)); err != nil {
		response.InternalError(w, err.Error())
		return
	}

	response.NoContent(w)
}

// AssignPermissions يعين صلاحيات لدور
func (c *RoleController) AssignPermissions(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "Invalid ID")
		return
	}

	var req struct {
		PermissionIDs []uint `json:"permission_ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	if err := c.roleService.AssignPermissions(uint(id), req.PermissionIDs); err != nil {
		response.InternalError(w, err.Error())
		return
	}

	response.SuccessMessage(w, "Permissions assigned successfully", nil)
}
