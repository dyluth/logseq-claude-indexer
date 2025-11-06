.PHONY: build test install clean lint demo help

# Build for current platform
build:
	@echo "Building logseq-claude-indexer..."
	@mkdir -p bin
	go build -o bin/logseq-claude-indexer ./cmd/logseq-claude-indexer

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p bin
	GOOS=darwin GOARCH=amd64 go build -o bin/logseq-claude-indexer-darwin-amd64 ./cmd/logseq-claude-indexer
	GOOS=darwin GOARCH=arm64 go build -o bin/logseq-claude-indexer-darwin-arm64 ./cmd/logseq-claude-indexer
	GOOS=linux GOARCH=amd64 go build -o bin/logseq-claude-indexer-linux-amd64 ./cmd/logseq-claude-indexer
	GOOS=windows GOARCH=amd64 go build -o bin/logseq-claude-indexer-windows-amd64.exe ./cmd/logseq-claude-indexer

# Run tests
test:
	@echo "Running tests..."
	go test -v -race -cover ./...

# Run tests with coverage report
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Install to GOPATH/bin
install:
	@echo "Installing logseq-claude-indexer..."
	go install ./cmd/logseq-claude-indexer

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	find testdata/fixtures -type d -name .claude -exec rm -rf {} + 2>/dev/null || true
	go clean

# Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: https://golangci-lint.run/usage/install/"; \
		exit 1; \
	fi

# Run on example repo (test fixtures)
demo:
	@echo "Running on test fixtures..."
	@make build
	./bin/logseq-claude-indexer generate --repo ./testdata/fixtures --verbose
	@echo "\nGenerated files:"
	@ls -lh testdata/fixtures/.claude/indexes/

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	go mod tidy

# Run all checks (fmt, lint, test)
check: fmt lint test
	@echo "All checks passed!"

# Show help
help:
	@echo "Available targets:"
	@echo "  build         - Build for current platform"
	@echo "  build-all     - Build for multiple platforms"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  install       - Install to GOPATH/bin"
	@echo "  clean         - Clean build artifacts"
	@echo "  lint          - Run linter"
	@echo "  demo          - Run on test fixtures"
	@echo "  fmt           - Format code"
	@echo "  tidy          - Tidy dependencies"
	@echo "  check         - Run all checks (fmt, lint, test)"
	@echo "  help          - Show this help message"
