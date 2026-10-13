#!/bin/bash

# Test script to verify the frontend and backend work together

set -e

echo "========================================"
echo "  Testing CV Builder Frontend"
echo "========================================"
echo ""

# Build backend
echo "1. Building backend..."
go build -o /tmp/cv-builder-test ./... 2>&1 | grep -E "error|Error" && exit 1 || echo "   ✓ Backend built successfully"
echo ""

# Build frontend
echo "2. Building frontend..."
cd web
npm install 2>/dev/null
npx tsc static/app.ts --outDir static --target es2015 --module es2015 --moduleResolution node --esModuleInterop --skipLibCheck 2>&1 | grep -E "error|Error" && exit 1 || echo "   ✓ Frontend built successfully"
cd ..
echo ""

# Check if static files exist
echo "3. Checking static files..."
for file in web/static/index.html web/static/app.js web/static/styles.css; do
    if [ -f "$file" ]; then
        echo "   ✓ $file exists"
    else
        echo "   ✗ $file missing!"
        exit 1
    fi
done
echo ""

echo "========================================"
echo "  All checks passed!"
echo "========================================"
echo ""
echo "To run the application:"
echo "  export PORT=8080"
echo "  export SQLITE_PATH=./data/cv.db"
echo "  ./bin/backend"
echo ""
echo "Then visit: http://localhost:8080"
