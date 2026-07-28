package generator

import (
	"fmt"
)

// StaticGenerator يولد ملفات ثابتة
type StaticGenerator struct {
	generator *Generator
}

// NewStaticGenerator ينشئ مولد ملفات ثابتة جديد
func NewStaticGenerator(generator *Generator) *StaticGenerator {
	return &StaticGenerator{
		generator: generator,
	}
}

// GenerateCSS يولد ملف CSS
func (sg *StaticGenerator) GenerateCSS(name string) error {
	outputPath := fmt.Sprintf("static/css/%s.css", name)
	return sg.generator.Generate("static.css.tpl", outputPath, map[string]string{
		"Name": name,
	})
}

// GenerateJS يولد ملف JavaScript
func (sg *StaticGenerator) GenerateJS(name string) error {
	outputPath := fmt.Sprintf("static/js/%s.js", name)
	return sg.generator.Generate("static.js.tpl", outputPath, map[string]string{
		"Name": name,
	})
}

// GenerateTheme يولد ثيمة
func (sg *StaticGenerator) GenerateTheme(name string) error {
	outputPath := fmt.Sprintf("static/css/themes/%s.css", name)
	return sg.generator.Generate("theme.css.tpl", outputPath, map[string]string{
		"Name": name,
	})
}
