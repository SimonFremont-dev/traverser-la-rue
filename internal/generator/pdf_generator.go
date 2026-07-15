package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/SimonFremont-dev/cv-builder/internal/models"
	"github.com/signintech/gopdf"
)

// PDFGenerator generates PDF format resumes
type PDFGenerator struct {
	Generator
}

// NewPDFGenerator creates a new PDFGenerator
func NewPDFGenerator() *PDFGenerator {
	return &PDFGenerator{
		Generator: Generator{
			TemplateDir: "templates",
		},
	}
}

// GeneratePDF generates a resume in PDF format
func (g *PDFGenerator) GeneratePDF(resume *models.Resume, outputPath string) error {
	// Create new PDF document
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	// Add a page
	pdf.AddPage()

	// Draw resume content
	g.drawResume(pdf, resume)

	// Write PDF to file
	if err := pdf.WritePdf(outputPath); err != nil {
		return fmt.Errorf("failed to write PDF: %w", err)
	}

	return nil
}

// drawResume draws the resume content on the PDF
func (g *PDFGenerator) drawResume(pdf *gopdf.GoPdf, resume *models.Resume) {
	// Set up the page
	leftMargin := 50.0
	topMargin := 50.0
	lineHeight := 8.0
	pageWidth := 595.0
	pageHeight := 842.0
	contentWidth := pageWidth - leftMargin*2

	// Set initial font
	pdf.SetFont("Helvetica", "", 12)

	// Current position
	x := leftMargin
	y := topMargin

	// Draw header (name)
	pdf.SetFont("Helvetica", "B", 24)
	pdf.SetXY(x, y)
	pdf.Cell(nil, fmt.Sprintf("%s", resume.Contact.Name))
	y += lineHeight * 2

	// Draw contact info
	pdf.SetFont("Helvetica", "", 12)
	pdf.SetTextColor(128, 128, 128)

	if resume.Contact.Email != "" {
		pdf.SetXY(x, y)
		pdf.Cell(nil, fmt.Sprintf("Email: %s", resume.Contact.Email))
		y += lineHeight
	}

	if resume.Contact.Phone != "" {
		pdf.SetXY(x, y)
		pdf.Cell(nil, fmt.Sprintf("Phone: %s", resume.Contact.Phone))
		y += lineHeight
	}

	if resume.Contact.Location != "" {
		pdf.SetXY(x, y)
		pdf.Cell(nil, fmt.Sprintf("Location: %s", resume.Contact.Location))
		y += lineHeight
	}

	y += lineHeight

	// Draw summary
	pdf.SetFont("Helvetica", "B", 14)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetXY(x, y)
	pdf.Cell(nil, "SUMMARY")
	y += lineHeight

	pdf.SetFont("Helvetica", "", 11)
	wrapText(pdf, resume.Summary, &x, &y, contentWidth, lineHeight)
	y += lineHeight

	// Draw skills
	if len(resume.Skills) > 0 {
		pdf.SetFont("Helvetica", "B", 14)
		pdf.SetTextColor(0, 0, 0)
		pdf.SetXY(x, y)
		pdf.Cell(nil, "SKILLS")
		y += lineHeight

		pdf.SetFont("Helvetica", "", 11)
		for _, skill := range resume.Skills {
			if y > pageHeight-50 {
				// New page
				pdf.AddPage()
				pdf.SetXY(leftMargin, topMargin)
				x = leftMargin
				y = topMargin
			}

			pdf.SetXY(x, y)
			pdf.Cell(nil, fmt.Sprintf("• %s", skill.Name))
			y += lineHeight
		}
		y += lineHeight
	}

	// Draw experience
	if len(resume.Experience) > 0 {
		pdf.SetFont("Helvetica", "B", 14)
		pdf.SetTextColor(0, 0, 0)
		pdf.SetXY(x, y)
		pdf.Cell(nil, "EXPERIENCE")
		y += lineHeight

		pdf.SetFont("Helvetica", "", 11)
		for _, exp := range resume.Experience {
			if y > pageHeight-50 {
				pdf.AddPage()
				pdf.SetXY(leftMargin, topMargin)
				x = leftMargin
				y = topMargin
			}

			// Position and company
			pdf.SetFont("Helvetica", "B", 12)
			pdf.SetTextColor(0, 0, 0)
			pdf.SetXY(x, y)
			pdf.Cell(nil, fmt.Sprintf("%s at %s", exp.Position, exp.Company))
			y += lineHeight

			// Dates
			pdf.SetFont("Helvetica", "", 10)
			pdf.SetTextColor(128, 128, 128)
			pdf.SetXY(x, y)
			pdf.Cell(nil, formatDateRangeForPDF(exp.StartDate, exp.EndDate, exp.Current))
			y += lineHeight

			// Description
			pdf.SetFont("Helvetica", "", 11)
			pdf.SetTextColor(0, 0, 0)
			xIndented := x + 20
			wrapText(pdf, exp.Description, &xIndented, &y, contentWidth-20, lineHeight)
			y += lineHeight
		}
	}

	// Draw education
	if len(resume.Education) > 0 {
		pdf.SetFont("Helvetica", "B", 14)
		pdf.SetTextColor(0, 0, 0)
		pdf.SetXY(x, y)
		pdf.Cell(nil, "EDUCATION")
		y += lineHeight

		pdf.SetFont("Helvetica", "", 11)
		for _, edu := range resume.Education {
			if y > pageHeight-50 {
				pdf.AddPage()
				pdf.SetXY(leftMargin, topMargin)
				x = leftMargin
				y = topMargin
			}

			pdf.SetXY(x, y)
			pdf.Cell(nil, fmt.Sprintf("%s in %s - %s", edu.Degree, edu.FieldOfStudy, edu.Institution))
			y += lineHeight
		}
	}
}

