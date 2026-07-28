package generator

import (
	"fmt"
	"strings"
)

// ControllerGenerator يولد متحكمات
type ControllerGenerator struct {
	generator *Generator
}

// NewControllerGenerator ينشئ مولد متحكمات جديد
func NewControllerGenerator(generator *Generator) *ControllerGenerator {
	return &ControllerGenerator{
		generator: generator,
	}
}

// Generate يولد متحكماً جديداً
func (cg *ControllerGenerator) Generate(appName, modelName string, crud bool) error {
	data := map[string]interface{}{
		"AppName": appName,
		"Model": map[string]string{
			"Name":  modelName,
			"Lower": strings.ToLower(modelName),
		},
		"CRUD": crud,
	}

	// إنشاء ملف المتحكم
	outputPath := fmt.Sprintf("apps/%s/controllers/%s_controller.go", appName, strings.ToLower(modelName))
	return cg.generator.Generate("controller.go.tpl", outputPath, data)
}

// GenerateCRUD يولد دوال CRUD
func (cg *ControllerGenerator) GenerateCRUD(appName, modelName string) error {
	data := map[string]interface{}{
		"AppName": appName,
		"Model": map[string]string{
			"Name":  modelName,
			"Lower": strings.ToLower(modelName),
		},
	}

	outputPath := fmt.Sprintf("apps/%s/controllers/%s_crud.go", appName, strings.ToLower(modelName))
	return cg.generator.Generate("crud.go.tpl", outputPath, data)
}
