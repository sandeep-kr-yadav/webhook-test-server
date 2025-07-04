#!/bin/bash

echo "=== Webhook Test Environment ==="
echo "Starting webhook test server..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go 1.21 or later."
    exit 1
fi

# Install dependencies
echo "Installing dependencies..."
go mod tidy

# Create sample files if they don't exist
if [ ! -f "test-files/sample.pdf" ]; then
    echo "Creating sample files..."
    make samples
fi

# Start the server
echo "Starting server on http://localhost:8080"
echo "Press Ctrl+C to stop"
echo ""
go run cmd/webhook-test-server.go 