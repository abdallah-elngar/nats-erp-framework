package migration

import (
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

// GormMigrator يدير هجرات GORM
type GormMigrator struct {
	db *gorm.DB
}

// NewGormMigrator ينشئ مدير هجرات GORM جديد
func NewGormMigrator(db *gorm.DB) *GormMigrator {
	return &GormMigrator{db: db}
}

// AutoMigrate يقوم بإنشاء الجداول تلقائياً
func (gm *GormMigrator) AutoMigrate(models ...interface{}) error {
	return gm.db.AutoMigrate(models...)
}

// MigrateTable يقوم بإنشاء جدول
func (gm *GormMigrator) MigrateTable(model interface{}) error {
	return gm.db.Migrator().CreateTable(model)
}

// DropTable يقوم بحذف جدول
func (gm *GormMigrator) DropTable(model interface{}) error {
	return gm.db.Migrator().DropTable(model)
}

// AddColumn يضيف عموداً
func (gm *GormMigrator) AddColumn(model interface{}, field string) error {
	return gm.db.Migrator().AddColumn(model, field)
}

// DropColumn يحذف عموداً
func (gm *GormMigrator) DropColumn(model interface{}, field string) error {
	return gm.db.Migrator().DropColumn(model, field)
}

// AlterColumn يعدل عموداً
func (gm *GormMigrator) AlterColumn(model interface{}, field string) error {
	return gm.db.Migrator().AlterColumn(model, field)
}

// AddIndex يضيف فهرساً
func (gm *GormMigrator) AddIndex(model interface{}, field string) error {
	return gm.db.Migrator().CreateIndex(model, field)
}

// DropIndex يحذف فهرساً
func (gm *GormMigrator) DropIndex(model interface{}, field string) error {
	return gm.db.Migrator().DropIndex(model, field)
}

// AddForeignKey يضيف مفتاحاً خارجياً
func (gm *GormMigrator) AddForeignKey(model interface{}, field string, ref string, onDelete string) error {
	return gm.db.Migrator().CreateConstraint(model, field)
}

// DropForeignKey يحذف مفتاحاً خارجياً
func (gm *GormMigrator) DropForeignKey(model interface{}, field string) error {
	return gm.db.Migrator().DropConstraint(model, field)
}

// HasTable يتحقق من وجود جدول
func (gm *GormMigrator) HasTable(model interface{}) bool {
	return gm.db.Migrator().HasTable(model)
}

// HasColumn يتحقق من وجود عمود
func (gm *GormMigrator) HasColumn(model interface{}, field string) bool {
	return gm.db.Migrator().HasColumn(model, field)
}

// HasIndex يتحقق من وجود فهرس
func (gm *GormMigrator) HasIndex(model interface{}, field string) bool {
	return gm.db.Migrator().HasIndex(model, field)
}

// GetTableName يحصل على اسم الجدول
func (gm *GormMigrator) GetTableName(model interface{}) string {
	stmt := &gorm.Statement{DB: gm.db}
	stmt.Parse(model)
	return stmt.Schema.Table
}

// GenerateMigrationSQL يولد SQL لهجرة
func (gm *GormMigrator) GenerateMigrationSQL(models ...interface{}) (string, error) {
	// تنفيذ الهجرة في وضع جاف
	var sql strings.Builder

	// إنشاء جدول مؤقت
	for _, model := range models {
		tableName := gm.GetTableName(model)
		sql.WriteString(fmt.Sprintf("-- Table: %s\n", tableName))
		sql.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n", tableName))

		// الحصول على حقول النموذج
		t := reflect.TypeOf(model)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			tag := field.Tag.Get("gorm")

			if tag == "-" {
				continue
			}

			// تحليل العلامات
			sql.WriteString(fmt.Sprintf("    %s %s,\n", field.Name, getSQLType(field.Type)))
		}

		sql.WriteString(");\n\n")
	}

	return sql.String(), nil
}

// getSQLType يحصل على نوع SQL
func getSQLType(t reflect.Type) string {
	switch t.Kind() {
	case reflect.Bool:
		return "BOOLEAN"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return "INTEGER"
	case reflect.Int64, reflect.Uint64:
		return "BIGINT"
	case reflect.Float32:
		return "REAL"
	case reflect.Float64:
		return "DOUBLE PRECISION"
	case reflect.String:
		return "VARCHAR(255)"
	case reflect.Struct:
		if t.Name() == "Time" {
			return "TIMESTAMP"
		}
		return "JSON"
	default:
		return "TEXT"
	}
}
