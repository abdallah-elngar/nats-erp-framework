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

// PermissionController متحكم الصلاحيات
type PermissionController struct {
	permissionService *services.PermissionService
}

// NewPermissionController ينشئ متحكم صلاحيات جديداً
func NewPermissionController(db ...*gorm.DB) *PermissionController {
	var permissionService *services.PermissionService
	if len(db) > 0 && db[0] != nil {
		permissionService = services.NewPermissionServiceWithDB(db[0])
	} else {
		permissionService = services.NewPermissionService()
	}
	return &PermissionController{
		permissionService: permissionService,
	}
}

// Index يعيد قائمة الصلاحيات
func (c *PermissionController) Index(w http.ResponseWriter, r *http.Request) {
	resource := r.URL.Query().Get("resource")
	var permissions interface{}
	var err error

	if resource != "" {
		permissions, err = c.permissionService.GetByResource(resource)
	} else {
		permissions, err = c.permissionService.GetAll()
	}

	if err != nil {
		response.InternalError(w, err.Error())
		return
	}

	response.Success(w, permissions)
}

// Show يعيد صلاحية واحدة
func (c *PermissionController) Show(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "Invalid ID")
		return
	}

	permission, err := c.permissionService.GetByID(uint(id))
	if err != nil {
		response.NotFound(w, "Permission not found")
		return
	}

	response.Success(w, permission)
}

// Create ينشئ صلاحية جديدة
func (c *PermissionController) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreatePermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	permission, err := c.permissionService.Create(req)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Created(w, permission)
}

// Delete يحذف صلاحية
func (c *PermissionController) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "Invalid ID")
		return
	}

	if err := c.permissionService.Delete(uint(id)); err != nil {
		response.InternalError(w, err.Error())
		return
	}

	response.NoContent(w)
}
