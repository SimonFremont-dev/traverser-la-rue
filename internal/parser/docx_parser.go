package parser

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/SimonFremont-dev/cv-builder/internal/generator"
	"github.com/SimonFremont-dev/cv-builder/internal/models"
)

// DOCXParser handles parsing of DOCX files
type DOCXParser struct {
	// Configuration
	ExtractText bool
}

// NewDOCXParser creates a new DOCXParser
func NewDOCXParser() *DOCXParser {
	return &DOCXParser{
		ExtractText: true,
	}
}

// ParseDOCX parses a DOCX file and extracts text content
func (p *DOCXParser) ParseDOCX(filePath string) (*models.Resume, error) {
	// Open the DOCX file (which is a ZIP archive)
	r, err := zip.OpenReader(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open DOCX file: %w", err)
	}
	defer r.Close()

	// Find and read document.xml
	var docXML []byte
	for _, f := range r.File {
		if f.Name == "word/document.xml" {
			rc, err := f.Open()
			if err != nil {
				return nil, fmt.Errorf("failed to open document.xml: %w", err)
			}
			defer rc.Close()
			docXML, err = io.ReadAll(rc)
			if err != nil {
				return nil, fmt.Errorf("failed to read document.xml: %w", err)
			}
			break
		}
	}

	if len(docXML) == 0 {
		return nil, fmt.Errorf("document.xml not found in DOCX file")
	}

	// Parse the XML to extract text
	text, err := extractTextFromDOCX(string(docXML))
	if err != nil {
		return nil, fmt.Errorf("failed to extract text from DOCX: %w", err)
	}

	// Create a basic resume from the extracted text
	resume := models.NewResume()
	resume.Metadata.ID = filepath.Base(filePath)
	resume.Summary = text

	return resume, nil
}

// ParseDOCXFromBytes parses DOCX content from a byte slice
func (p *DOCXParser) ParseDOCXFromBytes(data []byte) (*models.Resume, error) {
	// Try to read as ZIP
	r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("failed to read DOCX from bytes: %w", err)
	}

	// Find and read document.xml
	var docXML []byte
	for _, f := range r.File {
		if f.Name == "word/document.xml" {
			rc, err := f.Open()
			if err != nil {
				return nil, fmt.Errorf("failed to open document.xml: %w", err)
			}
			defer rc.Close()
			docXML, err = io.ReadAll(rc)
			if err != nil {
				return nil, fmt.Errorf("failed to read document.xml: %w", err)
			}
			break
		}
	}

	if len(docXML) == 0 {
		return nil, fmt.Errorf("document.xml not found in DOCX file")
	}

	// Parse the XML to extract text
	text, err := extractTextFromDOCX(string(docXML))
	if err != nil {
		return nil, fmt.Errorf("failed to extract text from DOCX: %w", err)
	}

	// Create a basic resume from the extracted text
	resume := models.NewResume()
	resume.Summary = text

	return resume, nil
}

// extractTextFromDOCX extracts plain text from DOCX XML
func extractTextFromDOCX(xmlContent string) (string, error) {
	// Simple approach: extract text from <w:t> elements
	// DOCX XML uses the w: namespace
	var textParts []string
	
	// Parse using a simple approach - look for <w:t> tags
	// This is a simplified parser that extracts text content
	
	// Use xml.Decoder to parse the document
	decoder := xml.NewDecoder(strings.NewReader(xmlContent))
	
	var currentText string
	inText := false
	
	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
		
		switch t := token.(type) {
		case xml.StartElement:
			// Check if this is a text element
			if t.Name.Local == "t" && t.Name.Space == "http://schemas.openxmlformats.org/wordprocessingml/2006/main" {
				inText = true
				currentText = ""
			}
		case xml.EndElement:
			if t.Name.Local == "t" && t.Name.Space == "http://schemas.openxmlformats.org/wordprocessingml/2006/main" {
				if currentText != "" {
					textParts = append(textParts, currentText)
				}
				inText = false
				currentText = ""
			}
		case xml.CharData:
			if inText {
				currentText += string(t)
			}
		}
	}
	
	// Join all text parts with spaces
	return strings.Join(textParts, " "), nil
}

// PType represents a paragraph in DOCX
type PType struct {
	R []RType `xml:"r"`
}

// RType represents a run in DOCX
type RType struct {
	T []TType `xml:"t"`
}

// TType represents text in DOCX
type TType struct {
	Text string `xml:",chardata"`
	Space string `xml:"space,attr,omitempty"`
}

// TblType represents a table in DOCX
type TblType struct {
	Tr []TrType `xml:"tr"`
}

// TrType represents a table row
type TrType struct {
	Tc []TcType `xml:"tc"`
}

// TcType represents a table cell
type TcType struct {
	P []PType `xml:"p"`
}

// ParseDOCXAndAdapt parses a DOCX resume and adapts it to match a job offer
// This is a convenience function that combines parsing and adaptation
func (p *DOCXParser) ParseDOCXAndAdapt(docxPath string, jobOfferPath string) (*models.Resume, *models.MatchResult, error) {
	// Parse the DOCX resume
	resume, err := p.ParseDOCX(docxPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse DOCX: %w", err)
	}

	// Parse the job offer
	jobParser := NewJobParser()
	jobOffer, err := jobParser.ParseJobOffer(jobOfferPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse job offer: %w", err)
	}

	// Create matcher and match
	// Note: This is a simplified version - in practice you'd want to use the full matcher
	// For now, just return the parsed resume without adaptation
	_ = jobOffer // jobOffer is parsed but not used in this simplified version
	return resume, nil, nil
}

// ParseDOCXAndGenerate parses a DOCX file and generates output in specified format
func (p *DOCXParser) ParseDOCXAndGenerate(docxPath string, outputPath string, format string) error {
	// Parse the DOCX
	resume, err := p.ParseDOCX(docxPath)
	if err != nil {
		return fmt.Errorf("failed to parse DOCX: %w", err)
	}

	// Generate output
	gen := generator.NewGenerator()
	if err := gen.GenerateResume(resume, format, outputPath); err != nil {
		return fmt.Errorf("failed to generate output: %w", err)
	}

	return nil
}

// ParseDOCXAndGenerateAll parses a DOCX file and generates DOCX and PDF output
func (p *DOCXParser) ParseDOCXAndGenerateAll(docxPath string, outputDir string) error {
	// Parse the DOCX
	resume, err := p.ParseDOCX(docxPath)
	if err != nil {
		return fmt.Errorf("failed to parse DOCX: %w", err)
	}

	// Generate multiple outputs
	pdfGen := generator.NewPDFGenerator()
	if err := pdfGen.GenerateMultipleOutputs(resume, outputDir); err != nil {
		return fmt.Errorf("failed to generate outputs: %w", err)
	}

	return nil
}
