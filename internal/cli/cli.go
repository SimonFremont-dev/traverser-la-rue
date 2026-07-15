package cli

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/SimonFremont-dev/cv-builder/internal/generator"
	"github.com/SimonFremont-dev/cv-builder/internal/matcher"
	"github.com/SimonFremont-dev/cv-builder/internal/models"
	"github.com/SimonFremont-dev/cv-builder/internal/parser"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// NewRootCommand creates a new root command
func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "resume-adapter",
		Short: "Adapt your resume to match job offers",
		Long: `resume-adapter is a tool that helps you adapt your resume content
by reading job offers and setting your matching skills and experience
to the forefront of your resume.

Features:
- Parse job offers from multiple formats (text, JSON, YAML, HTML, PDF)
- Match your resume against job requirements
- Generate adapted resumes optimized for specific job offers
- Get suggestions for improving your resume match
- Output in multiple formats (JSON, YAML, Markdown, HTML, Text)

Examples:
  # Match resume against a job offer
  resume-adapter match --resume my-resume.json --offer job-posting.txt

  # Generate an adapted resume
  resume-adapter adapt --resume my-resume.yaml --offer job-posting.json --output adapted-resume.md

  # Get match analysis with suggestions
  resume-adapter analyze --resume my-resume.json --offer job-posting.txt

  # Convert resume between formats
  resume-adapter convert --input my-resume.json --output my-resume.yaml --to yaml`,
	}

	// Add subcommands
	rootCmd.AddCommand(newMatchCommand())
	rootCmd.AddCommand(newAdaptCommand())
	rootCmd.AddCommand(newAnalyzeCommand())
	rootCmd.AddCommand(newConvertCommand())
	rootCmd.AddCommand(newGenerateCommand())
	rootCmd.AddCommand(newGenerateAllCommand())
	rootCmd.AddCommand(newVersionCommand())

	// Configure persistent flags
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default is $HOME/.resume-adapter.yaml)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")

	return rootCmd
}

// newMatchCommand creates the match subcommand
func newMatchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "match",
		Short: "Match resume against a job offer",
		Long: `Match your resume against a job offer to see how well you fit.

This command compares your resume skills, experience, and education
against the job requirements and provides a match score.`,
		Example: `resume-adapter match --resume my-resume.json --offer job.txt`,
		Run: func(cmd *cobra.Command, args []string) {
			verbose, _ := cmd.Flags().GetBool("verbose")
			resumePath, _ := cmd.Flags().GetString("resume")
			offerPath, _ := cmd.Flags().GetString("offer")
			format, _ := cmd.Flags().GetString("format")
			output, _ := cmd.Flags().GetString("output")
			criteriaFile, _ := cmd.Flags().GetString("criteria")

			if err := runMatch(verbose, resumePath, offerPath, format, output, criteriaFile); err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().StringP("resume", "r", "", "path to resume file (required)")
	cmd.Flags().StringP("offer", "o", "", "path to job offer file (required)")
	cmd.Flags().StringP("format", "f", "text", "output format: json, yaml, text")
	cmd.Flags().StringP("output", "O", "", "output file path (default: stdout)")
	cmd.Flags().StringP("criteria", "C", "", "path to custom matching criteria file")
	cmd.MarkFlagRequired("resume")
	cmd.MarkFlagRequired("offer")

	return cmd
}

// newAdaptCommand creates the adapt subcommand
func newAdaptCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "adapt",
		Short: "Create an adapted resume for a job offer",
		Long: `Generate a new version of your resume optimized for a specific job offer.

This command analyzes the job offer, matches it against your resume,
and creates a new resume with skills and experience reordered to
match the job requirements.`,
		Example: `resume-adapter adapt --resume my-resume.yaml --offer job.json --output adapted-resume.md`,
		Run: func(cmd *cobra.Command, args []string) {
			verbose, _ := cmd.Flags().GetBool("verbose")
			resumePath, _ := cmd.Flags().GetString("resume")
			offerPath, _ := cmd.Flags().GetString("offer")
			output, _ := cmd.Flags().GetString("output")
			outputFormat, _ := cmd.Flags().GetString("format")
			criteriaFile, _ := cmd.Flags().GetString("criteria")

			if err := runAdapt(verbose, resumePath, offerPath, output, outputFormat, criteriaFile); err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().StringP("resume", "r", "", "path to resume file (required)")
	cmd.Flags().StringP("offer", "o", "", "path to job offer file (required)")
	cmd.Flags().StringP("output", "O", "", "output file path (required)")
	cmd.Flags().StringP("format", "f", "markdown", "output format: json, yaml, markdown, html, text")
	cmd.Flags().StringP("criteria", "C", "", "path to custom matching criteria file")
	cmd.MarkFlagRequired("resume")
	cmd.MarkFlagRequired("offer")
	cmd.MarkFlagRequired("output")

	return cmd
}

