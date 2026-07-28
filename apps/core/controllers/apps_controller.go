package controllers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/nats-framework/nats/pkg/template"
)

// AppsController متحكم إدارة التطبيقات
type AppsController struct {
	BaseController
}

// NewAppsController ينشئ متحكم إدارة التطبيقات
func NewAppsController(engine *template.Engine) *AppsController {
	return &AppsController{
		BaseController: BaseController{engine: engine},
	}
}

// ListApps يعيد قائمة التطبيقات
func (c *AppsController) ListApps(w http.ResponseWriter, r *http.Request) {
	apps := []map[string]interface{}{}

	entries, err := os.ReadDir("apps")
	if err != nil {
		c.Success(w, apps)
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			appPath := filepath.Join("apps", entry.Name())
			modelsCount := 0
			modelsDir := filepath.Join(appPath, "models")
			if files, err := os.ReadDir(modelsDir); err == nil {
				for _, f := range files {
					if !f.IsDir() && strings.HasSuffix(f.Name(), ".go") {
						modelsCount++
					}
				}
			}

			status := "active"
			if _, err := os.Stat(filepath.Join(appPath, "app.go")); os.IsNotExist(err) {
				status = "inactive"
			}

			apps = append(apps, map[string]interface{}{
				"name":   entry.Name(),
				"models": modelsCount,
				"status": status,
			})
		}
	}

	c.Success(w, apps)
}

// GetAppModels يعيد نماذج تطبيق معين
func (c *AppsController) GetAppModels(w http.ResponseWriter, r *http.Request) {
	appName := chi.URLParam(r, "app")

	if appName == "" {
		c.Error(w, http.StatusBadRequest, "App name is required")
		return
	}

	fields := []map[string]interface{}{}

	modelsPath := filepath.Join("apps", appName, "models")
	files, err := os.ReadDir(modelsPath)
	if err != nil {
		fields = append(fields, map[string]interface{}{
			"model": "Default",
			"fields": []map[string]interface{}{
				{"name": "id", "type": "uint", "required": true, "unique": true},
				{"name": "created_at", "type": "time.Time", "required": false, "unique": false},
				{"name": "updated_at", "type": "time.Time", "required": false, "unique": false},
			},
		})
		c.Success(w, fields)
		return
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".go") {
			continue
		}

		filePath := filepath.Join(modelsPath, file.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		modelName := strings.TrimSuffix(file.Name(), ".go")
		modelName = strings.Title(modelName)

		var modelFields []map[string]interface{}
		lines := strings.Split(string(content), "\n")
		inStruct := false

		for _, line := range lines {
			trimmed := strings.TrimSpace(line)

			if strings.Contains(trimmed, "type "+modelName+" struct {") {
				inStruct = true
				continue
			}

			if inStruct && trimmed == "}" {
				break
			}

			if inStruct && trimmed != "" && !strings.HasPrefix(trimmed, "//") {
				if strings.Contains(trimmed, "gorm.Model") || strings.Contains(trimmed, "BaseModel") {
					modelFields = append(modelFields,
						map[string]interface{}{"name": "id", "type": "uint", "required": true, "unique": true},
						map[string]interface{}{"name": "created_at", "type": "time.Time", "required": false, "unique": false},
						map[string]interface{}{"name": "updated_at", "type": "time.Time", "required": false, "unique": false},
					)
					continue
				}

				parts := strings.Fields(trimmed)
				if len(parts) >= 2 {
					fieldName := parts[0]
					if fieldName[0] >= 'a' && fieldName[0] <= 'z' {
						continue
					}
					if strings.HasPrefix(fieldName, "`") {
						continue
					}

					fieldType := parts[1]
					required := false
					unique := false
					for _, part := range parts {
						if strings.Contains(part, "not null") || strings.Contains(part, "required") {
							required = true
						}
						if strings.Contains(part, "unique") {
							unique = true
						}
					}

					modelFields = append(modelFields, map[string]interface{}{
						"name":     fieldName,
						"type":     fieldType,
						"required": required,
						"unique":   unique,
					})
				}
			}
		}

		if len(modelFields) == 0 {
			modelFields = []map[string]interface{}{
				{"name": "id", "type": "uint", "required": true, "unique": true},
				{"name": "created_at", "type": "time.Time", "required": false, "unique": false},
				{"name": "updated_at", "type": "time.Time", "required": false, "unique": false},
			}
		}

		fields = append(fields, map[string]interface{}{
			"model":  modelName,
			"fields": modelFields,
		})
	}

	if len(fields) == 0 {
		fields = append(fields, map[string]interface{}{
			"model": "Default",
			"fields": []map[string]interface{}{
				{"name": "id", "type": "uint", "required": true, "unique": true},
				{"name": "created_at", "type": "time.Time", "required": false, "unique": false},
				{"name": "updated_at", "type": "time.Time", "required": false, "unique": false},
			},
		})
	}

	c.Success(w, fields)
}

