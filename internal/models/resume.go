package models

import (
	"time"
)

// Resume represents a professional resume/CV
type Resume struct {
	Metadata       ResumeMetadata                 `json:"metadata" yaml:"metadata"`
	Contact        ContactInfo                    `json:"contact" yaml:"contact"`
	Summary        string                         `json:"summary" yaml:"summary"`
	Experience     []WorkExperience               `json:"experience" yaml:"experience"`
	Education      []Education                    `json:"education" yaml:"education"`
	Skills         []Skill                        `json:"skills" yaml:"skills"`
	Projects       []Project                      `json:"projects" yaml:"projects"`
	Certifications []Certification                `json:"certifications" yaml:"certifications"`
	Languages      []Language                     `json:"languages" yaml:"languages"`
	CustomSections map[string][]CustomSectionItem `json:"custom_sections,omitempty" yaml:"custom_sections,omitempty"`
}

// ResumeMetadata contains metadata about the resume
type ResumeMetadata struct {
	ID          string    `json:"id" yaml:"id"`
	Title       string    `json:"title" yaml:"title"`
	Version     string    `json:"version" yaml:"version"`
	LastUpdated time.Time `json:"last_updated" yaml:"last_updated"`
	Template    string    `json:"template" yaml:"template"`
}

// ContactInfo contains personal contact information
type ContactInfo struct {
	Name      string `json:"name" yaml:"name"`
	Email     string `json:"email" yaml:"email"`
	Phone     string `json:"phone,omitempty" yaml:"phone,omitempty"`
	LinkedIn  string `json:"linkedin,omitempty" yaml:"linkedin,omitempty"`
	GitHub    string `json:"github,omitempty" yaml:"github,omitempty"`
	Portfolio string `json:"portfolio,omitempty" yaml:"portfolio,omitempty"`
	Location  string `json:"location,omitempty" yaml:"location,omitempty"`
}

// WorkExperience represents a job position
type WorkExperience struct {
	ID           string     `json:"id" yaml:"id"`
	Company      string     `json:"company" yaml:"company"`
	Position     string     `json:"position" yaml:"position"`
	StartDate    time.Time  `json:"start_date" yaml:"start_date"`
	EndDate      *time.Time `json:"end_date,omitempty" yaml:"end_date,omitempty"`
	Current      bool       `json:"current" yaml:"current"`
	Description  string     `json:"description" yaml:"description"`
	Achievements []string   `json:"achievements,omitempty" yaml:"achievements,omitempty"`
	Skills       []string   `json:"skills,omitempty" yaml:"skills,omitempty"`
	Location     string     `json:"location,omitempty" yaml:"location,omitempty"`
}

// Education represents educational background
type Education struct {
	ID           string     `json:"id" yaml:"id"`
	Institution  string     `json:"institution" yaml:"institution"`
	Degree       string     `json:"degree" yaml:"degree"`
	FieldOfStudy string     `json:"field_of_study" yaml:"field_of_study"`
	StartDate    time.Time  `json:"start_date" yaml:"start_date"`
	EndDate      *time.Time `json:"end_date,omitempty" yaml:"end_date,omitempty"`
	GPA          *float64   `json:"gpa,omitempty" yaml:"gpa,omitempty"`
	Courses      []string   `json:"courses,omitempty" yaml:"courses,omitempty"`
	Location     string     `json:"location,omitempty" yaml:"location,omitempty"`
}

// Skill represents a professional skill
type Skill struct {
	ID          string   `json:"id" yaml:"id"`
	Name        string   `json:"name" yaml:"name"`
	Category    string   `json:"category" yaml:"category"`
	Proficiency string   `json:"proficiency" yaml:"proficiency"` // Beginner, Intermediate, Advanced, Expert
	Keywords    []string `json:"keywords,omitempty" yaml:"keywords,omitempty"`
	Years       int      `json:"years,omitempty" yaml:"years,omitempty"`
}

// Project represents a project in the portfolio
type Project struct {
	ID           string     `json:"id" yaml:"id"`
	Name         string     `json:"name" yaml:"name"`
	Description  string     `json:"description" yaml:"description"`
	StartDate    time.Time  `json:"start_date" yaml:"start_date"`
	EndDate      *time.Time `json:"end_date,omitempty" yaml:"end_date,omitempty"`
	URL          string     `json:"url,omitempty" yaml:"url,omitempty"`
	Technologies []string   `json:"technologies,omitempty" yaml:"technologies,omitempty"`
	Role         string     `json:"role,omitempty" yaml:"role,omitempty"`
}

// Certification represents a professional certification
type Certification struct {
	ID           string     `json:"id" yaml:"id"`
	Name         string     `json:"name" yaml:"name"`
	Issuer       string     `json:"issuer" yaml:"issuer"`
	IssueDate    time.Time  `json:"issue_date" yaml:"issue_date"`
	ExpiryDate   *time.Time `json:"expiry_date,omitempty" yaml:"expiry_date,omitempty"`
	CredentialID string     `json:"credential_id,omitempty" yaml:"credential_id,omitempty"`
	URL          string     `json:"url,omitempty" yaml:"url,omitempty"`
}

// Language represents language proficiency
type Language struct {
	ID          string `json:"id" yaml:"id"`
	Language    string `json:"language" yaml:"language"`
	Proficiency string `json:"proficiency" yaml:"proficiency"` // Native, Fluent, Proficient, Basic
}

// CustomSectionItem represents an item in a custom resume section
type CustomSectionItem struct {
	ID          string   `json:"id" yaml:"id"`
	Title       string   `json:"title" yaml:"title"`
	Description string   `json:"description,omitempty" yaml:"description,omitempty"`
	Details     []string `json:"details,omitempty" yaml:"details,omitempty"`
}

// NewResume creates a new empty resume with default metadata
func NewResume() *Resume {
	return &Resume{
		Metadata: ResumeMetadata{
			Title:       "My Resume",
			Version:     "1.0.0",
			LastUpdated: time.Now(),
			Template:    "default",
		},
		Contact:        ContactInfo{},
		Experience:     []WorkExperience{},
		Education:      []Education{},
		Skills:         []Skill{},
		Projects:       []Project{},
		Certifications: []Certification{},
		Languages:      []Language{},
		CustomSections: make(map[string][]CustomSectionItem),
	}
}
