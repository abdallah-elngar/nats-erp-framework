package sync

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"

    "gorm.io/gorm"
    "github.com/nats-framework/nats/pkg/metadata"
)

type SyncManager struct {
    db         *gorm.DB
    metadataMgr *metadata.MetadataManager
}

func NewSyncManager(db *gorm.DB) *SyncManager {
    return &SyncManager{
        db:         db,
        metadataMgr: metadata.NewMetadataManager(db),
    }
}

// SyncFromDB يستعيد التغييرات من قاعدة البيانات إلى الملفات
func (s *SyncManager) SyncFromDB(appName string) error {
    metadata, err := s.metadataMgr.GetByApp(appName)
    if err != nil {
        return err
    }

    for _, meta := range metadata {
        if meta.Status == "applied" {
            continue
        }

        // تحديث ملف النموذج
        modelPath := filepath.Join("apps", meta.AppName, "models", strings.ToLower(meta.ModelName)+".go")
        if err := s.updateModelFile(modelPath, meta); err != nil {
            return err
        }

        // تحديث الحالة
        meta.Status = "applied"
        now := time.Now()
        meta.AppliedAt = &now
        s.metadataMgr.db.Save(&meta)
    }

    return nil
}

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
                    AppName:   appName,
                    ModelName: modelName,
                    FieldName: field.Name,
                    FieldType: field.Type,
                    IsRequired: field.Required,
                    IsUnique:  field.Unique,
                    Status:    "synced",
                }
                s.db.Create(newMeta)
            }
        }
    }

    return nil
}