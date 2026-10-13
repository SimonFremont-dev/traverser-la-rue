# Frontend Setup for CV Builder

This document describes the frontend UI setup for the CV Builder project.

## Overview

A TypeScript-based frontend has been created that allows users to:

1. **Load Job Offers** - From URL, file upload, or direct text input
2. **Extract Keywords** - Automatically extract technical keywords from job offers
3. **Match Skills** - Match extracted keywords against skills in the PostgreSQL database
4. **Add Custom Skills** - Add unmatched keywords as custom skills to their profile
5. **Generate Resume** - Generate a customized resume with matched skills

## Directory Structure

```
backend/
├── web/
│   ├── static/
│   │   ├── index.html      # Main HTML page
│   │   ├── app.ts          # TypeScript source
│   │   ├── app.js          # Compiled JavaScript
│   │   └── styles.css      # CSS styles
│   └── package.json        # npm dependencies
├── internal/
│   └── api/
│       └── handlers.go     # API handlers (updated to serve frontend)
└── main.go               # Main entry point
```

## API Endpoints Used

The frontend uses the following API endpoints:

- `GET /api/profiles` - List all profiles
- `GET /api/profiles/{id}/skills` - Get profile skills
- `POST /api/profiles/{id}/skills` - Add skill to profile
- `POST /api/extract-keywords` - Extract keywords from text
- `POST /api/match` - Match profile against keywords
- `POST /api/resume/generate` - Generate resume
- `GET /api/skills` - Get all skills from database
- `GET /api/health` - Health check

## How to Run

### 1. Build and Start the Server

```bash
# Build the backend and frontend
make all

# Or manually:
make build-frontend  # Build TypeScript frontend
make build           # Build Go backend

# Start the server
export PORT=8080
export SQLITE_PATH=./data/cv.db
./bin/backend
```

### 2. Access the Application

Open your browser and navigate to: `http://localhost:8080`

### 3. Using the Application

#### Step 1: Load Job Offer
- Enter a job offer URL and click "Load from URL"
- Or upload a text file containing the job offer
- Or paste the job offer text directly into the textarea
- Click "Process Job Offer"

#### Step 2: Review Extracted Keywords
- View the keywords extracted from the job offer
- Click "Continue to Skill Matching"

#### Step 3: Match Skills
- Select your profile from the dropdown
- View skills from the database
- View matched keywords (in green)
- View unmatched keywords (in red) with option to add as custom skill
- Click "Generate Resume"

#### Step 4: Review Resume
- View the generated resume
- Download as text file
- Copy to clipboard
- Start over with a new job offer

## TypeScript Development

### Install Dependencies

```bash
cd web
npm install
```

### Build Frontend

```bash
npx tsc static/app.ts --outDir static --target es2015 --module es2015 --moduleResolution node --esModuleInterop --skipLibCheck
```

### Watch Mode (for development)

```bash
npx tsc static/app.ts --outDir static --target es2015 --module es2015 --moduleResolution node --esModuleInterop --skipLibCheck --watch
```

## Backend Changes

The following changes were made to the backend to support the frontend:

1. **Added static file serving** in `internal/api/handlers.go`:
   - Serve static files from `/static/` path
   - Serve `index.html` for all other GET requests (SPA support)

2. **Added skills endpoint** in `internal/api/handlers.go`:
   - `GET /api/skills` - Returns all skills from PostgreSQL database

3. **Updated service** in `internal/service/service.go`:
   - Added `sqlRepo` field to Service struct
   - Added `GetAllSkills()` method to fetch skills from SQL database
   - Updated NewService to accept SQLRepository parameter

4. **Updated main.go**:
   - Pass sqlRepo to NewService

## Notes

- The frontend assumes PostgreSQL is configured with the skill database
- If PostgreSQL is not configured, the skills section will show "No skills in database"
- Skills can still be added as custom skills to profiles even without PostgreSQL
- The frontend uses local storage for temporary state during the session

## Troubleshooting

### Frontend not loading
- Ensure static files are in `web/static/` directory
- Check that the server is running on the correct port
- Verify that the static file handler is registered in the router

### Skills not showing
- Check that PostgreSQL is running and configured
- Verify that migrations have been applied
- Check that the SQLDSN environment variable is set correctly

### TypeScript compilation errors
- Run `npm install` in the web directory
- Ensure TypeScript is installed globally or in the project
- Check for syntax errors in app.ts

## Future Enhancements

- Add authentication and user sessions
- Save job offers to database
- Save match results to database
- Add rich text editing for resume
- Export resume in multiple formats (PDF, DOCX)
- Add search and filtering for skills
- Add profile management UI
