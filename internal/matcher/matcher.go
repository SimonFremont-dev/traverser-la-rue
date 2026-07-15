package matcher

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/SimonFremont-dev/cv-builder/internal/models"
)

// Matcher handles matching resumes to job offers
type Matcher struct {
	Criteria models.MatchCriteria
}

// NewMatcher creates a new Matcher instance with default criteria
func NewMatcher() *Matcher {
	return &Matcher{
		Criteria: models.DefaultMatchCriteria(),
	}
}

// NewMatcherWithCriteria creates a new Matcher with custom criteria
func NewMatcherWithCriteria(criteria models.MatchCriteria) *Matcher {
	return &Matcher{
		Criteria: criteria,
	}
}

// Match performs the matching between a resume and a job offer
func (m *Matcher) Match(resume *models.Resume, jobOffer *models.JobOffer) *models.MatchResult {
	if resume == nil || jobOffer == nil {
		return nil
	}

	result := &models.MatchResult{
		JobOfferID:          jobOffer.Metadata.ID,
		ResumeID:            resume.Metadata.ID,
		SkillMatches:        make([]models.SkillMatch, 0),
		MissingRequirements: make([]string, 0),
		SuggestedChanges:    make([]models.SuggestedChange, 0),
	}

	// Match skills
	result.SkillMatches = m.matchSkills(resume, jobOffer)

	// Match experience
	result.ExperienceMatch = m.matchExperience(resume, jobOffer)

	// Match education
	result.EducationMatch = m.matchEducation(resume, jobOffer)

	// Identify missing requirements
	result.MissingRequirements = m.identifyMissingRequirements(resume, jobOffer, result)

	// Calculate overall score
	result.OverallScore = m.calculateOverallScore(result)

	// Assign grade
	result.MatchGrade = m.assignGrade(result.OverallScore)

	// Generate suggestions
	result.SuggestedChanges = m.generateSuggestions(resume, jobOffer, result)

	return result
}

// MatchWithAdaptation performs matching and returns an adapted resume
func (m *Matcher) MatchWithAdaptation(resume *models.Resume, jobOffer *models.JobOffer) *models.MatchResult {
	result := m.Match(resume, jobOffer)

	// Create adapted resume
	adaptedResume := m.adaptResume(resume, jobOffer, result)
	result.AdaptedResume = adaptedResume

	return result
}

// matchSkills matches resume skills with job requirements
func (m *Matcher) matchSkills(resume *models.Resume, jobOffer *models.JobOffer) []models.SkillMatch {
	var matches []models.SkillMatch

	// Combine all required and preferred skills from job offer
	allJobSkills := append(jobOffer.Requirements.RequiredSkills, jobOffer.Requirements.PreferredSkills...)

	for _, resumeSkill := range resume.Skills {
		for _, jobSkill := range allJobSkills {
			match := m.calculateSkillMatch(&resumeSkill, &jobSkill)
			if match.MatchScore > 0 {
				matches = append(matches, *match)
			}
		}
	}

	// Also check for keyword matches in descriptions
	for _, keyword := range jobOffer.Requirements.Keywords {
		keywordLower := strings.ToLower(keyword)

		// Check if keyword appears in resume
		for _, exp := range resume.Experience {
			if strings.Contains(strings.ToLower(exp.Description), keywordLower) ||
				strings.Contains(strings.ToLower(exp.Position), keywordLower) {

				// Create a pseudo-skill match
				match := &models.SkillMatch{
					JobRequirement: &models.SkillRequirement{
						Name:     keyword,
						Required: true,
					},
					MatchScore: 80.0, // High score for keyword match
					MatchType:  "Keyword",
				}
				matches = append(matches, *match)
				break
			}
		}
	}

	// Sort matches by score
	matches = models.SortSkillMatchesByScore(matches)

	return matches
}

