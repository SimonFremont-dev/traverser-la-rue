# Makefile for cv-builder

BINARY_NAME := resume-adapter
BINARY_PATH := ./bin/$(BINARY_NAME)
SOURCE_PATH := ./cmd/resume-adapter

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p bin
	go build -o $(BINARY_PATH) $(SOURCE_PATH)
	@echo "Build complete: $(BINARY_PATH)"

# Install the application
install:
	@echo "Installing $(BINARY_NAME)..."
	go install $(SOURCE_PATH)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/

# Run the application
run: build
	@$(BINARY_PATH)

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	go mod tidy

# Generate a new resume template
generate:
	@$(BINARY_PATH) generate --output my-resume.yaml --format yaml

# Match resume against job offer
match:
	@$(BINARY_PATH) match --resume examples/resume.json --offer examples/job_offer.txt

# Analyze resume against job offer
analyze:
	@$(BINARY_PATH) analyze --resume examples/resume.json --offer examples/job_offer.txt

# Adapt resume for job offer
adapt:
	@$(BINARY_PATH) adapt --resume examples/resume.json --offer examples/job_offer.txt --output adapted-resume.md --format markdown

# Build all
all: deps build test

.PHONY: build install clean run test deps tidy generate match analyze adapt all