// newAnalyzeCommand creates the analyze subcommand
func newAnalyzeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "analyze",
		Short: "Detailed analysis of resume vs job offer",
		Long: `Perform a detailed analysis of how well your resume matches a job offer.

This provides comprehensive feedback including:
- Overall match score and grade
- Skill by skill matching
- Experience and education matching
- Missing requirements
- Suggestions for improvement`,
		Example: `resume-adapter analyze --resume my-resume.json --offer job.txt`,
		Run: func(cmd *cobra.Command, args []string) {
			verbose, _ := cmd.Flags().GetBool("verbose")
			resumePath, _ := cmd.Flags().GetString("resume")
			offerPath, _ := cmd.Flags().GetString("offer")
			format, _ := cmd.Flags().GetString("format")
			criteriaFile, _ := cmd.Flags().GetString("criteria")

			if err := runAnalyze(verbose, resumePath, offerPath, format, criteriaFile); err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().StringP("resume", "r", "", "path to resume file (required)")
	cmd.Flags().StringP("offer", "o", "", "path to job offer file (required)")
	cmd.Flags().StringP("format", "f", "text", "output format: json, yaml, text")
	cmd.Flags().StringP("criteria", "C", "", "path to custom matching criteria file")
	cmd.MarkFlagRequired("resume")
	cmd.MarkFlagRequired("offer")

	return cmd
}

// newConvertCommand creates the convert subcommand
func newConvertCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "convert",
		Short: "Convert resume between formats",
		Long: `Convert your resume from one format to another.

Supports: JSON, YAML, Markdown, HTML, Text`,
		Example: `resume-adapter convert --input my-resume.json --output my-resume.yaml --to yaml`,
		Run: func(cmd *cobra.Command, args []string) {
			verbose, _ := cmd.Flags().GetBool("verbose")
			input, _ := cmd.Flags().GetString("input")
			output, _ := cmd.Flags().GetString("output")
			toFormat, _ := cmd.Flags().GetString("to")

			if err := runConvert(verbose, input, output, toFormat); err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().StringP("input", "i", "", "input file path (required)")
	cmd.Flags().StringP("output", "o", "", "output file path (required)")
	cmd.Flags().StringP("to", "t", "", "target format: json, yaml, markdown, html, text (required)")
	cmd.MarkFlagRequired("input")
	cmd.MarkFlagRequired("output")
	cmd.MarkFlagRequired("to")

	return cmd
}

// newGenerateCommand creates the generate subcommand
func newGenerateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a new resume template",
		Long: `Create a new empty resume template in the specified format.

This is useful for creating a starting point for your resume.`,
		Example: `resume-adapter generate --output my-resume.yaml --format yaml`,
		Run: func(cmd *cobra.Command, args []string) {
			verbose, _ := cmd.Flags().GetBool("verbose")
			output, _ := cmd.Flags().GetString("output")
			format, _ := cmd.Flags().GetString("format")

			if err := runGenerate(verbose, output, format); err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().StringP("output", "o", "", "output file path (required)")
	cmd.Flags().StringP("format", "f", "yaml", "output format: json, yaml, markdown, html, text")
	cmd.MarkFlagRequired("output")

	return cmd
}

// newVersionCommand creates the version subcommand
func newVersionCommand() *cobra.Command {
	version := "0.1.0"

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Long:  `Print the version of resume-adapter`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("resume-adapter version %s\n", version)
		},
	}

	return cmd
}

