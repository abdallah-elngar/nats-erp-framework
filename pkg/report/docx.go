package report

import (
	"bytes"
	"io"

	"github.com/nguyenthenguyen/docx"
)

// DOCXGenerator يولد تقارير DOCX
type DOCXGenerator struct {
	config *Config
}

// NewDOCXGenerator ينشئ مولد DOCX جديد
func NewDOCXGenerator(config *Config) *DOCXGenerator {
	return &DOCXGenerator{config: config}
}

// Generate يولد تقرير DOCX
func (d *DOCXGenerator) Generate(report *Report) ([]byte, error) {
	var buf bytes.Buffer
	if err := d.GenerateToWriter(report, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GenerateToWriter يولد تقرير DOCX إلى كاتب
func (d *DOCXGenerator) GenerateToWriter(report *Report, w io.Writer) error {
	// ✅ إنشاء مستند جديد
	doc := new(docx.Docx)

	// ✅ إضافة عنوان (باستخدام الطريقة الصحيحة للمكتبة)
	doc.AddHeading(report.Title, 1)

	// ✅ إضافة تاريخ
	doc.AddParagraph("Date: " + report.Date.Format("2006-01-02"))

	// ✅ إضافة بيانات التقرير (إذا كانت موجودة)
	if report.Data != nil {
		doc.AddParagraph("")
		doc.AddParagraph("Report Data:")
	}

	// ✅ حفظ المستند إلى الكاتب
	return doc.WriteTo(w)
}

// AddHeading يضيف عنواناً
func (d *DOCXGenerator) AddHeading(doc *docx.Docx, text string, level int) {
	doc.AddHeading(text, level)
}

// AddParagraph يضيف فقرة
func (d *DOCXGenerator) AddParagraph(doc *docx.Docx, text string) {
	doc.AddParagraph(text)
}

// AddTable يضيف جدولاً
func (d *DOCXGenerator) AddTable(doc *docx.Docx, headers []string, data [][]string) {
	// إنشاء جدول
	table := doc.AddTable()

	// إضافة رؤوس
	row := table.AddRow()
	for _, header := range headers {
		cell := row.AddCell()
		cell.AddParagraph(header)
	}

	// إضافة بيانات
	for _, rowData := range data {
		row := table.AddRow()
		for _, cellData := range rowData {
			cell := row.AddCell()
			cell.AddParagraph(cellData)
		}
	}
}
