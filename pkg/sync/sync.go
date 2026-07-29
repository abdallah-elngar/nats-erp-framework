package sync

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/nats-framework/nats/pkg/metadata"
	"gorm.io/gorm"
)

// SyncManager مدير المزامنة بين الملفات وقاعدة البيانات
type SyncManager struct {
	db          *gorm.DB
	metadataMgr *metadata.MetadataManager
}

// NewSyncManager ينشئ مدير مزامنة جديد
func NewSyncManager(db *gorm.DB) *SyncManager {
	return &SyncManager{
		db:          db,
		metadataMgr: metadata.NewMetadataManager(db),
	}
}

// ============================================
// المزامنة من قاعدة البيانات إلى الملفات
// ============================================

// SyncFromDB يستعيد التغييرات من قاعدة البيانات إلى الملفات
func (s *SyncManager) SyncFromDB(appName string) error {
	metadataList, err := s.metadataMgr.GetByApp(appName)
	if err != nil {
		return err
	}

	for _, meta := range metadataList {
		if meta.Status == "applied" {
			continue
		}

		// تحديث ملف النموذج
		modelPath := filepath.Join("apps", meta.AppName, "models", strings.ToLower(meta.ModelName)+".go")
		if err := s.updateModelFile(modelPath, meta); err != nil {
			return err
		}

		// ✅ استخدام GetDB() بدلاً من db مباشرة
		if err := s.metadataMgr.UpdateStatus(meta.MigrationID, "applied"); err != nil {
			return err
		}
	}

	return nil
}

// updateModelFile يحدث ملف النموذج بإضافة حقل جديد
func (s *SyncManager) updateModelFile(modelPath string, meta metadata.AppMetadata) error {
	// قراءة الملف الحالي
	content, err := os.ReadFile(modelPath)
	if err != nil {
		return err
	}

	// التحقق من وجود الحقل مسبقاً
	if strings.Contains(string(content), meta.FieldName+" ") {
		return nil // الحقل موجود بالفعل
	}

	// إضافة الحقل الجديد
	goType := s.getGoType(meta.FieldType)
	fieldLine := fmt.Sprintf("    %s %s `gorm:\"column:%s\"`",
		strings.Title(meta.FieldName),
		goType,
		strings.ToLower(meta.FieldName))

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

	// حفظ الملف
	newContent := strings.Join(newLines, "\n")
	return os.WriteFile(modelPath, []byte(newContent), 0644)
}

// getGoType يحول نوع الحقل إلى نوع Go
func (s *SyncManager) getGoType(fieldType string) string {
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

// ============================================
// المزامنة من الملفات إلى قاعدة البيانات
// ============================================

// SyncToDB يزامن التغييرات من الملفات إلى قاعدة البيانات
func (s *SyncManager) SyncToDB(appName string) error {
	// قراءة جميع النماذج من الملفات
	modelsDir := filepath.Join("apps", appName, "models")
	files, err := os.ReadDir(modelsDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".go") {
			continue
		}

		modelName := strings.TrimSuffix(file.Name(), ".go")
		modelName = strings.Title(modelName)

		// استخراج الحقول من ملف النموذج
		fields, err := s.extractFieldsFromModel(filepath.Join(modelsDir, file.Name()), modelName)
		if err != nil {
			continue
		}

		// مزامنة الحقول مع قاعدة البيانات
		for _, field := range fields {
			// التحقق من وجود الحقل في قاعدة البيانات
			var existing metadata.AppMetadata
			err := s.db.Where("app_name = ? AND model_name = ? AND field_name = ?",
				appName, modelName, field.Name).First(&existing).Error

			if err == gorm.ErrRecordNotFound {
				// حقل جديد - إضافته إلى قاعدة البيانات
				newMeta := &metadata.AppMetadata{
					AppName:    appName,
					ModelName:  modelName,
					FieldName:  field.Name,
					FieldType:  field.Type,
					IsRequired: field.Required,
					IsUnique:   field.Unique,
					Status:     "synced",
				}
				s.db.Create(newMeta)
			}
		}
	}

	return nil
}