// calculateSkillMatch calculates the match score between a resume skill and a job requirement
func (m *Matcher) calculateSkillMatch(resumeSkill *models.Skill, jobSkill *models.SkillRequirement) *models.SkillMatch {
	match := &models.SkillMatch{
		ResumeSkill:    resumeSkill,
		JobRequirement: jobSkill,
		MatchScore:     0.0,
		MatchType:      "None",
	}

	// Check for exact match
	if strings.EqualFold(resumeSkill.Name, jobSkill.Name) {
		match.MatchScore = 100.0
		match.MatchType = "Exact"

		// Bonus for proficiency match
		if m.proficiencyMatch(resumeSkill.Proficiency, jobSkill.Proficiency) {
			match.ProficiencyMatch = "Perfect"
		} else {
			match.ProficiencyMatch = "Partial"
			match.MatchScore *= 0.9
		}

		// Bonus for years of experience
		if resumeSkill.Years > 0 && jobSkill.Years != nil {
			if resumeSkill.Years >= *jobSkill.Years {
				match.YearsMatch = 100.0
			} else {
				match.YearsMatch = float64(resumeSkill.Years) / float64(*jobSkill.Years) * 100
				match.MatchScore *= 0.8
			}
		}

		return match
	}

	// Check for partial matches (substring or similar)
	similarity := m.calculateSimilarity(resumeSkill.Name, jobSkill.Name)
	if similarity >= 0.8 {
		match.MatchScore = similarity * 100
		match.MatchType = "Partial"
		return match
	}

	// Check for alternative matches
	for _, alt := range jobSkill.Alternatives {
		if strings.EqualFold(resumeSkill.Name, alt) {
			match.MatchScore = 90.0
			match.MatchType = "Alternative"
			return match
		}
	}

	// Check if resume skill name appears in job skill name or vice versa
	if strings.Contains(strings.ToLower(jobSkill.Name), strings.ToLower(resumeSkill.Name)) ||
		strings.Contains(strings.ToLower(resumeSkill.Name), strings.ToLower(jobSkill.Name)) {
		match.MatchScore = 70.0
		match.MatchType = "Substring"
		return match
	}

	// Check keyword matches
	for _, kw := range resumeSkill.Keywords {
		if strings.Contains(strings.ToLower(jobSkill.Name), strings.ToLower(kw)) {
			match.MatchScore = 60.0
			match.MatchType = "Keyword"
			return match
		}
	}

	return match
}

// proficiencyMatch checks if resume proficiency meets job requirement
func (m *Matcher) proficiencyMatch(resumeProf, jobProf string) bool {
	if resumeProf == "" || jobProf == "" {
		return true
	}

	// Define proficiency hierarchy
	proficiencyLevels := map[string]int{
		"Beginner":     1,
		"Intermediate": 2,
		"Advanced":     3,
		"Expert":       4,
		"Proficient":   3,
		"Fluent":       4,
		"Native":       5,
	}

	resumeLevel, ok1 := proficiencyLevels[strings.ToLower(resumeProf)]
	jobLevel, ok2 := proficiencyLevels[strings.ToLower(jobProf)]

	if !ok1 || !ok2 {
		return true
	}

	return resumeLevel >= jobLevel
}

// calculateSimilarity calculates string similarity (0-1)
func (m *Matcher) calculateSimilarity(s1, s2 string) float64 {
	// Use Levenshtein distance for similarity
	return 1.0 - float64(levenshteinDistance(s1, s2))/float64(max(len(s1), len(s2)))
}

// levenshteinDistance calculates the Levenshtein distance between two strings
func levenshteinDistance(s, t string) int {
	s = strings.ToLower(s)
	t = strings.ToLower(t)

	// Create a matrix to store distances
	d := make([][]int, len(s)+1)
	for i := range d {
		d[i] = make([]int, len(t)+1)
	}

	// Initialize the matrix
	for i := 0; i <= len(s); i++ {
		d[i][0] = i
	}
	for j := 0; j <= len(t); j++ {
		d[0][j] = j
	}

	// Fill in the matrix
	for i := 1; i <= len(s); i++ {
		for j := 1; j <= len(t); j++ {
			if s[i-1] == t[j-1] {
				d[i][j] = d[i-1][j-1]
			} else {
				d[i][j] = min(
					d[i-1][j]+1,   // Deletion
					d[i][j-1]+1,   // Insertion
					d[i-1][j-1]+1, // Substitution
				)
			}
		}
	}

	return d[len(s)][len(t)]
}