// newGenerateAllCommand creates the generate-all subcommand
func newGenerateAllCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate-all",
		Short: "Generate resume in multiple formats (DOCX and PDF) from input",
		Long: `Generate a resume in multiple formats from a single input file.

This command takes a resume in any supported format (JSON, YAML, Markdown, Text, DOCX)
and generates output in DOCX and PDF formats, plus HTML and Markdown for reference.

If the input is a DOCX file, it will extract the text and generate all formats.

Examples:
  # Generate all formats from a DOCX resume
  resume-adapter generate-all --input my-resume.docx --output ./output/

  # Generate all formats from a JSON resume
  resume-adapter generate-all --input my-resume.json --output ./output/

  # Generate all formats with a specific job offer for adaptation
  resume-adapter generate-all --input my-resume.json --offer job.txt --output ./output/ --adapt`,
		Run: func(cmd *cobra.Command, args []string) {
			verbose, _ := cmd.Flags().GetBool("verbose")
			inputPath, _ := cmd.Flags().GetString("input")
			outputDir, _ := cmd.Flags().GetString("output")
			offerPath, _ := cmd.Flags().GetString("offer")
			adapt, _ := cmd.Flags().GetBool("adapt")

			if err := runGenerateAll(verbose, inputPath, outputDir, offerPath, adapt); err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().StringP("input", "i", "", "input file path (required, supports: json, yaml, md, txt, docx)")
	cmd.Flags().StringP("output", "o", "", "output directory (required)")
	cmd.Flags().StringP("offer", "j", "", "optional job offer file for adaptation")
	cmd.Flags().BoolP("adapt", "a", false, "adapt resume to match job offer")
	cmd.MarkFlagRequired("input")
	cmd.MarkFlagRequired("output")

	return cmd
}

// runMatch executes the match command
func runMatch(verbose bool, resumePath, offerPath, format, output, criteriaFile string) error {
	// Load resume
	resumeParser := parser.NewResumeParser()
	resume, err := resumeParser.ParseResume(resumePath)
	if err != nil {
		return fmt.Errorf("failed to parse resume: %w", err)
	}

	// Load job offer
	jobParser := parser.NewJobParser()
	jobOffer, err := jobParser.ParseJobOffer(offerPath)
	if err != nil {
		return fmt.Errorf("failed to parse job offer: %w", err)
	}

	// Load custom criteria if specified
	var criteria models.MatchCriteria
	if criteriaFile != "" {
		data, err := os.ReadFile(criteriaFile)
		if err != nil {
			return fmt.Errorf("failed to read criteria file: %w", err)
		}
		if err := json.Unmarshal(data, &criteria); err != nil {
			return fmt.Errorf("failed to parse criteria: %w", err)
		}
	} else {
		criteria = models.DefaultMatchCriteria()
	}

	// Create matcher and perform matching
	matcher := matcher.NewMatcherWithCriteria(criteria)
	result := matcher.Match(resume, jobOffer)

	// Output result
	if output == "" {
		// Print to stdout
		return printMatchResult(result, format)
	}

	// Write to file
	return writeMatchResult(result, output, format)
}

// runAdapt executes the adapt command
func runAdapt(verbose bool, resumePath, offerPath, outputPath, outputFormat, criteriaFile string) error {
	// Load resume
	resumeParser := parser.NewResumeParser()
	resume, err := resumeParser.ParseResume(resumePath)
	if err != nil {
		return fmt.Errorf("failed to parse resume: %w", err)
	}

	// Load job offer
	jobParser := parser.NewJobParser()
	jobOffer, err := jobParser.ParseJobOffer(offerPath)
	if err != nil {
		return fmt.Errorf("failed to parse job offer: %w", err)
	}

	// Load custom criteria if specified
	var criteria models.MatchCriteria
	if criteriaFile != "" {
		data, err := os.ReadFile(criteriaFile)
		if err != nil {
			return fmt.Errorf("failed to read criteria file: %w", err)
		}
		if err := json.Unmarshal(data, &criteria); err != nil {
			return fmt.Errorf("failed to parse criteria: %w", err)
		}
	} else {
		criteria = models.DefaultMatchCriteria()
	}

	// Create matcher and perform matching with adaptation
	matcher := matcher.NewMatcherWithCriteria(criteria)
	result := matcher.MatchWithAdaptation(resume, jobOffer)

	if result.AdaptedResume == nil {
		return fmt.Errorf("failed to create adapted resume")
	}

	// Generate output
	gen := generator.NewGenerator()
	if err := gen.GenerateResume(result.AdaptedResume, outputFormat, outputPath); err != nil {
		return fmt.Errorf("failed to generate adapted resume: %w", err)
	}

	// Print match summary
	fmt.Printf("Adapted resume created: %s\n", outputPath)
	fmt.Printf("Match score: %.1f%% (Grade: %s)\n", result.OverallScore, result.MatchGrade)

	if verbose {
		fmt.Printf("\nMissing requirements:\n")
		for _, req := range result.MissingRequirements {
			fmt.Printf("  - %s\n", req)
		}

		fmt.Printf("\nSuggestions:\n")
		for _, sug := range result.SuggestedChanges {
			fmt.Printf("  [%s] %s: %s\n", sug.Priority, sug.Type, sug.Reason)
		}
	}

	return nil
}

