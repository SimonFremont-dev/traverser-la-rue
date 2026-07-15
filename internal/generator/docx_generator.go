package generator

import (
	"archive/zip"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/SimonFremont-dev/cv-builder/internal/models"
)

// DOCXGenerator generates DOCX format resumes
type DOCXGenerator struct {
	Generator
}

// NewDOCXGenerator creates a new DOCXGenerator
func NewDOCXGenerator() *DOCXGenerator {
	return &DOCXGenerator{
		Generator: Generator{
			TemplateDir: "templates",
		},
	}
}

// GenerateDOCX generates a resume in DOCX format
func (g *DOCXGenerator) GenerateDOCX(resume *models.Resume, outputPath string) error {
	// First generate markdown content
	markdown := g.resumeToMarkdown(resume)

	// Convert markdown to DOCX
	docxContent, err := markdownToDOCX(markdown, resume.Contact.Name)
	if err != nil {
		return fmt.Errorf("failed to convert markdown to DOCX: %w", err)
	}

	// Write the DOCX file
	if err := os.WriteFile(outputPath, docxContent, 0644); err != nil {
		return fmt.Errorf("failed to write DOCX file: %w", err)
	}

	return nil
}

// markdownToDOCX converts markdown to DOCX format
func markdownToDOCX(markdown string, resumeName string) ([]byte, error) {
	// Create a basic DOCX file structure
	// DOCX is a ZIP archive containing XML files
	
	// Create the document.xml content
	docXML := createDocumentXML(markdown)
	
	// Create a temporary directory to build the DOCX
	tmpDir, err := os.MkdirTemp("", "docx-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)
	
	// Create required directories
	dirs := []string{
		"word",
		"word/_rels",
		"_rels",
		"docProps",
	}
	
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755); err != nil {
			return nil, err
		}
	}
	
	// Write document.xml
	if err := os.WriteFile(filepath.Join(tmpDir, "word", "document.xml"), []byte(docXML), 0644); err != nil {
		return nil, err
	}
	
	// Write [Content_Types].xml
	contentTypes := `<?xml version="1.0" encoding="UTF-8"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Default Extension="png" ContentType="image/png"/>
  <Default Extension="jpeg" ContentType="image/jpeg"/>
  <Default Extension="jpg" ContentType="image/jpeg"/>
  <Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>
</Types>`
	if err := os.WriteFile(filepath.Join(tmpDir, "[Content_Types].xml"), []byte(contentTypes), 0644); err != nil {
		return nil, err
	}
	
	// Write _rels/.rels
	rels := `<?xml version="1.0" encoding="UTF-8"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId3" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/extended-properties" Target="docProps/app.xml"/>
  <Relationship Id="rId2" Type="http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties" Target="docProps/core.xml"/>
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/>
</Relationships>`
	if err := os.WriteFile(filepath.Join(tmpDir, "_rels", ".rels"), []byte(rels), 0644); err != nil {
		return nil, err
	}
	
	// Write docProps/core.xml
	coreProps := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<cp:coreProperties xmlns:cp="http://schemas.openxmlformats.org/package/2006/metadata/core-properties" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:dcterms="http://purl.org/dc/terms/" xmlns:dcmitype="http://purl.org/dc/dcmitype/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <dc:title>%s Resume</dc:title>
  <dc:subject>Resume</dc:subject>
  <dc:creator>%s</dc:creator>
  <cp:lastModifiedBy>%s</cp:lastModifiedBy>
  <dcterms:modified xsi:type="dcterms:dateTime">%s</dcterms:modified>
  <dcterms:created xsi:type="dcterms:dateTime">%s</dcterms:created>
</cp:coreProperties>`,
		resumeName,
		resumeName,
		resumeName,
		time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		time.Now().UTC().Format("2006-01-02T15:04:05Z"),
	)
	if err := os.WriteFile(filepath.Join(tmpDir, "docProps", "core.xml"), []byte(coreProps), 0644); err != nil {
		return nil, err
	}
	
	// Write docProps/app.xml
	appProps := `<?xml version="1.0" encoding="UTF-8"?>
<Properties xmlns="http://schemas.openxmlformats.org/officeDocument/2006/extended-properties" xmlns:vt="http://schemas.openxmlformats.org/officeDocument/2006/docPropsVTypes">
  <Template>Normal.dotm</Template>
  <TotalTime>0</TotalTime>
  <Pages>1</Pages>
  <Words>0</Words>
  <Characters>0</Characters>
  <Application>cv-builder</Application>
  <DocSecurity>0</DocSecurity>
  <ScaleCrop>false</ScaleCrop>
  <HeadingPairs>
    <vt:vector size="9" baseType="variant">
    </vt:vector>
  </HeadingPairs>
  <TitlesOfParts>
    <vt:vector size="1" baseType="lpstr">
      <vt:lpstr>Resume</vt:lpstr>
    </vt:vector>
  </TitlesOfParts>
  <Company></Company>
  <LinksUpToDate>false</LinksUpToDate>
  <SharedDoc>false</SharedDoc>
  <HyperlinksChanged>false</HyperlinksChanged>
  <AppVersion>1.0.0.0</AppVersion>
</Properties>`
	if err := os.WriteFile(filepath.Join(tmpDir, "docProps", "app.xml"), []byte(appProps), 0644); err != nil {
		return nil, err
	}
	
	// Write word/_rels/document.xml.rels
	docRels := `<?xml version="1.0" encoding="UTF-8"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
