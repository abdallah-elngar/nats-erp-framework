package generator

import (
	"fmt"
	"strings"
)

// ModelGenerator يولد نماذج
type ModelGenerator struct {
	generator *Generator
}

// NewModelGenerator ينشئ مولد نماذج جديد
func NewModelGenerator(generator *Generator) *ModelGenerator {
	return &ModelGenerator{
		generator: generator,
	}
}

// Generate يولد نموذجاً جديداً
func (mg *ModelGenerator) Generate(appName, modelName string, fields []Field, relations []Relation) error {
	data := map[string]interface{}{
		"AppName": appName,
		"Model": map[string]interface{}{
			"Name":      modelName,
			"Lower":     strings.ToLower(modelName),
			"Fields":    fields,
			"Relations": relations,
		},
	}

	outputPath := fmt.Sprintf("apps/%s/models/%s.go", appName, strings.ToLower(modelName))
	return mg.generator.Generate("model.go.tpl", outputPath, data)
}

// GenerateDTO يولد DTO جديداً
func (mg *ModelGenerator) GenerateDTO(appName, modelName string, fields []Field) error {
	data := map[string]interface{}{
		"AppName": appName,
		"Model": map[string]interface{}{
			"Name":   modelName,
			"Lower":  strings.ToLower(modelName),
			"Fields": fields,
		},
	}

	outputPath := fmt.Sprintf("apps/%s/dto/%s_dto.go", appName, strings.ToLower(modelName))
	return mg.generator.Generate("dto.go.tpl", outputPath, data)
}