// min returns the minimum of integers
func min(nums ...int) int {
	if len(nums) == 0 {
		return 0
	}
	min := nums[0]
	for _, n := range nums[1:] {
		if n < min {
			min = n
		}
	}
	return min
}

// max returns the maximum of integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// matchExperience matches resume experience with job requirements
func (m *Matcher) matchExperience(resume *models.Resume, jobOffer *models.JobOffer) models.ExperienceMatch {
	match := models.ExperienceMatch{
		RelevantExperience: make([]models.RelevantExperienceItem, 0),
		MissingExperience:  make([]string, 0),
	}

	// Calculate total experience from resume
	totalYears := m.calculateTotalExperience(resume)

	if jobOffer.Requirements.Experience.TotalYears != nil {
		requiredYears := float64(*jobOffer.Requirements.Experience.TotalYears)
		if totalYears >= requiredYears {
			match.TotalYearsMatch = 100.0
		} else {
			match.TotalYearsMatch = (totalYears / requiredYears) * 100
		}
	}

	// Check for in-field experience
	if jobOffer.Requirements.Experience.InField != nil {
		inFieldYears := m.calculateFieldExperience(resume, jobOffer)
		requiredYears := float64(*jobOffer.Requirements.Experience.InField)
		if inFieldYears >= requiredYears {
			match.InFieldMatch = 100.0
		} else {
			match.InFieldMatch = (inFieldYears / requiredYears) * 100
		}
	}

	// Check for management experience
	if jobOffer.Requirements.Experience.Management {
		match.ManagementMatch = m.hasManagementExperience(resume)
	}

	// Find relevant experience items
	for _, exp := range resume.Experience {
		relevance := m.calculateExperienceRelevance(exp, jobOffer)
		if relevance > 0 {
			match.RelevantExperience = append(match.RelevantExperience, models.RelevantExperienceItem{
				WorkExperience:  &exp,
				RelevanceScore:  relevance,
				MatchedKeywords: m.findMatchedKeywords(exp, jobOffer),
			})
		}
	}

	// Sort by relevance
	sort.Slice(match.RelevantExperience, func(i, j int) bool {
		return match.RelevantExperience[i].RelevanceScore > match.RelevantExperience[j].RelevanceScore
	})

	// Identify missing experience
	for _, req := range jobOffer.Requirements.Experience.Specific {
		if !m.hasExperienceType(resume, req) {
			match.MissingExperience = append(match.MissingExperience, req)
		}
	}

	return match
}

// calculateTotalExperience calculates total work experience in years
func (m *Matcher) calculateTotalExperience(resume *models.Resume) float64 {
	totalMonths := 0

	for _, exp := range resume.Experience {
		if exp.Current {
			// Count from start date to now
			totalMonths += int(m.monthsBetween(exp.StartDate, time.Now()))
		} else if exp.EndDate != nil {
			totalMonths += int(m.monthsBetween(exp.StartDate, *exp.EndDate))
		}
	}

	return float64(totalMonths) / 12.0
}

// calculateFieldExperience calculates experience in the field relevant to the job
func (m *Matcher) calculateFieldExperience(resume *models.Resume, jobOffer *models.JobOffer) float64 {
	totalMonths := 0

	for _, exp := range resume.Experience {
		// Check if experience is relevant to the job
		if m.isRelevantExperience(exp, jobOffer) {
			if exp.Current {
				totalMonths += int(m.monthsBetween(exp.StartDate, time.Now()))
			} else if exp.EndDate != nil {
				totalMonths += int(m.monthsBetween(exp.StartDate, *exp.EndDate))
			}
		}
	}

	return float64(totalMonths) / 12.0
}

// hasManagementExperience checks if resume has management experience
func (m *Matcher) hasManagementExperience(resume *models.Resume) bool {
	for _, exp := range resume.Experience {
		lowerPosition := strings.ToLower(exp.Position)
		if strings.Contains(lowerPosition, "manager") ||
			strings.Contains(lowerPosition, "director") ||
			strings.Contains(lowerPosition, "lead") ||
			strings.Contains(lowerPosition, "head") {
			return true
		}
	}
	return false
}

