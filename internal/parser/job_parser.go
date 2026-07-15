package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/SimonFremont-dev/cv-builder/internal/models"
	"gopkg.in/yaml.v3"
)

// JobParser handles parsing of job offer files
type JobParser struct {
	// Configuration for parsing
	ExtractKeywords bool
	StrictMode      bool
}

// NewJobParser creates a new JobParser instance
func NewJobParser() *JobParser {
	return &JobParser{
		ExtractKeywords: true,
		StrictMode:      false,
	}
}

// ParseJobOffer parses a job offer from a file based on its extension
func (p *JobParser) ParseJobOffer(filePath string) (*models.JobOffer, error) {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".json":
		return p.ParseJobOfferJSON(filePath)
	case ".yaml", ".yml":
		return p.ParseJobOfferYAML(filePath)
	case ".txt", ".text":
		return p.ParseJobOfferText(filePath)
	case ".html", ".htm":
		return p.ParseJobOfferHTML(filePath)
	case ".pdf":
		return p.ParseJobOfferPDF(filePath)
	default:
		// Try to detect format from content
		return p.ParseJobOfferAuto(filePath)
	}
}

// ParseJobOfferJSON parses a job offer from a JSON file
func (p *JobParser) ParseJobOfferJSON(filePath string) (*models.JobOffer, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON job offer file: %w", err)
	}

	var jobOffer models.JobOffer
	jobOffer.Metadata.Format = "json"
	jobOffer.Source = filePath
	jobOffer.RawContent = string(data)

	if err := json.Unmarshal(data, &jobOffer); err != nil {
		return nil, fmt.Errorf("failed to parse JSON job offer: %w", err)
	}

	// Extract keywords if enabled
	if p.ExtractKeywords {
		p.extractKeywords(&jobOffer)
	}

	return &jobOffer, nil
}

// ParseJobOfferYAML parses a job offer from a YAML file
func (p *JobParser) ParseJobOfferYAML(filePath string) (*models.JobOffer, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML job offer file: %w", err)
	}

	var jobOffer models.JobOffer
	jobOffer.Metadata.Format = "yaml"
	jobOffer.Source = filePath
	jobOffer.RawContent = string(data)

	if err := yaml.Unmarshal(data, &jobOffer); err != nil {
		return nil, fmt.Errorf("failed to parse YAML job offer: %w", err)
	}

	// Extract keywords if enabled
	if p.ExtractKeywords {
		p.extractKeywords(&jobOffer)
	}

	return &jobOffer, nil
}

// ParseJobOfferText parses a job offer from a plain text file
func (p *JobParser) ParseJobOfferText(filePath string) (*models.JobOffer, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read text job offer file: %w", err)
	}

	content := string(data)
	jobOffer := p.parseTextContent(content)
	jobOffer.Metadata.Format = "text"
	jobOffer.Source = filePath
	jobOffer.RawContent = content

	// Extract keywords if enabled
	if p.ExtractKeywords {
		p.extractKeywords(jobOffer)
	}

	return jobOffer, nil
}

// ParseJobOfferHTML parses a job offer from an HTML file
// This is a simplified parser - in production you'd use a proper HTML parser
func (p *JobParser) ParseJobOfferHTML(filePath string) (*models.JobOffer, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read HTML job offer file: %w", err)
	}

	content := string(data)
	// Remove HTML tags
	textContent := stripHTMLTags(content)

	jobOffer := p.parseTextContent(textContent)
	jobOffer.Metadata.Format = "html"
	jobOffer.Source = filePath
	jobOffer.RawContent = content

	// Extract keywords if enabled
	if p.ExtractKeywords {
		p.extractKeywords(jobOffer)
	}

	return jobOffer, nil
}

// ParseJobOfferPDF parses a job offer from a PDF file
// Note: This requires the pdfcpu library
func (p *JobParser) ParseJobOfferPDF(filePath string) (*models.JobOffer, error) {
	// For now, we'll just read the file and treat it as text
	// In a production system, you'd use a PDF parsing library
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read PDF job offer file: %w", err)
	}

	content := string(data)
	jobOffer := p.parseTextContent(content)
	jobOffer.Metadata.Format = "pdf"
	jobOffer.Source = filePath
	jobOffer.RawContent = content

	// Extract keywords if enabled
	if p.ExtractKeywords {
		p.extractKeywords(jobOffer)
	}

	return jobOffer, nil
}

