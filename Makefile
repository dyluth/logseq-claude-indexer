.PHONY: build test install install-user setup-git-hook test-datasets clean lint demo help

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
	@echo "✓ Installed to GOPATH/bin"
	@echo ""
	@echo "To set up git hooks in your Logseq repo:"
	@echo "  cd /path/to/your/logseq/repo"
	@echo "  make -C /path/to/logseq-claude-indexer setup-git-hook LOGSEQ_REPO=."

# Install to ~/.local/bin (user-local installation)
install-user: build
	@echo "Installing logseq-claude-indexer to ~/.local/bin..."
	@mkdir -p ~/.local/bin
	@cp bin/logseq-claude-indexer ~/.local/bin/
	@chmod +x ~/.local/bin/logseq-claude-indexer
	@echo "✓ Installed to ~/.local/bin/logseq-claude-indexer"
	@echo ""
	@echo "Make sure ~/.local/bin is in your PATH. Add this to ~/.bashrc or ~/.zshrc:"
	@echo "  export PATH=\"\$$HOME/.local/bin:\$$PATH\""
	@echo ""
	@echo "Next steps:"
	@echo ""
	@echo "1. Set up git hook in your Logseq repo:"
	@echo ""
	@echo "   cd /path/to/your/logseq/repo"
	@echo "   cat > .git/hooks/post-commit << 'EOF'"
	@echo "   #!/bin/bash"
	@echo "   logseq-claude-indexer generate --repo . --output .claude/indexes --quiet"
	@echo "   EOF"
	@echo "   chmod +x .git/hooks/post-commit"
	@echo ""
	@echo "2. Generate initial indexes:"
	@echo ""
	@echo "   cd /path/to/your/logseq/repo"
	@echo "   logseq-claude-indexer generate --repo . --output .claude/indexes"
	@echo ""
	@echo "For detailed instructions, see INSTALL.md"

# Setup git hook in a Logseq repository
# Usage: make setup-git-hook LOGSEQ_REPO=/path/to/logseq/repo
setup-git-hook:
	@if [ -z "$(LOGSEQ_REPO)" ]; then \
		echo "Error: LOGSEQ_REPO not specified."; \
		echo "Usage: make setup-git-hook LOGSEQ_REPO=/path/to/your/logseq/repo"; \
		exit 1; \
	fi
	@if [ ! -d "$(LOGSEQ_REPO)/.git" ]; then \
		echo "Error: $(LOGSEQ_REPO) is not a git repository"; \
		exit 1; \
	fi
	@echo "Setting up post-commit hook in $(LOGSEQ_REPO)..."
	@echo '#!/bin/bash' > $(LOGSEQ_REPO)/.git/hooks/post-commit
	@echo '# Auto-generate Claude indexes after each commit' >> $(LOGSEQ_REPO)/.git/hooks/post-commit
	@echo 'logseq-claude-indexer generate --repo . --output .claude/indexes --quiet' >> $(LOGSEQ_REPO)/.git/hooks/post-commit
	@chmod +x $(LOGSEQ_REPO)/.git/hooks/post-commit
	@echo "✓ Git hook installed at $(LOGSEQ_REPO)/.git/hooks/post-commit"
	@echo ""
	@echo "The indexer will now run automatically after each commit."
	@echo "To run manually: cd $(LOGSEQ_REPO) && logseq-claude-indexer generate --repo . --output .claude/indexes"

# Test on all datasets
test-datasets: build
	@echo "Testing on all datasets..."
	@echo "\n=== Testing on fixtures ==="
	@./bin/logseq-claude-indexer generate --repo testdata/fixtures --output .claude/indexes
	@echo "\n=== Testing on synthetic ==="
	@./bin/logseq-claude-indexer generate --repo testdata/synthetic --output .claude/indexes
	@echo "\n=== Testing on user-provided ==="
	@./bin/logseq-claude-indexer generate --repo testdata/user-provided --output .claude/indexes
	@echo "\n✓ All datasets tested successfully"
	@echo "\nGenerated index files:"
	@ls testdata/user-provided/.claude/indexes/

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
	@echo "Logseq Claude Indexer - Makefile Commands"
	@echo ""
	@echo "Build & Installation:"
	@echo "  build              - Build for current platform"
	@echo "  build-all          - Build for multiple platforms"
	@echo "  install            - Install to GOPATH/bin"
	@echo "  install-user       - Install to ~/.local/bin (recommended for users)"
	@echo "  clean              - Clean build artifacts"
	@echo ""
	@echo "Git Hook Setup:"
	@echo "  setup-git-hook     - Install post-commit hook in Logseq repo"
	@echo "                       Usage: make setup-git-hook LOGSEQ_REPO=/path/to/repo"
	@echo ""
	@echo "Testing:"
	@echo "  test               - Run unit tests"
	@echo "  test-coverage      - Run tests with coverage report"
	@echo "  test-datasets      - Test on all datasets (fixtures, synthetic, user-provided)"
	@echo "  demo               - Run on test fixtures"
	@echo ""
	@echo "Code Quality:"
	@echo "  fmt                - Format code"
	@echo "  lint               - Run linter (requires golangci-lint)"
	@echo "  tidy               - Tidy dependencies"
	@echo "  check              - Run all checks (fmt, lint, test)"
	@echo ""
	@echo "Quick Start:"
	@echo "  1. Build:          make build"
	@echo "  2. Test:           make test"
	@echo "  3. Install:        make install-user"
	@echo "  4. Setup hook:     make setup-git-hook LOGSEQ_REPO=/path/to/your/logseq"
	@echo ""
	@echo "For more information, see README.md"
