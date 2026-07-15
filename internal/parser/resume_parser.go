package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/SimonFremont-dev/cv-builder/internal/models"
	"gopkg.in/yaml.v3"
)

// ResumeParser handles parsing of resume files
type ResumeParser struct {
	// Configuration for parsing
	StrictMode bool
}

// NewResumeParser creates a new ResumeParser instance
func NewResumeParser() *ResumeParser {
	return &ResumeParser{
		StrictMode: false,
	}
}

// ParseResume parses a resume from a file based on its extension
func (p *ResumeParser) ParseResume(filePath string) (*models.Resume, error) {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".json":
		return p.ParseResumeJSON(filePath)
	case ".yaml", ".yml":
		return p.ParseResumeYAML(filePath)
	case ".txt", ".text":
		return p.ParseResumeText(filePath)
	default:
		// Try to detect format from content
		return p.ParseResumeAuto(filePath)
	}
}

// ParseResumeJSON parses a resume from a JSON file
func (p *ResumeParser) ParseResumeJSON(filePath string) (*models.Resume, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON resume file: %w", err)
	}

	var resume models.Resume
	if err := json.Unmarshal(data, &resume); err != nil {
		return nil, fmt.Errorf("failed to parse JSON resume: %w", err)
	}

	// Validate the parsed resume
	if p.StrictMode {
		if err := validateResume(&resume); err != nil {
			return nil, fmt.Errorf("resume validation failed: %w", err)
		}
	}

	return &resume, nil
}

// ParseResumeYAML parses a resume from a YAML file
func (p *ResumeParser) ParseResumeYAML(filePath string) (*models.Resume, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML resume file: %w", err)
	}

	var resume models.Resume
	if err := yaml.Unmarshal(data, &resume); err != nil {
		return nil, fmt.Errorf("failed to parse YAML resume: %w", err)
	}

	// Validate the parsed resume
	if p.StrictMode {
		if err := validateResume(&resume); err != nil {
			return nil, fmt.Errorf("resume validation failed: %w", err)
		}
	}

	return &resume, nil
}

// ParseResumeText parses a resume from a plain text file
// This is a simple parser that tries to extract basic information
func (p *ResumeParser) ParseResumeText(filePath string) (*models.Resume, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read text resume file: %w", err)
	}

	resume := models.NewResume()
	content := string(data)

	// Simple text parsing - this would be enhanced with NLP in a production system
	resume.Summary = extractSummary(content)

	// Set a default ID based on filename
	resume.Metadata.ID = strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))

	return resume, nil
}

// ParseResumeAuto tries to automatically detect and parse the resume format
func (p *ResumeParser) ParseResumeAuto(filePath string) (*models.Resume, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read resume file: %w", err)
	}

	// Try JSON first
	var resumeJSON models.Resume
	if err := json.Unmarshal(data, &resumeJSON); err == nil {
		return &resumeJSON, nil
	}

	// Try YAML
	var resumeYAML models.Resume
	if err := yaml.Unmarshal(data, &resumeYAML); err == nil {
		return &resumeYAML, nil
	}

	// Fall back to text parsing
	return p.ParseResumeText(filePath)
}

// ParseResumeFromString parses a resume from a string based on detected format
func (p *ResumeParser) ParseResumeFromString(content string, format string) (*models.Resume, error) {
	format = strings.ToLower(format)

	switch format {
	case "json", "":
		var resume models.Resume
		if err := json.Unmarshal([]byte(content), &resume); err != nil {
			// If JSON fails, try YAML
			if err := yaml.Unmarshal([]byte(content), &resume); err != nil {
				return nil, fmt.Errorf("failed to parse resume: %w", err)
			}
		}
		return &resume, nil
	case "yaml", "yml":
		var resume models.Resume
		if err := yaml.Unmarshal([]byte(content), &resume); err != nil {
			return nil, fmt.Errorf("failed to parse YAML resume: %w", err)
		}
		return &resume, nil
	case "text", "txt":
		resume := models.NewResume()
		resume.Summary = extractSummary(content)
		return resume, nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// validateResume performs basic validation on the resume
func validateResume(resume *models.Resume) error {
	if resume.Contact.Name == "" {
		return fmt.Errorf("resume must have a contact name")
	}
	if resume.Contact.Email == "" {
		return fmt.Errorf("resume must have a contact email")
	}
	return nil
}

// extractSummary extracts a summary from text content
func extractSummary(content string) string {
	lines := strings.Split(content, "\n")
	var summaryLines []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && len(trimmed) > 10 && len(trimmed) < 200 {
			summaryLines = append(summaryLines, trimmed)
		}
		if len(summaryLines) >= 3 {
			break
		}
	}

	return strings.Join(summaryLines, " ")
}
