# Build stage
FROM golang:1.21-alpine AS builder

# Install gcc for CGO
RUN apk add --no-cache gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with CGO enabled
RUN CGO_ENABLED=1 GOOS=linux go build -o /backend

# Runtime stage
FROM alpine:3.19

# Install CA certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Set working directory
WORKDIR /app

# Create data directory for SQLite database
RUN mkdir -p /app/data

# Copy binary from builder
COPY --from=builder /backend /app/backend

# Copy static files if any (for future use)
COPY --from=builder /app/web/static /app/static

# Set environment variables
ENV PORT=8080
ENV SQLITE_PATH=/app/data/personal.db
ENV SQL_DSN=""

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
	CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/health || exit 1

# Run the application
CMD ["/app/backend"]