// wrapText wraps text to fit within a width for gopdf
func wrapText(pdf *gopdf.GoPdf, text string, x *float64, y *float64, maxWidth float64, lineHeight float64) {
	words := strings.Fields(text)
	currentLine := ""

	for _, word := range words {
		testLine := currentLine
		if testLine != "" {
			testLine += " "
		}
		testLine += word

		// Estimate width - gopdf doesn't have easy text measurement
		// Use character count as approximation (roughly 6 chars per mm)
		if len(testLine) > int(maxWidth/2) {
			if currentLine != "" {
				pdf.SetXY(*x, *y)
				pdf.Cell(nil, currentLine)
				*y += lineHeight
			}
			currentLine = word
		} else {
			currentLine = testLine
		}
	}

	// Draw the last line
	if currentLine != "" {
		pdf.SetXY(*x, *y)
		pdf.Cell(nil, currentLine)
		*y += lineHeight
	}
}

// formatDateRangeForPDF formats a date range for PDF display
func formatDateRangeForPDF(start time.Time, end *time.Time, current bool) string {
	startStr := start.Format("Jan 2006")

	if current || end == nil {
		return fmt.Sprintf("%s - Present", startStr)
	}

	endStr := end.Format("Jan 2006")
	return fmt.Sprintf("%s - %s", startStr, endStr)
}

// GeneratePDFFromHTML generates a PDF from HTML content
func (g *PDFGenerator) GeneratePDFFromHTML(htmlContent string, outputPath string) error {
	// For production, use wkhtmltopdf or a Go HTML to PDF library
	// For now, we'll generate a simple PDF from the HTML text content

	// As a fallback, generate a simple PDF from the HTML text content
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()
	pdf.SetFont("Helvetica", "", 12)

	// Simple text extraction from HTML (remove tags)
	text := stripHTMLTags(htmlContent)

	// Write text to PDF
	leftMargin := 50.0
	topMargin := 50.0
	lineHeight := 8.0
	contentWidth := 595.0 - leftMargin*2

	x := leftMargin
	y := topMargin

	wrapText(pdf, text, &x, &y, contentWidth, lineHeight)

	if err := pdf.WritePdf(outputPath); err != nil {
		return fmt.Errorf("failed to write PDF from HTML: %w", err)
	}

	return nil
}

// stripHTMLTags removes HTML tags from a string
func stripHTMLTags(html string) string {
	inTag := false
	var result strings.Builder

	for _, ch := range html {
		if ch == '<' {
			inTag = true
		} else if ch == '>' {
			inTag = false
		} else if !inTag {
			result.WriteRune(ch)
		}
	}

	return result.String()
}

// GenerateFromResume generates PDF directly from resume
func (g *PDFGenerator) GenerateFromResume(resume *models.Resume, outputPath string) error {
	return g.GeneratePDF(resume, outputPath)
}

// GenerateMultipleOutputs generates DOCX, PDF, and other formats
func (g *PDFGenerator) GenerateMultipleOutputs(resume *models.Resume, outputDir string) error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	baseName := sanitizeFilename(resume.Contact.Name)

	// Generate DOCX
	docxGen := NewDOCXGenerator()
	docxPath := filepath.Join(outputDir, fmt.Sprintf("%s_resume.docx", baseName))
	if err := docxGen.GenerateDOCX(resume, docxPath); err != nil {
		return fmt.Errorf("failed to generate DOCX: %w", err)
	}

	// Generate PDF
	pdfPath := filepath.Join(outputDir, fmt.Sprintf("%s_resume.pdf", baseName))
	if err := g.GeneratePDF(resume, pdfPath); err != nil {
		return fmt.Errorf("failed to generate PDF: %w", err)
	}

	// Generate HTML for reference
	htmlPath := filepath.Join(outputDir, fmt.Sprintf("%s_resume.html", baseName))
	gen := NewGenerator()
	if err := gen.GenerateResume(resume, "html", htmlPath); err != nil {
		return fmt.Errorf("failed to generate HTML: %w", err)
	}

	// Generate Markdown for reference
	mdPath := filepath.Join(outputDir, fmt.Sprintf("%s_resume.md", baseName))
	if err := gen.GenerateResume(resume, "markdown", mdPath); err != nil {
		return fmt.Errorf("failed to generate Markdown: %w", err)
	}

	return nil
}