// calculateExperienceRelevance calculates how relevant an experience is to the job
func (m *Matcher) calculateExperienceRelevance(exp models.WorkExperience, jobOffer *models.JobOffer) float64 {
	score := 0.0

	// Check position match
	positionScore := m.calculatePositionSimilarity(exp.Position, jobOffer.Position.Title)
	score += positionScore * 0.4

	// Check description match with job requirements
	if strings.Contains(strings.ToLower(exp.Description), strings.ToLower(jobOffer.Position.Title)) {
		score += 20.0
	}

	// Check for skill matches
	skillCount := 0
	for _, skill := range exp.Skills {
		for _, reqSkill := range jobOffer.Requirements.RequiredSkills {
			if strings.EqualFold(skill, reqSkill.Name) {
				skillCount++
				break
			}
		}
	}
	score += float64(skillCount) * 5.0

	// Check for keyword matches
	keywords := m.findMatchedKeywords(exp, jobOffer)
	score += float64(len(keywords)) * 2.0

	// Cap at 100
	if score > 100 {
		score = 100
	}

	return score
}

// calculatePositionSimilarity calculates similarity between position titles
func (m *Matcher) calculatePositionSimilarity(pos1, pos2 string) float64 {
	if pos1 == "" || pos2 == "" {
		return 0.0
	}

	similarity := m.calculateSimilarity(pos1, pos2)

	// Bonus for same words
	words1 := strings.Fields(strings.ToLower(pos1))
	words2 := strings.Fields(strings.ToLower(pos2))

	commonCount := 0
	for _, w1 := range words1 {
		for _, w2 := range words2 {
			if w1 == w2 {
				commonCount++
				break
			}
		}
	}

	// Calculate Jaccard similarity
	if len(words1) > 0 && len(words2) > 0 {
		union := len(words1) + len(words2) - commonCount
		jaccard := float64(commonCount) / float64(union)
		// Weighted average
		similarity = (similarity*0.6 + jaccard*0.4)
	}

	return similarity * 100
}

// findMatchedKeywords finds keywords from job offer that match the experience
func (m *Matcher) findMatchedKeywords(exp models.WorkExperience, jobOffer *models.JobOffer) []string {
	var keywords []string

	lowerDesc := strings.ToLower(exp.Description + " " + exp.Position)

	// Check against job keywords
	for _, kw := range jobOffer.Requirements.Keywords {
		if strings.Contains(lowerDesc, strings.ToLower(kw)) {
			keywords = append(keywords, kw)
		}
	}

	// Check against required skills
	for _, req := range jobOffer.Requirements.RequiredSkills {
		if strings.Contains(lowerDesc, strings.ToLower(req.Name)) {
			keywords = append(keywords, req.Name)
		}
	}

	return keywords
}

// hasExperienceType checks if resume has a specific type of experience
func (m *Matcher) hasExperienceType(resume *models.Resume, expType string) bool {
	lowerType := strings.ToLower(expType)

	for _, exp := range resume.Experience {
		if strings.Contains(strings.ToLower(exp.Position), lowerType) ||
			strings.Contains(strings.ToLower(exp.Description), lowerType) {
			return true
		}
	}

	return false
}

// matchEducation matches resume education with job requirements
func (m *Matcher) matchEducation(resume *models.Resume, jobOffer *models.JobOffer) models.EducationMatch {
	match := models.EducationMatch{
		RelevantEducation: make([]*models.Education, 0),
		MissingEducation:  make([]string, 0),
	}

	// Check degree match
	for _, edu := range resume.Education {
		for _, req := range jobOffer.Requirements.Education {
			if req.Degree == "" || strings.EqualFold(edu.Degree, req.Degree) {
				match.DegreeMatch = true
				match.RelevantEducation = append(match.RelevantEducation, &edu)
				break
			}
		}
	}

	// Check field of study match
	for _, edu := range resume.Education {
		for _, req := range jobOffer.Requirements.Education {
			if req.FieldOfStudy == "" || strings.EqualFold(edu.FieldOfStudy, req.FieldOfStudy) {
				match.FieldMatch = true
				break
			}
		}
	}

	// Check GPA if available
	for _, edu := range resume.Education {
		if edu.GPA != nil {
			// Simple check - in reality, you'd need more context
			if *edu.GPA >= 3.0 {
				match.GPAMeet = true
				break
			}
		}
	}

	// Identify missing education
	for _, req := range jobOffer.Requirements.Education {
		if req.Required {
			found := false
			for _, edu := range resume.Education {
				if strings.EqualFold(edu.Degree, req.Degree) {
					found = true
					break
				}
			}
			if !found {
				match.MissingEducation = append(match.MissingEducation, fmt.Sprintf("%s degree in %s", req.Degree, req.FieldOfStudy))
			}
		}
	}

	return match
}

