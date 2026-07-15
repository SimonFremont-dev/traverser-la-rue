package models

import (
	"sort"
)

// MatchResult represents the result of matching a resume against a job offer
type MatchResult struct {
	JobOfferID          string            `json:"job_offer_id" yaml:"job_offer_id"`
	ResumeID            string            `json:"resume_id" yaml:"resume_id"`
	OverallScore        float64           `json:"overall_score" yaml:"overall_score"` // 0-100
	MatchGrade          string            `json:"match_grade" yaml:"match_grade"`     // A, B, C, D, F
	SkillMatches        []SkillMatch      `json:"skill_matches" yaml:"skill_matches"`
	ExperienceMatch     ExperienceMatch   `json:"experience_match" yaml:"experience_match"`
	EducationMatch      EducationMatch    `json:"education_match" yaml:"education_match"`
	MissingRequirements []string          `json:"missing_requirements" yaml:"missing_requirements"`
	SuggestedChanges    []SuggestedChange `json:"suggested_changes" yaml:"suggested_changes"`
	AdaptedResume       *Resume           `json:"adapted_resume,omitempty" yaml:"adapted_resume,omitempty"`
}

// SkillMatch represents a match between a resume skill and a job requirement
type SkillMatch struct {
	ResumeSkill      *Skill            `json:"resume_skill" yaml:"resume_skill"`
	JobRequirement   *SkillRequirement `json:"job_requirement" yaml:"job_requirement"`
	MatchScore       float64           `json:"match_score" yaml:"match_score"` // 0-100
	MatchType        string            `json:"match_type" yaml:"match_type"`   // Exact, Partial, Alternative, None
	YearsMatch       float64           `json:"years_match,omitempty" yaml:"years_match,omitempty"`
	ProficiencyMatch string            `json:"proficiency_match,omitempty" yaml:"proficiency_match,omitempty"`
}

// ExperienceMatch represents the match between resume experience and job requirements
type ExperienceMatch struct {
	TotalYearsMatch    float64                  `json:"total_years_match" yaml:"total_years_match"`
	InFieldMatch       float64                  `json:"in_field_match,omitempty" yaml:"in_field_match,omitempty"`
	ManagementMatch    bool                     `json:"management_match" yaml:"management_match"`
	RelevantExperience []RelevantExperienceItem `json:"relevant_experience" yaml:"relevant_experience"`
	MissingExperience  []string                 `json:"missing_experience,omitempty" yaml:"missing_experience,omitempty"`
}

// RelevantExperienceItem represents a work experience that matches job requirements
type RelevantExperienceItem struct {
	WorkExperience  *WorkExperience `json:"work_experience" yaml:"work_experience"`
	RelevanceScore  float64         `json:"relevance_score" yaml:"relevance_score"`
	MatchedKeywords []string        `json:"matched_keywords" yaml:"matched_keywords"`
}

// EducationMatch represents the match between resume education and job requirements
type EducationMatch struct {
	DegreeMatch       bool         `json:"degree_match" yaml:"degree_match"`
	FieldMatch        bool         `json:"field_match" yaml:"field_match"`
	GPAMeet           bool         `json:"gpa_meet,omitempty" yaml:"gpa_meet,omitempty"`
	RelevantEducation []*Education `json:"relevant_education,omitempty" yaml:"relevant_education,omitempty"`
	MissingEducation  []string     `json:"missing_education,omitempty" yaml:"missing_education,omitempty"`
}

// SuggestedChange represents a suggested change to improve resume match
type SuggestedChange struct {
	Type          string `json:"type" yaml:"type"`         // AddSkill, RemoveSkill, Reword, Reorder, Emphasize
	Priority      string `json:"priority" yaml:"priority"` // High, Medium, Low
	Section       string `json:"section" yaml:"section"`
	ItemID        string `json:"item_id,omitempty" yaml:"item_id,omitempty"`
	CurrentText   string `json:"current_text,omitempty" yaml:"current_text,omitempty"`
	SuggestedText string `json:"suggested_text,omitempty" yaml:"suggested_text,omitempty"`
	Reason        string `json:"reason" yaml:"reason"`
}

// MatchCriteria represents criteria for matching
type MatchCriteria struct {
	SkillWeight           float64 `json:"skill_weight" yaml:"skill_weight"`
	ExperienceWeight      float64 `json:"experience_weight" yaml:"experience_weight"`
	EducationWeight       float64 `json:"education_weight" yaml:"education_weight"`
	KeywordWeight         float64 `json:"keyword_weight" yaml:"keyword_weight"`
	MinScoreForSuggestion float64 `json:"min_score_for_suggestion" yaml:"min_score_for_suggestion"`
	SuggestAlternatives   bool    `json:"suggest_alternatives" yaml:"suggest_alternatives"`
}

// DefaultMatchCriteria returns default matching criteria
func DefaultMatchCriteria() MatchCriteria {
	return MatchCriteria{
		SkillWeight:           0.4,
		ExperienceWeight:      0.3,
		EducationWeight:       0.2,
		KeywordWeight:         0.1,
		MinScoreForSuggestion: 70.0,
		SuggestAlternatives:   true,
	}
}

// SortSkillMatchesByScore sorts skill matches by match score in descending order
func SortSkillMatchesByScore(matches []SkillMatch) []SkillMatch {
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].MatchScore > matches[j].MatchScore
	})
	return matches
}

// SortSuggestedChangesByPriority sorts suggested changes by priority
func SortSuggestedChangesByPriority(changes []SuggestedChange) []SuggestedChange {
	priorityOrder := map[string]int{
		"High":   0,
		"Medium": 1,
		"Low":    2,
	}
	sort.Slice(changes, func(i, j int) bool {
		return priorityOrder[changes[i].Priority] < priorityOrder[changes[j].Priority]
	})
	return changes
}
