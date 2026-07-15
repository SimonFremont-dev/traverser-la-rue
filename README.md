# CV Builder - Resume Adaptation Tool

A Go-based tool that helps you adapt your resume content by reading job offers and automatically highlighting your matching skills and experience.

## Features

- **Multi-format Parsing**: Parse job offers from text, JSON, YAML, HTML, and PDF formats
- **Intelligent Matching**: Match your resume against job requirements with fuzzy string matching
- **Adaptive Resume Generation**: Generate optimized resumes tailored to specific job offers
- **Comprehensive Analysis**: Get detailed reports on how well your resume matches job requirements
- **Suggestions**: Receive actionable suggestions to improve your resume match
- **Multiple Output Formats**: Generate resumes in JSON, YAML, Markdown, HTML, or plain text

## Installation

### Prerequisites

- Go 1.21 or higher
- Git

### Install from Source

```bash
# Clone the repository
git clone https://github.com/SimonFremont-dev/cv-builder.git
cd cv-builder

# Build the application
make build

# Or install directly
go install ./cmd/resume-adapter
```

### Using Makefile

```bash
# Build the binary
make build

# Install globally
make install

# Run all tests
make test

# Clean up
make clean
```

## Usage

### Basic Commands

```bash
# Get help
resume-adapter --help

# Check version
resume-adapter version

# Generate a new resume template
resume-adapter generate --output my-resume.yaml --format yaml

# Convert resume between formats
resume-adapter convert --input my-resume.json --output my-resume.yaml --to yaml
```

### Matching Commands

```bash
# Match resume against a job offer
resume-adapter match --resume my-resume.json --offer job-posting.txt

# Get detailed analysis
resume-adapter analyze --resume my-resume.yaml --offer job-offer.json --verbose

# Save match results to file
resume-adapter match --resume my-resume.json --offer job.txt --format json --output results.json
```

### Adaptation Commands

```bash
# Create an adapted resume for a job offer
resume-adapter adapt --resume my-resume.json --offer job-posting.txt --output adapted-resume.md

# Adapt with custom output format
resume-adapter adapt --resume my-resume.yaml --offer job.json --output adapted.html --format html

# Use custom matching criteria
resume-adapter adapt --resume my-resume.json --offer job.txt --output adapted.yaml --criteria my-criteria.json
```

## Resume Format

The tool supports structured resumes in JSON and YAML format. Here's an example structure:

```yaml
metadata:
  id: my-resume
  title: My Professional Resume
  version: 1.0.0
  template: default

contact:
  name: John Doe
  email: john@example.com
  phone: +1 234-567-8900
  linkedin: https://linkedin.com/in/johndoe
  github: https://github.com/johndoe
  location: San Francisco, CA

summary: Senior Software Engineer with 5+ years of experience...

experience:
  - id: exp-1
    company: Tech Company
    position: Senior Backend Engineer
    start_date: 2020-01-15T00:00:00Z
    current: true
    description: Lead backend development...
    skills: [Go, Kubernetes, Docker]
    achievements:
      - Improved API performance by 60%
      - Migrated to microservices

education:
  - id: edu-1
    institution: Stanford University
    degree: Master of Science
    field_of_study: Computer Science
    start_date: 2018-09-01T00:00:00Z
    end_date: 2020-05-15T00:00:00Z
    gpa: 3.8

skills:
  - id: skill-1
    name: Go
    category: Programming Languages
    proficiency: Expert
    years: 5
  - id: skill-2
    name: Kubernetes
    category: DevOps
    proficiency: Advanced
    years: 4
```

## Job Offer Format

Job offers can be provided in various formats:

### Text Format

Simple text with job description:

```
Position: Senior Backend Engineer
Company: Tech Corp
Location: Remote

Requirements:
- 5+ years of Go experience
- Knowledge of Kubernetes and Docker
- Experience with RESTful APIs

Preferred Skills:
- AWS experience
- GraphQL knowledge
```

### JSON/YAML Format

Structured job offer format:

```json
{
  "metadata": {
    "title": "Senior Backend Engineer",
    "format": "json"
  },
  "company": {
    "name": "Tech Corp",
    "location": "San Francisco, CA"
  },
  "position": {
    "title": "Senior Backend Engineer",
    "type": "Full-time",
    "salary": {
      "min": 120000,
      "max": 150000,
      "currency": "USD",
      "period": "yearly"
    }
  },
  "requirements": {
    "required_skills": [
      {"name": "Go", "years": 5, "required": true},
      {"name": "Kubernetes", "required": true}
    ],
    "preferred_skills": [
      {"name": "AWS", "nice_to_have": true}
    ],
    "experience": {
      "total_years": 5
    },
    "education": [
      {"degree": "Bachelor", "required": true}
    ]
  }
}
```

## Matching Criteria

You can customize the matching algorithm by providing a criteria file:

```json
{
  "skill_weight": 0.4,
  "experience_weight": 0.3,
  "education_weight": 0.2,
  "keyword_weight": 0.1,
  "min_score_for_suggestion": 70.0,
  "suggest_alternatives": true
}
```

## How It Works

1. **Parsing**: The tool parses both your resume and the job offer, extracting skills, experience, education, and other relevant information.

2. **Matching**: Using fuzzy string matching and Levenshtein distance, the tool compares your skills and experience against job requirements.

3. **Scoring**: Each match is assigned a score based on:
   - Exact matches (100%)
   - Partial matches (80-99%)
   - Alternative matches (90%)
   - Keyword matches (80%)

4. **Adaptation**: The tool reorders your skills and experience to highlight the most relevant ones for the job.

5. **Suggestions**: Based on missing requirements, the tool provides actionable suggestions to improve your resume.

## Project Structure

```
cv-builder/
├── cmd/
│   └── resume-adapter/
│       └── main.go          # CLI entry point
├── internal/
│   ├── cli/                 # CLI commands and handlers
│   ├── generator/           # Resume output generators
│   ├── matcher/             # Matching algorithms
│   ├── models/              # Data models
│   └── parser/              # Parsers for various formats
├── examples/                # Example files
├── Makefile                 # Build automation
├── go.mod                   # Go module file
└── README.md                # This file
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/your-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin feature/your-feature`)
5. Create a new Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