// runAnalyze executes the analyze command
func runAnalyze(verbose bool, resumePath, offerPath, format, criteriaFile string) error {
	// Load resume
	resumeParser := parser.NewResumeParser()
	resume, err := resumeParser.ParseResume(resumePath)
	if err != nil {
		return fmt.Errorf("failed to parse resume: %w", err)
	}

	// Load job offer
	jobParser := parser.NewJobParser()
	jobOffer, err := jobParser.ParseJobOffer(offerPath)
	if err != nil {
		return fmt.Errorf("failed to parse job offer: %w", err)
	}

	// Load custom criteria if specified
	var criteria models.MatchCriteria
	if criteriaFile != "" {
		data, err := os.ReadFile(criteriaFile)
		if err != nil {
			return fmt.Errorf("failed to read criteria file: %w", err)
		}
		if err := json.Unmarshal(data, &criteria); err != nil {
			return fmt.Errorf("failed to parse criteria: %w", err)
		}
	} else {
		criteria = models.DefaultMatchCriteria()
	}

	// Create matcher and perform matching
	matcher := matcher.NewMatcherWithCriteria(criteria)
	result := matcher.Match(resume, jobOffer)

	// Output result
	if format == "" {
		format = "text"
	}

	// Print detailed analysis to stdout
	return printAnalysisResult(result, format, verbose)
}

// runConvert executes the convert command
func runConvert(verbose bool, input, output, toFormat string) error {
	// Determine input format from extension
	ext := strings.ToLower(filepath.Ext(input))
	var format string
	switch ext {
	case ".json":
		format = "json"
	case ".yaml", ".yml":
		format = "yaml"
	case ".md", ".markdown":
		format = "markdown"
	case ".html":
		format = "html"
	case ".txt", ".text":
		format = "text"
	default:
		format = "text"
	}

	// Parse based on format
	var resume *models.Resume
	var err error

	switch format {
	case "json", "yaml", "yml":
		resumeParser := parser.NewResumeParser()
		resume, err = resumeParser.ParseResume(input)
		if err != nil {
			return fmt.Errorf("failed to parse resume: %w", err)
		}
	case "markdown", "md", "html", "text", "txt":
		// For these formats, we need to create a basic resume
		// Read the file and create a minimal resume
		data, err := os.ReadFile(input)
		if err != nil {
			return fmt.Errorf("failed to read input file: %w", err)
		}

		// Create a basic resume
		resume = models.NewResume()
		resume.Metadata.ID = strings.TrimSuffix(filepath.Base(input), filepath.Ext(input))

		// For markdown/html/text, try to extract some info
		content := string(data)
		if strings.Contains(content, "# ") || strings.Contains(content, "<h1>") {
			// Extract name from first heading
			lines := strings.Split(content, "\n")
			if len(lines) > 0 {
				name := strings.TrimSpace(lines[0])
				name = strings.TrimPrefix(name, "#")
				name = strings.TrimSpace(name)
				if name != "" {
					resume.Contact.Name = name
				}
			}
		}
	default:
		return fmt.Errorf("unsupported input format: %s", format)
	}

	// Generate output
	gen := generator.NewGenerator()
	if err := gen.GenerateResume(resume, toFormat, output); err != nil {
		return fmt.Errorf("failed to generate output: %w", err)
	}

	fmt.Printf("Converted %s to %s format: %s\n", input, toFormat, output)
	return nil
}

