# Makefile for CV Builder Backend

.PHONY: all build test run clean migrate migrate-down lint fmt frontend build-frontend

# Go parameters
BINARY_NAME := backend
BINARY_PATH := ./bin/$(BINARY_NAME)
SOURCE_PATH := .

# Frontend parameters
FRONTEND_DIR := web
STATIC_DIR := $(FRONTEND_DIR)/static

# Database parameters
DATABASE_URL ?= postgres://cv_user:cv_password@localhost:5432/cv_db?sslmode=disable
MIGRATIONS_DIR ?= migrations

# Atlas parameters
ATLAS_URL ?= file://$(MIGRATIONS_DIR)

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p bin
	go build -o $(BINARY_PATH) $(SOURCE_PATH)
	@echo "Build complete: $(BINARY_PATH)"

# Run the application
run: build build-frontend
	@$(BINARY_PATH)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Download Go dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download

# Tidy Go dependencies
tidy:
	@echo "Tidying dependencies..."
	go mod tidy

# Lint code
lint:
	@echo "Running linter..."
	golangci-lint run ./...

# Format code
fmt:
	@echo "Formatting code..."
	gofmt -w .

# ============================================
# DATABASE MIGRATIONS (Atlas)
# ============================================

# Show migration status
migrate-status:
	@echo "Migration status:"
	@atlas schema apply --url $(ATLAS_URL) --to $(DATABASE_URL) --dry-run

# Apply all pending migrations
migrate:
	@echo "Applying migrations..."
	@atlas schema apply --url $(ATLAS_URL) --to $(DATABASE_URL)

# Apply specific migration (use: make migrate-to VERSION=xxx)
migrate-to:
	@echo "Applying migrations up to $(VERSION)..."
	@atlas schema apply --url $(ATLAS_URL) --to $(DATABASE_URL) --to $(VERSION)

# Create a new migration file (use: make migrate-new NAME=xxx)
migrate-new:
	@echo "Creating new migration: $(NAME)"
	@atlas migrate diff --url $(ATLAS_URL) --to $(DATABASE_URL) $(NAME)

# Revert the last migration
migrate-down:
	@echo "Reverting last migration..."
	@atlas schema apply --url $(ATLAS_URL) --to $(DATABASE_URL) --down 1

# Revert all migrations
migrate-reset:
	@echo "Reverting all migrations..."
	@atlas schema apply --url $(ATLAS_URL) --to $(DATABASE_URL) --down -1

# Validate migrations
migrate-validate:
	@echo "Validating migrations..."
	@atlas schema validate --url $(ATLAS_URL) --to $(DATABASE_URL)

# ============================================
# DATABASE OPERATIONS
# ============================================

# Connect to database
db-connect:
	@echo "Connecting to database..."
	psql $(DATABASE_URL)

# Reset database (DANGEROUS!)
db-reset:
	@echo "WARNING: This will drop and recreate the database!"
	@psql $(DATABASE_URL) -c "DROP DATABASE IF EXISTS cv_db; CREATE DATABASE cv_db;"
	@$(MAKE) migrate

# Seed database with test data
db-seed:
	@echo "Seeding database..."
	@psql $(DATABASE_URL) -f $(MIGRATIONS_DIR)/seed_data.sql

# ============================================
# DOCKER OPERATIONS
# ============================================

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t cv-builder-backend .

# Run Docker container
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 --rm cv-builder-backend

# Run Docker with database
docker-run-full:
	@echo "Running Docker with database..."
	cd .. && docker-compose up

# Stop Docker containers
docker-stop:
	@echo "Stopping Docker containers..."
	cd .. && docker-compose down

# ============================================
# ALL TARGETS
# ============================================

all: deps build-frontend build test

# Build frontend (TypeScript)
build-frontend:
	@echo "Building frontend..."
	cd $(FRONTEND_DIR) && npm install
	cd $(FRONTEND_DIR) && npx tsc static/app.ts --outDir static --target es2015 --module es2015 --moduleResolution node --esModuleInterop --skipLibCheck

