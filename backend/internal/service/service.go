package service

import (
	"fmt"
	"strings"

	"github.com/SimonFremont-dev/cv-builder/backend/internal/models"
	"github.com/SimonFremont-dev/cv-builder/backend/internal/repository"
	"github.com/SimonFremont-dev/cv-builder/backend/internal/keyword"
)

// Service handles business logic
type Service struct {
	repo        *repository.SQLiteRepository
	sqlRepo    *repository.SQLRepository
	extractor *keyword.Extractor
}

// NewService creates a new service
func NewService(repo *repository.SQLiteRepository, sqlRepo *repository.SQLRepository) *Service {
	// Initialize keyword extractor with technical keywords
	categories := map[string][]string{}
	extractor := keyword.NewExtractor(categories)
	
	return &Service{
		repo:     repo,
		sqlRepo: sqlRepo,
		extractor: extractor,
	}
}

// ExtractKeywords extracts keywords from text
func (s *Service) ExtractKeywords(text string) ([]string, error) {
	keywords := s.extractor.ExtractTechnicalKeywords(text)
	return keywords, nil
}

// ExtractAndStandardizeKeywords extracts keywords and tries to map them to standardized skills
func (s *Service) ExtractAndStandardizeKeywords(text string) ([]models.StandardizedSkill, error) {
	keywords := s.extractor.ExtractTechnicalKeywords(text)
	
	// For now, just return as standardized skills with empty skill info
	// This will be enhanced when we have the SQL repository connected
	var result []models.StandardizedSkill
	for _, kw := range keywords {
		result = append(result, models.StandardizedSkill{
			Keyword: kw,
			SkillType: "custom",
			SkillName: kw,
			Confidence: 0.8,
		})
	}
	
	return result, nil
}

// MatchProfileToKeywords matches a profile against extracted keywords
func (s *Service) MatchProfileToKeywords(profileID int, keywords []string) (*models.MatchResponse, error) {
	// Get profile
	profile, err := s.repo.GetProfile(profileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}

	var matchedSkills []models.MatchedSkill
	var unmatchedKeywords []string
	var missingSkills []models.StandardizedSkill

	// Get all custom skills from profile
	profileCustomSkills := make(map[string]models.CustomSkill)
	for _, skill := range profile.CustomSkills {
		profileCustomSkills[strings.ToLower(skill.Name)] = skill
	}

	// Check each keyword
	for _, kw := range keywords {
		kwLower := strings.ToLower(kw)
		
		// Check if keyword matches any profile skill
		if skill, exists := profileCustomSkills[kwLower]; exists {
			matchedSkills = append(matchedSkills, models.MatchedSkill{
				ProfileSkillID: skill.ID,
				SkillType:      "custom",
				SkillID:        skill.ID,
				SkillName:      skill.Name,
				Proficiency:    skill.Proficiency,
				YearsExperience: skill.YearsExperience,
				MatchScore:     100.0,
				MatchType:      "exact",
			})
		} else {
			// Check if any profile skill contains the keyword
			found := false
			for _, skill := range profile.CustomSkills {
				if strings.Contains(strings.ToLower(skill.Name), kwLower) {
					matchedSkills = append(matchedSkills, models.MatchedSkill{
						ProfileSkillID: skill.ID,
						SkillType:      "custom",
						SkillID:        skill.ID,
						SkillName:      skill.Name,
						Proficiency:    skill.Proficiency,
						YearsExperience: skill.YearsExperience,
						MatchScore:     80.0,
						MatchType:      "partial",
					})
					found = true
					break
				}
			}
			if !found {
				unmatchedKeywords = append(unmatchedKeywords, kw)
				missingSkills = append(missingSkills, models.StandardizedSkill{
					Keyword: kw,
					SkillType: "unknown",
					SkillName: kw,
					Confidence: 0.0,
				})
			}
		}
	}

	return &models.MatchResponse{
		Profile:        profile,
		MatchedSkills: matchedSkills,
		UnmatchedKeywords: unmatchedKeywords,
		MissingSkills: missingSkills,
	}, nil
}

// AddMissingSkill adds a skill to a profile after user confirmation
func (s *Service) AddMissingSkill(profileID int, skillName string, skillData *models.SkillUpdate) error {
	// Check if skill already exists
	exists, err := s.repo.CheckCustomSkillExists(profileID, skillName)
	if err != nil {
		return fmt.Errorf("failed to check skill: %w", err)
	}
	
	if exists {
		return fmt.Errorf("skill already exists")
	}

	// Create custom skill from update data
	skill := &models.CustomSkill{
		Name:           skillName,
		Category:       skillData.Category,
		Proficiency:    skillData.Proficiency,
		YearsExperience: skillData.Years,
	}

	// Add skill to profile
	if err := s.repo.AddCustomSkill(profileID, skill); err != nil {
		return fmt.Errorf("failed to add skill: %w", err)
	}

	return nil
}

