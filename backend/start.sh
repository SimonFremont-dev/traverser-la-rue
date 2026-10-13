#!/bin/bash

# CV Builder Start Script
# This script builds and launches the CV Builder application with frontend

set -e

echo "========================================"
echo "  CV Builder - Starting Application"
echo "========================================"
echo ""

# Build the backend
echo "Building backend..."
go build -o bin/backend .
echo "✓ Backend built"
echo ""

# Build the frontend
echo "Building frontend..."
cd web
npm install 2>/dev/null || true
npx tsc static/app.ts --outDir static --target es2015 --module es2015 --moduleResolution node --esModuleInterop --skipLibCheck 2>/dev/null || true
cd ..
echo "✓ Frontend built"
echo ""

# Set default configuration
export PORT=${PORT:-8080}
export SQLITE_PATH=${SQLITE_PATH:-***}

# Start the server
echo "Starting server on port $PORT..."
echo "Access the application at: http://localhost:$PORT"
echo ""

./bin/backend
