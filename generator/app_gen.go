package generator

import (
	"fmt"
	"strings"
	"time"
)

// Model يمثل نموذجاً في النظام
type Model struct {
	Name      string
	Fields    []Field
	Relations []Relation
}

// Field يمثل حقلاً في النموذج
type Field struct {
	Name      string
	Type      string
	Required  bool
	Unique    bool
	Max       int
	Min       int
	Relations []Relation
}

// Relation يمثل علاقة بين النماذج
type Relation struct {
	Type       string // belongs_to, has_many, has_one, many_to_many
	Model      string
	ForeignKey string
}

// AppGenerator يولد تطبيقات
type AppGenerator struct {
	generator *Generator
}

// NewAppGenerator ينشئ مولد تطبيقات جديد
func NewAppGenerator(generator *Generator) *AppGenerator {
	return &AppGenerator{
		generator: generator,
	}
}

// GenerateApp يولد تطبيقاً جديداً
func (ag *AppGenerator) GenerateApp(appName string, models []Model) error {
	data := map[string]interface{}{
		"AppName":   appName,
		"AppTitle":  strings.Title(appName),
		"Models":    models,
		"Timestamp": time.Now().Format("2006-01-02 15:04:05"),
	}

	templates := map[string]string{
		"app.go.tpl":      fmt.Sprintf("apps/%s/app.go", appName),
		"register.go.tpl": fmt.Sprintf("apps/%s/register.go", appName),
		"router.go.tpl":   fmt.Sprintf("apps/%s/routes/router.go", appName),
	}

	return ag.generator.GenerateMultiple(templates, data)
}

// GenerateModel يولد نموذجاً
func (ag *AppGenerator) GenerateModel(appName string, model Model) error {
	data := map[string]interface{}{
		"AppName":   appName,
		"Model":     model,
		"Timestamp": time.Now().Format("2006-01-02 15:04:05"),
	}

	outputPath := fmt.Sprintf("apps/%s/models/%s.go", appName, strings.ToLower(model.Name))
	return ag.generator.Generate("model.go.tpl", outputPath, data)
}

// GenerateController يولد متحكماً
func (ag *AppGenerator) GenerateController(appName string, model Model) error {
	data := map[string]interface{}{
		"AppName":   appName,
		"Model":     model,
		"Timestamp": time.Now().Format("2006-01-02 15:04:05"),
	}

	outputPath := fmt.Sprintf("apps/%s/controllers/%s_controller.go", appName, strings.ToLower(model.Name))
	return ag.generator.Generate("controller.go.tpl", outputPath, data)
}

// GenerateMigration يولد هجرة
func (ag *AppGenerator) GenerateMigration(appName string, model Model) error {
	timestamp := time.Now().Format("20060102150405")
	data := map[string]interface{}{
		"AppName":   appName,
		"Model":     model,
		"Timestamp": time.Now().Format("2006-01-02 15:04:05"),
		"ID":        timestamp,
	}

	gormPath := fmt.Sprintf("apps/%s/migrations/gorm/%s_create_%s_table.go",
		appName, data["ID"], strings.ToLower(model.Name))
	if err := ag.generator.Generate("gorm_migration.go.tpl", gormPath, data); err != nil {
		return err
	}

	sqlPath := fmt.Sprintf("apps/%s/migrations/sql/%s_create_%s_table.sql",
		appName, data["ID"], strings.ToLower(model.Name))
	return ag.generator.Generate("sql_migration.go.tpl", sqlPath, data)
}

// GenerateCRUD يولد دوال CRUD
func (ag *AppGenerator) GenerateCRUD(appName string, model Model) error {
	data := map[string]interface{}{
		"AppName": appName,
		"Model":   model,
	}

	outputPath := fmt.Sprintf("apps/%s/controllers/%s_crud.go", appName, strings.ToLower(model.Name))
	return ag.generator.Generate("crud.go.tpl", outputPath, data)
}