// ParseJobOfferAuto tries to automatically detect and parse the job offer format
func (p *JobParser) ParseJobOfferAuto(filePath string) (*models.JobOffer, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read job offer file: %w", err)
	}

	content := string(data)

	// Try JSON first
	var jobOfferJSON models.JobOffer
	if err := json.Unmarshal(data, &jobOfferJSON); err == nil {
		jobOfferJSON.Metadata.Format = "json"
		jobOfferJSON.Source = filePath
		jobOfferJSON.RawContent = content
		if p.ExtractKeywords {
			p.extractKeywords(&jobOfferJSON)
		}
		return &jobOfferJSON, nil
	}

	// Try YAML
	var jobOfferYAML models.JobOffer
	if err := yaml.Unmarshal(data, &jobOfferYAML); err == nil {
		jobOfferYAML.Metadata.Format = "yaml"
		jobOfferYAML.Source = filePath
		jobOfferYAML.RawContent = content
		if p.ExtractKeywords {
			p.extractKeywords(&jobOfferYAML)
		}
		return &jobOfferYAML, nil
	}

	// Fall back to text parsing
	return p.ParseJobOfferText(filePath)
}

// ParseJobOfferFromString parses a job offer from a string based on detected format
func (p *JobParser) ParseJobOfferFromString(content string, format string) (*models.JobOffer, error) {
	format = strings.ToLower(format)

	var jobOffer *models.JobOffer

	switch format {
	case "json", "":
		var jo models.JobOffer
		if err := json.Unmarshal([]byte(content), &jo); err != nil {
			// If JSON fails, try YAML
			if err := yaml.Unmarshal([]byte(content), &jo); err != nil {
				// Fall back to text parsing
				jobOffer = p.parseTextContent(content)
			} else {
				jobOffer = &jo
			}
		} else {
			jobOffer = &jo
		}
	case "yaml", "yml":
		var jo models.JobOffer
		if err := yaml.Unmarshal([]byte(content), &jo); err != nil {
			return nil, fmt.Errorf("failed to parse YAML job offer: %w", err)
		}
		jobOffer = &jo
	case "text", "txt":
		jobOffer = p.parseTextContent(content)
	case "html":
		textContent := stripHTMLTags(content)
		jobOffer = p.parseTextContent(textContent)
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}

	jobOffer.Metadata.Format = format
	jobOffer.RawContent = content

	// Extract keywords if enabled
	if p.ExtractKeywords && jobOffer != nil {
		p.extractKeywords(jobOffer)
	}

	return jobOffer, nil
}

// parseTextContent parses text content into a JobOffer structure
func (p *JobParser) parseTextContent(content string) *models.JobOffer {
	jobOffer := models.NewJobOffer()
	jobOffer.Description = content

	// Extract company info
	if company := extractCompany(content); company != "" {
		jobOffer.Company.Name = company
	}

	// Extract position title
	if title := extractPositionTitle(content); title != "" {
		jobOffer.Position.Title = title
		jobOffer.Metadata.Title = title
	}

	// Extract location
	if location := extractLocation(content); location != "" {
		jobOffer.Position.Location = location
		jobOffer.Company.Location = location
	}

	// Extract salary information
	if salary := extractSalary(content); salary != nil {
		jobOffer.Position.Salary = salary
	}

	// Extract requirements
	requirements := extractRequirements(content)
	if requirements != nil {
		jobOffer.Requirements = *requirements
	}

	// Extract keywords
	keywords := extractKeywordsFromText(content)
	if len(keywords) > 0 {
		jobOffer.Requirements.Keywords = keywords
	}

	return jobOffer
}

// extractCompany extracts company name from content
func extractCompany(content string) string {
	// Simple pattern matching for company name
	patterns := []string{
		"Company:", "Company Name:", "Employer:",
		"About us:", "About the company:",
	}

	for _, pattern := range patterns {
		if idx := strings.Index(strings.ToUpper(content), strings.ToUpper(pattern)); idx != -1 {
			start := idx + len(pattern)
			end := strings.Index(content[start:], "\n")
			if end == -1 {
				end = len(content)
			} else {
				end += start
			}
			return strings.TrimSpace(content[start:end])
		}
	}

	// Try to find common company indicators
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.Contains(strings.ToUpper(line), "COMPANY") ||
			strings.Contains(strings.ToUpper(line), "EMPLOYER") {
			return strings.TrimSpace(line)
		}
	}

	return ""
}

// extractPositionTitle extracts job position title from content
func extractPositionTitle(content string) string {
	patterns := []string{
		"Position:", "Title:", "Job Title:", "Role:",
		"We are looking for a", "We need a", "Hiring ",
	}

	for _, pattern := range patterns {
		if idx := strings.Index(strings.ToUpper(content), strings.ToUpper(pattern)); idx != -1 {
			start := idx + len(pattern)
			end := strings.Index(content[start:], "\n")
			if end == -1 {
				end = len(content)
			} else {
				end += start
			}
			return strings.TrimSpace(content[start:end])
		}
	}

	// Check first line for title
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) > 10 && len(trimmed) < 100 &&
			!strings.Contains(strings.ToUpper(trimmed), "COMPANY") &&
			!strings.Contains(strings.ToUpper(trimmed), "ABOUT") {
			return trimmed
		}
	}

	return ""
}

