# CV Builder Backend API

A modular backend API for the CV Builder project that matches profiles to job offers and generates resumes.

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        Backend API                              │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────────┐ │
│  │   Handlers  │    │  Services   │    │  Repositories    │ │
│  │  (HTTP API) │───▶│ (Business   │───▶│  (Data Access)   │ │
│  │             │    │   Logic)    │    │                 │ │
│  └─────────────┘    └─────────────┘    └─────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      Data Storage                               │
│  ┌─────────────────────┐    ┌─────────────────────────────┐ │
│  │    SQLite            │    │          SQL Database          │ │
│  │  (Personal Data)     │    │    (Job Offers, etc.)         │ │
│  │  - Profiles          │    │  - Job Offers                  │ │
│  │  - Skills            │    │  - Keywords                   │ │
│  │  - Experience        │    │  - Matches                    │ │
│  │  - Education         │    │                             │ │
│  └─────────────────────┘    └─────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

## Features

- **Profile Management**: Create, read, update profiles with skills, experience, and education
- **Keyword Extraction**: Extract technical keywords from job descriptions
- **Skill Matching**: Match extracted keywords against profile skills
- **Resume Generation**: Generate text-based resumes with matched skills
- **Interactive Skill Addition**: Ask users about missing skills and add to profile

## API Endpoints

### Profile Endpoints
- `GET /api/profiles` - List all profiles
- `POST /api/profiles` - Create a new profile
- `GET /api/profiles/{id}` - Get a specific profile
- `PUT /api/profiles/{id}` - Update a profile

### Skill Endpoints
- `GET /api/profiles/{id}/skills` - List skills for a profile
- `POST /api/profiles/{id}/skills` - Add a skill to a profile

### Keyword Endpoints
- `POST /api/extract-keywords` - Extract keywords from text

### Matching Endpoints
- `POST /api/match` - Match a profile against keywords

### Resume Endpoints
- `POST /api/resume/generate` - Generate a resume

### Job Offer Endpoints
- `POST /api/job-offers/process` - Process a job offer (extract, match, generate)

### Health
- `GET /api/health` - Health check

## Quick Start

### Prerequisites
- Docker
- Docker Compose
- Go 1.21+ (for local development)

### Using Docker Compose

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f backend

# Stop all services
docker-compose down
```

### Local Development

```bash
# Navigate to backend directory
cd backend

# Initialize Go module
go mod download

# Run the server
SQLITE_PATH=./data/personal.db PORT=8080 go run .

# Or build and run
go build -o backend
SQLITE_PATH=./data/personal.db PORT=8080 ./backend
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | HTTP server port | 8080 |
| `SQLITE_PATH` | Path to SQLite database | ./data/personal.db |
| `SQL_DSN` | SQL database connection string | (empty) |

## Database Structure

### SQLite (Personal Data)
- `profiles` - User profiles
- `skills` - User skills
- `experience` - Work experience
- `education` - Education history
- `job_offers` - Processed job offers
- `keywords` - Extracted keywords
- `resumes` - Generated resumes

### SQL Database (Job Data)
- Will be designed later for job offers and matching data

## Configuration

Create a `config.json` file:

```json
{
  "port": "8080",
  "sqlite_path": "./data/personal.db",
  "sql_dsn": "postgres://user:password@localhost:5432/db?sslmode=disable",
  "keyword_categories": ["programming", "tools", "frameworks", "databases"]
}
```

## Testing the API

### Extract Keywords

```bash
curl -X POST http://localhost:8080/api/extract-keywords \
  -H "Content-Type: application/json" \
  -d '{"text": "We need a Go developer with Docker and Kubernetes experience"}'
```

### Create Profile

```bash
curl -X POST http://localhost:8080/api/profiles \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "skills": [
      {"name": "Go", "category": "programming", "proficiency": "Advanced", "years": 5},
      {"name": "Docker", "category": "tools", "proficiency": "Intermediate", "years": 3}
    ],
    "experience": [
      {"company": "Tech Corp", "position": "Senior Developer", "start_date": "2020-01-01", "current": true}
    ]
  }'
```

### Process Job Offer

```bash
curl -X POST "http://localhost:8080/api/job-offers/process?profile_id=1" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Backend Developer",
    "company": "Awesome Inc",
    "description": "We need a Go developer with Kubernetes experience",
    "requirements": "5+ years Go, Docker, Kubernetes, PostgreSQL"
  }'
```

## Project Structure

```
backend/
├── Dockerfile              # Docker build configuration
├── go.mod                  # Go module definition
├── main.go                 # Application entry point
├── README.md               # This file
└── internal/
    ├── api/                # HTTP handlers
    │   └── handlers.go
    ├── config/            # Configuration management
    │   └── config.go
    ├── keyword/           # Keyword extraction
    │   └── extractor.go
    ├── models/            # Data models
    │   └── models.go
    ├── repository/        # Database access
    │   └── sqlite_repository.go
    └── service/           # Business logic
        └── service.go
```

## Future Enhancements

- [ ] Add SQL database repository implementation
- [ ] Add authentication/authorization
- [ ] Add rate limiting
- [ ] Add OpenAPI/Swagger documentation
- [ ] Add unit and integration tests
- [ ] Add logging middleware
- [ ] Add metrics and monitoring
- [ ] Add API versioning
- [ ] Add request/response validation
