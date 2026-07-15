package parser

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/SimonFremont-dev/cv-builder/internal/models"
)

// Helper function for string containment check
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// TestResumeParser_ParseJSON tests parsing a JSON resume
func TestResumeParser_ParseJSON(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    *models.Resume
		wantErr bool
	}{
		{
			name:    "valid JSON resume",
			content: createValidJSONResume(),
			want:    createExpectedResume(),
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			content: `this is not valid json or yaml`,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty JSON",
			content: `{}`,
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewResumeParser()
			got, err := p.ParseResumeFromString(tt.content, "json")

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseResumeFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want == nil {
				return
			}

			// Compare key fields
			if got.Contact.Name != tt.want.Contact.Name {
				t.Errorf("Contact.Name = %v, want %v", got.Contact.Name, tt.want.Contact.Name)
			}
			if got.Contact.Email != tt.want.Contact.Email {
				t.Errorf("Contact.Email = %v, want %v", got.Contact.Email, tt.want.Contact.Email)
			}
			if len(got.Skills) != len(tt.want.Skills) {
				t.Errorf("len(Skills) = %v, want %v", len(got.Skills), len(tt.want.Skills))
			}
		})
	}
}

// TestResumeParser_ParseYAML tests parsing a YAML resume
func TestResumeParser_ParseYAML(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{
			name:    "valid YAML resume",
			content: createValidYAMLResume(),
			wantErr: false,
		},
		{
			name:    "invalid YAML",
			content: "invalid: yaml: content:",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewResumeParser()
			_, err := p.ParseResumeFromString(tt.content, "yaml")

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseResumeFromString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestResumeParser_ParseFile tests parsing from a file
func TestResumeParser_ParseFile(t *testing.T) {
	// Create a temporary file
	tmpDir := t.TempDir()
	jsonFile := filepath.Join(tmpDir, "resume.json")
	yamlFile := filepath.Join(tmpDir, "resume.yaml")

	// Write JSON file
	jsonContent := createValidJSONResume()
	if err := os.WriteFile(jsonFile, []byte(jsonContent), 0644); err != nil {
		t.Fatalf("Failed to write JSON file: %v", err)
	}

	// Write YAML file
	yamlContent := createValidYAMLResume()
	if err := os.WriteFile(yamlFile, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to write YAML file: %v", err)
	}

	t.Run("JSON file", func(t *testing.T) {
		p := NewResumeParser()
		got, err := p.ParseResume(jsonFile)
		if err != nil {
			t.Errorf("ParseResume() error = %v", err)
			return
		}
		if got.Contact.Name != "John Doe" {
			t.Errorf("Contact.Name = %v, want John Doe", got.Contact.Name)
		}
	})

	t.Run("YAML file", func(t *testing.T) {
		p := NewResumeParser()
		got, err := p.ParseResume(yamlFile)
		if err != nil {
			t.Errorf("ParseResume() error = %v", err)
			return
		}
		if got.Contact.Name != "John Doe" {
			t.Errorf("Contact.Name = %v, want John Doe", got.Contact.Name)
		}
	})

	t.Run("Text file", func(t *testing.T) {
		textFile := filepath.Join(tmpDir, "resume.txt")
		if err := os.WriteFile(textFile, []byte("Test resume content"), 0644); err != nil {
			t.Fatalf("Failed to write text file: %v", err)
		}
		p := NewResumeParser()
		got, err := p.ParseResume(textFile)
		if err != nil {
			t.Errorf("ParseResume() error = %v", err)
			return
		}
		if got != nil {
			// Text parsing should create a basic resume
			t.Logf("Text parsing successful, summary length: %d", len(got.Summary))
		}
	})
}

// TestResumeParser_NewResumeParser tests the constructor
func TestResumeParser_NewResumeParser(t *testing.T) {
	p := NewResumeParser()
	if p.StrictMode {
		t.Error("StrictMode should be false by default")
	}
}

// TestResumeParser_StrictMode tests strict mode validation
func TestResumeParser_StrictMode(t *testing.T) {
	// Create a temporary file with incomplete resume
	tmpDir := t.TempDir()
	incompleteFile := filepath.Join(tmpDir, "incomplete.json")

	// Write incomplete JSON (missing Contact.Name and Contact.Email)
	incompleteJSON := `{"metadata": {"title": "Test"}, "contact": {"phone": "123"}, "skills": [], "experience": [], "education": []}`
	if err := os.WriteFile(incompleteFile, []byte(incompleteJSON), 0644); err != nil {
		t.Fatalf("Failed to write incomplete JSON file: %v", err)
	}

	p := NewResumeParser()
	p.StrictMode = true

	_, err := p.ParseResumeJSON(incompleteFile)
	if err == nil {
		t.Error("Expected error for incomplete resume in strict mode (missing Contact.Name and Contact.Email)")
		return
	}
	// Check error message contains expected text
	if !contains(err.Error(), "resume must have a contact name") && !contains(err.Error(), "resume must have a contact email") {
		t.Logf("Got error but unexpected message: %s", err.Error())
	}
}

// TestResumeParser_AutoDetect tests automatic format detection
func TestResumeParser_AutoDetect(t *testing.T) {
	p := NewResumeParser()

	t.Run("JSON auto-detect", func(t *testing.T) {
		jsonContent := createValidJSONResume()
		got, err := p.ParseResumeFromString(jsonContent, "")
		if err != nil {
			t.Errorf("Failed to auto-detect JSON: %v", err)
			return
		}
		if got.Contact.Name != "John Doe" {
			t.Errorf("Contact.Name = %v, want John Doe", got.Contact.Name)
		}
	})

	t.Run("YAML auto-detect", func(t *testing.T) {
		yamlContent := createValidYAMLResume()
		got, err := p.ParseResumeFromString(yamlContent, "")
		if err != nil {
			t.Errorf("Failed to auto-detect YAML: %v", err)
			return
		}
		if got.Contact.Name != "John Doe" {
			t.Errorf("Contact.Name = %v, want John Doe", got.Contact.Name)
		}
	})
}

// Helper functions

func createValidJSONResume() string {
	resume := models.Resume{
		Metadata: models.ResumeMetadata{
			Title:       "Test Resume",
			Version:     "1.0.0",
			LastUpdated: time.Now(),
			Template:    "default",
		},
		Contact: models.ContactInfo{
			Name:  "John Doe",
			Email: "john@example.com",
		},
		Summary: "Test summary",
		Skills: []models.Skill{
			{
				ID:       "skill-1",
				Name:     "Go",
				Category: "Programming",
			},
		},
		Experience: []models.WorkExperience{
			{
				ID:       "exp-1",
				Company:  "Test Corp",
				Position: "Developer",
			},
		},
		Education: []models.Education{
			{
				ID:     "edu-1",
				Degree: "Bachelor",
			},
		},
	}

	data, _ := json.Marshal(resume)
	return string(data)
}

func createValidYAMLResume() string {
	return `metadata:
  title: Test Resume
  version: 1.0.0
contact:
  name: John Doe
  email: john@example.com
summary: Test summary
skills:
  - id: skill-1
    name: Go
    category: Programming
experience:
  - id: exp-1
    company: Test Corp
    position: Developer
education:
  - id: edu-1
    degree: Bachelor
`
}

func createExpectedResume() *models.Resume {
	return &models.Resume{
		Metadata: models.ResumeMetadata{
			Title:       "Test Resume",
			Version:     "1.0.0",
			LastUpdated: time.Now(),
			Template:    "default",
		},
		Contact: models.ContactInfo{
			Name:  "John Doe",
			Email: "john@example.com",
		},
		Summary: "Test summary",
		Skills: []models.Skill{
			{
				ID:       "skill-1",
				Name:     "Go",
				Category: "Programming",
			},
		},
		Experience: []models.WorkExperience{
			{
				ID:       "exp-1",
				Company:  "Test Corp",
				Position: "Developer",
			},
		},
		Education: []models.Education{
			{
				ID:     "edu-1",
				Degree: "Bachelor",
			},
		},
	}
}