// runGenerate executes the generate command
func runGenerate(verbose bool, output, format string) error {
	// Create a new empty resume
	resume := models.NewResume()
	resume.Metadata.ID = strings.TrimSuffix(filepath.Base(output), filepath.Ext(output))
	resume.Contact.Name = "Your Name"
	resume.Contact.Email = "your.email@example.com"
	resume.Summary = "Professional summary highlighting your key skills and experience."

	// Generate output
	gen := generator.NewGenerator()
	if err := gen.GenerateResume(resume, format, output); err != nil {
		return fmt.Errorf("failed to generate resume template: %w", err)
	}

	fmt.Printf("Resume template created: %s\n", output)
	return nil
}

// printMatchResult prints the match result in the specified format
func printMatchResult(result *models.MatchResult, format string) error {
	if result == nil {
		return fmt.Errorf("no match result to print")
	}

	switch strings.ToLower(format) {
	case "json":
		data, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal result to JSON: %w", err)
		}
		fmt.Println(string(data))
	case "yaml", "yml":
		data, err := yaml.Marshal(result)
		if err != nil {
			return fmt.Errorf("failed to marshal result to YAML: %w", err)
		}
		fmt.Println(string(data))
	default:
		// Text format
		fmt.Printf("Match Score: %.1f%%\n", result.OverallScore)
		fmt.Printf("Grade: %s\n", result.MatchGrade)
		fmt.Printf("\nMatched Skills:\n")
		for _, match := range result.SkillMatches {
			fmt.Printf("  - %s: %.0f%% (%s)\n",
				match.JobRequirement.Name, match.MatchScore, match.MatchType)
		}

		fmt.Printf("\nExperience Match:\n")
		fmt.Printf("  Total Years: %.1f%%\n", result.ExperienceMatch.TotalYearsMatch)
		if result.ExperienceMatch.InFieldMatch > 0 {
			fmt.Printf("  In Field: %.1f%%\n", result.ExperienceMatch.InFieldMatch)
		}
		if result.ExperienceMatch.ManagementMatch {
			fmt.Printf("  Management: Yes\n")
		}

		fmt.Printf("\nEducation Match:\n")
		fmt.Printf("  Degree: %v\n", result.EducationMatch.DegreeMatch)
		fmt.Printf("  Field: %v\n", result.EducationMatch.FieldMatch)

		fmt.Printf("\nMissing Requirements:\n")
		if len(result.MissingRequirements) == 0 {
			fmt.Println("  None")
		} else {
			for _, req := range result.MissingRequirements {
				fmt.Printf("  - %s\n", req)
			}
		}
	}

	return nil
}

