package report

import (
	"bytes"
	"fmt"
	"io"

	"github.com/jung-kurt/gofpdf/v2"
)

// PDFGenerator يولد تقارير PDF
type PDFGenerator struct {
	config *Config
}

// NewPDFGenerator ينشئ مولد PDF جديد
func NewPDFGenerator(config *Config) *PDFGenerator {
	return &PDFGenerator{config: config}
}

// Generate يولد تقرير PDF
func (p *PDFGenerator) Generate(report *Report) ([]byte, error) {
	var buf bytes.Buffer
	if err := p.GenerateToWriter(report, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GenerateToWriter يولد تقرير PDF إلى كاتب
func (p *PDFGenerator) GenerateToWriter(report *Report, w io.Writer) error {
	pdf := gofpdf.New("P", "mm", "A4", "")

	// إضافة صفحة
	pdf.AddPage()

	// إعداد الخط
	pdf.SetFont("Arial", "B", 16)

	// إضافة عنوان
	pdf.Cell(0, 10, report.Title)
	pdf.Ln(10)

	// إضافة تاريخ
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 10, fmt.Sprintf("Date: %s", report.Date.Format("2006-01-02")))
	pdf.Ln(10)

	// إضافة بيانات التقرير
	if report.Data != nil {
		pdf.SetFont("Arial", "", 10)
		// يمكن إضافة المزيد من المحتوى حسب نوع البيانات
	}

	// حفظ PDF إلى الكاتب
	return pdf.Output(w)
}

// AddHeader يضيف رأساً للصفحة
func (p *PDFGenerator) AddHeader(pdf *gofpdf.Fpdf, text string) {
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, text)
	pdf.Ln(8)
}

// AddParagraph يضيف فقرة
func (p *PDFGenerator) AddParagraph(pdf *gofpdf.Fpdf, text string) {
	pdf.SetFont("Arial", "", 10)
	pdf.MultiCell(0, 5, text, "", "", false)
	pdf.Ln(4)
}

// AddTable يضيف جدولاً
func (p *PDFGenerator) AddTable(pdf *gofpdf.Fpdf, headers []string, data [][]string) {
	// حساب عرض الأعمدة
	colWidth := 190.0 / float64(len(headers))

	// إضافة رأس الجدول
	pdf.SetFont("Arial", "B", 10)
	for _, header := range headers {
		pdf.CellFormat(colWidth, 7, header, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(7)

	// إضافة بيانات الجدول
	pdf.SetFont("Arial", "", 9)
	for _, row := range data {
		for _, cell := range row {
			pdf.CellFormat(colWidth, 7, cell, "1", 0, "L", false, 0, "")
		}
		pdf.Ln(7)
	}
	pdf.Ln(4)
}
