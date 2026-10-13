package models

// Source represents a data source for skills (ESCO, O*NET, Personal)
type Source struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	IsOfficial  bool   `json:"is_official"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

// Type represents a skill type category (Programming Language, Tool, etc.)
type Type struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
}

// Skill represents a skill with source, type, and name
type Skill struct {
	ID          int    `json:"id"`
	SourceID    int    `json:"source_id"`
	TypeID      int    `json:"type_id"`
	Name        string `json:"name"`
	ExternalID  string `json:"external_id,omitempty"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
	
	// Transient fields (populated from joins)
	SourceName string `json:"source_name,omitempty"`
	TypeName   string `json:"type_name,omitempty"`
	IsOfficial bool   `json:"is_official,omitempty"`
}

// SkillWithDetails represents a skill with full source and type information
type SkillWithDetails struct {
	ID          int    `json:"id"`
	SourceID    int    `json:"source_id"`
	TypeID      int    `json:"type_id"`
	Source      Source `json:"source"`
	Type        Type   `json:"type"`
	Name        string `json:"name"`
	ExternalID  string `json:"external_id,omitempty"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

// Profile represents personal data stored in SQLite
type Profile struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Skills       []ProfileSkill `json:"skills"`
	CustomSkills []CustomSkill  `json:"custom_skills,omitempty"`
	Experience   []ProfileExperience `json:"experience"`
	Education    []ProfileEducation `json:"education"`
	CreatedAt    string `json:"created_at,omitempty"`
	UpdatedAt    string `json:"updated_at,omitempty"`
}

// ProfileSkill represents a link between a profile and a skill
type ProfileSkill struct {
	ID             int    `json:"id"`
	SkillID        int    `json:"skill_id"`
	SourceID       int    `json:"source_id"`
	TypeID         int    `json:"type_id"`
	Proficiency    string `json:"proficiency,omitempty"` // Beginner, Intermediate, Advanced, Expert
	YearsExperience int    `json:"years_experience,omitempty"`
	LastVerified   string `json:"last_verified,omitempty"`
	Notes          string `json:"notes,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	
	// Transient fields
	SkillName     string `json:"skill_name,omitempty"`
	SourceName    string `json:"source_name,omitempty"`
	TypeName      string `json:"type_name,omitempty"`
	IsOfficial    bool   `json:"is_official,omitempty"`
}

// ProfileExperience represents work experience
type ProfileExperience struct {
	ID          int    `json:"id"`
	Company     string `json:"company"`
	Position    string `json:"position"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date,omitempty"`
	Current     bool   `json:"current"`
	Description string `json:"description,omitempty"`
}

// ProfileEducation represents education history
type ProfileEducation struct {
	ID           int    `json:"id"`
	Institution  string `json:"institution"`
	Degree       string `json:"degree"`
	FieldOfStudy string `json:"field_of_study,omitempty"`
	StartDate    string `json:"start_date,omitempty"`
	EndDate      string `json:"end_date,omitempty"`
	GPA          float64 `json:"gpa,omitempty"`
}

// JobOffer represents a job offer to match against
type JobOffer struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Company     string `json:"company,omitempty"`
	Description string `json:"description,omitempty"`
	Requirements string `json:"requirements,omitempty"`
	Location    string `json:"location,omitempty"`
	URL         string `json:"url,omitempty"`
	Source      string `json:"source,omitempty"`
	RawContent  string `json:"raw_content,omitempty"`
}

// KeywordExtractionRequest represents a request to extract keywords from text
type KeywordExtractionRequest struct {
	Text string `json:"text"`
}

// KeywordExtractionResponse contains extracted keywords
type KeywordExtractionResponse struct {
	Keywords     []string          `json:"keywords"`
	Standardized []StandardizedSkill `json:"standardized,omitempty"`
}

// StandardizedSkill represents a keyword mapped to a standardized skill
type StandardizedSkill struct {
	Keyword     string `json:"keyword"`
	Skill       Skill  `json:"skill,omitempty"`
	SkillType   string `json:"skill_type,omitempty"`
	SkillID     int    `json:"skill_id,omitempty"`
	SkillName   string `json:"skill_name,omitempty"`
	Confidence  float64 `json:"confidence"`
}

