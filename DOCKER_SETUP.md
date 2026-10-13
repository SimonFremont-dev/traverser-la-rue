# Docker Setup for CV Builder Backend

This document describes the Docker setup for the CV Builder backend API.

## What Was Created

### Directory Structure

```
cv-builder/
├── backend/
│   ├── Dockerfile              # Multi-stage build for Go application
│   ├── .dockerignore           # Files to exclude from Docker build
│   ├── go.mod                  # Go module dependencies
│   ├── go.sum                  # Go dependency checksums
│   ├── main.go                 # Application entry point
│   ├── README.md               # Backend documentation
│   └── internal/
│       ├── api/
│       │   └── handlers.go     # HTTP API handlers
│       ├── config/
│       │   └── config.go       # Configuration management
│       ├── keyword/
│       │   └── extractor.go    # Keyword extraction logic
│       ├── models/
│       │   └── models.go       # Data models
│       ├── repository/
│       │   └── sqlite_repository.go  # SQLite database operations
│       └── service/
│           └── service.go       # Business logic
└── docker-compose.yml          # Multi-container orchestration
```

## Architecture

### Containers

1. **Backend API** (`cv-builder-backend`)
   - Go 1.21 application
   - Exposes REST API on port 8080
   - Uses SQLite for personal data storage
   - Connects to PostgreSQL for job data with ESCO & O*NET dictionaries
   - Built with multi-stage Docker build for minimal image size

2. **PostgreSQL Database** (`cv-builder-sql-db`)
   - PostgreSQL 15 Alpine
   - Stores job offers, ESCO skills, O*NET skills, and matching data
   - Persistent volume for data storage
   - Health checks configured

### Database Structure

The PostgreSQL database contains:

- **ESCO Skills**: `esco_skills` - European Skills/Competences taxonomy
- **O*NET Skills**: `onet_skills` - US Occupational Information Network taxonomy
- **Skill Categories**: `skill_categories` - Classification system
- **Keyword Mappings**: `keyword_skill_map` - Maps extracted keywords to standardized skills
- **Skill Synonyms**: `skill_synonyms` - Alternative names for skills
- **Profiles**: `profiles` - User profiles
- **Job Offers**: `job_offers` - Job postings to match against
- **Matching Results**: `match_results` - Results of profile-job matching

See `backend/migrations/001_initial_schema.sql` for full schema definition.

### Data Flow

```
Frontend (Future) → HTTP Requests → Backend API
                                      │
                                      ├── SQLite (local file)
                                      │     └── Personal data (profiles, skills, experience, education)
                                      │
                                      └── PostgreSQL (networked)
                                            └── Job offers, keywords, matches
```

## Features Implemented

### 1. Profile Management
- Create, read, update profiles
- Store skills, experience, and education
- Persistent in SQLite database

### 2. Keyword Extraction
- Extracts technical keywords from text
- Identifies programming languages, tools, frameworks, databases, etc.
- Filters out common non-technical words
- Returns clean list of keywords

### 3. Skill Matching
- Compares extracted keywords against profile skills
- Identifies matched skills
- Identifies missing skills
- Returns comprehensive matching results

### 4. Resume Generation
- Generates simple text-based resumes
- Includes matched skills prominently
- Shows other profile skills
- Optionally includes unmatched keywords as suggestions

### 5. Job Offer Processing
- Accepts job offer text
- Extracts keywords automatically
- Matches against profile
- Generates resume with matched skills

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/profiles` | List all profiles |
| POST | `/api/profiles` | Create a new profile |
| GET | `/api/profiles/{id}` | Get a specific profile |
| PUT | `/api/profiles/{id}` | Update a profile |
| GET | `/api/profiles/{id}/skills` | List skills for a profile |
| POST | `/api/profiles/{id}/skills` | Add a skill to a profile |
| POST | `/api/extract-keywords` | Extract keywords from text |
| POST | `/api/match` | Match profile against keywords |
| POST | `/api/resume/generate` | Generate a resume |
| POST | `/api/job-offers/process` | Process job offer (extract, match, generate) |
| GET | `/api/health` | Health check |

## Usage

### Starting the Services

```bash
# Navigate to project root
cd cv-builder

# Start all containers
docker-compose up -d

# View logs
docker-compose logs -f backend

# Stop all containers
docker-compose down
```

### Testing the API

```bash
# Health check
curl http://localhost:8080/api/health

# Extract keywords
curl -X POST http://localhost:8080/api/extract-keywords \
  -H "Content-Type: application/json" \
  -d '{"text": "We need a Go developer with Docker and Kubernetes experience"}'

# Process a job offer
curl -X POST "http://localhost:8080/api/job-offers/process?profile_id=1" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Backend Developer",
    "company": "Tech Corp",
    "description": "We need a Go developer with Kubernetes experience",
    "requirements": "5+ years Go, Docker, Kubernetes"
  }'
