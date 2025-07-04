.PHONY: build run test clean help

# Default target
help:
	@echo "Webhook Test Environment - Available commands:"
	@echo "  build    - Build the webhook server"
	@echo "  run      - Run the webhook server"
	@echo "  test     - Run all test cases"
	@echo "  clean    - Clean build artifacts"
	@echo "  help     - Show this help message"

# Build the webhook server
build:
	@echo "Building webhook server..."
	go build -o webhook-server cmd/webhook-test-server.go
	@echo "Build complete: webhook-server"

# Run the webhook server
run:
	@echo "Starting webhook server..."
	go run cmd/webhook-test-server.go

# Run all test cases
test:
	@echo "Running test cases..."
	@cd test-files && ./test-multipart-cases.sh

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f webhook-server
	rm -f webhook-server.log
	@echo "Clean complete"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	@echo "Dependencies installed"

# Create sample files for testing
samples:
	@echo "Creating sample files for testing..."
	@cd test-files && \
	echo "Sample PDF content" > sample.pdf && \
	echo "Sample PNG content" > sample.png && \
	echo "name,value\ntest,123" > sample.csv && \
	echo "Sample Excel content" > sample.xlsx
	@echo "Sample files created"

# Setup complete environment
setup: deps samples
	@echo "Environment setup complete!"
	@echo "Run 'make run' to start the server"
	@echo "Run 'make test' to run test cases" 