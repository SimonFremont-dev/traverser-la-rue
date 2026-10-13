# CV Builder - User Guide

## Overview

CV Builder is a tool that helps you create a customized resume by matching your skills against job offers. It extracts keywords from job descriptions, matches them against skills in our database (ESCO and O*NET taxonomies), and generates a tailored resume.

## Quick Start

1. **Build and run the application:**
   ```bash
   cd /path/to/cv-builder/backend
   make all
   export PORT=8080
   export SQLITE_PATH=./data/cv.db
   ./bin/backend
   ```

2. **Access the web interface:**
   Open your browser and go to: `http://localhost:8080`

## Using the Application

### Step 1: Load a Job Offer

You have three options to load a job offer:

**Option A: Enter a URL**
- Paste the job offer URL in the "Job Offer URL" field
- Click "Load from URL"
- **Note:** LinkedIn and some other job sites have anti-scraping protection. If you get an error with LinkedIn URLs, use Option B or C instead.

**Option B: Upload a File**
- Click "Choose File" and select a text file (.txt, .md, etc.) containing the job description
- Click "Upload File"

**Option C: Paste Text Directly**
- Copy the job description text from any source
- Paste it into the large text area
- This works for all job offers, including LinkedIn

Once loaded, click the **"Process Job Offer"** button.

### Step 2: Review Extracted Keywords

The application will:
- Extract technical keywords from the job description
- Display them as tags on the page

Review the keywords to ensure they're relevant. If you notice missing keywords, you can go back and edit the text.

Click **"Continue to Skill Matching"** to proceed.

### Step 3: Match Skills

1. **Select Your Profile**
   - Choose your profile from the dropdown menu
   - If no profiles exist, you'll need to create one via the API

2. **View Database Skills**
   - The application will display skills from our database that match the keywords
   - Skills are categorized by source (ESCO, O*NET) and type

3. **Review Matched Keywords**
   - **Green (Matched):** Keywords found in our skills database
   - **Red (Unmatched):** Keywords not found in our database

4. **Add Custom Skills**
   - For unmatched keywords, click the "Add to Profile" button
   - This adds the keyword as a custom skill to your profile
   - You can set the proficiency level and category later via the API

### Step 4: Generate Resume

Click **"Generate Resume"** to create your customized resume.

The resume will include:
- Your profile information
- Matched skills from the job offer
- Custom skills you added
- Relevant experience and education

**Options:**
- **Download Resume:** Save the resume as a text file
- **Copy to Clipboard:** Copy the resume text for pasting elsewhere
- **Start Over:** Begin a new job offer analysis

## Managing Profiles

### Creating a Profile

Currently, profiles need to be created via the API. Here's how:

```bash
curl -X POST http://localhost:8080/api/profiles \
  -H "Content-Type: application/json" \
  -d '{"name": "Your Name", "email": "your@email.com"}'
```

### Adding Skills to Profile

Skills are added automatically when you click "Add to Profile" in the UI. You can also add them via API:

```bash
curl -X POST http://localhost:8080/api/profiles/1/skills \
  -H "Content-Type: application/json" \
  -d '{"name": "Skill Name", "category": "Programming", "proficiency": "Advanced", "years": 5}'
```

## Database Configuration

### SQLite (Personal Data)
- Stores your profiles, custom skills, and generated resumes
- Configured via `SQLITE_PATH` environment variable
- Default: `./data/cv.db`

### PostgreSQL (Skills Database)
- Stores standardized skills from ESCO and O*NET taxonomies
- Configured via `SQLDSN` environment variable
- Example: `postgres://user:password@localhost:5432/cv_db?sslmode=disable`

### Running Migrations

If you're using PostgreSQL, apply the database migrations:

```bash
# Using Atlas
atlas schema apply --url file://migrations --to postgres://user:password@localhost:5432/cv_db

# Or use Make
make migrate DATABASE_URL=postgres://user:password@localhost:5432/cv_db
```

## Troubleshooting

### LinkedIn URLs Don't Work

LinkedIn actively blocks automated fetching. **Solution:** Copy the job description text directly and paste it into the text area.

### No Skills Showing

1. Check that PostgreSQL is running and configured
2. Verify the `SQLDSN` environment variable is set correctly
3. Ensure migrations have been applied
4. If still not working, the frontend will show "No skills in database" - this is expected if PostgreSQL isn't configured

### URL Fetching Fails

Some websites block automated requests. For these sites:
- Copy the job description text
- Paste it directly into the text area
- Click "Process Job Offer"

### Error: "No profiles found"

You need to create a profile first:
```bash
curl -X POST http://localhost:8080/api/profiles \
  -H "Content-Type: application/json" \
  -d '{"name": "Your Name", "email": "your@email.com"}'
```

Then refresh the page.

## Tips for Best Results

1. **Use Clean Text:** When pasting job descriptions, remove any formatting or HTML tags
2. **Check Keywords:** Review the extracted keywords and add any missing ones as custom skills
3. **Keep Profiles Updated:** Regularly add new skills to your profile
4. **Use Multiple Profiles:** Create different profiles for different job types

## API Reference

### Profiles
- `GET /api/profiles` - List all profiles
- `GET /api/profiles/{id}` - Get a specific profile
- `POST /api/profiles` - Create a new profile
- `PUT /api/profiles/{id}` - Update a profile

### Skills
- `GET /api/profiles/{id}/skills` - Get profile skills
- `POST /api/profiles/{id}/skills` - Add skill to profile
- `GET /api/skills` - Get all skills from database

### Processing
- `POST /api/extract-keywords` - Extract keywords from text
- `POST /api/match` - Match profile against keywords
- `POST /api/resume/generate` - Generate resume
- `POST /api/job-offers/process` - Process a job offer

### Utilities
- `GET /api/fetch-url?url={url}` - Fetch content from a URL
- `GET /api/health` - Health check
- `GET /` or `/static/` - Serve frontend

## Example Workflow

```bash
# Start server
./bin/backend &

# Create a profile
curl -X POST http://localhost:8080/api/profiles \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com"}'

# Add some skills
curl -X POST http://localhost:8080/api/profiles/1/skills \
  -H "Content-Type: application/json" \
  -d '{"name": "Go", "category": "Programming", "proficiency": "Advanced", "years": 5}'

curl -X POST http://localhost:8080/api/profiles/1/skills \
  -H "Content-Type: application/json" \
  -d '{"name": "Docker", "category": "DevOps", "proficiency": "Intermediate", "years": 3}'

# Now use the web interface at http://localhost:8080
# Load a job offer, process it, and generate your resume!
```
