package controllers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"github.com/nats-framework/nats/apps/users/models"
	"github.com/nats-framework/nats/apps/users/services"
	"github.com/nats-framework/nats/pkg/response"
)

// UserController متحكم المستخدمين
type UserController struct {
	userService *services.UserService
}

// NewUserController ينشئ متحكم مستخدمين جديداً
func NewUserController(db *gorm.DB) *UserController {
	return &UserController{
		userService: services.NewUserService(db),
	}
}

// Index يعيد قائمة المستخدمين
func (c *UserController) Index(w http.ResponseWriter, r *http.Request) {
	users, err := c.userService.GetAll()
	if err != nil {
		response.InternalError(w, err.Error())
		return
	}

	response.Success(w, users)
}

// Show يعيد مستخدم واحداً
func (c *UserController) Show(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "Invalid ID")
		return
	}

	user, err := c.userService.GetByID(uint(id))
	if err != nil {
		response.NotFound(w, "User not found")
		return
	}

	response.Success(w, user)
}

// Create ينشئ مستخدماً جديداً
func (c *UserController) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username" validate:"required,min=3,max=50"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
		FullName string `json:"full_name"`
		Role     string `json:"role"`
	}

	if err := response.BindJSON(r, &req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	user, err := c.userService.Create(req.Username, req.Email, req.Password, req.FullName, req.Role)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Created(w, user)
}

// Update يحدث مستخدماً
func (c *UserController) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "Invalid ID")
		return
	}

	var req struct {
		Email    string `json:"email"`
		FullName string `json:"full_name"`
		Role     string `json:"role"`
		Status   string `json:"status"`
	}

	if err := response.BindJSON(r, &req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	user, err := c.userService.Update(uint(id), req.Email, req.FullName, req.Role, req.Status)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Success(w, user)
}

// Delete يحذف مستخدماً
func (c *UserController) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "Invalid ID")
		return
	}

	if err := c.userService.Delete(uint(id)); err != nil {
		response.InternalError(w, err.Error())
		return
	}

	response.NoContent(w)
}

// ✅ Profile يعيد الملف الشخصي للمستخدم الحالي
func (c *UserController) Profile(w http.ResponseWriter, r *http.Request) {
	// الحصول على المستخدم من السياق
	user, ok := getUserFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "User not authenticated")
		return
	}

	// جلب بيانات المستخدم الكاملة
	userData, err := c.userService.GetByID(user.ID)
	if err != nil {
		response.NotFound(w, "User not found")
		return
	}

	response.Success(w, userData)
}

// ✅ UpdateProfile يحدث الملف الشخصي للمستخدم الحالي
func (c *UserController) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	// الحصول على المستخدم من السياق
	user, ok := getUserFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "User not authenticated")
		return
	}

	var req struct {
		Email    string `json:"email"`
		FullName string `json:"full_name"`
		Avatar   string `json:"avatar"`
	}

	if err := response.BindJSON(r, &req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	// تحديث الملف الشخصي
	updatedUser, err := c.userService.UpdateProfile(user.ID, req.Email, req.FullName, req.Avatar)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Success(w, updatedUser)
}

// ✅ ChangePassword يغير كلمة مرور المستخدم
func (c *UserController) ChangePassword(w http.ResponseWriter, r *http.Request) {
	// الحصول على المستخدم من السياق
	user, ok := getUserFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "User not authenticated")
		return
	}

	var req struct {
		CurrentPassword string `json:"current_password" validate:"required"`
		NewPassword     string `json:"new_password" validate:"required,min=8"`
		ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
	}

	if err := response.BindJSON(r, &req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	// تغيير كلمة المرور
	if err := c.userService.ChangePassword(user.ID, req.CurrentPassword, req.NewPassword); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.SuccessMessage(w, "Password changed successfully", nil)
}

// getUserFromContext يحصل على المستخدم من السياق
func getUserFromContext(ctx context.Context) (*models.User, bool) {
	user, ok := ctx.Value("user").(*models.User)
	return user, ok
}