```

## Environment Variables

### Backend Container

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | HTTP server port | 8080 |
| `SQLITE_PATH` | SQLite database path | /app/data/personal.db |
| `SQL_DSN` | PostgreSQL connection string | postgres://cv_user:cv_password@sql-db:5432/cv_db?sslmode=disable |

### PostgreSQL Container

| Variable | Description | Default |
|----------|-------------|---------|
| `POSTGRES_DB` | Database name | cv_db |
| `POSTGRES_USER` | Database user | cv_user |
| `POSTGRES_PASSWORD` | Database password | cv_password |

## Volumes

1. **SQLite Data**: `./backend/data:/app/data`
   - Persists personal data across container restarts
   - Contains the SQLite database file

2. **Output Files**: `./backend/output:/app/output`
   - Stores generated resumes
   - Accessible from host machine

3. **PostgreSQL Data**: `postgres_data:/var/lib/postgresql/data`
   - Named volume for PostgreSQL data
   - Managed by Docker
   - Persists across container restarts

## Network

- **cv-network**: Custom bridge network
  - Connects backend and sql-db containers
  - Allows containers to communicate using service names
  - Isolated from host network

## Building the Backend Locally

```bash
cd backend

# Initialize dependencies
go mod download

# Build the application
go build -o backend

# Run with environment variables
SQLITE_PATH=./data/personal.db PORT=8080 ./backend
```

## Dockerfile Details

The Dockerfile uses a **multi-stage build**:

1. **Builder Stage** (`golang:1.21-alpine`)
   - Downloads Go dependencies
   - Compiles the application with `CGO_ENABLED=0` for static binary
   - Produces a minimal binary

2. **Runtime Stage** (`alpine:3.19`)
   - Uses lightweight Alpine Linux
   - Copies only the compiled binary
   - Adds CA certificates for HTTPS support
   - Creates necessary directories
   - Sets up health checks

This results in a **small (~15MB) production-ready container image**.

## Database Schema

### SQLite Tables

- **profiles**: User profiles (id, name, email, timestamps)
- **skills**: Profile skills (id, profile_id, name, category, proficiency, years)
- **experience**: Work experience (id, profile_id, company, position, dates, description)
- **education**: Education history (id, profile_id, institution, degree, field, dates, gpa)
- **job_offers**: Processed job offers (id, title, company, description, requirements, location)
- **keywords**: Extracted keywords (id, job_offer_id, keyword, category)
- **resumes**: Generated resumes (id, profile_id, job_offer_id, content, format, created_at)

### Indexes

Indexes are created for:
- Profile email (unique)
- Skill names
- Profile IDs on all related tables
- Keywords for fast searching

## Future Evolution

The backend is designed to be easily extended:

### Adding SQL Database Support
- Create `repository/sql_repository.go`
- Implement the same interface as SQLite repository
- Update service to use both repositories
- Configure connection in `main.go`

### Adding New API Endpoints
- Add new handler functions in `internal/api/handlers.go`
- Add new routes in `SetupRoutes` method
- Add corresponding service methods

### Adding Authentication
- Add middleware for JWT validation
- Protect sensitive endpoints
- Add user authentication endpoints

### Adding More Output Formats
- Extend `GenerateResume` to support PDF, DOCX, HTML
- Add new format-specific generators
- Update API to accept format parameter

## Troubleshooting

### Common Issues

1. **Docker Build Fails**
   - Check Go module dependencies: `go mod tidy`
   - Verify Go version (1.21 required)
   - Check Dockerfile syntax

2. **Database Connection Fails**
   - Verify PostgreSQL container is running: `docker-compose ps`
   - Check PostgreSQL logs: `docker-compose logs sql-db`
   - Test connection manually

3. **SQLite Permission Issues**
   - Ensure volume mount has correct permissions
   - Check container file permissions: `docker exec cv-builder-backend ls -la /app/data`
   - May need to adjust host directory permissions

4. **Port Conflicts**
   - Port 8080 or 5432 may already be in use
   - Change ports in docker-compose.yml
   - Or stop conflicting services

### Debug Commands

```bash
# View running containers
docker-compose ps

# View container logs
docker-compose logs backend

# Enter backend container
docker exec -it cv-builder-backend sh

# Test database connection from backend container
docker exec -it cv-builder-backend nc -zv sql-db 5432

# Check SQLite database
docker exec -it cv-builder-backend sqlite3 /app/data/personal.db ".tables"
```

## Next Steps

1. **Test the backend**: Start containers and verify API endpoints
2. **Add frontend**: Create frontend service in docker-compose
3. **Implement SQL repository**: Add PostgreSQL data access
4. **Add authentication**: Secure API endpoints
5. **Add more features**: PDF generation, better matching algorithms
6. **Add tests**: Unit and integration tests
7. **Add monitoring**: Prometheus metrics, health checks