// identifyMissingRequirements identifies requirements that are not met
func (m *Matcher) identifyMissingRequirements(resume *models.Resume, jobOffer *models.JobOffer, result *models.MatchResult) []string {
	var missing []string

	// Check required skills
	for _, reqSkill := range jobOffer.Requirements.RequiredSkills {
		found := false
		for _, match := range result.SkillMatches {
			if match.JobRequirement != nil && match.JobRequirement.Name == reqSkill.Name && match.MatchScore >= 80 {
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, fmt.Sprintf("Required skill: %s", reqSkill.Name))
		}
	}

	// Check experience requirements
	if jobOffer.Requirements.Experience.TotalYears != nil {
		totalYears := m.calculateTotalExperience(resume)
		requiredYears := float64(*jobOffer.Requirements.Experience.TotalYears)
		if totalYears < requiredYears {
			missing = append(missing, fmt.Sprintf("Total experience: %d+ years", *jobOffer.Requirements.Experience.TotalYears))
		}
	}

	// Check education requirements
	for _, req := range jobOffer.Requirements.Education {
		if req.Required {
			found := false
			for _, edu := range resume.Education {
				if strings.EqualFold(edu.Degree, req.Degree) {
					found = true
					break
				}
			}
			if !found {
				missing = append(missing, fmt.Sprintf("Education: %s", req.Degree))
			}
		}
	}

	// Check certifications
	for _, reqCert := range jobOffer.Requirements.Certifications {
		found := false
		for _, cert := range resume.Certifications {
			if strings.EqualFold(cert.Name, reqCert) {
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, fmt.Sprintf("Certification: %s", reqCert))
		}
	}

	// Check languages
	for _, reqLang := range jobOffer.Requirements.Languages {
		if reqLang.Required {
			found := false
			for _, lang := range resume.Languages {
				if strings.EqualFold(lang.Language, reqLang.Language) {
					found = true
					break
				}
			}
			if !found {
				missing = append(missing, fmt.Sprintf("Language: %s", reqLang.Language))
			}
		}
	}

	return missing
}

// calculateOverallScore calculates the overall match score
func (m *Matcher) calculateOverallScore(result *models.MatchResult) float64 {
	if result == nil {
		return 0.0
	}

	// Calculate skill score (average of top matches)
	skillScore := 0.0
	if len(result.SkillMatches) > 0 {
		total := 0.0
		for _, match := range result.SkillMatches {
			total += match.MatchScore
		}
		skillScore = total / float64(len(result.SkillMatches))
	}

	// Calculate experience score
	experienceScore := (result.ExperienceMatch.TotalYearsMatch +
		result.ExperienceMatch.InFieldMatch) / 2.0
	if result.ExperienceMatch.ManagementMatch {
		experienceScore = experienceScore * 1.1
	}
	if experienceScore > 100 {
		experienceScore = 100
	}

	// Calculate education score
	educationScore := 0.0
	if result.EducationMatch.DegreeMatch {
		educationScore += 50
	}
	if result.EducationMatch.FieldMatch {
		educationScore += 50
	}
	if educationScore > 100 {
		educationScore = 100
	}

	// Weighted average
	weightedScore := (skillScore*m.Criteria.SkillWeight +
		experienceScore*m.Criteria.ExperienceWeight +
		educationScore*m.Criteria.EducationWeight)

	// Penalty for missing requirements
	missingPenalty := float64(len(result.MissingRequirements)) * 5.0
	if missingPenalty > 30 {
		missingPenalty = 30
	}

	finalScore := weightedScore - missingPenalty
	if finalScore < 0 {
		finalScore = 0
	}
	if finalScore > 100 {
		finalScore = 100
	}

	return math.Round(finalScore*100) / 100
}

