package models

// Config represents the application configuration
type Config struct {
	Resume    ResumeConfig    `json:"resume" yaml:"resume"`
	JobOffer  JobOfferConfig  `json:"job_offer" yaml:"job_offer"`
	Output    OutputConfig    `json:"output" yaml:"output"`
	Generator GeneratorConfig `json:"generator,omitempty" yaml:"generator,omitempty"`
	Matcher   MatcherConfig   `json:"matcher,omitempty" yaml:"matcher,omitempty"`
}

// ResumeConfig contains resume configuration
type ResumeConfig struct {
	Path     string `json:"path" yaml:"path"`
	Filename string `json:"filename,omitempty" yaml:"filename,omitempty"`
}

// JobOfferConfig contains job offer configuration
type JobOfferConfig struct {
	Path     string `json:"path" yaml:"path"`
	Filename string `json:"filename,omitempty" yaml:"filename,omitempty"`
}

// OutputConfig contains output configuration
type OutputConfig struct {
	Directory         string `json:"directory" yaml:"directory"`
	UseCompanyAndRole bool   `json:"use_company_and_role" yaml:"use_company_and_role"`
	Format            string `json:"format,omitempty" yaml:"format,omitempty"`
}

// GeneratorConfig contains generator configuration
type GeneratorConfig struct {
	Formats []string `json:"formats,omitempty" yaml:"formats,omitempty"`
}

// MatcherConfig contains matcher configuration
type MatcherConfig struct {
	CriteriaPath string `json:"criteria_path,omitempty" yaml:"criteria_path,omitempty"`
}

// NewDefaultConfig creates a new default configuration
func NewDefaultConfig() *Config {
	return &Config{
		Resume: ResumeConfig{
			Path: "input/resume",
		},
		JobOffer: JobOfferConfig{
			Path: "input/job-offers",
		},
		Output: OutputConfig{
			Directory:         "output",
			UseCompanyAndRole: true,
			Format:            "pdf",
		},
		Generator: GeneratorConfig{
			Formats: []string{"pdf", "docx"},
		},
	}
}
