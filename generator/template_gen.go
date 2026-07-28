package generator

import (
	"fmt"
	"strings"
)

// TemplateGenerator يولد قوالب
type TemplateGenerator struct {
	generator *Generator
}

// NewTemplateGenerator ينشئ مولد قوالب جديد
func NewTemplateGenerator(generator *Generator) *TemplateGenerator {
	return &TemplateGenerator{
		generator: generator,
	}
}

// Generate يولد قالباً جديداً
func (tg *TemplateGenerator) Generate(appName, templateName, templateType string) error {
	data := map[string]string{
		"AppName":      appName,
		"TemplateName": templateName,
		"Title":        strings.Title(strings.ReplaceAll(templateName, "/", " ")),
		"Type":         templateType,
	}

	outputPath := fmt.Sprintf("apps/%s/templates/%s.html", appName, templateName)
	return tg.generator.Generate("template_"+templateType+".html.tpl", outputPath, data)
}

// GenerateLayout يولد تخطيطاً
func (tg *TemplateGenerator) GenerateLayout(appName, layoutName string) error {
	outputPath := fmt.Sprintf("apps/%s/templates/layouts/%s.html", appName, layoutName)
	return tg.generator.Generate("layout.html.tpl", outputPath, map[string]string{
		"AppName":    appName,
		"LayoutName": layoutName,
	})
}
