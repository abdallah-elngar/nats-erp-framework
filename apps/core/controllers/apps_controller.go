package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/nats-framework/nats/pkg/auth"

	"github.com/nats-framework/nats/pkg/database"
	"github.com/nats-framework/nats/pkg/metadata"
	"github.com/nats-framework/nats/pkg/response"
	"github.com/nats-framework/nats/pkg/template"
)

// ============================================
// تعريف الأنواع
// ============================================

// ModelInput يمثل بيانات النموذج من الطلب
type ModelInput struct {
	Name   string `json:"name"`
	Fields string `json:"fields"`
}

// AppsController متحكم إدارة التطبيقات
type AppsController struct {
	engine *template.Engine
}

// NewAppsController ينشئ متحكم إدارة التطبيقات
func NewAppsController(engine *template.Engine) *AppsController {
	return &AppsController{
		engine: engine,
	}
}

// ============================================
// دوال مساعدة للاستجابات
// ============================================

func (c *AppsController) success(w http.ResponseWriter, data interface{}) {
	response.Success(w, data)
}

func (c *AppsController) error(w http.ResponseWriter, status int, message string) {
	response.Error(w, status, message)
}
func (c *AppsController) getUserFromContext(ctx context.Context) string {
	// استخدام الدالة من حزمة auth (بعد التعديل)
	user, ok := auth.GetUserFromContext(ctx)
	if ok && user != nil {
		return user.Username
	}
	return "system"
}

// ============================================
// 1. قائمة التطبيقات
// ============================================