// assignGrade assigns a letter grade based on the score
func (m *Matcher) assignGrade(score float64) string {
	if score >= 90 {
		return "A"
	} else if score >= 80 {
		return "B"
	} else if score >= 70 {
		return "C"
	} else if score >= 60 {
		return "D"
	} else {
		return "F"
	}
}

// generateSuggestions generates suggestions to improve the resume
func (m *Matcher) generateSuggestions(resume *models.Resume, jobOffer *models.JobOffer, result *models.MatchResult) []models.SuggestedChange {
	var suggestions []models.SuggestedChange

	// Suggest adding missing required skills
	for _, missing := range result.MissingRequirements {
		if strings.HasPrefix(missing, "Required skill:") {
			skill := strings.TrimPrefix(missing, "Required skill: ")
			suggestions = append(suggestions, models.SuggestedChange{
				Type:          "AddSkill",
				Priority:      "High",
				Section:       "skills",
				Reason:        fmt.Sprintf("Job requires %s which is missing from your resume", skill),
				SuggestedText: skill,
			})
		}
	}

	// Suggest emphasizing matched skills
	for _, match := range result.SkillMatches {
		if match.MatchScore >= 80 && match.ResumeSkill != nil {
			// Check if skill is mentioned in experience
			foundInExperience := false
			for _, exp := range resume.Experience {
				if containsSkill(exp.Skills, match.ResumeSkill.Name) {
					foundInExperience = true
					break
				}
			}

			if !foundInExperience {
				suggestions = append(suggestions, models.SuggestedChange{
					Type:     "Emphasize",
					Priority: "Medium",
					Section:  "experience",
					ItemID:   match.ResumeSkill.ID,
					Reason:   fmt.Sprintf("Skill %s matches job requirement but could be emphasized in experience", match.ResumeSkill.Name),
				})
			}
		}
	}

	// Suggest reordering to put most relevant experience first
	if len(result.ExperienceMatch.RelevantExperience) > 0 {
		suggestions = append(suggestions, models.SuggestedChange{
			Type:     "Reorder",
			Priority: "Medium",
			Section:  "experience",
			Reason:   "Move most relevant experience to the top of your resume",
		})
	}

	// Sort suggestions by priority
	suggestions = models.SortSuggestedChangesByPriority(suggestions)

	return suggestions
}

// adaptResume creates an adapted version of the resume based on the job offer
func (m *Matcher) adaptResume(resume *models.Resume, jobOffer *models.JobOffer, result *models.MatchResult) *models.Resume {
	// Create a copy of the resume
	adapted := *resume
	adapted.Metadata.Version = fmt.Sprintf("%s-adapted", resume.Metadata.Version)
	adapted.Metadata.Title = fmt.Sprintf("%s - Adapted for %s", resume.Metadata.Title, jobOffer.Position.Title)

	// Adapt summary to mention key requirements
	if len(result.SkillMatches) > 0 {
		keySkills := make([]string, 0, min(5, len(result.SkillMatches)))
		for i := 0; i < min(5, len(result.SkillMatches)); i++ {
			if result.SkillMatches[i].ResumeSkill != nil {
				keySkills = append(keySkills, result.SkillMatches[i].ResumeSkill.Name)
			}
		}
		if len(keySkills) > 0 {
			adapted.Summary = fmt.Sprintf("%s With expertise in %s.", resume.Summary, strings.Join(keySkills, ", "))
		}
	}

	// Reorder skills to put matched ones first
	m.reorderSkills(&adapted, result)

	// Reorder experience to put relevant ones first
	m.reorderExperience(&adapted, result)

	// Add keywords to skills if missing
	m.addMissingKeywords(&adapted, jobOffer)

	return &adapted
}

