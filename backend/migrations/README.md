# Database Migrations

This directory contains SQL migration files for the CV Builder backend database.

## Migration System

Migrations are automatically applied when the backend starts if a PostgreSQL database is configured.

### Files

1. **`001_initial_schema.sql`** - Initial database schema with:
   - ESCO skills table
   - O*NET skills table
   - Skill categories
   - Keyword to skill mappings
   - Profile data
   - Job offers
   - Matching results
   - Resumes

2. **`002_sample_data.sql`** - Sample data including:
   - Programming languages (Go, Python, Java, etc.)
   - Tools (Docker, Kubernetes, Git, etc.)
   - Databases (PostgreSQL, MySQL, MongoDB, etc.)
   - Cloud platforms (AWS, Azure, GCP, etc.)
   - Frameworks (React, Angular, Spring Boot, etc.)
   - Keyword mappings (e.g., "golang" → "Go", "k8s" → "Kubernetes")

## Database Schema Overview

### ESCO Skills (`esco_skills`)

European Skills/Competences, Qualifications and Occupations taxonomy.

Fields:
- `esco_id` - Unique ESCO identifier (e.g., "1.1.1.1.1")
- `uri` - ESCO URI
- `concept_type` - Type: skill, knowledge, or competence
- `preferred_label` - Primary skill name
- `description` - Skill description
- `alt_labels` - Alternative labels
- `category_id` - Reference to skill category
- `parent_esco_id` - Parent in ESCO hierarchy
- `level` - Hierarchy level
- `relevance_score` - Relevance weight (0-1)

### O*NET Skills (`onet_skills`)

US Occupational Information Network taxonomy.

Fields:
- `onet_code` - O*NET-SOC code
- `element_id` - Specific element identifier
- `content_model` - Content model (Skills, Knowledge, Abilities)
- `name` - Skill name
- `description` - Skill description
- `category` - O*NET category
- `category_id` - Reference to skill category
- `example` - Usage example
- `relevance_score` - Relevance weight (0-1)

### Skill Categories (`skill_categories`)

Classification system for organizing skills.

Fields:
- `name` - Category name (Programming, Tools, Databases, etc.)
- `description` - Category description
- `parent_id` - Parent category (for hierarchical organization)

### Keyword to Skill Mapping (`keyword_skill_map`)

Maps extracted keywords to standardized skills.

Fields:
- `keyword` - The keyword extracted from text
- `skill_type` - 'esco' or 'onet'
- `skill_id` - ID in the corresponding skills table
- `confidence` - Confidence score (0-1)
- `source` - How the mapping was created (manual, automated, etc.)

### Skill Synonyms (`skill_synonyms`)

Alternative names for skills.

Fields:
- `skill_type` - 'esco' or 'onet'
- `skill_id` - ID in the skills table
- `synonym` - Alternative name
- `language` - Language code (e.g., 'en')

### Profiles (`profiles`)

User profiles.

Fields:
- `name` - User's name
- `email` - User's email (unique)
- `uuid` - Unique identifier
- `created_at` - Creation timestamp
- `updated_at` - Last update timestamp

### Profile Skills (`profile_skills`)

Links profiles to standardized skills.

Fields:
- `profile_id` - Reference to profile
- `skill_type` - 'esco' or 'onet'
- `skill_id` - ID in the skills table
- `proficiency` - Proficiency level (Beginner, Intermediate, Advanced, Expert)
- `years_experience` - Years of experience
- `last_verified` - When the skill was last verified
- `notes` - Additional notes

### Profile Custom Skills (`profile_custom_skills`)

Custom skills not in standardized dictionaries.

Fields:
- `profile_id` - Reference to profile
- `name` - Skill name
- `category` - Skill category
- `proficiency` - Proficiency level
- `years_experience` - Years of experience

### Job Offers (`job_offers`)

Job postings to match against.

Fields:
- `title` - Job title
- `company` - Company name
- `location` - Job location
- `description` - Job description
- `requirements` - Job requirements
- `url` - Job posting URL
- `source` - Source (linkedin, indeed, etc.)
- `raw_content` - Original scraped content
- `processed_at` - When the job was processed

### Job Keywords (`job_keywords`)

Keywords extracted from job offers.

Fields:
- `job_offer_id` - Reference to job offer
- `keyword` - Extracted keyword
- `context` - Surrounding text
- `frequency` - How many times it appears
- `relevance` - Relevance score
- `category` - Detected category

### Job Requirements (`job_requirements`)

Structured requirements from job offers.

