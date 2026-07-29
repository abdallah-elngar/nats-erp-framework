package metadata

import (
	"time"

	"gorm.io/gorm"
)

// AppMetadata يمثل البيانات الوصفية للتطبيق
type AppMetadata struct {
	ID           uint   `gorm:"primaryKey"`
	AppName      string `gorm:"index;not null;size:100"`
	ModelName    string `gorm:"index;not null;size:100"`
	FieldName    string `gorm:"index;not null;size:100"`
	FieldType    string `gorm:"size:50"`
	IsRequired   bool   `gorm:"default:false"`
	IsUnique     bool   `gorm:"default:false"`
	DefaultValue string `gorm:"type:text"`
	CreatedBy    string `gorm:"size:100"`
	MigrationID  string `gorm:"size:50"`
	Status       string `gorm:"default:'pending';size:20"` // pending, applied, failed, rolled_back
	AppliedAt    *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// TableName يعيد اسم الجدول
func (AppMetadata) TableName() string {
	return "app_metadata"
}

// MetadataManager يدير البيانات الوصفية
type MetadataManager struct {
	db *gorm.DB
}

// NewMetadataManager ينشئ مدير بيانات وصفية جديد
func NewMetadataManager(db *gorm.DB) *MetadataManager {
	if db != nil {
		db.AutoMigrate(&AppMetadata{})
	}
	return &MetadataManager{db: db}
}

// Save يحفظ بيانات وصفية جديدة
func (m *MetadataManager) Save(metadata *AppMetadata) error {
	if m.db == nil {
		return nil
	}
	return m.db.Create(metadata).Error
}

// GetByApp يعيد البيانات الوصفية لتطبيق معين
func (m *MetadataManager) GetByApp(appName string) ([]AppMetadata, error) {
	if m.db == nil {
		return []AppMetadata{}, nil
	}
	var metadata []AppMetadata
	err := m.db.Where("app_name = ?", appName).Order("created_at DESC").Find(&metadata).Error
	return metadata, err
}

// GetByAppAndModel يعيد البيانات الوصفية لنموذج معين
func (m *MetadataManager) GetByAppAndModel(appName, modelName string) ([]AppMetadata, error) {
	if m.db == nil {
		return []AppMetadata{}, nil
	}
	var metadata []AppMetadata
	err := m.db.Where("app_name = ? AND model_name = ?", appName, modelName).
		Order("created_at DESC").Find(&metadata).Error
	return metadata, err
}

// GetPending يعيد التغييرات المعلقة
func (m *MetadataManager) GetPending(appName string) ([]AppMetadata, error) {
	if m.db == nil {
		return []AppMetadata{}, nil
	}
	var metadata []AppMetadata
	query := m.db.Where("status = ?", "pending")
	if appName != "" {
		query = query.Where("app_name = ?", appName)
	}
	err := query.Order("created_at ASC").Find(&metadata).Error
	return metadata, err
}

// UpdateStatus يحدث حالة تغيير
func (m *MetadataManager) UpdateStatus(migrationID, status string) error {
	if m.db == nil {
		return nil
	}
	now := time.Now()
	return m.db.Model(&AppMetadata{}).
		Where("migration_id = ?", migrationID).
		Updates(map[string]interface{}{
			"status":     status,
			"applied_at": &now,
			"updated_at": now,
		}).Error
}

// UpdateStatusByApp يحدث حالة جميع تغييرات تطبيق معين
func (m *MetadataManager) UpdateStatusByApp(appName, status string) error {
	if m.db == nil {
		return nil
	}
	now := time.Now()
	return m.db.Model(&AppMetadata{}).
		Where("app_name = ? AND status = ?", appName, "pending").
		Updates(map[string]interface{}{
			"status":     status,
			"applied_at": &now,
			"updated_at": now,
		}).Error
}

// GetByMigrationID يعيد البيانات الوصفية حسب معرف الهجرة
func (m *MetadataManager) GetByMigrationID(migrationID string) (*AppMetadata, error) {
	if m.db == nil {
		return nil, nil
	}
	var metadata AppMetadata
	err := m.db.Where("migration_id = ?", migrationID).First(&metadata).Error
	if err != nil {
		return nil, err
	}
	return &metadata, nil
}

// DeleteByApp يحذف جميع البيانات الوصفية لتطبيق معين
func (m *MetadataManager) DeleteByApp(appName string) error {
	if m.db == nil {
		return nil
	}
	return m.db.Where("app_name = ?", appName).Delete(&AppMetadata{}).Error
}

// DeleteField يحذف حقل معين
func (m *MetadataManager) DeleteField(appName, modelName, fieldName string) error {
	if m.db == nil {
		return nil
	}
	return m.db.Where("app_name = ? AND model_name = ? AND field_name = ?",
		appName, modelName, fieldName).Delete(&AppMetadata{}).Error
}

// CountPending يعيد عدد التغييرات المعلقة
func (m *MetadataManager) CountPending(appName string) (int64, error) {
	if m.db == nil {
		return 0, nil
	}
	var count int64
	query := m.db.Model(&AppMetadata{}).Where("status = ?", "pending")
	if appName != "" {
		query = query.Where("app_name = ?", appName)
	}
	err := query.Count(&count).Error
	return count, err
}

// GetHistory يعيد تاريخ تغييرات تطبيق معين
func (m *MetadataManager) GetHistory(appName string, limit int) ([]AppMetadata, error) {
	if m.db == nil {
		return []AppMetadata{}, nil
	}
	var metadata []AppMetadata
	err := m.db.Where("app_name = ?", appName).
		Order("created_at DESC").
		Limit(limit).
		Find(&metadata).Error
	return metadata, err
}