// MatchRequest represents a request to match a profile against keywords
type MatchRequest struct {
	ProfileID  int    `json:"profile_id"`
	Keywords   []string `json:"keywords"`
	JobOfferID int    `json:"job_offer_id,omitempty"`
}

// MatchedSkill represents a skill that matched
type MatchedSkill struct {
	ProfileSkillID int    `json:"profile_skill_id"`
	Skill         Skill  `json:"skill,omitempty"`
	SkillType     string `json:"skill_type,omitempty"`
	SkillID       int    `json:"skill_id,omitempty"`
	SkillName     string `json:"skill_name,omitempty"`
	Proficiency   string `json:"proficiency,omitempty"`
	YearsExperience int   `json:"years_experience,omitempty"`
	MatchScore    float64 `json:"match_score"`
	MatchType     string `json:"match_type"` // exact, partial, fuzzy
}

// MatchResponse contains matching results
type MatchResponse struct {
	Profile        *Profile            `json:"profile"`
	MatchedSkills []MatchedSkill      `json:"matched_skills"`
	UnmatchedKeywords []string        `json:"unmatched_keywords"`
	MissingSkills []StandardizedSkill `json:"missing_skills"`
	Suggestions   []string            `json:"suggestions,omitempty"`
	MatchScore    float64             `json:"match_score,omitempty"`
	MatchGrade    string              `json:"match_grade,omitempty"`
}

// SkillUpdate represents a request to add/update a skill
type SkillUpdate struct {
	Name        string `json:"name"`
	Category    string `json:"category,omitempty"`
	TypeName    string `json:"type_name,omitempty"` // Type name (not ID)
	Proficiency string `json:"proficiency,omitempty"`
	Years       int    `json:"years,omitempty"`
}

// ResumeGenerationRequest represents a request to generate a resume
type ResumeGenerationRequest struct {
	ProfileID  int `json:"profile_id"`
	JobOfferID int `json:"job_offer_id,omitempty"`
	Keywords   []string `json:"keywords,omitempty"`
	IncludeUnmatched bool `json:"include_unmatched,omitempty"`
}

// Resume represents the generated resume
type Resume struct {
	ID          int    `json:"id"`
	ProfileID   int    `json:"profile_id"`
	JobOfferID  int    `json:"job_offer_id,omitempty"`
	Content     string `json:"content"`
	Format      string `json:"format"` // text, markdown, json, html, pdf, docx
	Filename    string `json:"filename,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
}

// DatabaseConfig represents configuration for database connections
type DatabaseConfig struct {
	SQLitePath string `json:"sqlite_path"`
	SQLDSN    string `json:"sql_dsn,omitempty"`
}

// AppConfig represents application configuration
type AppConfig struct {
	Port          string          `json:"port"`
	SQLitePath    string          `json:"sqlite_path"`
	SQLDSN        string          `json:"sql_dsn,omitempty"`
	KeywordConfig KeywordConfig   `json:"keyword_config"`
	AtlasConfig   AtlasConfig     `json:"atlas_config,omitempty"`
}

// AtlasConfig represents Atlas migration configuration
type AtlasConfig struct {
	URL string `json:"url,omitempty"`
}

// KeywordConfig represents configuration for keyword extraction
type KeywordConfig struct {
	MinLength      int     `json:"min_length"`
	MaxLength      int     `json:"max_length"`
	MatchThreshold float64 `json:"match_threshold"`
}

// CustomSkill represents a user-defined custom skill
type CustomSkill struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Category       string `json:"category,omitempty"`
	Proficiency    string `json:"proficiency,omitempty"`
	YearsExperience int    `json:"years_experience,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
}

// SkillFilter represents filters for skill queries
type SkillFilter struct {
	SourceIDs []int  `json:"source_ids,omitempty"`
	TypeIDs   []int  `json:"type_ids,omitempty"`
	IsOfficial *bool `json:"is_official,omitempty"`
	Search     string `json:"search,omitempty"`
	Limit      int    `json:"limit,omitempty"`
}

// SkillQueryResult represents a paginated skill query result
type SkillQueryResult struct {
	Skills     []SkillWithDetails `json:"skills"`
	TotalCount int                `json:"total_count"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
}