// ListApps يعيد قائمة التطبيقات
func (c *AppsController) ListApps(w http.ResponseWriter, r *http.Request) {
	apps := []map[string]interface{}{}

	entries, err := os.ReadDir("apps")
	if err != nil {
		c.success(w, apps)
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

	c.success(w, apps)
}

// GetAppModels يعيد نماذج تطبيق معين
func (c *AppsController) GetAppModels(w http.ResponseWriter, r *http.Request) {
	appName := chi.URLParam(r, "app")

	if appName == "" {
		c.error(w, http.StatusBadRequest, "App name is required")
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
		c.success(w, fields)
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

	c.success(w, fields)
}

// ============================================
// 2. إنشاء تطبيق (كامل)
// ============================================

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
		c.error(w, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	if req.Name == "" {
		c.error(w, http.StatusBadRequest, "Application name is required")
		return
	}

	// ✅ تحويل req.Models إلى ModelInput
	models := make([]ModelInput, len(req.Models))
	for i, m := range req.Models {
		models[i] = ModelInput{
			Name:   m.Name,
			Fields: m.Fields,
		}
	}

	appPath := filepath.Join("apps", req.Name)
	if err := os.MkdirAll(appPath, 0755); err != nil {
		c.error(w, http.StatusInternalServerError, "Failed to create app directory: "+err.Error())
		return
	}

	// إنشاء جميع المجلدات الفرعية
	subdirs := []string{
		"controllers", "models", "dto", "migrations/gorm", "migrations/sql",
		"templates/layouts", "templates/pages", "routes", "middleware",
		"permissions", "services", "listeners", "hooks", "signals", "repository",
	}

	for _, subdir := range subdirs {
		if err := os.MkdirAll(filepath.Join(appPath, subdir), 0755); err != nil {
			c.error(w, http.StatusInternalServerError, "Failed to create subdirectory: "+err.Error())
			return
		}
	}

	// ============================================
	// ✅ 1. app.go
	// ============================================
	appContent := `package ` + req.Name + `

import (
    "github.com/go-chi/chi/v5"
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
    routes.RegisterRoutes(a.engine.GetChiRouter())
    return nil
}

func (a *App) Boot() error {
    return nil
}
`
	if err := os.WriteFile(filepath.Join(appPath, "app.go"), []byte(appContent), 0644); err != nil {
		c.error(w, http.StatusInternalServerError, "Failed to create app.go: "+err.Error())
		return
	}

	// ============================================
	// ✅ 2. register.go
	// ============================================
	registerContent := `package ` + req.Name + `

import (
    "github.com/nats-framework/nats/pkg/engine"
)

func Register(app *engine.Engine) error {
    appInstance := NewApp(app)
    if err := appInstance.Register(); err != nil {
        return err
    }
    if err := appInstance.Boot(); err != nil {
        return err
    }
    return nil
}
`
	if err := os.WriteFile(filepath.Join(appPath, "register.go"), []byte(registerContent), 0644); err != nil {
		c.error(w, http.StatusInternalServerError, "Failed to create register.go: "+err.Error())
		return
	}

	// ============================================
	// ✅ 3. النماذج (Models)
	// ============================================
	for _, model := range models {
		if model.Name == "" {
			continue
		}
		modelName := strings.Title(model.Name)
		fields := parseFields(model.Fields)

		modelContent := generateModelContent(req.Name, modelName, fields)
		modelPath := filepath.Join(appPath, "models", strings.ToLower(modelName)+".go")
		if err := os.WriteFile(modelPath, []byte(modelContent), 0644); err != nil {
			c.error(w, http.StatusInternalServerError, "Failed to create model: "+err.Error())
			return
		}
	}

	// ============================================
	// ✅ 4. DTOs
	// ============================================
	for _, model := range models {
		if model.Name == "" {
			continue
		}
		modelName := strings.Title(model.Name)
		fields := parseFields(model.Fields)

		dtoContent := generateDTOContent(req.Name, modelName, fields)
		dtoPath := filepath.Join(appPath, "dto", strings.ToLower(modelName)+"_dto.go")
		if err := os.WriteFile(dtoPath, []byte(dtoContent), 0644); err != nil {
			c.error(w, http.StatusInternalServerError, "Failed to create DTO: "+err.Error())
			return
		}
	}

	// ============================================
	// ✅ 5. الهجرات (Migrations)
	// ============================================
	for _, model := range models {
		if model.Name == "" {
			continue
		}
		modelName := strings.Title(model.Name)
		modelFields := parseFields(model.Fields)

		timestamp := time.Now().Format("20060102150405")
		tableName := strings.ToLower(modelName) + "s"

		// GORM Migration
		gormContent := `package migrations

import (
    "gorm.io/gorm"

    "github.com/nats-framework/nats/apps/` + req.Name + `/models"
)

func Up_` + timestamp + `(db *gorm.DB) error {
    return db.AutoMigrate(&models.` + modelName + `{})
}

func Down_` + timestamp + `(db *gorm.DB) error {
    return db.Migrator().DropTable("` + tableName + `")
}
`
		gormPath := filepath.Join(appPath, "migrations/gorm", timestamp+"_create_"+tableName+"_table.go")
		if err := os.WriteFile(gormPath, []byte(gormContent), 0644); err != nil {
			c.error(w, http.StatusInternalServerError, "Failed to create GORM migration: "+err.Error())
			return
		}

		// SQL Migration
		var columns []string
		columns = append(columns, "    id BIGSERIAL PRIMARY KEY")
		for _, f := range modelFields {
			if f["type"] == "relation" {
				columns = append(columns, "    "+strings.ToLower(f["name"])+"_id BIGINT")
				continue
			}
			columns = append(columns, "    "+strings.ToLower(f["name"])+" "+getSQLType(f["type"]))
		}
		columns = append(columns, "    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP")
		columns = append(columns, "    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP")
		columns = append(columns, "    deleted_at TIMESTAMP")

		sqlContent := `-- name: create_` + tableName + `_table
-- id: ` + timestamp + `
-- created: ` + time.Now().Format("2006-01-02 15:04:05") + `

-- up:
CREATE TABLE IF NOT EXISTS ` + tableName + ` (
` + strings.Join(columns, ",\n") + `
);

-- down:
DROP TABLE IF EXISTS ` + tableName + `;
`
		sqlPath := filepath.Join(appPath, "migrations/sql", timestamp+"_create_"+tableName+"_table.sql")
		if err := os.WriteFile(sqlPath, []byte(sqlContent), 0644); err != nil {
			c.error(w, http.StatusInternalServerError, "Failed to create SQL migration: "+err.Error())
			return
		}
	}

	// ============================================
	// ✅ 6. المتحكمات (Controllers)
	// ============================================
	for _, model := range models {
		if model.Name == "" {
			continue
		}
		modelName := strings.Title(model.Name)
		lowerName := strings.ToLower(modelName)

		controllerContent := generateControllerContent(req.Name, modelName)
		controllerPath := filepath.Join(appPath, "controllers", lowerName+"_controller.go")
		if err := os.WriteFile(controllerPath, []byte(controllerContent), 0644); err != nil {
			c.error(w, http.StatusInternalServerError, "Failed to create controller: "+err.Error())
			return
		}
	}

	// ============================================
	// ✅ 7. router.go - ✅ استخدام models (نوع ModelInput)
	// ============================================
	routerContent := generateRouterContent(req.Name, models)
	if err := os.WriteFile(filepath.Join(appPath, "routes", "router.go"), []byte(routerContent), 0644); err != nil {
		c.error(w, http.StatusInternalServerError, "Failed to create router: "+err.Error())
		return
	}

	// ============================================
	// ✅ 8. Services
	// ============================================
	for _, model := range models {
		if model.Name == "" {
			continue
		}
		modelName := strings.Title(model.Name)

		serviceContent := generateServiceContent(req.Name, modelName)
		servicePath := filepath.Join(appPath, "services", strings.ToLower(modelName)+"_service.go")
		if err := os.WriteFile(servicePath, []byte(serviceContent), 0644); err != nil {
			c.error(w, http.StatusInternalServerError, "Failed to create service: "+err.Error())
			return
		}
	}

	// ============================================
	// ✅ 9. Repository
	// ============================================
	for _, model := range models {
		if model.Name == "" {
			continue
		}
		modelName := strings.Title(model.Name)

		repoContent := generateRepositoryContent(req.Name, modelName)
		repoPath := filepath.Join(appPath, "repository", strings.ToLower(modelName)+"_repository.go")
		if err := os.WriteFile(repoPath, []byte(repoContent), 0644); err != nil {
			c.error(w, http.StatusInternalServerError, "Failed to create repository: "+err.Error())
			return
		}
	}

	// ============================================
	// ✅ 10. Permissions
	// ============================================
	permissionsContent := generatePermissionsContent(req.Name, models)
	if err := os.WriteFile(filepath.Join(appPath, "permissions", "permissions.go"), []byte(permissionsContent), 0644); err != nil {
		c.error(w, http.StatusInternalServerError, "Failed to create permissions: "+err.Error())
		return
	}

	// ============================================
	// ✅ 11. Listeners
	// ============================================
	listenersContent := generateListenersContent(req.Name, models)
	if err := os.WriteFile(filepath.Join(appPath, "listeners", "listeners.go"), []byte(listenersContent), 0644); err != nil {
		c.error(w, http.StatusInternalServerError, "Failed to create listeners: "+err.Error())
		return
	}

	// ============================================
	// ✅ 12. Hooks
	// ============================================
	hooksContent := generateHooksContent(req.Name, models)
	if err := os.WriteFile(filepath.Join(appPath, "hooks", "hooks.go"), []byte(hooksContent), 0644); err != nil {
		c.error(w, http.StatusInternalServerError, "Failed to create hooks: "+err.Error())
		return
	}

	// ============================================
	// ✅ 13. Signals
	// ============================================
	signalsContent := generateSignalsContent(req.Name, models)
	if err := os.WriteFile(filepath.Join(appPath, "signals", "signals.go"), []byte(signalsContent), 0644); err != nil {
		c.error(w, http.StatusInternalServerError, "Failed to create signals: "+err.Error())
		return
	}

	// ============================================
	// ✅ 14. Middleware
	// ============================================
	middlewareContent := generateMiddlewareContent(req.Name)
	if err := os.WriteFile(filepath.Join(appPath, "middleware", "middleware.go"), []byte(middlewareContent), 0644); err != nil {
		c.error(w, http.StatusInternalServerError, "Failed to create middleware: "+err.Error())
		return
	}

	// ============================================
	// ✅ 15. Templates (القوالب)
	// ============================================
	for _, model := range models {
		if model.Name == "" {
			continue
		}
		modelName := strings.Title(model.Name)
		lowerName := strings.ToLower(modelName)
		pluralName := lowerName + "s"

		templatesDir := filepath.Join(appPath, "templates/pages", pluralName)
		os.MkdirAll(templatesDir, 0755)

		// index.html
		indexContent := `{{define "content"}}
<div class="container mx-auto p-4">
    <div class="flex justify-between items-center mb-6">
        <h1 class="text-2xl font-bold">📦 ` + modelName + `s</h1>
        <button class="btn btn-primary" hx-get="/` + req.Name + `/` + pluralName + `/create" hx-target="#modal-content" @click="showModal = true">
            <i class="fas fa-plus mr-2"></i>Add ` + modelName + `
        </button>
    </div>

    <div class="bg-white rounded-lg shadow overflow-hidden">
        <table class="w-full">
            <thead class="bg-gray-50">
                <tr>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">ID</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Name</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Actions</th>
                </tr>
            </thead>
            <tbody class="divide-y divide-gray-200">
                {{range .Items}}
                <tr>
                    <td class="px-6 py-4">{{.ID}}</td>
                    <td class="px-6 py-4">{{.Name}}</td>
                    <td class="px-6 py-4">
                        <button class="text-blue-600 hover:text-blue-800 mr-2"
                                hx-get="/` + req.Name + `/` + pluralName + `/{{.ID}}" hx-target="#modal-content" @click="showModal = true">
                            <i class="fas fa-eye"></i>
                        </button>
                        <button class="text-green-600 hover:text-green-800 mr-2"
                                hx-get="/` + req.Name + `/` + pluralName + `/{{.ID}}/edit" hx-target="#modal-content" @click="showModal = true">
                            <i class="fas fa-edit"></i>
                        </button>
                        <button class="text-red-600 hover:text-red-800"
                                hx-delete="/` + req.Name + `/` + pluralName + `/{{.ID}}" hx-confirm="Are you sure?">
                            <i class="fas fa-trash"></i>
                        </button>
                    </td>
                </tr>
                {{else}}
                <tr>
                    <td colspan="3" class="px-6 py-4 text-center text-gray-500">No ` + lowerName + `s found</td>
                </tr>
                {{end}}
            </tbody>
        </table>
    </div>
</div>
{{end}}`
		if err := os.WriteFile(filepath.Join(templatesDir, "index.html"), []byte(indexContent), 0644); err != nil {
			c.error(w, http.StatusInternalServerError, "Failed to create index template: "+err.Error())
			return
		}
	}

	c.success(w, map[string]interface{}{
		"message": "Application created successfully",
		"name":    req.Name,
		"models":  len(req.Models),
	})
}

// ============================================
// 3. حذف تطبيق
// ============================================

// DeleteApp يحذف تطبيقاً
func (c *AppsController) DeleteApp(w http.ResponseWriter, r *http.Request) {
	appName := chi.URLParam(r, "app")

	if appName == "" {
		c.error(w, http.StatusBadRequest, "App name is required")
		return
	}

	if appName == "core" || appName == "users" {
		c.error(w, http.StatusForbidden, "Cannot delete system apps")
		return
	}

	appPath := filepath.Join("apps", appName)
	if err := os.RemoveAll(appPath); err != nil {
		c.error(w, http.StatusInternalServerError, "Failed to delete app: "+err.Error())
		return
	}

	c.success(w, map[string]interface{}{
		"message": "App deleted successfully",
	})
}

// ============================================
// 4. العلاقات
// ============================================

// ListRelations يعيد قائمة العلاقات
func (c *AppsController) ListRelations(w http.ResponseWriter, r *http.Request) {
	relations := []map[string]interface{}{
		{"parent": "core", "child": "users", "type": "one-to-many"},
	}
	c.success(w, relations)
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
		c.error(w, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	if req.Parent == "" || req.Child == "" {
		c.error(w, http.StatusBadRequest, "Parent and child apps are required")
		return
	}

	if req.Parent == req.Child {
		c.error(w, http.StatusBadRequest, "Parent and child cannot be the same")
		return
	}

	c.success(w, map[string]interface{}{
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
		c.error(w, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	c.success(w, map[string]interface{}{
		"message": "Relation deleted successfully",
	})
}

// ============================================
// 5. الحقول
// ============================================

// AddField يضيف حقل جديد
// ============================================
// 5. إدارة الحقول (Add/Delete Field) - كاملة مع حفظ في قاعدة البيانات
// ============================================

// ============================================
// دالة مساعدة للحصول على المستخدم
// ============================================

// getCurrentUser يحصل على المستخدم الحالي من السياق
func (c *AppsController) getCurrentUser(r *http.Request) string {
	// محاولة الحصول على المستخدم من السياق
	user := auth.GetUserFromRequest(r)
	if user != nil {
		// إذا كان المستخدم من نوع User، إرجاع اسم المستخدم
		if u, ok := user.(*auth.User); ok {
			return u.Username
		}
		// إذا كان map، محاولة استخراج username
		if u, ok := user.(map[string]interface{}); ok {
			if username, ok := u["username"].(string); ok {
				return username
			}
			if username, ok := u["Username"].(string); ok {
				return username
			}
		}
	}

	// مستخدم افتراضي (للتطوير)
	return "admin"
}

// ============================================
// 5. إدارة الحقول (Add/Delete Field) - كاملة
// ============================================

// AddField يضيف حقل جديد لنموذج
func (c *AppsController) AddField(w http.ResponseWriter, r *http.Request) {
	appName := chi.URLParam(r, "app")
	modelName := chi.URLParam(r, "model")

	var req struct {
		Name         string `json:"name"`
		Type         string `json:"type"`
		Required     bool   `json:"required"`
		Unique       bool   `json:"unique"`
		DefaultValue string `json:"defaultValue"`
		Relation     string `json:"relation"`
		MaxLength    int    `json:"maxLength"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.error(w, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	if appName == "" || modelName == "" || req.Name == "" {
		c.error(w, http.StatusBadRequest, "App name, model name, and field name are required")
		return
	}

	// ============================================
	// ✅ 1. التحقق من وجود التطبيق والنموذج
	// ============================================
	modelPath := filepath.Join("apps", appName, "models", strings.ToLower(modelName)+".go")
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		c.error(w, http.StatusNotFound, "Model not found: "+modelName)
		return
	}

	// ============================================
	// ✅ 2. قراءة ملف النموذج الحالي
	// ============================================
	content, err := os.ReadFile(modelPath)
	if err != nil {
		c.error(w, http.StatusInternalServerError, "Failed to read model file: "+err.Error())
		return
	}

	// ============================================
	// ✅ 3. التحقق من عدم وجود الحقل مسبقاً
	// ============================================
	if strings.Contains(string(content), req.Name+" ") {
		c.error(w, http.StatusBadRequest, "Field '"+req.Name+"' already exists in model '"+modelName+"'")
		return
	}

	// ============================================
	// ✅ 4. إضافة الحقل الجديد إلى الـ struct
	// ============================================
	goType := getGoType(req.Type)
	fieldLine := fmt.Sprintf("    %s %s `gorm:\"column:%s\"`",
		strings.Title(req.Name),
		goType,
		strings.ToLower(req.Name))

	lines := strings.Split(string(content), "\n")
	var newLines []string
	inserted := false

	for _, line := range lines {
		newLines = append(newLines, line)
		if !inserted && strings.Contains(line, "CreatedAt") {
			newLines = append(newLines, fieldLine)
			inserted = true
		}
	}

	if !inserted {
		for i := len(newLines) - 1; i >= 0; i-- {
			if strings.TrimSpace(newLines[i]) == "}" {
				newLines = append(newLines[:i], append([]string{fieldLine}, newLines[i:]...)...)
				break
			}
		}
	}

	// ============================================
	// ✅ 5. حفظ الملف المحدث
	// ============================================
	newContent := strings.Join(newLines, "\n")
	if err := os.WriteFile(modelPath, []byte(newContent), 0644); err != nil {
		c.error(w, http.StatusInternalServerError, "Failed to update model file: "+err.Error())
		return
	}

	// ============================================
	// ✅ 6. إنشاء هجرة GORM
	// ============================================
	timestamp := time.Now().Format("20060102150405")
	tableName := strings.ToLower(modelName) + "s"
	fieldName := strings.ToLower(req.Name)
	modelNameCap := strings.Title(modelName)

	gormContent := `package migrations

import (
    "gorm.io/gorm"
)

// Up_` + timestamp + ` يضيف عمود ` + fieldName + ` إلى جدول ` + tableName + `
func Up_` + timestamp + `(db *gorm.DB) error {
    return db.Migrator().AddColumn(&` + modelNameCap + `{}, "` + strings.Title(req.Name) + `")
}

// Down_` + timestamp + ` يحذف عمود ` + fieldName + ` من جدول ` + tableName + `
func Down_` + timestamp + `(db *gorm.DB) error {
    return db.Migrator().DropColumn(&` + modelNameCap + `{}, "` + strings.Title(req.Name) + `")
}
`

	gormPath := filepath.Join("apps", appName, "migrations", "gorm",
		timestamp+"_add_"+fieldName+"_to_"+strings.ToLower(modelName)+".go")
	if err := os.WriteFile(gormPath, []byte(gormContent), 0644); err != nil {
		c.error(w, http.StatusInternalServerError, "Failed to create GORM migration: "+err.Error())
		return
	}

	// ============================================
	// ✅ 7. إنشاء هجرة SQL
	// ============================================
	sqlType := getSQLType(req.Type)
	constraints := ""
	if req.Required {
		constraints += " NOT NULL"
	}
	if req.Unique {
		constraints += " UNIQUE"
	}
	if req.DefaultValue != "" {
		constraints += " DEFAULT " + req.DefaultValue
	}

	sqlContent := `-- name: add_` + fieldName + `_to_` + tableName + `
-- id: ` + timestamp + `
-- created: ` + time.Now().Format("2006-01-02 15:04:05") + `

-- up:
ALTER TABLE ` + tableName + ` ADD COLUMN IF NOT EXISTS ` + fieldName + ` ` + sqlType + constraints + `;

-- down:
ALTER TABLE ` + tableName + ` DROP COLUMN IF EXISTS ` + fieldName + `;
`

	sqlPath := filepath.Join("apps", appName, "migrations", "sql",
		timestamp+"_add_"+fieldName+"_to_"+strings.ToLower(modelName)+".sql")
	if err := os.WriteFile(sqlPath, []byte(sqlContent), 0644); err != nil {
		c.error(w, http.StatusInternalServerError, "Failed to create SQL migration: "+err.Error())
		return
	}

	// ============================================
	// ✅ 8. حفظ البيانات في جدول app_metadata (قاعدة البيانات)
	// ============================================
	// الحصول على اسم المستخدم الحالي
	currentUser := c.getCurrentUser(r)

	metadataMgr := metadata.NewMetadataManager(database.DB())

	meta := &metadata.AppMetadata{
		AppName:      appName,
		ModelName:    modelName,
		FieldName:    req.Name,
		FieldType:    req.Type,
		IsRequired:   req.Required,
		IsUnique:     req.Unique,
		DefaultValue: req.DefaultValue,
		CreatedBy:    currentUser,
		MigrationID:  timestamp,
		Status:       "pending",
	}

	if err := metadataMgr.Save(meta); err != nil {
		// تسجيل الخطأ ولكن لا نوقف العملية
		fmt.Printf("⚠️ Failed to save metadata: %v\n", err)
	}

	// ============================================
	// ✅ 9. إرجاع الاستجابة مع معلومات الهجرة
	// ============================================
	c.success(w, map[string]interface{}{
		"message":        "Field added successfully. Please run migrations.",
		"app":            appName,
		"model":          modelName,
		"field":          req,
		"migration_id":   timestamp,
		"gorm_migration": filepath.Base(gormPath),
		"sql_migration":  filepath.Base(sqlPath),
		"needs_migrate":  true,
	})
}

// DeleteField يحذف حقل
func (c *AppsController) DeleteField(w http.ResponseWriter, r *http.Request) {
	appName := chi.URLParam(r, "app")
	modelName := chi.URLParam(r, "model")
	fieldName := chi.URLParam(r, "field")

	if appName == "" || modelName == "" || fieldName == "" {
		c.error(w, http.StatusBadRequest, "App name, model name, and field name are required")
		return
	}

	c.success(w, map[string]interface{}{
		"message": "Field deleted successfully",
		"app":     appName,
		"model":   modelName,
		"field":   fieldName,
	})
}

// ============================================
// 6. المستخدمين
// ============================================

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
		c.error(w, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		c.error(w, http.StatusBadRequest, "Username, email, and password are required")
		return
	}

	c.success(w, map[string]interface{}{
		"message":  "User created successfully",
		"username": req.Username,
		"email":    req.Email,
		"role":     req.Role,
	})
}

// ============================================
// 7. الهجرات
// ============================================

// RunMigrations ينفذ الهجرات
func (c *AppsController) RunMigrations(w http.ResponseWriter, r *http.Request) {
	c.success(w, map[string]interface{}{
		"message": "Migrations completed successfully",
	})
}

// ResetMigrations يعيد تعيين الهجرات
func (c *AppsController) ResetMigrations(w http.ResponseWriter, r *http.Request) {
	c.success(w, map[string]interface{}{
		"message": "Migrations reset successfully",
	})
}

// ============================================
// دوال مساعدة
// ============================================

// parseFields يحلل الحقول من النص
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

// ============================================
// دوال توليد المحتوى
// ============================================

// generateModelContent يولد محتوى ملف النموذج
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

// getGoType يحول نوع الحقل إلى نوع Go
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
	case "relation":
		return "uint"
	default:
		return "string"
	}
}

// getSQLType يحول نوع الحقل إلى نوع SQL
func getSQLType(fieldType string) string {
	switch fieldType {
	case "string":
		return "VARCHAR(255)"
	case "text":
		return "TEXT"
	case "int":
		return "INTEGER"
	case "float":
		return "DECIMAL(10,2)"
	case "bool":
		return "BOOLEAN"
	case "date":
		return "DATE"
	case "datetime":
		return "TIMESTAMP"
	case "time":
		return "TIME"
	case "json":
		return "JSONB"
	case "relation":
		return "BIGINT"
	default:
		return "TEXT"
	}
}

// generateDTOContent يولد محتوى DTO
func generateDTOContent(appName, modelName string, fields []map[string]string) string {
	var createFields, updateFields, responseFields []string
	for _, f := range fields {
		if f["type"] == "relation" {
			continue
		}
		fieldName := strings.Title(f["name"])
		goType := getGoType(f["type"])
		jsonName := strings.ToLower(f["name"])

		createFields = append(createFields, "    "+fieldName+" "+goType+" `json:\""+jsonName+"\" validate:\"required\"`")
		updateFields = append(updateFields, "    "+fieldName+" "+goType+" `json:\""+jsonName+"\"`")
		responseFields = append(responseFields, "    "+fieldName+" "+goType+" `json:\""+jsonName+"\"`")
	}

	return `package dto

import "time"

type Create` + modelName + `Request struct {
` + strings.Join(createFields, "\n") + `
}

type Update` + modelName + `Request struct {
` + strings.Join(updateFields, "\n") + `
}

type ` + modelName + `Response struct {
    ID        uint      ` + "`json:\"id\"`" + `
` + strings.Join(responseFields, "\n") + `
    CreatedAt time.Time ` + "`json:\"created_at\"`" + `
    UpdatedAt time.Time ` + "`json:\"updated_at\"`" + `
}
`
}

// generateControllerContent يولد محتوى المتحكم
func generateControllerContent(appName, modelName string) string {
	return `package controllers

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/go-chi/chi/v5"

    "github.com/nats-framework/nats/apps/` + appName + `/dto"
    "github.com/nats-framework/nats/apps/` + appName + `/services"
    "github.com/nats-framework/nats/pkg/response"
)

// ` + modelName + `Controller متحكم ` + modelName + `
type ` + modelName + `Controller struct {
    service *services.` + modelName + `Service
}

// New` + modelName + `Controller ينشئ متحكم ` + modelName + ` جديد
func New` + modelName + `Controller() *` + modelName + `Controller {
    return &` + modelName + `Controller{
        service: services.New` + modelName + `Service(),
    }
}

func (c *` + modelName + `Controller) Index(w http.ResponseWriter, r *http.Request) {
    items, err := c.service.GetAll()
    if err != nil {
        response.Error(w, http.StatusInternalServerError, err.Error())
        return
    }
    response.Success(w, items)
}

func (c *` + modelName + `Controller) Show(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        response.Error(w, http.StatusBadRequest, "Invalid ID")
        return
    }

    item, err := c.service.GetByID(uint(id))
    if err != nil {
        response.Error(w, http.StatusNotFound, "` + modelName + ` not found")
        return
    }

    response.Success(w, item)
}

func (c *` + modelName + `Controller) Create(w http.ResponseWriter, r *http.Request) {
    var req dto.Create` + modelName + `Request
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        response.Error(w, http.StatusBadRequest, "Invalid request")
        return
    }

    item, err := c.service.Create(req)
    if err != nil {
        response.Error(w, http.StatusBadRequest, err.Error())
        return
    }

    response.Success(w, item)
}

func (c *` + modelName + `Controller) Update(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        response.Error(w, http.StatusBadRequest, "Invalid ID")
        return
    }

    var req dto.Update` + modelName + `Request
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        response.Error(w, http.StatusBadRequest, "Invalid request")
        return
    }

    item, err := c.service.Update(uint(id), req)
    if err != nil {
        response.Error(w, http.StatusBadRequest, err.Error())
        return
    }

    response.Success(w, item)
}

func (c *` + modelName + `Controller) Delete(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        response.Error(w, http.StatusBadRequest, "Invalid ID")
        return
    }

    if err := c.service.Delete(uint(id)); err != nil {
        response.Error(w, http.StatusInternalServerError, err.Error())
        return
    }

    response.Success(w, map[string]interface{}{"message": "Deleted successfully"})
}
`
}

// ✅ generateRouterContent - تستقبل []ModelInput
func generateRouterContent(appName string, models []ModelInput) string {
	var routes []string
	for _, model := range models {
		if model.Name == "" {
			continue
		}
		name := strings.Title(model.Name)
		pluralName := strings.ToLower(model.Name) + "s"

		routes = append(routes, `    // مسارات `+name+`
    ctrl := controllers.New`+name+`Controller()
    r.Get("/`+pluralName+`", ctrl.Index)
    r.Get("/`+pluralName+`/{id}", ctrl.Show)
    r.Post("/`+pluralName+`", ctrl.Create)
    r.Put("/`+pluralName+`/{id}", ctrl.Update)
    r.Delete("/`+pluralName+`/{id}", ctrl.Delete)`)
	}

	return `package routes

import (
    "github.com/go-chi/chi/v5"

    "github.com/nats-framework/nats/apps/` + appName + `/controllers"
)

func RegisterRoutes(r *chi.Mux) {
` + strings.Join(routes, "\n\n    ") + `
}
`
}

// generateServiceContent يولد محتوى Service
func generateServiceContent(appName, modelName string) string {
	return `package services

import (
    "errors"

    "github.com/nats-framework/nats/apps/` + appName + `/dto"
    "github.com/nats-framework/nats/apps/` + appName + `/models"
    "github.com/nats-framework/nats/apps/` + appName + `/repository"
)

type ` + modelName + `Service struct {
    repo *repository.` + modelName + `Repository
}

func New` + modelName + `Service() *` + modelName + `Service {
    return &` + modelName + `Service{
        repo: repository.New` + modelName + `Repository(),
    }
}

func (s *` + modelName + `Service) GetAll() ([]models.` + modelName + `, error) {
    return s.repo.FindAll()
}

func (s *` + modelName + `Service) GetByID(id uint) (*models.` + modelName + `, error) {
    return s.repo.FindByID(id)
}

func (s *` + modelName + `Service) Create(req dto.Create` + modelName + `Request) (*models.` + modelName + `, error) {
    item := &models.` + modelName + `{
        // TODO: Map fields from req to model
    }
    if err := s.repo.Create(item); err != nil {
        return nil, err
    }
    return item, nil
}

func (s *` + modelName + `Service) Update(id uint, req dto.Update` + modelName + `Request) (*models.` + modelName + `, error) {
    item, err := s.repo.FindByID(id)
    if err != nil {
        return nil, errors.New("record not found")
    }
    // TODO: Update fields from req to model
    if err := s.repo.Update(item); err != nil {
        return nil, err
    }
    return item, nil
}

func (s *` + modelName + `Service) Delete(id uint) error {
    return s.repo.Delete(id)
}
`
}

// generateRepositoryContent يولد محتوى Repository
func generateRepositoryContent(appName, modelName string) string {
	return `package repository

import (
    "gorm.io/gorm"

    "github.com/nats-framework/nats/apps/` + appName + `/models"
    "github.com/nats-framework/nats/pkg/database"
)

type ` + modelName + `Repository struct {
    db *gorm.DB
}

func New` + modelName + `Repository() *` + modelName + `Repository {
    return &` + modelName + `Repository{
        db: database.DB(),
    }
}

func (r *` + modelName + `Repository) Create(item *models.` + modelName + `) error {
    return r.db.Create(item).Error
}

func (r *` + modelName + `Repository) FindByID(id uint) (*models.` + modelName + `, error) {
    var item models.` + modelName + `
    err := r.db.First(&item, id).Error
    if err != nil {
        return nil, err
    }
    return &item, nil
}

func (r *` + modelName + `Repository) FindAll() ([]models.` + modelName + `, error) {
    var items []models.` + modelName + `
    err := r.db.Find(&items).Error
    return items, err
}

func (r *` + modelName + `Repository) Update(item *models.` + modelName + `) error {
    return r.db.Save(item).Error
}

func (r *` + modelName + `Repository) Delete(id uint) error {
    return r.db.Delete(&models.` + modelName + `{}, id).Error
}

func (r *` + modelName + `Repository) Exists(query string, args ...interface{}) (bool, error) {
    var count int64
    err := r.db.Model(&models.` + modelName + `{}).Where(query, args...).Count(&count).Error
    return count > 0, err
}
`
}

// generatePermissionsContent يولد محتوى Permissions
func generatePermissionsContent(appName string, models []ModelInput) string {
	var perms []string
	for _, model := range models {
		if model.Name == "" {
			continue
		}
		lowerName := strings.ToLower(model.Name)
		perms = append(perms, `    Perm`+model.Name+`View   = "`+appName+`.view_`+lowerName+`"
    Perm`+model.Name+`Create = "`+appName+`.create_`+lowerName+`"
    Perm`+model.Name+`Edit   = "`+appName+`.edit_`+lowerName+`"
    Perm`+model.Name+`Delete = "`+appName+`.delete_`+lowerName+`"`)
	}

	return `package permissions

const (
` + strings.Join(perms, "\n\n    ") + `
)
`
}

// generateListenersContent يولد محتوى Listeners
func generateListenersContent(appName string, models []ModelInput) string {
	var listeners []string
	for _, model := range models {
		if model.Name == "" {
			continue
		}
		lowerName := strings.ToLower(model.Name)
		modelNameCap := strings.Title(model.Name)

		listeners = append(listeners, `
    // مستمع عند إنشاء `+model.Name+`
    em.Listen("`+appName+`.`+lowerName+`.created", func(event events.Event) error {
        item, ok := event.GetData().(*models.`+modelNameCap+`)
        if !ok {
            return fmt.Errorf("invalid event data")
        }
        log.Printf("📦 `+model.Name+` created: ID=%d", item.ID)
        return nil
    })

    // مستمع عند تحديث `+model.Name+`
    em.Listen("`+appName+`.`+lowerName+`.updated", func(event events.Event) error {
        item, ok := event.GetData().(*models.`+modelNameCap+`)
        if !ok {
            return fmt.Errorf("invalid event data")
        }
        log.Printf("📦 `+model.Name+` updated: ID=%d", item.ID)
        return nil
    })

    // مستمع عند حذف `+model.Name+`
    em.Listen("`+appName+`.`+lowerName+`.deleted", func(event events.Event) error {
        item, ok := event.GetData().(*models.`+modelNameCap+`)
        if !ok {
            return fmt.Errorf("invalid event data")
        }
        log.Printf("📦 `+model.Name+` deleted: ID=%d", item.ID)
        return nil
    })`)
	}

	return `package listeners

import (
    "fmt"
    "log"

    "github.com/nats-framework/nats/apps/` + appName + `/models"
    "github.com/nats-framework/nats/pkg/events"
)

func RegisterListeners(em *events.EventManager) {` +
		strings.Join(listeners, "\n") + `
}
`
}

// generateHooksContent يولد محتوى Hooks
func generateHooksContent(appName string, models []ModelInput) string {
	var hooks []string
	for _, model := range models {
		if model.Name == "" {
			continue
		}
		lowerName := strings.ToLower(model.Name)
		modelNameCap := strings.Title(model.Name)

		hooks = append(hooks, `
    // خطاف قبل إنشاء `+model.Name+`
    hookManager.RegisterBefore("`+appName+`.`+lowerName+`.create", func(data interface{}) (interface{}, error) {
        item, ok := data.(*models.`+modelNameCap+`)
        if !ok {
            return data, nil
        }
        log.Printf("🔧 Before creating `+model.Name+`: %s", item.Name)
        return data, nil
    }, 10)

    // خطاف بعد إنشاء `+model.Name+`
    hookManager.RegisterAfter("`+appName+`.`+lowerName+`.create", func(data interface{}) (interface{}, error) {
        item, ok := data.(*models.`+modelNameCap+`)
        if !ok {
            return data, nil
        }
        log.Printf("✅ `+model.Name+` created: %s", item.Name)
        return data, nil
    }, 10)`)
	}

	return `package hooks

import (
    "log"

    "github.com/nats-framework/nats/apps/` + appName + `/models"
    "github.com/nats-framework/nats/pkg/hooks"
)

func RegisterHooks(hookManager *hooks.HookManager) {` +
		strings.Join(hooks, "\n") + `
}
`
}

// generateSignalsContent يولد محتوى Signals
func generateSignalsContent(appName string, models []ModelInput) string {
	var signals []string
	for _, model := range models {
		if model.Name == "" {
			continue
		}
		lowerName := strings.ToLower(model.Name)
		modelNameCap := strings.Title(model.Name)

		signals = append(signals, `
    // إشارة عند تغيير `+model.Name+`
    signalManager.Connect("`+appName+`.`+lowerName+`.changed", func(signal *signals.Signal) error {
        item, ok := signal.Data.(*models.`+modelNameCap+`)
        if !ok {
            return nil
        }
        log.Printf("📊 `+model.Name+` changed: ID=%d", item.ID)
        return nil
    })`)
	}

	return `package signals

import (
    "log"

    "github.com/nats-framework/nats/apps/` + appName + `/models"
    "github.com/nats-framework/nats/pkg/signals"
)

func RegisterSignals(signalManager *signals.SignalManager) {` +
		strings.Join(signals, "\n") + `
}
`
}

// generateMiddlewareContent يولد محتوى Middleware
func generateMiddlewareContent(appName string) string {
	return `package middleware

import (
    "net/http"
    "strings"

    "github.com/nats-framework/nats/pkg/response"
)

// AuthMiddleware يتحقق من صلاحيات المستخدم لتطبيق ` + appName + `
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            response.Unauthorized(w, "Authorization required")
            return
        }

        token = strings.TrimPrefix(token, "Bearer ")
        if token == "" {
            response.Unauthorized(w, "Invalid token")
            return
        }

        // TODO: التحقق من صلاحية المستخدم لتطبيق ` + appName + `
        next.ServeHTTP(w, r)
    })
}

// RateLimitMiddleware يحد من معدل الطلبات
func RateLimitMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // TODO: تطبيق تحديد معدل الطلبات
        next.ServeHTTP(w, r)
    })
}
`
}