// printAnalysisResult prints the analysis result in the specified format
func printAnalysisResult(result *models.MatchResult, format string, verbose bool) error {
	if result == nil {
		return fmt.Errorf("no match result to print")
	}

	switch strings.ToLower(format) {
	case "json":
		data, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal result to JSON: %w", err)
		}
		fmt.Println(string(data))
	case "yaml", "yml":
		data, err := yaml.Marshal(result)
		if err != nil {
			return fmt.Errorf("failed to marshal result to YAML: %w", err)
		}
		fmt.Println(string(data))
	default:
		// Text format - detailed analysis
		fmt.Printf("=== RESUME ANALYSIS REPORT ===\n\n")

		fmt.Printf("Overall Match Score: %.1f%% (Grade: %s)\n\n", result.OverallScore, result.MatchGrade)

		// Job offer info
		fmt.Printf("Job Offer: %s\n", result.JobOfferID)
		fmt.Printf("Resume: %s\n\n", result.ResumeID)

		// Skill matches
		fmt.Printf("SKILL MATCHES (%d total):\n", len(result.SkillMatches))
		for i, match := range result.SkillMatches {
			var resumeSkillName string
			if match.ResumeSkill != nil {
				resumeSkillName = match.ResumeSkill.Name
			} else {
				resumeSkillName = "N/A"
			}
			fmt.Printf("  %d. %s (Job: %s) - %.1f%% [%s]\n",
				i+1, resumeSkillName, match.JobRequirement.Name, match.MatchScore, match.MatchType)
			if verbose && match.ProficiencyMatch != "" {
				fmt.Printf("     Proficiency: %s\n", match.ProficiencyMatch)
			}
			if verbose && match.YearsMatch > 0 {
				fmt.Printf("     Years Match: %.1f%%\n", match.YearsMatch)
			}
		}
		fmt.Println()

		// Experience match
		fmt.Println("EXPERIENCE MATCH:")
		fmt.Printf("  Total Years Match: %.1f%%\n", result.ExperienceMatch.TotalYearsMatch)
		if result.ExperienceMatch.InFieldMatch > 0 {
			fmt.Printf("  In Field Match: %.1f%%\n", result.ExperienceMatch.InFieldMatch)
		}
		fmt.Printf("  Management Experience: %v\n", result.ExperienceMatch.ManagementMatch)

		if len(result.ExperienceMatch.RelevantExperience) > 0 {
			fmt.Println("  Relevant Experience:")
			for i, exp := range result.ExperienceMatch.RelevantExperience {
				if exp.WorkExperience != nil {
					fmt.Printf("    %d. %s at %s - %.1f%% relevance\n",
						i+1, exp.WorkExperience.Position, exp.WorkExperience.Company, exp.RelevanceScore)
					if verbose && len(exp.MatchedKeywords) > 0 {
						fmt.Printf("       Matched Keywords: %s\n", strings.Join(exp.MatchedKeywords, ", "))
					}
				}
			}
		}

		if len(result.ExperienceMatch.MissingExperience) > 0 {
			fmt.Println("  Missing Experience:")
			for _, exp := range result.ExperienceMatch.MissingExperience {
				fmt.Printf("    - %s\n", exp)
			}
		}
		fmt.Println()

		// Education match
		fmt.Println("EDUCATION MATCH:")
		fmt.Printf("  Degree Match: %v\n", result.EducationMatch.DegreeMatch)
		fmt.Printf("  Field Match: %v\n", result.EducationMatch.FieldMatch)
		fmt.Printf("  GPA Meets Requirement: %v\n", result.EducationMatch.GPAMeet)

		if len(result.EducationMatch.MissingEducation) > 0 {
			fmt.Println("  Missing Education:")
			for _, edu := range result.EducationMatch.MissingEducation {
				fmt.Printf("    - %s\n", edu)
			}
		}
		fmt.Println()

		// Missing requirements
		fmt.Println("MISSING REQUIREMENTS:")
		if len(result.MissingRequirements) == 0 {
			fmt.Println("  None - You meet all requirements!")
		} else {
			for i, req := range result.MissingRequirements {
				fmt.Printf("  %d. %s\n", i+1, req)
			}
		}
		fmt.Println()

		// Suggestions
		if len(result.SuggestedChanges) > 0 {
			fmt.Println("SUGGESTIONS FOR IMPROVEMENT:")
			for i, sug := range result.SuggestedChanges {
				fmt.Printf("  %d. [%s Priority] %s\n", i+1, sug.Priority, sug.Reason)
				if sug.SuggestedText != "" {
					fmt.Printf("     Suggested: %s\n", sug.SuggestedText)
				}
				if sug.Section != "" {
					fmt.Printf("     Section: %s\n", sug.Section)
				}
			}
		}
	}

	return nil
}

// writeMatchResult writes the match result to a file
func writeMatchResult(result *models.MatchResult, outputPath, format string) error {
	if result == nil {
		return fmt.Errorf("no match result to write")
	}

	var data []byte
	var err error

	switch strings.ToLower(format) {
	case "json":
		data, err = json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal result to JSON: %w", err)
		}
	case "yaml", "yml":
		data, err = yaml.Marshal(result)
		if err != nil {
			return fmt.Errorf("failed to marshal result to YAML: %w", err)
		}
	default:
		// Text format
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Match Score: %.1f%%\n", result.OverallScore))
		sb.WriteString(fmt.Sprintf("Grade: %s\n\n", result.MatchGrade))
		sb.WriteString("Matched Skills:\n")
		for _, match := range result.SkillMatches {
			sb.WriteString(fmt.Sprintf("  - %s: %.0f%% (%s)\n",
				match.JobRequirement.Name, match.MatchScore, match.MatchType))
		}
		sb.WriteString("\nMissing Requirements:\n")
		for _, req := range result.MissingRequirements {
			sb.WriteString(fmt.Sprintf("  - %s\n", req))
		}
		data = []byte(sb.String())
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}