</Relationships>`
	if err := os.WriteFile(filepath.Join(tmpDir, "word", "_rels", "document.xml.rels"), []byte(docRels), 0644); err != nil {
		return nil, err
	}
	
	// Create the DOCX file (ZIP archive)
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)
	
	// Add all files to the ZIP
	filesToAdd := []struct {
		name string
		path string
	}{
		{"[Content_Types].xml", "[Content_Types].xml"},
		{"_rels/.rels", "_rels/.rels"},
		{"docProps/app.xml", "docProps/app.xml"},
		{"docProps/core.xml", "docProps/core.xml"},
		{"word/_rels/document.xml.rels", "word/_rels/document.xml.rels"},
		{"word/document.xml", "word/document.xml"},
	}
	
	for _, f := range filesToAdd {
		data, err := os.ReadFile(filepath.Join(tmpDir, f.path))
		if err != nil {
			return nil, err
		}
		w, err := w.Create(f.name)
		if err != nil {
			return nil, err
		}
		if _, err := w.Write(data); err != nil {
			return nil, err
		}
	}
	
	if err := w.Close(); err != nil {
		return nil, err
	}
	
	return buf.Bytes(), nil
}

// createDocumentXML creates the main document XML for DOCX
func createDocumentXML(markdown string) string {
	// Convert markdown to HTML-like structure for DOCX
	// This is a simplified conversion
	lines := strings.Split(markdown, "\n")
	
	var xmlBuilder strings.Builder
	
	// Write XML header
	xmlBuilder.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main" xmlns:wne="http://schemas.microsoft.com/office/word/2006/wordml">
<w:body>
`)
	
	// Process each line
	inList := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		
		if trimmed == "" {
			xmlBuilder.WriteString("<w:p><w:r><w:t> </w:t></w:r></w:p>\n")
			continue
		}
		
		// Check for headers
		if strings.HasPrefix(trimmed, "# ") {
			// H1
			title := strings.TrimPrefix(trimmed, "# ")
			xmlBuilder.WriteString(fmt.Sprintf("<w:p><w:pPr><w:pStyle w:val=\"Title\"/></w:pPr><w:r><w:t>%s</w:t></w:r></w:p>\n", escapeXML(title)))
		} else if strings.HasPrefix(trimmed, "## ") {
			// H2
			title := strings.TrimPrefix(trimmed, "## ")
			xmlBuilder.WriteString(fmt.Sprintf("<w:p><w:pPr><w:pStyle w:val=\"Heading1\"/></w:pPr><w:r><w:t>%s</w:t></w:r></w:p>\n", escapeXML(title)))
		} else if strings.HasPrefix(trimmed, "### ") {
			// H3
			title := strings.TrimPrefix(trimmed, "### ")
			xmlBuilder.WriteString(fmt.Sprintf("<w:p><w:pPr><w:pStyle w:val=\"Heading2\"/></w:pPr><w:r><w:t>%s</w:t></w:r></w:p>\n", escapeXML(title)))
		} else if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") {
			// List item
			if !inList {
				xmlBuilder.WriteString("<w:p><w:pPr><w:numPr><w:ilvl w:val=\"0\"/></w:numPr></w:pPr>")
				inList = true
			}
			content := strings.TrimPrefix(strings.TrimPrefix(trimmed, "- "), "* ")
			xmlBuilder.WriteString(fmt.Sprintf("<w:r><w:t>%s</w:t></w:r></w:p>\n", escapeXML(content)))
		} else if strings.HasPrefix(trimmed, "```") {
			// Code block - skip for now
			xmlBuilder.WriteString("<w:p><w:r><w:t> </w:t></w:r></w:p>\n")
		} else {
			// Regular paragraph
			if inList {
				inList = false
			}
			xmlBuilder.WriteString(fmt.Sprintf("<w:p><w:r><w:t>%s</w:t></w:r></w:p>\n", escapeXML(trimmed)))
		}
	}
	
	// Close XML
	xmlBuilder.WriteString(`</w:body>
</w:document>`)
	
	return xmlBuilder.String()
}

// escapeXML escapes special XML characters
func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}

// GenerateMultipleFormats generates resume in multiple formats (DOCX and PDF)
func (g *DOCXGenerator) GenerateMultipleFormats(resume *models.Resume, outputDir string) error {
	// Generate DOCX
	docxPath := filepath.Join(outputDir, fmt.Sprintf("%s_resume.docx", sanitizeFilename(resume.Contact.Name)))
	if err := g.GenerateDOCX(resume, docxPath); err != nil {
		return fmt.Errorf("failed to generate DOCX: %w", err)
	}
	
	// Generate PDF by first creating HTML, then converting
	// For now, we'll generate HTML and note that PDF conversion requires external tools
	htmlPath := filepath.Join(outputDir, fmt.Sprintf("%s_resume.html", sanitizeFilename(resume.Contact.Name)))
	gen := NewGenerator()
	if err := gen.GenerateResume(resume, "html", htmlPath); err != nil {
		return fmt.Errorf("failed to generate HTML: %w", err)
	}
	
	// Note: PDF generation requires external tools like wkhtmltopdf or a PDF library
	// For production use, consider:
	// 1. Using wkhtmltopdf: wkhtmltopdf input.html output.pdf
	// 2. Using a Go PDF library like github.com/go-pdf/fpdf
	// 3. Using a cloud service
	
	return nil
}

// sanitizeFilename creates a safe filename from a string
func sanitizeFilename(name string) string {
	// Remove invalid characters
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "/", "_")
	name = strings.ReplaceAll(name, "\\", "_")
	name = strings.ReplaceAll(name, ":", "_")
	name = strings.ReplaceAll(name, "*", "_")
	name = strings.ReplaceAll(name, "?", "_")
	name = strings.ReplaceAll(name, "\"", "_")
	name = strings.ReplaceAll(name, "<", "_")
	name = strings.ReplaceAll(name, ">", "_")
	name = strings.ReplaceAll(name, "|", "_")
	return strings.ToLower(name)
}