// GenerateResume generates a simple text resume based on profile and matched keywords
func (s *Service) GenerateResume(profileID int, keywords []string, includeUnmatched bool) (string, error) {
	// Get profile
	profile, err := s.repo.GetProfile(profileID)
	if err != nil {
		return "", fmt.Errorf("failed to get profile: %w", err)
	}

	// Perform matching
	matchResult, err := s.MatchProfileToKeywords(profileID, keywords)
	if err != nil {
		return "", fmt.Errorf("failed to match: %w", err)
	}

	// Build resume content
	var sb strings.Builder

	// Header
	sb.WriteString(fmt.Sprintf("=%s\n\n", strings.Repeat("=", 50)))
	sb.WriteString(fmt.Sprintf("Name: %s\n", profile.Name))
	sb.WriteString(fmt.Sprintf("Email: %s\n", profile.Email))
	sb.WriteString("\n")

	// Experience
	if len(profile.Experience) > 0 {
		sb.WriteString("EXPERIENCE\n")
		sb.WriteString(strings.Repeat("-", 50) + "\n")
		for _, exp := range profile.Experience {
			sb.WriteString(fmt.Sprintf("\n%s at %s\n", exp.Position, exp.Company))
			sb.WriteString(fmt.Sprintf("  %s - %s\n", exp.StartDate, exp.EndDate))
			if exp.Description != "" {
				sb.WriteString(fmt.Sprintf("  %s\n", exp.Description))
			}
		}
		sb.WriteString("\n")
	}

	// Education
	if len(profile.Education) > 0 {
		sb.WriteString("EDUCATION\n")
		sb.WriteString(strings.Repeat("-", 50) + "\n")
		for _, edu := range profile.Education {
			sb.WriteString(fmt.Sprintf("\n%s in %s\n", edu.Degree, edu.FieldOfStudy))
			sb.WriteString(fmt.Sprintf("  %s, %s - %s\n", edu.Institution, edu.StartDate, edu.EndDate))
		}
		sb.WriteString("\n")
	}

	// Skills section
	sb.WriteString("SKILLS\n")
	sb.WriteString(strings.Repeat("-", 50) + "\n")

	// Matched skills
	if len(matchResult.MatchedSkills) > 0 {
		sb.WriteString("\nMatched Skills:\n")
		for _, skill := range matchResult.MatchedSkills {
			proficiency := skill.Proficiency
			if proficiency == "" {
				proficiency = "Proficient"
			}
			sb.WriteString(fmt.Sprintf("  - %s (%s, %d years)\n", skill.SkillName, proficiency, skill.YearsExperience))
		}
	}

	// Other custom skills from profile
	if len(profile.CustomSkills) > len(matchResult.MatchedSkills) {
		sb.WriteString("\nOther Skills:\n")
		matchedSkillNames := make(map[string]bool)
		for _, skill := range matchResult.MatchedSkills {
			matchedSkillNames[skill.SkillName] = true
		}
		for _, skill := range profile.CustomSkills {
			if !matchedSkillNames[skill.Name] {
				proficiency := skill.Proficiency
				if proficiency == "" {
					proficiency = "Proficient"
				}
				sb.WriteString(fmt.Sprintf("  - %s (%s, %d years)\n", skill.Name, proficiency, skill.YearsExperience))
			}
		}
	}

	// If includeUnmatched, show unmatched keywords as potential skills
	if includeUnmatched && len(matchResult.UnmatchedKeywords) > 0 {
		sb.WriteString("\nPotential Skills to Add:\n")
		for _, kw := range matchResult.UnmatchedKeywords {
			sb.WriteString(fmt.Sprintf("  - %s\n", kw))
		}
	}

	sb.WriteString("\n" + strings.Repeat("=", 50) + "\n")

	return sb.String(), nil
}

// ProcessJobOffer processes a job offer: extracts keywords, matches to profile, and generates resume
func (s *Service) ProcessJobOffer(profileID int, jobOffer *models.JobOffer) (*models.MatchResponse, string, error) {
	// Save job offer to SQLite
	if err := s.repo.SaveJobOffer(jobOffer); err != nil {
		return nil, "", fmt.Errorf("failed to save job offer: %w", err)
	}

	// Extract keywords from job offer
	keywords := s.extractor.ExtractTechnicalKeywords(jobOffer.Description + " " + jobOffer.Requirements)
	
	// Save keywords
	if err := s.repo.SaveKeywords(jobOffer.ID, keywords); err != nil {
		return nil, "", fmt.Errorf("failed to save keywords: %w", err)
	}

	// Match profile to keywords
	matchResult, err := s.MatchProfileToKeywords(profileID, keywords)
	if err != nil {
		return nil, "", fmt.Errorf("failed to match: %w", err)
	}

	// Generate resume
	resumeContent, err := s.GenerateResume(profileID, keywords, true)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate resume: %w", err)
	}

	return matchResult, resumeContent, nil
}

// SaveResume saves a resume
func (s *Service) SaveResume(resume *models.Resume) error {
	if err := s.repo.SaveResume(resume); err != nil {
		return fmt.Errorf("failed to save resume: %w", err)
	}
	return nil
}

// GetAllProfiles returns all profiles
func (s *Service) GetAllProfiles() ([]*models.Profile, error) {
	return s.repo.GetAllProfiles()
}

// GetProfileSkills returns all custom skills for a profile
func (s *Service) GetProfileSkills(profileID int) ([]models.CustomSkill, error) {
	return s.repo.GetProfileCustomSkills(profileID)
}

// GetAllSkills returns all skills from the SQL database
func (s *Service) GetAllSkills() ([]models.SkillWithDetails, error) {
	if s.sqlRepo == nil {
		// Return empty array if SQL repository is not configured
		return []models.SkillWithDetails{}, nil
	}
	
	// Create skill repository
	skillRepo := repository.NewSkillSQLRepository(s.sqlRepo.DB())
	skills, err := skillRepo.GetAllSkills()
	if err != nil {
		return nil, err
	}
	
	// Convert []*models.SkillWithDetails to []models.SkillWithDetails
	result := make([]models.SkillWithDetails, len(skills))
	for i, skill := range skills {
		if skill != nil {
			result[i] = *skill
		}
	}
	return result, nil
}

// CreateProfile creates a new profile in SQLite
func (s *Service) CreateProfile(profile *models.Profile) error {
	return s.repo.CreateProfile(profile)
}

// UpdateProfile updates a profile in SQLite
func (s *Service) UpdateProfile(profile *models.Profile) error {
	return s.repo.UpdateProfile(profile)
}
