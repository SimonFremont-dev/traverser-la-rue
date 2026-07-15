package models

import (
	"time"
)

// JobOffer represents a job posting or offer
type JobOffer struct {
	Metadata     JobMetadata  `json:"metadata" yaml:"metadata"`
	Company      CompanyInfo  `json:"company" yaml:"company"`
	Position     PositionInfo `json:"position" yaml:"position"`
	Description  string       `json:"description" yaml:"description"`
	Requirements Requirements `json:"requirements" yaml:"requirements"`
	Benefits     []string     `json:"benefits,omitempty" yaml:"benefits,omitempty"`
	RawContent   string       `json:"-" yaml:"-"` // Original raw content
	Source       string       `json:"source,omitempty" yaml:"source,omitempty"`
}

// JobMetadata contains metadata about the job offer
type JobMetadata struct {
	ID          string     `json:"id" yaml:"id"`
	Title       string     `json:"title" yaml:"title"`
	URL         string     `json:"url,omitempty" yaml:"url,omitempty"`
	PostedDate  *time.Time `json:"posted_date,omitempty" yaml:"posted_date,omitempty"`
	ExpiresDate *time.Time `json:"expires_date,omitempty" yaml:"expires_date,omitempty"`
	Format      string     `json:"format" yaml:"format"` // text, html, pdf, json, yaml
	Language    string     `json:"language,omitempty" yaml:"language,omitempty"`
}

// CompanyInfo contains information about the hiring company
type CompanyInfo struct {
	Name        string `json:"name" yaml:"name"`
	Industry    string `json:"industry,omitempty" yaml:"industry,omitempty"`
	Size        string `json:"size,omitempty" yaml:"size,omitempty"`
	Location    string `json:"location,omitempty" yaml:"location,omitempty"`
	Website     string `json:"website,omitempty" yaml:"website,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
}

// PositionInfo contains information about the job position
type PositionInfo struct {
	Title      string  `json:"title" yaml:"title"`
	Type       string  `json:"type,omitempty" yaml:"type,omitempty"`   // Full-time, Part-time, Contract, Internship
	Level      string  `json:"level,omitempty" yaml:"level,omitempty"` // Entry, Mid, Senior, Lead, Executive
	Department string  `json:"department,omitempty" yaml:"department,omitempty"`
	Salary     *Salary `json:"salary,omitempty" yaml:"salary,omitempty"`
	Location   string  `json:"location,omitempty" yaml:"location,omitempty"`
	Remote     bool    `json:"remote,omitempty" yaml:"remote,omitempty"`
}

// Salary represents compensation information
type Salary struct {
	Min      *float64 `json:"min,omitempty" yaml:"min,omitempty"`
	Max      *float64 `json:"max,omitempty" yaml:"max,omitempty"`
	Currency string   `json:"currency,omitempty" yaml:"currency,omitempty"`
	Period   string   `json:"period,omitempty" yaml:"period,omitempty"` // hourly, monthly, yearly
	Equity   bool     `json:"equity,omitempty" yaml:"equity,omitempty"`
	Bonuses  bool     `json:"bonuses,omitempty" yaml:"bonuses,omitempty"`
}

// Requirements contains the skills and experience required for the job
type Requirements struct {
	RequiredSkills  []SkillRequirement     `json:"required_skills" yaml:"required_skills"`
	PreferredSkills []SkillRequirement     `json:"preferred_skills,omitempty" yaml:"preferred_skills,omitempty"`
	Education       []EducationRequirement `json:"education,omitempty" yaml:"education,omitempty"`
	Experience      ExperienceRequirement  `json:"experience" yaml:"experience"`
	Certifications  []string               `json:"certifications,omitempty" yaml:"certifications,omitempty"`
	Languages       []LanguageRequirement  `json:"languages,omitempty" yaml:"languages,omitempty"`
	SoftSkills      []string               `json:"soft_skills,omitempty" yaml:"soft_skills,omitempty"`
	Keywords        []string               `json:"keywords,omitempty" yaml:"keywords,omitempty"`
}

// SkillRequirement represents a required or preferred skill
type SkillRequirement struct {
	Name         string   `json:"name" yaml:"name"`
	Category     string   `json:"category,omitempty" yaml:"category,omitempty"`
	Years        *int     `json:"years,omitempty" yaml:"years,omitempty"`
	Proficiency  string   `json:"proficiency,omitempty" yaml:"proficiency,omitempty"`
	Required     bool     `json:"required" yaml:"required"`
	NiceToHave   bool     `json:"nice_to_have,omitempty" yaml:"nice_to_have,omitempty"`
	Alternatives []string `json:"alternatives,omitempty" yaml:"alternatives,omitempty"`
}

// EducationRequirement represents education requirements
type EducationRequirement struct {
	Degree       string `json:"degree,omitempty" yaml:"degree,omitempty"`
	FieldOfStudy string `json:"field_of_study,omitempty" yaml:"field_of_study,omitempty"`
	Required     bool   `json:"required" yaml:"required"`
}

// ExperienceRequirement represents work experience requirements
type ExperienceRequirement struct {
	TotalYears *int     `json:"total_years,omitempty" yaml:"total_years,omitempty"`
	MinYears   *int     `json:"min_years,omitempty" yaml:"min_years,omitempty"`
	InField    *int     `json:"in_field,omitempty" yaml:"in_field,omitempty"`
	Management bool     `json:"management,omitempty" yaml:"management,omitempty"`
	Specific   []string `json:"specific,omitempty" yaml:"specific,omitempty"` // Specific experience types
}

// LanguageRequirement represents language requirements
type LanguageRequirement struct {
	Language    string `json:"language" yaml:"language"`
	Proficiency string `json:"proficiency,omitempty" yaml:"proficiency,omitempty"`
	Required    bool   `json:"required" yaml:"required"`
}

// NewJobOffer creates a new empty job offer
func NewJobOffer() *JobOffer {
	return &JobOffer{
		Metadata: JobMetadata{
			Format: "text",
		},
		Company:  CompanyInfo{},
		Position: PositionInfo{},
		Requirements: Requirements{
			RequiredSkills:  []SkillRequirement{},
			PreferredSkills: []SkillRequirement{},
			Education:       []EducationRequirement{},
			Experience:      ExperienceRequirement{},
			Certifications:  []string{},
			Languages:       []LanguageRequirement{},
			SoftSkills:      []string{},
			Keywords:        []string{},
		},
		Benefits: []string{},
	}
}
