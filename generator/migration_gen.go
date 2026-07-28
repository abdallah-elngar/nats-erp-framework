package generator

import (
	"fmt"
	"strings"
	"time"
)

// MigrationGenerator يولد هجرات
type MigrationGenerator struct {
	generator *Generator
}

// NewMigrationGenerator ينشئ مولد هجرات جديد
func NewMigrationGenerator(generator *Generator) *MigrationGenerator {
	return &MigrationGenerator{
		generator: generator,
	}
}

// Generate يولد هجرة جديدة
func (mg *MigrationGenerator) Generate(appName, modelName string) error {
	timestamp := time.Now().Format("20060102150405")
	data := map[string]interface{}{
		"AppName": appName,
		"Model": map[string]string{
			"Name":  modelName,
			"Lower": strings.ToLower(modelName),
		},
		"ID":        timestamp,
		"Timestamp": time.Now().Format("2006-01-02 15:04:05"),
	}

	// إنشاء هجرة GORM
	gormPath := fmt.Sprintf("apps/%s/migrations/gorm/%s_create_%s_table.go",
		appName, timestamp, strings.ToLower(modelName))
	if err := mg.generator.Generate("gorm_migration.go.tpl", gormPath, data); err != nil {
		return err
	}

	// إنشاء هجرة SQL
	sqlPath := fmt.Sprintf("apps/%s/migrations/sql/%s_create_%s_table.sql",
		appName, timestamp, strings.ToLower(modelName))
	return mg.generator.Generate("sql_migration.go.tpl", sqlPath, data)
}
