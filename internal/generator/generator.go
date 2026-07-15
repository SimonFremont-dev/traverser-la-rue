package generator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/SimonFremont-dev/cv-builder/internal/models"
	"gopkg.in/yaml.v3"
)

// Generator handles generation of resume output in various formats
type Generator struct {
	TemplateDir string
}

// NewGenerator creates a new Generator instance
func NewGenerator() *Generator {
	return &Generator{
		TemplateDir: "templates",
	}
}

// NewGeneratorWithTemplateDir creates a new Generator with a custom template directory
func NewGeneratorWithTemplateDir(templateDir string) *Generator {
	return &Generator{
		TemplateDir: templateDir,
	}
}

// GenerateResume generates a resume in the specified format
func (g *Generator) GenerateResume(resume *models.Resume, format string, outputPath string) error {
	format = strings.ToLower(format)

	switch format {
	case "json":
		return g.generateJSON(resume, outputPath)
	case "yaml", "yml":
		return g.generateYAML(resume, outputPath)
	case "md", "markdown":
		return g.generateMarkdown(resume, outputPath)
	case "txt", "text":
		return g.generateText(resume, outputPath)
	case "html":
		return g.generateHTML(resume, outputPath)
	case "docx":
		docxGen := NewDOCXGenerator()
		return docxGen.GenerateDOCX(resume, outputPath)
	case "pdf":
		pdfGen := NewPDFGenerator()
		return pdfGen.GeneratePDF(resume, outputPath)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

// generateJSON generates a resume in JSON format
func (g *Generator) generateJSON(resume *models.Resume, outputPath string) error {
	data, err := json.MarshalIndent(resume, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal resume to JSON: %w", err)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write JSON resume: %w", err)
	}

	return nil
}

// generateYAML generates a resume in YAML format
func (g *Generator) generateYAML(resume *models.Resume, outputPath string) error {
	data, err := yaml.Marshal(resume)
	if err != nil {
		return fmt.Errorf("failed to marshal resume to YAML: %w", err)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write YAML resume: %w", err)
	}

	return nil
}

// generateMarkdown generates a resume in Markdown format
func (g *Generator) generateMarkdown(resume *models.Resume, outputPath string) error {
	markdown := g.resumeToMarkdown(resume)

	if err := os.WriteFile(outputPath, []byte(markdown), 0644); err != nil {
		return fmt.Errorf("failed to write Markdown resume: %w", err)
	}

	return nil
}

// generateText generates a resume in plain text format
func (g *Generator) generateText(resume *models.Resume, outputPath string) error {
	text := g.resumeToText(resume)

	if err := os.WriteFile(outputPath, []byte(text), 0644); err != nil {
		return fmt.Errorf("failed to write text resume: %w", err)
	}

	return nil
}

// generateHTML generates a resume in HTML format
func (g *Generator) generateHTML(resume *models.Resume, outputPath string) error {
	// Try to use template if available
	if g.TemplateDir != "" {
		templatePath := filepath.Join(g.TemplateDir, "resume.html")
		if _, err := os.Stat(templatePath); err == nil {
			return g.generateHTMLWithTemplate(resume, templatePath, outputPath)
		}
	}

	// Fall back to default HTML generation
	html := g.resumeToHTML(resume)

	if err := os.WriteFile(outputPath, []byte(html), 0644); err != nil {
		return fmt.Errorf("failed to write HTML resume: %w", err)
	}

	return nil
}

// generateHTMLWithTemplate generates HTML using a template file
func (g *Generator) generateHTMLWithTemplate(resume *models.Resume, templatePath, outputPath string) error {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("failed to parse HTML template: %w", err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, resume); err != nil {
		return fmt.Errorf("failed to execute HTML template: %w", err)
	}

	return nil
}

// resumeToMarkdown converts a resume to Markdown format
func (g *Generator) resumeToMarkdown(resume *models.Resume) string {
	var sb strings.Builder

	// Header
	sb.WriteString("# ")
	sb.WriteString(resume.Contact.Name)
	sb.WriteString("\n\n")

	// Contact info
	sb.WriteString("## Contact Information\n\n")
	sb.WriteString(fmt.Sprintf("- **Email:** %s\n", resume.Contact.Email))
	if resume.Contact.Phone != "" {
		sb.WriteString(fmt.Sprintf("- **Phone:** %s\n", resume.Contact.Phone))
	}
	if resume.Contact.Location != "" {
		sb.WriteString(fmt.Sprintf("- **Location:** %s\n", resume.Contact.Location))
	}
	if resume.Contact.LinkedIn != "" {
		sb.WriteString(fmt.Sprintf("- **LinkedIn:** %s\n", resume.Contact.LinkedIn))
	}
	if resume.Contact.GitHub != "" {
		sb.WriteString(fmt.Sprintf("- **GitHub:** %s\n", resume.Contact.GitHub))
	}
	if resume.Contact.Portfolio != "" {
		sb.WriteString(fmt.Sprintf("- **Portfolio:** %s\n", resume.Contact.Portfolio))
	}
	sb.WriteString("\n")

	// Summary
	if resume.Summary != "" {
		sb.WriteString("## Summary\n\n")
		sb.WriteString(resume.Summary)
		sb.WriteString("\n\n")
	}

	// Skills
	if len(resume.Skills) > 0 {
		sb.WriteString("## Skills\n\n")

		// Group skills by category
		categories := make(map[string][]models.Skill)
		for _, skill := range resume.Skills {
			category := skill.Category
			if category == "" {
				category = "General"
			}
			categories[category] = append(categories[category], skill)
		}

		for category, skills := range categories {
			sb.WriteString(fmt.Sprintf("### %s\n\n", category))
			for _, skill := range skills {
				sb.WriteString(fmt.Sprintf("- %s", skill.Name))
				if skill.Proficiency != "" {
					sb.WriteString(fmt.Sprintf(" (%s)", skill.Proficiency))
				}
				if skill.Years > 0 {
					sb.WriteString(fmt.Sprintf(" - %d years", skill.Years))
				}
				sb.WriteString("\n")
			}
			sb.WriteString("\n")
		}
	}

	// Experience
	if len(resume.Experience) > 0 {
		sb.WriteString("## Experience\n\n")
		for _, exp := range resume.Experience {
			sb.WriteString(fmt.Sprintf("### %s\n", exp.Position))
			sb.WriteString(fmt.Sprintf("**%s** | %s", exp.Company, formatDateRange(exp.StartDate, exp.EndDate, exp.Current)))
			if exp.Location != "" {
				sb.WriteString(fmt.Sprintf(" | %s", exp.Location))
			}
			sb.WriteString("\n\n")

			if exp.Description != "" {
				sb.WriteString(exp.Description)
				sb.WriteString("\n\n")
			}

			if len(exp.Achievements) > 0 {
				sb.WriteString("**Achievements:**\n")
				for _, achievement := range exp.Achievements {
					sb.WriteString(fmt.Sprintf("- %s\n", achievement))
				}
				sb.WriteString("\n")
			}

			if len(exp.Skills) > 0 {
				sb.WriteString("**Skills:** ")
				sb.WriteString(strings.Join(exp.Skills, ", "))
				sb.WriteString("\n\n")
			}
		}
	}

	// Education
	if len(resume.Education) > 0 {
		sb.WriteString("## Education\n\n")
		for _, edu := range resume.Education {
			sb.WriteString(fmt.Sprintf("### %s in %s\n", edu.Degree, edu.FieldOfStudy))
			sb.WriteString(fmt.Sprintf("**%s** | %s", edu.Institution, formatDateRange(edu.StartDate, edu.EndDate, false)))
			if edu.Location != "" {
				sb.WriteString(fmt.Sprintf(" | %s", edu.Location))
			}
			if edu.GPA != nil {
				sb.WriteString(fmt.Sprintf(" | GPA: %.2f", *edu.GPA))
			}
			sb.WriteString("\n\n")
		}
	}

	// Projects
	if len(resume.Projects) > 0 {
		sb.WriteString("## Projects\n\n")
		for _, project := range resume.Projects {
			sb.WriteString(fmt.Sprintf("### %s\n", project.Name))
			if project.URL != "" {
				sb.WriteString(fmt.Sprintf("[%s](%s) | ", project.Name, project.URL))
			}
			sb.WriteString(fmt.Sprintf("%s\n", formatDateRange(project.StartDate, project.EndDate, project.EndDate == nil)))
			sb.WriteString(fmt.Sprintf("%s\n", project.Description))
			if len(project.Technologies) > 0 {
				sb.WriteString(fmt.Sprintf("**Technologies:** %s\n", strings.Join(project.Technologies, ", ")))
			}
			sb.WriteString("\n")
		}
	}

	// Certifications
	if len(resume.Certifications) > 0 {
		sb.WriteString("## Certifications\n\n")
		for _, cert := range resume.Certifications {
			sb.WriteString(fmt.Sprintf("- **%s** - %s", cert.Name, cert.Issuer))
			if cert.IssueDate != (time.Time{}) {
				sb.WriteString(fmt.Sprintf(" (%s)", cert.IssueDate.Format("Jan 2006")))
			}
			if cert.CredentialID != "" {
				sb.WriteString(fmt.Sprintf(" - ID: %s", cert.CredentialID))
			}
			if cert.URL != "" {
				sb.WriteString(fmt.Sprintf(" - [%s](%s)", cert.URL, cert.URL))
			}
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	// Languages
	if len(resume.Languages) > 0 {
		sb.WriteString("## Languages\n\n")
		for _, lang := range resume.Languages {
			sb.WriteString(fmt.Sprintf("- **%s**: %s\n", lang.Language, lang.Proficiency))
		}
		sb.WriteString("\n")
	}

	// Custom sections
	for sectionName, items := range resume.CustomSections {
		sb.WriteString(fmt.Sprintf("## %s\n\n", sectionName))
		for _, item := range items {
			sb.WriteString(fmt.Sprintf("### %s\n", item.Title))
			if item.Description != "" {
				sb.WriteString(fmt.Sprintf("%s\n", item.Description))
			}
			for _, detail := range item.Details {
				sb.WriteString(fmt.Sprintf("- %s\n", detail))
			}
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// resumeToText converts a resume to plain text format
func (g *Generator) resumeToText(resume *models.Resume) string {
	var sb strings.Builder

	// Header
	sb.WriteString(fmt.Sprintf("%s\n", resume.Contact.Name))
	sb.WriteString(strings.Repeat("=", len(resume.Contact.Name)))
	sb.WriteString("\n\n")

	// Contact info
	sb.WriteString("CONTACT\n")
	sb.WriteString(fmt.Sprintf("Email: %s\n", resume.Contact.Email))
	if resume.Contact.Phone != "" {
		sb.WriteString(fmt.Sprintf("Phone: %s\n", resume.Contact.Phone))
	}
	if resume.Contact.Location != "" {
		sb.WriteString(fmt.Sprintf("Location: %s\n", resume.Contact.Location))
	}
	if resume.Contact.LinkedIn != "" {
		sb.WriteString(fmt.Sprintf("LinkedIn: %s\n", resume.Contact.LinkedIn))
	}
	if resume.Contact.GitHub != "" {
		sb.WriteString(fmt.Sprintf("GitHub: %s\n", resume.Contact.GitHub))
	}
	if resume.Contact.Portfolio != "" {
		sb.WriteString(fmt.Sprintf("Portfolio: %s\n", resume.Contact.Portfolio))
	}
	sb.WriteString("\n")

	// Summary
	if resume.Summary != "" {
		sb.WriteString("SUMMARY\n")
		sb.WriteString(strings.Repeat("-", 7))
		sb.WriteString("\n")
		sb.WriteString(resume.Summary)
		sb.WriteString("\n\n")
	}

	// Skills
	if len(resume.Skills) > 0 {
		sb.WriteString("SKILLS\n")
		sb.WriteString(strings.Repeat("-", 6))
		sb.WriteString("\n")

		for _, skill := range resume.Skills {
			sb.WriteString(fmt.Sprintf("- %s", skill.Name))
			if skill.Proficiency != "" {
				sb.WriteString(fmt.Sprintf(" (%s)", skill.Proficiency))
			}
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	// Experience
	if len(resume.Experience) > 0 {
		sb.WriteString("EXPERIENCE\n")
		sb.WriteString(strings.Repeat("-", 10))
		sb.WriteString("\n\n")

		for _, exp := range resume.Experience {
			sb.WriteString(fmt.Sprintf("%s\n", exp.Position))
			sb.WriteString(fmt.Sprintf("%s | %s\n", exp.Company, formatDateRange(exp.StartDate, exp.EndDate, exp.Current)))
			if exp.Location != "" {
				sb.WriteString(fmt.Sprintf("Location: %s\n", exp.Location))
			}
			if exp.Description != "" {
				sb.WriteString(fmt.Sprintf("\n%s\n", exp.Description))
			}
			if len(exp.Achievements) > 0 {
				sb.WriteString("\nAchievements:\n")
				for _, achievement := range exp.Achievements {
					sb.WriteString(fmt.Sprintf("  * %s\n", achievement))
				}
			}
			sb.WriteString("\n")
		}
	}

	// Education
	if len(resume.Education) > 0 {
		sb.WriteString("EDUCATION\n")
		sb.WriteString(strings.Repeat("-", 9))
		sb.WriteString("\n\n")

		for _, edu := range resume.Education {
			sb.WriteString(fmt.Sprintf("%s in %s\n", edu.Degree, edu.FieldOfStudy))
			sb.WriteString(fmt.Sprintf("%s | %s\n", edu.Institution, formatDateRange(edu.StartDate, edu.EndDate, false)))
			if edu.Location != "" {
				sb.WriteString(fmt.Sprintf("Location: %s\n", edu.Location))
			}
			if edu.GPA != nil {
				sb.WriteString(fmt.Sprintf("GPA: %.2f\n", *edu.GPA))
			}
			sb.WriteString("\n")
		}
	}

	// Projects
	if len(resume.Projects) > 0 {
		sb.WriteString("PROJECTS\n")
		sb.WriteString(strings.Repeat("-", 8))
		sb.WriteString("\n\n")

		for _, project := range resume.Projects {
			sb.WriteString(fmt.Sprintf("%s\n", project.Name))
			sb.WriteString(fmt.Sprintf("  %s\n", project.Description))
			if project.URL != "" {
				sb.WriteString(fmt.Sprintf("  URL: %s\n", project.URL))
			}
			if len(project.Technologies) > 0 {
				sb.WriteString(fmt.Sprintf("  Technologies: %s\n", strings.Join(project.Technologies, ", ")))
			}
			sb.WriteString("\n")
		}
	}

	// Certifications
	if len(resume.Certifications) > 0 {
		sb.WriteString("CERTIFICATIONS\n")
		sb.WriteString(strings.Repeat("-", 13))
		sb.WriteString("\n\n")

		for _, cert := range resume.Certifications {
			sb.WriteString(fmt.Sprintf("%s - %s\n", cert.Name, cert.Issuer))
			if cert.IssueDate != (time.Time{}) {
				sb.WriteString(fmt.Sprintf("  Issued: %s\n", cert.IssueDate.Format("Jan 2006")))
			}
			sb.WriteString("\n")
		}
	}

	// Languages
	if len(resume.Languages) > 0 {
		sb.WriteString("LANGUAGES\n")
		sb.WriteString(strings.Repeat("-", 9))
		sb.WriteString("\n\n")

		for _, lang := range resume.Languages {
			sb.WriteString(fmt.Sprintf("%s: %s\n", lang.Language, lang.Proficiency))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// resumeToHTML converts a resume to HTML format
func (g *Generator) resumeToHTML(resume *models.Resume) string {
	var sb strings.Builder

	// HTML header
	sb.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>`)
	sb.WriteString(resume.Contact.Name)
	sb.WriteString(`</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; margin: 0; padding: 20px; color: #333; }
        .container { max-width: 800px; margin: 0 auto; }
        h1 { color: #2c3e50; border-bottom: 2px solid #3498db; padding-bottom: 10px; }
        h2 { color: #3498db; margin-top: 30px; }
        .section { margin-bottom: 30px; }
        .job { margin-bottom: 20px; padding-bottom: 20px; border-bottom: 1px solid #eee; }
        .job:last-child { border-bottom: none; }
        .job-title { font-weight: bold; font-size: 1.2em; color: #2c3e50; }
        .job-company { color: #7f8c8d; }
        .date-range { color: #95a5a6; font-size: 0.9em; }
        .skills-list { display: flex; flex-wrap: wrap; gap: 5px; }
        .skill-tag { background: #e0f7fa; color: #006064; padding: 3px 8px; border-radius: 3px; font-size: 0.9em; }
        .contact-info { display: flex; flex-wrap: wrap; gap: 15px; margin-bottom: 15px; }
        .contact-item { display: flex; align-items: center; }
        .summary { background: #f8f9fa; padding: 15px; border-radius: 5px; margin-bottom: 20px; }
    </style>
</head>
<body>
    <div class="container">
        <h1>`)
	sb.WriteString(resume.Contact.Name)
	sb.WriteString(`</h1>

        <div class="section">
            <h2>Contact Information</h2>
            <div class="contact-info">
                <div class="contact-item">
                    <strong>Email:</strong> `)
	sb.WriteString(resume.Contact.Email)
	sb.WriteString(`
                </div>
`)

	if resume.Contact.Phone != "" {
		sb.WriteString(`                <div class="contact-item">
                    <strong>Phone:</strong> `)
		sb.WriteString(resume.Contact.Phone)
		sb.WriteString(`
                </div>
`)
	}

	if resume.Contact.Location != "" {
		sb.WriteString(`                <div class="contact-item">
                    <strong>Location:</strong> `)
		sb.WriteString(resume.Contact.Location)
		sb.WriteString(`
                </div>
`)
	}

	if resume.Contact.LinkedIn != "" {
		sb.WriteString(`                <div class="contact-item">
                    <strong>LinkedIn:</strong> <a href="`)
		sb.WriteString(resume.Contact.LinkedIn)
		sb.WriteString(`" target="_blank">`)
		sb.WriteString(resume.Contact.LinkedIn)
		sb.WriteString(`</a>
                </div>
`)
	}

	if resume.Contact.GitHub != "" {
		sb.WriteString(`                <div class="contact-item">
                    <strong>GitHub:</strong> <a href="`)
		sb.WriteString(resume.Contact.GitHub)
		sb.WriteString(`" target="_blank">`)
		sb.WriteString(resume.Contact.GitHub)
		sb.WriteString(`</a>
                </div>
`)
	}

	if resume.Contact.Portfolio != "" {
		sb.WriteString(`                <div class="contact-item">
                    <strong>Portfolio:</strong> <a href="`)
		sb.WriteString(resume.Contact.Portfolio)
		sb.WriteString(`" target="_blank">`)
		sb.WriteString(resume.Contact.Portfolio)
		sb.WriteString(`</a>
                </div>
`)
	}

	sb.WriteString(`            </div>
        </div>
`)

	// Summary
	if resume.Summary != "" {
		sb.WriteString(`
        <div class="section">
            <h2>Summary</h2>
            <div class="summary">`)
		sb.WriteString(template.HTMLEscapeString(resume.Summary))
		sb.WriteString(`</div>
        </div>
`)
	}

	// Skills
	if len(resume.Skills) > 0 {
		sb.WriteString(`
        <div class="section">
            <h2>Skills</h2>
            <div class="skills-list">
`)
		for _, skill := range resume.Skills {
			sb.WriteString(`                <span class="skill-tag">`)
			sb.WriteString(template.HTMLEscapeString(skill.Name))
			if skill.Proficiency != "" {
				sb.WriteString(` (<em>`)
				sb.WriteString(template.HTMLEscapeString(skill.Proficiency))
				sb.WriteString(`</em>`)
			}
			sb.WriteString(`</span>
`)
		}
		sb.WriteString(`            </div>
        </div>
`)
	}

	// Experience
	if len(resume.Experience) > 0 {
		sb.WriteString(`
        <div class="section">
            <h2>Experience</h2>
`)
		for _, exp := range resume.Experience {
			sb.WriteString(`
            <div class="job">
                <div class="job-title">`)
			sb.WriteString(template.HTMLEscapeString(exp.Position))
			sb.WriteString(`</div>
                <div class="job-company">`)
			sb.WriteString(template.HTMLEscapeString(exp.Company))
			sb.WriteString(`</div>
                <div class="date-range">`)
			sb.WriteString(formatDateRange(exp.StartDate, exp.EndDate, exp.Current))
			if exp.Location != "" {
				sb.WriteString(` | `)
				sb.WriteString(template.HTMLEscapeString(exp.Location))
			}
			sb.WriteString(`</div>
`)
			if exp.Description != "" {
				sb.WriteString(`
                <p>`)
				sb.WriteString(template.HTMLEscapeString(exp.Description))
				sb.WriteString(`</p>
`)
			}
			sb.WriteString(`            </div>
`)
		}
		sb.WriteString(`        </div>
`)
	}

	// Education
	if len(resume.Education) > 0 {
		sb.WriteString(`
        <div class="section">
            <h2>Education</h2>
`)
		for _, edu := range resume.Education {
			sb.WriteString(`
            <div class="job">
                <div class="job-title">`)
			sb.WriteString(template.HTMLEscapeString(edu.Degree))
			if edu.FieldOfStudy != "" {
				sb.WriteString(` in `)
				sb.WriteString(template.HTMLEscapeString(edu.FieldOfStudy))
			}
			sb.WriteString(`</div>
                <div class="job-company">`)
			sb.WriteString(template.HTMLEscapeString(edu.Institution))
			sb.WriteString(`</div>
                <div class="date-range">`)
			sb.WriteString(formatDateRange(edu.StartDate, edu.EndDate, false))
			if edu.Location != "" {
				sb.WriteString(` | `)
				sb.WriteString(template.HTMLEscapeString(edu.Location))
			}
			if edu.GPA != nil {
				sb.WriteString(fmt.Sprintf(` | GPA: %.2f`, *edu.GPA))
			}
			sb.WriteString(`</div>
            </div>
`)
		}
		sb.WriteString(`        </div>
`)
	}

	// Projects
	if len(resume.Projects) > 0 {
		sb.WriteString(`
        <div class="section">
            <h2>Projects</h2>
`)
		for _, project := range resume.Projects {
			sb.WriteString(`
            <div class="job">
                <div class="job-title">`)
			sb.WriteString(template.HTMLEscapeString(project.Name))
			sb.WriteString(`</div>
                <div class="date-range">`)
			sb.WriteString(formatDateRange(project.StartDate, project.EndDate, project.EndDate == nil))
			sb.WriteString(`</div>
                <p>`)
			sb.WriteString(template.HTMLEscapeString(project.Description))
			sb.WriteString(`</p>
`)
			if len(project.Technologies) > 0 {
				sb.WriteString(`
                <div><strong>Technologies:</strong> `)
				sb.WriteString(template.HTMLEscapeString(strings.Join(project.Technologies, ", ")))
				sb.WriteString(`</div>
`)
			}
			sb.WriteString(`            </div>
`)
		}
		sb.WriteString(`        </div>
`)
	}

	// Certifications
	if len(resume.Certifications) > 0 {
		sb.WriteString(`
        <div class="section">
            <h2>Certifications</h2>
            <ul>
`)
		for _, cert := range resume.Certifications {
			sb.WriteString(`                <li><strong>`)
			sb.WriteString(template.HTMLEscapeString(cert.Name))
			sb.WriteString(`</strong> - `)
			sb.WriteString(template.HTMLEscapeString(cert.Issuer))
			if cert.IssueDate != (time.Time{}) {
				sb.WriteString(fmt.Sprintf(` (%s)`, cert.IssueDate.Format("Jan 2006")))
			}
			sb.WriteString(`</li>
`)
		}
		sb.WriteString(`            </ul>
        </div>
`)
	}

	// Languages
	if len(resume.Languages) > 0 {
		sb.WriteString(`
        <div class="section">
            <h2>Languages</h2>
            <ul>
`)
		for _, lang := range resume.Languages {
			sb.WriteString(`                <li><strong>`)
			sb.WriteString(template.HTMLEscapeString(lang.Language))
			sb.WriteString(`</strong>: `)
			sb.WriteString(template.HTMLEscapeString(lang.Proficiency))
			sb.WriteString(`</li>
`)
		}
		sb.WriteString(`            </ul>
        </div>
`)
	}

	// HTML footer
	sb.WriteString(`
    </div>
</body>
</html>`)

	return sb.String()
}

// formatDateRange formats a date range as a string
func formatDateRange(start time.Time, end *time.Time, current bool) string {
	startStr := start.Format("Jan 2006")

	if current || end == nil {
		return fmt.Sprintf("%s - Present", startStr)
	}

	endStr := end.Format("Jan 2006")
	return fmt.Sprintf("%s - %s", startStr, endStr)
}
