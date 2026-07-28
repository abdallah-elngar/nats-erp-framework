package orm

import (
	"gorm.io/gorm"
)

// GormAdapter محول لـ GORM
type GormAdapter struct {
	db *gorm.DB
}

// NewGormAdapter ينشئ محولاً جديداً
func NewGormAdapter(db *gorm.DB) *GormAdapter {
	return &GormAdapter{db: db}
}

// AutoMigrate يقوم بإنشاء الجداول تلقائياً
func (a *GormAdapter) AutoMigrate(models ...interface{}) error {
	return a.db.AutoMigrate(models...)
}

// GetTableName يحصل على اسم الجدول للنموذج
func (a *GormAdapter) GetTableName(model interface{}) string {
	stmt := &gorm.Statement{DB: a.db}
	if err := stmt.Parse(model); err != nil {
		return ""
	}
	return stmt.Schema.Table
}

// GetFields يحصل على حقول النموذج
func (a *GormAdapter) GetFields(model interface{}) []string {
	stmt := &gorm.Statement{DB: a.db}
	if err := stmt.Parse(model); err != nil {
		return nil
	}

	var fields []string
	for _, field := range stmt.Schema.Fields {
		if field.DBName != "" {
			fields = append(fields, field.DBName)
		}
	}
	return fields
}

// GetRelations يحصل على علاقات النموذج
func (a *GormAdapter) GetRelations(model interface{}) []string {
	stmt := &gorm.Statement{DB: a.db}
	if err := stmt.Parse(model); err != nil {
		return nil
	}

	var relations []string
	for _, rel := range stmt.Schema.Relationships.Relations {
		relations = append(relations, rel.Name)
	}
	return relations
}

// GenerateMigrationSQL يولد SQL لهجرة
func (a *GormAdapter) GenerateMigrationSQL(models ...interface{}) (string, error) {
	// تنفيذ الترحيل في وضع جاف
	// سيتم تنفيذها لاحقاً
	return "", nil
}