// reorderSkills reorders skills to prioritize matched ones
func (m *Matcher) reorderSkills(resume *models.Resume, result *models.MatchResult) {
	if len(result.SkillMatches) == 0 {
		return
	}

	// Create a map of matched skill IDs
	matchedIDs := make(map[string]bool)
	for _, match := range result.SkillMatches {
		if match.ResumeSkill != nil {
			matchedIDs[match.ResumeSkill.ID] = true
		}
	}

	// Sort skills: matched first, then by name
	sort.SliceStable(resume.Skills, func(i, j int) bool {
		iMatched := matchedIDs[resume.Skills[i].ID]
		jMatched := matchedIDs[resume.Skills[j].ID]

		if iMatched != jMatched {
			return iMatched
		}
		return resume.Skills[i].Name < resume.Skills[j].Name
	})
}

// reorderExperience reorders experience to prioritize relevant ones
func (m *Matcher) reorderExperience(resume *models.Resume, result *models.MatchResult) {
	if len(result.ExperienceMatch.RelevantExperience) == 0 {
		return
	}

	// Create a map of experience relevance
	relevanceMap := make(map[string]float64)
	for _, item := range result.ExperienceMatch.RelevantExperience {
		if item.WorkExperience != nil {
			relevanceMap[item.WorkExperience.ID] = item.RelevanceScore
		}
	}

	// Sort experience by relevance (descending), then by date (descending)
	sort.SliceStable(resume.Experience, func(i, j int) bool {
		iRel := relevanceMap[resume.Experience[i].ID]
		jRel := relevanceMap[resume.Experience[j].ID]

		if iRel != jRel {
			return iRel > jRel
		}
		// If same relevance, sort by date (newer first)
		return resume.Experience[i].StartDate.After(resume.Experience[j].StartDate)
	})
}

// addMissingKeywords adds missing keywords to the resume
func (m *Matcher) addMissingKeywords(resume *models.Resume, jobOffer *models.JobOffer) {
	// Check for keywords in job offer that are not in resume
	for _, keyword := range jobOffer.Requirements.Keywords {
		keywordLower := strings.ToLower(keyword)
		found := false

		// Check if keyword exists in resume
		for _, skill := range resume.Skills {
			if strings.EqualFold(skill.Name, keyword) {
				found = true
				break
			}
		}

		for _, exp := range resume.Experience {
			if strings.Contains(strings.ToLower(exp.Description), keywordLower) ||
				strings.Contains(strings.ToLower(exp.Position), keywordLower) {
				found = true
				break
			}
		}

		if !found {
			// Add as a skill
			resume.Skills = append(resume.Skills, models.Skill{
				ID:       fmt.Sprintf("kw-%d", len(resume.Skills)),
				Name:     keyword,
				Category: "Other",
			})
		}
	}
}

// containsSkill checks if a slice contains a skill
func containsSkill(skills []string, skill string) bool {
	for _, s := range skills {
		if strings.EqualFold(s, skill) {
			return true
		}
	}
	return false
}

// monthsBetween calculates the number of months between two dates
func (m *Matcher) monthsBetween(start, end time.Time) float64 {
	return float64((end.Year()-start.Year())*12 + int(end.Month()) - int(start.Month()))
}

// isRelevantExperience checks if an experience is relevant to the job
func (m *Matcher) isRelevantExperience(exp models.WorkExperience, jobOffer *models.JobOffer) bool {
	// Check if any skill matches
	for _, skill := range exp.Skills {
		for _, reqSkill := range jobOffer.Requirements.RequiredSkills {
			if strings.EqualFold(skill, reqSkill.Name) {
				return true
			}
		}
		for _, reqSkill := range jobOffer.Requirements.PreferredSkills {
			if strings.EqualFold(skill, reqSkill.Name) {
				return true
			}
		}
	}

	// Check if description contains keywords
	lowerDesc := strings.ToLower(exp.Description + " " + exp.Position)
	for _, keyword := range jobOffer.Requirements.Keywords {
		if strings.Contains(lowerDesc, strings.ToLower(keyword)) {
			return true
		}
	}

	// Check position similarity
	positionScore := m.calculatePositionSimilarity(exp.Position, jobOffer.Position.Title)
	return positionScore > 50
}