// extractFieldsFromModel يستخرج الحقول من ملف النموذج
func (s *SyncManager) extractFieldsFromModel(modelPath, modelName string) ([]FieldInfo, error) {
	var fields []FieldInfo

	file, err := os.Open(modelPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	inStruct := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// البحث عن بداية الـ struct
		if strings.Contains(line, "type "+modelName+" struct {") {
			inStruct = true
			continue
		}

		// نهاية الـ struct
		if inStruct && line == "}" {
			break
		}

		// استخراج الحقول من داخل الـ struct
		if inStruct && line != "" && !strings.HasPrefix(line, "//") {
			// تخطي الحقول المدمجة
			if strings.Contains(line, "gorm.Model") || strings.Contains(line, "BaseModel") {
				// إضافة الحقول الأساسية
				fields = append(fields, FieldInfo{
					Name:     "id",
					Type:     "uint",
					Required: true,
					Unique:   true,
				})
				fields = append(fields, FieldInfo{
					Name:     "created_at",
					Type:     "time.Time",
					Required: false,
					Unique:   false,
				})
				fields = append(fields, FieldInfo{
					Name:     "updated_at",
					Type:     "time.Time",
					Required: false,
					Unique:   false,
				})
				continue
			}

			// تحليل الحقل
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				fieldName := parts[0]
				// تخطي الحقول غير المصدرة (تبدأ بحرف صغير)
				if len(fieldName) > 0 && fieldName[0] >= 'a' && fieldName[0] <= 'z' {
					continue
				}
				// تخطي العلامات
				if strings.HasPrefix(fieldName, "`") {
					continue
				}

				fieldType := parts[1]
				required := false
				unique := false

				// استخراج العلامات
				for _, part := range parts {
					if strings.Contains(part, "not null") || strings.Contains(part, "required") {
						required = true
					}
					if strings.Contains(part, "unique") {
						unique = true
					}
				}

				// تخطي الحقول التي ليست من الحقول الفعلية
				if fieldName == "ID" || fieldName == "CreatedAt" || fieldName == "UpdatedAt" || fieldName == "DeletedAt" {
					continue
				}

				fields = append(fields, FieldInfo{
					Name:     fieldName,
					Type:     fieldType,
					Required: required,
					Unique:   unique,
				})
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return fields, nil
}

// ============================================
// هياكل البيانات المساعدة
// ============================================

// FieldInfo يمثل معلومات الحقل
type FieldInfo struct {
	Name     string
	Type     string
	Required bool
	Unique   bool
}

// ============================================
// دوال إضافية للمزامنة
// ============================================

// SyncAll يزامن جميع التطبيقات
func (s *SyncManager) SyncAll() error {
	appsDir := "apps"
	entries, err := os.ReadDir(appsDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			if err := s.SyncToDB(entry.Name()); err != nil {
				return err
			}
		}
	}

	return nil
}

// GetPendingChanges يعيد التغييرات المعلقة
func (s *SyncManager) GetPendingChanges(appName string) ([]metadata.AppMetadata, error) {
	return s.metadataMgr.GetPending(appName)
}

// ApplyPendingChanges يطبق التغييرات المعلقة
func (s *SyncManager) ApplyPendingChanges(appName string) error {
	pending, err := s.metadataMgr.GetPending(appName)
	if err != nil {
		return err
	}

	for _, meta := range pending {
		modelPath := filepath.Join("apps", meta.AppName, "models", strings.ToLower(meta.ModelName)+".go")
		if err := s.updateModelFile(modelPath, meta); err != nil {
			return err
		}

		// ✅ استخدام UpdateStatus بدلاً من الوصول المباشر إلى db
		if err := s.metadataMgr.UpdateStatus(meta.MigrationID, "applied"); err != nil {
			return err
		}
	}

	return nil
}

// RollbackPendingChanges يلغي التغييرات المعلقة
func (s *SyncManager) RollbackPendingChanges(appName string) error {
	pending, err := s.metadataMgr.GetPending(appName)
	if err != nil {
		return err
	}

	for _, meta := range pending {
		// حذف من قاعدة البيانات
		if err := s.metadataMgr.DeleteField(meta.AppName, meta.ModelName, meta.FieldName); err != nil {
			return err
		}
	}

	return nil
}

// GetSyncStatus يعيد حالة المزامنة
func (s *SyncManager) GetSyncStatus(appName string) (map[string]interface{}, error) {
	pending, err := s.metadataMgr.CountPending(appName)
	if err != nil {
		return nil, err
	}

	metadataList, err := s.metadataMgr.GetByApp(appName)
	if err != nil {
		return nil, err
	}

	status := map[string]interface{}{
		"app_name":      appName,
		"pending_count": pending,
		"total_changes": len(metadataList),
		"last_sync":     time.Now(),
	}

	return status, nil
}