Fields:
- `job_offer_id` - Reference to job offer
- `requirement_type` - Type (skill, experience, education, certification, language)
- `name` - Requirement name
- `min_years` - Minimum years required
- `is_required` - Whether it's a hard requirement
- `weight` - Importance weight

### Match Results (`match_results`)

Results of matching profiles to job offers.

Fields:
- `profile_id` - Reference to profile
- `job_offer_id` - Reference to job offer
- `overall_score` - Match score (0-100)
- `match_grade` - Letter grade (A, B, C, D, F)
- `total_matches` - Number of matches
- `total_requirements` - Total requirements

### Skill Matches (`skill_matches`)

Individual skill matches.

Fields:
- `match_result_id` - Reference to match result
- `profile_skill_type` - 'esco', 'onet', or 'custom'
- `profile_skill_id` - ID in profile skills or custom skills
- `job_requirement_id` - Reference to job requirement (optional)
- `match_score` - Match score (0-100)
- `match_type` - Type of match (exact, partial, fuzzy)
- `confidence` - Confidence score
- `notes` - Additional notes

### Missing Requirements (`missing_requirements`)

Requirements not matched by the profile.

Fields:
- `match_result_id` - Reference to match result
- `job_requirement_id` - Reference to job requirement
- `name` - Requirement name
- `requirement_type` - Type of requirement
- `priority` - Priority level (High, Medium, Low)
- `asked_user` - Whether user was asked about this
- `user_response` - User's response (yes, no, maybe)
- `user_response_at` - When user responded

### Resumes (`resumes`)

Generated resumes.

Fields:
- `profile_id` - Reference to profile
- `job_offer_id` - Reference to job offer (optional)
- `match_result_id` - Reference to match result (optional)
- `content` - Resume content
- `format` - Format (text, markdown, json, html, pdf, docx)
- `filename` - Filename
- `storage_path` - Storage path
- `created_at` - Creation timestamp

## Running Migrations

Migrations are automatically applied when the backend starts:

```go
// In main.go
sqlRepo, err := repository.NewSQLRepository(cfg.SQLDSN)
if err != nil {
    log.Fatalf("Failed to initialize SQL repository: %v", err)
}
defer sqlRepo.Close()

// Run migrations
if err := sqlRepo.RunMigrations("migrations"); err != nil {
    log.Printf("Warning: failed to run migrations: %v", err)
} else {
    log.Println("SQL migrations applied successfully")
}
```

## Adding New Migrations

1. Create a new file with a sequential number prefix (e.g., `003_new_feature.sql`)
2. Add your SQL statements to the file
3. The migration will be automatically applied on next startup

## Views

The schema includes several helpful views:

- **`all_skills`** - Combined view of ESCO and O*NET skills
- **`profile_skills_with_details`** - Profile skills with skill details
- **`keyword_to_skill_mapping`** - Keyword mappings with skill names

## Sample Data

The `002_sample_data.sql` file includes sample data for:

- Skill categories (Programming, Tools, Databases, Cloud, etc.)
- ESCO skills (Go, Python, Docker, Kubernetes, etc.)
- O*NET skills (Programming, Containerization, etc.)
- Keyword mappings (golang → Go, k8s → Kubernetes, etc.)
- Skill synonyms

This sample data allows you to test the system immediately after startup.

## Importing Real Data

To import real ESCO and O*NET data:

### ESCO

1. Download ESCO data from: https://ec.europa.eu/esco portal
2. Parse the JSON/CSV files
3. Insert into `esco_skills` table using `BulkInsertESCOSkills`

### O*NET

1. Download O*NET data from: https://www.onetcenter.org/taxonomy.html
2. Parse the XML/JSON files
3. Insert into `onet_skills` table using `BulkInsertONETSkills`

Example Go code for bulk insertion:

```go
// For ESCO
skills := []*models.ESCOSkill{...}
if err := repo.BulkInsertESCOSkills(skills); err != nil {
    log.Fatalf("Failed to import ESCO skills: %v", err)
}

// For O*NET
onetSkills := []*models.ONETSkill{...}
if err := repo.BulkInsertONETSkills(onetSkills); err != nil {
    log.Fatalf("Failed to import O*NET skills: %v", err)
}
```

## Next Steps

1. **Import full ESCO and O*NET datasets** - Populate the skill dictionaries
2. **Enhance keyword extraction** - Use the standardized skills for better matching
3. **Improve matching algorithm** - Leverage the hierarchical structure of ESCO/O*NET
4. **Add skill validation** - Validate user skills against standardized dictionaries
5. **Add skill suggestions** - Suggest relevant skills based on profile data
