package report

import (
	"bytes"
	"io"

	"github.com/xuri/excelize/v2"
)

// ExcelGenerator يولد تقارير Excel
type ExcelGenerator struct {
	config *Config
}

// NewExcelGenerator ينشئ مولد Excel جديد
func NewExcelGenerator(config *Config) *ExcelGenerator {
	return &ExcelGenerator{config: config}
}

// Generate يولد تقرير Excel
func (e *ExcelGenerator) Generate(report *Report) ([]byte, error) {
	var buf bytes.Buffer
	if err := e.GenerateToWriter(report, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GenerateToWriter يولد تقرير Excel إلى كاتب
func (e *ExcelGenerator) GenerateToWriter(report *Report, w io.Writer) error {
	f := excelize.NewFile()

	// إنشاء ورقة
	sheet := "Sheet1"
	f.SetSheetName("Sheet1", sheet)

	// إضافة عنوان
	f.SetCellValue(sheet, "A1", report.Title)
	f.SetCellValue(sheet, "A2", "Date: "+report.Date.Format("2006-01-02"))

	// إضافة بيانات
	// ...

	// حفظ Excel
	return f.Write(w)
}

// AddSheet يضيف ورقة
func (e *ExcelGenerator) AddSheet(f *excelize.File, name string) {
	f.NewSheet(name)
}

// AddHeaders يضيف رؤوساً للجدول
func (e *ExcelGenerator) AddHeaders(f *excelize.File, sheet string, headers []string, row int) {
	for i, header := range headers {
		col := string(rune('A' + i))
		cell := col + string(rune(row+'0'))
		f.SetCellValue(sheet, cell, header)
	}
}
