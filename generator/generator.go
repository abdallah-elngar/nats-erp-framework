package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Generator يولد الملفات
type Generator struct {
	templatesDir string
	outputDir    string
	funcMap      template.FuncMap
}

// NewGenerator ينشئ مولداً جديداً
func NewGenerator(templatesDir, outputDir string) *Generator {
	return &Generator{
		templatesDir: templatesDir,
		outputDir:    outputDir,
		funcMap: template.FuncMap{
			"upper":    strings.ToUpper,
			"lower":    strings.ToLower,
			"title":    strings.Title,
			"join":     strings.Join,
			"add":      func(a, b int) int { return a + b },
			"sub":      func(a, b int) int { return a - b },
			"mul":      func(a, b int) int { return a * b },
			"div":      func(a, b int) int { return a / b },
			"contains": strings.Contains,
			"replace":  strings.ReplaceAll,
			"trim":     strings.TrimSpace,
		},
	}
}

// Generate يقوم بإنشاء ملف من قالب
func (g *Generator) Generate(templateName, outputPath string, data interface{}) error {
	// قراءة القالب
	templatePath := filepath.Join(g.templatesDir, templateName)
	tmpl, err := template.New(templateName).Funcs(g.funcMap).ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// إنشاء الملف
	outputFile := filepath.Join(g.outputDir, outputPath)
	if err := os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
		return err
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// تنفيذ القالب
	if err := tmpl.Execute(file, data); err != nil {
		return err
	}

	return nil
}

// GenerateFromString ينشئ ملفاً من نص القالب
func (g *Generator) GenerateFromString(templateContent, outputPath string, data interface{}) error {
	tmpl, err := template.New("temp").Funcs(g.funcMap).Parse(templateContent)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	outputFile := filepath.Join(g.outputDir, outputPath)
	if err := os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
		return err
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		return err
	}

	return nil
}

// GenerateMultiple ينشئ ملفات متعددة
func (g *Generator) GenerateMultiple(templates map[string]string, data interface{}) error {
	for templateName, outputPath := range templates {
		if err := g.Generate(templateName, outputPath, data); err != nil {
			return err
		}
	}
	return nil
}