// extractLocation extracts location from content
func extractLocation(content string) string {
	patterns := []string{
		"Location:", "Where:", "Based in",
		"Remote", "On-site", "Hybrid",
	}

	for _, pattern := range patterns {
		if idx := strings.Index(strings.ToUpper(content), strings.ToUpper(pattern)); idx != -1 {
			start := idx + len(pattern)
			end := strings.Index(content[start:], "\n")
			if end == -1 {
				end = len(content)
			} else {
				end += start
			}
			return strings.TrimSpace(content[start:end])
		}
	}

	return ""
}

// extractSalary extracts salary information from content
func extractSalary(content string) *models.Salary {
	salary := &models.Salary{}

	// Pattern to match salary ranges
	re := regexp.MustCompile(`\$?\d{1,3}[,\d]{0,3}[\s]?[\-]\s?\$?\d{1,3}[,\d]{0,3}`)
	matches := re.FindAllString(content, -1)

	if len(matches) > 0 {
		// Parse first match
		salaryRange := strings.ReplaceAll(matches[0], "$", "")
		salaryRange = strings.ReplaceAll(salaryRange, ",", "")
		salaryRange = strings.ReplaceAll(strings.ToLower(salaryRange), "to", "-")
		salaryRange = strings.ReplaceAll(salaryRange, " ", "")

		// Extract min and max
		parts := strings.Split(salaryRange, "-")

		if len(parts) >= 2 {
			var min, max float64
			if _, err := fmt.Sscanf(parts[0], "%f", &min); err == nil {
				salary.Min = &min
			}
			if _, err := fmt.Sscanf(parts[len(parts)-1], "%f", &max); err == nil {
				salary.Max = &max
			}
		}

		// Try to detect currency
		if strings.Contains(content, "USD") || strings.Contains(content, "$") {
			salary.Currency = "USD"
		} else if strings.Contains(content, "EUR") || strings.Contains(content, "€") {
			salary.Currency = "EUR"
		} else if strings.Contains(content, "GBP") || strings.Contains(content, "£") {
			salary.Currency = "GBP"
		}

		// Try to detect period
		if strings.Contains(strings.ToLower(content), "year") || strings.Contains(strings.ToLower(content), "annual") {
			salary.Period = "yearly"
		} else if strings.Contains(strings.ToLower(content), "month") {
			salary.Period = "monthly"
		} else if strings.Contains(strings.ToLower(content), "hour") {
			salary.Period = "hourly"
		}
	}

	return salary
}

// extractRequirements extracts requirements from content
func extractRequirements(content string) *models.Requirements {
	reqs := models.NewJobOffer().Requirements

	// Extract required skills
	if skills := extractSkills(content, "required"); len(skills) > 0 {
		for _, skill := range skills {
			reqs.RequiredSkills = append(reqs.RequiredSkills, models.SkillRequirement{
				Name:     skill,
				Required: true,
			})
		}
	}

	// Extract preferred skills
	if skills := extractSkills(content, "preferred"); len(skills) > 0 {
		for _, skill := range skills {
			reqs.PreferredSkills = append(reqs.PreferredSkills, models.SkillRequirement{
				Name:       skill,
				Required:   false,
				NiceToHave: true,
			})
		}
	}

	// Extract experience requirements
	if exp := extractExperienceRequirements(content); exp != nil {
		reqs.Experience = *exp
	}

	// Extract education requirements
	if edu := extractEducationRequirements(content); len(edu) > 0 {
		reqs.Education = edu
	}

	return &reqs
}

// extractSkills extracts skills from content based on context
func extractSkills(content string, context string) []string {
	var skills []string

	// Look for sections containing the context
	lowerContent := strings.ToLower(content)
	contextIdx := strings.Index(lowerContent, context)

	if contextIdx == -1 {
		return skills
	}

	// Extract the section
	start := contextIdx
	end := strings.Index(lowerContent[start+len(context):], "\n\n")
	if end == -1 {
		end = len(lowerContent)
	} else {
		end += start + len(context)
	}

	section := content[start:end]

	// Simple skill extraction - look for words that look like skills
	// This would be enhanced with NLP in a production system
	words := strings.FieldsFunc(section, func(r rune) bool {
		return r == ' ' || r == ',' || r == ';' || r == '\n' || r == '\t' || r == '.'
	})

	for _, word := range words {
		cleaned := strings.TrimSpace(strings.Trim(word, "-,;."))
		if len(cleaned) > 2 && isValidSkill(cleaned) {
			skills = append(skills, cleaned)
		}
	}

	return skills
}

