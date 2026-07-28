package report

import (
	"io"
	"time"
)

// Report يمثل تقريراً
type Report struct {
	Title  string
	Author string
	Date   time.Time
	Data   interface{}
	Format string
}

// ReportGenerator يولد التقارير
type ReportGenerator interface {
	Generate(report *Report) ([]byte, error)
	GenerateToWriter(report *Report, w io.Writer) error
}

// Config يمثل إعدادات التقارير
type Config struct {
	FontPath     string
	TemplatePath string
	OutputPath   string
}

// DefaultConfig يعيد الإعدادات الافتراضية
func DefaultConfig() *Config {
	return &Config{
		FontPath:     "./static/fonts/",
		TemplatePath: "./templates/reports/",
		OutputPath:   "./storage/reports/",
	}
}

// NewReport ينشئ تقريراً جديداً
func NewReport(title, author, format string, data interface{}) *Report {
	return &Report{
		Title:  title,
		Author: author,
		Date:   time.Now(),
		Data:   data,
		Format: format,
	}
}