// runGenerateAll executes the generate-all command
func runGenerateAll(verbose bool, inputPath, outputDir, offerPath string, adapt bool) error {
	// Check if input file exists
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return fmt.Errorf("input file does not exist: %s", inputPath)
	}

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	var resume *models.Resume
	var err error

	// Parse input based on extension
	ext := strings.ToLower(filepath.Ext(inputPath))
	switch ext {
	case ".docx":
		// Parse DOCX file
		docxParser := parser.NewDOCXParser()
		resume, err = docxParser.ParseDOCX(inputPath)
		if err != nil {
			return fmt.Errorf("failed to parse DOCX: %w", err)
		}
	case ".json":
		// Parse JSON resume
		resumeParser := parser.NewResumeParser()
		resume, err = resumeParser.ParseResume(inputPath)
		if err != nil {
			return fmt.Errorf("failed to parse JSON: %w", err)
		}
	case ".yaml", ".yml":
		// Parse YAML resume
		resumeParser := parser.NewResumeParser()
		resume, err = resumeParser.ParseResume(inputPath)
		if err != nil {
			return fmt.Errorf("failed to parse YAML: %w", err)
		}
	case ".md", ".markdown":
		// Parse Markdown - read and create basic resume
		data, err := os.ReadFile(inputPath)
		if err != nil {
			return fmt.Errorf("failed to read Markdown file: %w", err)
		}
		resume = models.NewResume()
		resume.Summary = string(data)
		resume.Metadata.ID = strings.TrimSuffix(filepath.Base(inputPath), ext)
	case ".txt", ".text":
		// Parse Text - read and create basic resume
		data, err := os.ReadFile(inputPath)
		if err != nil {
			return fmt.Errorf("failed to read text file: %w", err)
		}
		resume = models.NewResume()
		resume.Summary = string(data)
		resume.Metadata.ID = strings.TrimSuffix(filepath.Base(inputPath), ext)
	default:
		return fmt.Errorf("unsupported input format: %s. Supported: .docx, .json, .yaml, .yml, .md, .markdown, .txt, .text", ext)
	}

	// If adaptation is requested, match against job offer
	if adapt && offerPath != "" {
		jobParser := parser.NewJobParser()
		jobOffer, err := jobParser.ParseJobOffer(offerPath)
		if err != nil {
			return fmt.Errorf("failed to parse job offer: %w", err)
		}

		matcher := matcher.NewMatcher()
		matchResult := matcher.MatchWithAdaptation(resume, jobOffer)
		if matchResult.AdaptedResume != nil {
			resume = matchResult.AdaptedResume
			if verbose {
				fmt.Printf("Adapted resume with match score: %.1f%%\n", matchResult.OverallScore)
			}
		}
	}

	// Generate all formats
	pdfGen := generator.NewPDFGenerator()
	if err := pdfGen.GenerateMultipleOutputs(resume, outputDir); err != nil {
		return fmt.Errorf("failed to generate outputs: %w", err)
	}

	fmt.Printf("Generated files in %s:\n", outputDir)
	fmt.Printf("  - %s_resume.docx\n", sanitizeFilename(resume.Contact.Name))
	fmt.Printf("  - %s_resume.pdf\n", sanitizeFilename(resume.Contact.Name))
	fmt.Printf("  - %s_resume.html\n", sanitizeFilename(resume.Contact.Name))
	fmt.Printf("  - %s_resume.md\n", sanitizeFilename(resume.Contact.Name))

	return nil
}

// sanitizeFilename creates a safe filename from a string
func sanitizeFilename(name string) string {
	if name == "" {
		name = "resume"
	}
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "/", "_")
	name = strings.ReplaceAll(name, "\\", "_")
	name = strings.ReplaceAll(name, ":", "_")
	name = strings.ReplaceAll(name, "*", "_")
	name = strings.ReplaceAll(name, "?", "_")
	name = strings.ReplaceAll(name, "\"", "_")
	name = strings.ReplaceAll(name, "<", "_")
	name = strings.ReplaceAll(name, ">", "_")
	name = strings.ReplaceAll(name, "|", "_")
	return strings.ToLower(name)
}