// isValidSkill checks if a word looks like a valid skill
func isValidSkill(word string) bool {
	// Simple check - skills are usually capitalized or contain hyphens
	if strings.Contains(word, "-") || strings.ToUpper(word) == word {
		return true
	}

	// Check for common skill patterns
	commonSkills := []string{
		"Go", "Golang", "Python", "Java", "C++", "JavaScript", "TypeScript",
		"React", "Angular", "Vue", "Node", "Docker", "Kubernetes",
		"AWS", "Azure", "GCP", "SQL", "NoSQL", "REST", "GraphQL",
		"Agile", "Scrum", "CI/CD", "DevOps", "Machine Learning",
	}

	for _, skill := range commonSkills {
		if strings.EqualFold(word, skill) {
			return true
		}
	}

	return len(word) > 3
}

// extractExperienceRequirements extracts experience requirements from content
func extractExperienceRequirements(content string) *models.ExperienceRequirement {
	req := &models.ExperienceRequirement{}

	// Look for experience patterns
	// Extract total years
	re := regexp.MustCompile(`(?i)(\d+)\s*year`)
	matches := re.FindAllStringSubmatch(content, -1)
	if len(matches) > 0 && len(matches[0]) > 1 {
		var years int
		if _, err := fmt.Sscanf(matches[0][1], "%d", &years); err == nil {
			req.TotalYears = &years
		}
	}

	re = regexp.MustCompile(`(?i)(\d+)\s*years`)
	matches = re.FindAllStringSubmatch(content, -1)
	if len(matches) > 0 && len(matches[0]) > 1 {
		var years int
		if _, err := fmt.Sscanf(matches[0][1], "%d", &years); err == nil {
			req.TotalYears = &years
		}
	}

	re = regexp.MustCompile(`(?i)(\d+)\+\s*year`)
	matches = re.FindAllStringSubmatch(content, -1)
	if len(matches) > 0 && len(matches[0]) > 1 {
		var years int
		if _, err := fmt.Sscanf(matches[0][1], "%d", &years); err == nil {
			req.MinYears = &years
		}
	}

	return req
}

// extractEducationRequirements extracts education requirements from content
func extractEducationRequirements(content string) []models.EducationRequirement {
	var reqs []models.EducationRequirement

	// Look for degree patterns
	degrees := []string{"Bachelor", "Master", "PhD", "Doctorate", "MBA", "Associate"}

	for _, degree := range degrees {
		if strings.Contains(strings.ToUpper(content), strings.ToUpper(degree)) {
			reqs = append(reqs, models.EducationRequirement{
				Degree:   degree,
				Required: true,
			})
		}
	}

	return reqs
}

// extractKeywordsFromText extracts keywords from text content
func extractKeywordsFromText(content string) []string {
	var keywords []string

	// Common keywords to look for
	keywordPatterns := []string{
		"required", "must have", "experience with",
		"knowledge of", "familiar with", "expert in",
		"proficient in", "skilled in", "ability to",
	}

	// Look for words that follow these patterns
	for _, pattern := range keywordPatterns {
		lowerContent := strings.ToLower(content)
		lowerPattern := strings.ToLower(pattern)

		idx := 0
		for {
			pos := strings.Index(lowerContent[idx:], lowerPattern)
			if pos == -1 {
				break
			}

			start := idx + pos + len(lowerPattern)
			// Find the next word
			wordStart := strings.IndexFunc(content[start:], func(r rune) bool {
				return r != ' ' && r != '\t' && r != '\n'
			})

			if wordStart == -1 {
				break
			}

			wordEnd := strings.IndexFunc(content[start+wordStart:], func(r rune) bool {
				return r == ' ' || r == '\t' || r == '\n' || r == ',' || r == '.' || r == ';'
			})

			if wordEnd == -1 {
				wordEnd = len(content)
			} else {
				wordEnd += start + wordStart
			}

			word := strings.TrimSpace(content[start+wordStart : wordEnd])
			if len(word) > 2 {
				keywords = append(keywords, word)
			}

			idx = start + wordStart + len(word)
		}
	}

	return keywords
}

// extractKeywords adds keywords to the job offer's requirements
func (p *JobParser) extractKeywords(jobOffer *models.JobOffer) {
	if jobOffer == nil {
		return
	}

	// Extract from description
	if jobOffer.Description != "" {
		keywords := extractKeywordsFromText(jobOffer.Description)
		for _, kw := range keywords {
			found := false
			for _, existing := range jobOffer.Requirements.Keywords {
				if strings.EqualFold(existing, kw) {
					found = true
					break
				}
			}
			if !found {
				jobOffer.Requirements.Keywords = append(jobOffer.Requirements.Keywords, kw)
			}
		}
	}
}

// stripHTMLTags removes HTML tags from content
func stripHTMLTags(content string) string {
	re := regexp.MustCompile("<[^>]*>")
	return re.ReplaceAllString(content, " ")
}
