# CV Builder Frontend

This directory contains the frontend UI for the CV Builder application.

## Structure

- `static/` - Static files served by the backend
  - `index.html` - Main HTML page
  - `styles.css` - CSS styles
  - `app.ts` - TypeScript source code
  - `app.js` - Compiled JavaScript (generated from app.ts)

## Development

### Install Dependencies

```bash
cd web
npm install
```

### Build Frontend

```bash
# Compile TypeScript to JavaScript
npx tsc static/app.ts --outDir static --target es2015 --module es2015 --moduleResolution node --esModuleInterop --skipLibCheck
```

Or use the Makefile target:

```bash
cd ..
make build-frontend
```

### Watch Mode

For development, you can run TypeScript in watch mode:

```bash
cd web
npx tsc static/app.ts --outDir static --target es2015 --module es2015 --moduleResolution node --esModuleInterop --skipLibCheck --watch
```

## Features

The frontend provides a 4-step workflow:

1. **Load Job Offer**: Load a job offer from URL, upload a text file, or paste text directly
2. **Extracted Keywords**: View keywords extracted from the job offer
3. **Skill Matching**: Match extracted keywords against skills in the database and profile
4. **Resume Generation**: Generate and download a customized resume

## API Endpoints Used

- `GET /api/profiles` - List all profiles
- `GET /api/profiles/{id}/skills` - Get profile skills
- `POST /api/profiles/{id}/skills` - Add skill to profile
- `POST /api/extract-keywords` - Extract keywords from text
- `POST /api/match` - Match profile against keywords
- `POST /api/resume/generate` - Generate resume
- `GET /api/skills` - Get all skills from database
- `GET /api/health` - Health check