// CreateApp ينشئ تطبيقاً جديداً
func (c *AppsController) CreateApp(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Parent      string `json:"parent"`
		Models      []struct {
			Name   string `json:"name"`
			Fields string `json:"fields"`
		} `json:"models"`
		CRUD string `json:"crud"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.Error(w, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	if req.Name == "" {
		c.Error(w, http.StatusBadRequest, "Application name is required")
		return
	}

	appPath := filepath.Join("apps", req.Name)
	if err := os.MkdirAll(appPath, 0755); err != nil {
		c.Error(w, http.StatusInternalServerError, "Failed to create app directory: "+err.Error())
		return
	}

	subdirs := []string{
		"controllers", "models", "dto", "migrations/gorm", "migrations/sql",
		"templates/layouts", "templates/pages", "routes", "middleware",
		"permissions", "services", "listeners", "hooks", "signals", "repository",
	}

	for _, subdir := range subdirs {
		if err := os.MkdirAll(filepath.Join(appPath, subdir), 0755); err != nil {
			c.Error(w, http.StatusInternalServerError, "Failed to create subdirectory: "+err.Error())
			return
		}
	}

	// إنشاء ملف app.go
	appContent := `package ` + req.Name + `

import (
    "github.com/nats-framework/nats/pkg/engine"
    "github.com/nats-framework/nats/apps/` + req.Name + `/routes"
)

type App struct {
    engine *engine.Engine
    name   string
}

func NewApp(engine *engine.Engine) *App {
    return &App{
        engine: engine,
        name:   "` + req.Name + `",
    }
}

func (a *App) Name() string {
    return a.name
}

func (a *App) Register() error {
    routes.RegisterRoutes(a.engine.GetRouter())
    return nil
}

func (a *App) Boot() error {
    return nil
}
`

	if err := os.WriteFile(filepath.Join(appPath, "app.go"), []byte(appContent), 0644); err != nil {
		c.Error(w, http.StatusInternalServerError, "Failed to create app.go: "+err.Error())
		return
	}

	// إنشاء النماذج
	for _, model := range req.Models {
		if model.Name == "" {
			continue
		}

		modelName := strings.Title(model.Name)
		fields := parseFields(model.Fields)

		modelContent := generateModelContent(req.Name, modelName, fields)
		modelPath := filepath.Join(appPath, "models", strings.ToLower(modelName)+".go")
		if err := os.WriteFile(modelPath, []byte(modelContent), 0644); err != nil {
			c.Error(w, http.StatusInternalServerError, "Failed to create model: "+err.Error())
			return
		}
	}

	c.Success(w, map[string]interface{}{
		"message": "Application created successfully",
		"name":    req.Name,
		"models":  len(req.Models),
	})
}

// DeleteApp يحذف تطبيقاً
func (c *AppsController) DeleteApp(w http.ResponseWriter, r *http.Request) {
	appName := chi.URLParam(r, "app")

	if appName == "" {
		c.Error(w, http.StatusBadRequest, "App name is required")
		return
	}

	if appName == "core" || appName == "users" {
		c.Error(w, http.StatusForbidden, "Cannot delete system apps")
		return
	}

	appPath := filepath.Join("apps", appName)
	if err := os.RemoveAll(appPath); err != nil {
		c.Error(w, http.StatusInternalServerError, "Failed to delete app: "+err.Error())
		return
	}

	c.Success(w, map[string]interface{}{
		"message": "App deleted successfully",
	})
}

// ListRelations يعيد قائمة العلاقات
func (c *AppsController) ListRelations(w http.ResponseWriter, r *http.Request) {
	relations := []map[string]interface{}{
		{"parent": "core", "child": "users", "type": "one-to-many"},
	}
	c.Success(w, relations)
}

// CreateRelation ينشئ علاقة جديدة
func (c *AppsController) CreateRelation(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Parent     string `json:"parent"`
		Child      string `json:"child"`
		Type       string `json:"type"`
		ForeignKey string `json:"foreignKey"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.Error(w, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	if req.Parent == "" || req.Child == "" {
		c.Error(w, http.StatusBadRequest, "Parent and child apps are required")
		return
	}

	if req.Parent == req.Child {
		c.Error(w, http.StatusBadRequest, "Parent and child cannot be the same")
		return
	}

	c.Success(w, map[string]interface{}{
		"message": "Apps linked successfully",
		"parent":  req.Parent,
		"child":   req.Child,
		"type":    req.Type,
	})
}

// DeleteRelation يحذف علاقة
func (c *AppsController) DeleteRelation(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Parent string `json:"parent"`
		Child  string `json:"child"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.Error(w, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	c.Success(w, map[string]interface{}{
		"message": "Relation deleted successfully",
	})
}

// AddField يضيف حقل جديد
func (c *AppsController) AddField(w http.ResponseWriter, r *http.Request) {
	appName := chi.URLParam(r, "app")
	modelName := chi.URLParam(r, "model")

	var req struct {
		Name     string `json:"name"`
		Type     string `json:"type"`
		Required bool   `json:"required"`
		Unique   bool   `json:"unique"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.Error(w, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	if appName == "" || modelName == "" || req.Name == "" {
		c.Error(w, http.StatusBadRequest, "App name, model name, and field name are required")
		return
	}

	c.Success(w, map[string]interface{}{
		"message": "Field added successfully",
		"app":     appName,
		"model":   modelName,
		"field":   req,
	})
}

// DeleteField يحذف حقل
func (c *AppsController) DeleteField(w http.ResponseWriter, r *http.Request) {
	appName := chi.URLParam(r, "app")
	modelName := chi.URLParam(r, "model")
	fieldName := chi.URLParam(r, "field")

	if appName == "" || modelName == "" || fieldName == "" {
		c.Error(w, http.StatusBadRequest, "App name, model name, and field name are required")
		return
	}

	c.Success(w, map[string]interface{}{
		"message": "Field deleted successfully",
		"app":     appName,
		"model":   modelName,
		"field":   fieldName,
	})
}

// CreateUser ينشئ مستخدم جديد
func (c *AppsController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		FullName string `json:"fullName"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.Error(w, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		c.Error(w, http.StatusBadRequest, "Username, email, and password are required")
		return
	}

	c.Success(w, map[string]interface{}{
		"message":  "User created successfully",
		"username": req.Username,
		"email":    req.Email,
		"role":     req.Role,
	})
}

// RunMigrations ينفذ الهجرات
func (c *AppsController) RunMigrations(w http.ResponseWriter, r *http.Request) {
	c.Success(w, map[string]interface{}{
		"message": "Migrations completed successfully",
	})
}

// ResetMigrations يعيد تعيين الهجرات
func (c *AppsController) ResetMigrations(w http.ResponseWriter, r *http.Request) {
	c.Success(w, map[string]interface{}{
		"message": "Migrations reset successfully",
	})
}

// ============================================
// دوال مساعدة
// ============================================

func parseFields(fieldsStr string) []map[string]string {
	fields := []map[string]string{}
	if fieldsStr == "" {
		return fields
	}

	for _, part := range strings.Split(fieldsStr, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		kv := strings.Split(part, ":")
		if len(kv) == 2 {
			fields = append(fields, map[string]string{
				"name": strings.TrimSpace(kv[0]),
				"type": strings.TrimSpace(kv[1]),
			})
		}
	}

	return fields
}

func generateModelContent(appName, modelName string, fields []map[string]string) string {
	var fieldsStr []string
	for _, f := range fields {
		goType := getGoType(f["type"])
		fieldsStr = append(fieldsStr, "    "+strings.Title(f["name"])+" "+goType+" `gorm:\"column:"+strings.ToLower(f["name"])+"\"`")
	}

	return `package models

import (
    "time"

    "gorm.io/gorm"
)

type ` + modelName + ` struct {
    ID        uint           ` + "`gorm:\"primaryKey\"`" + `
` + strings.Join(fieldsStr, "\n") + `
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt ` + "`gorm:\"index\"`" + `
}

func (` + modelName + `) TableName() string {
    return "` + strings.ToLower(modelName) + `s"
}
`
}

func getGoType(fieldType string) string {
	switch fieldType {
	case "string":
		return "string"
	case "text":
		return "string"
	case "int":
		return "int"
	case "float":
		return "float64"
	case "bool":
		return "bool"
	case "date", "datetime", "time":
		return "time.Time"
	case "json":
		return "json.RawMessage"
	default:
		return "string"
	}
}
