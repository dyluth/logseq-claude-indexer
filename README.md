# logseq-claude-indexer

A production-ready Golang CLI tool that scans Logseq knowledge bases and generates Claude Code-optimized index files.
This allows claude code to be run directly inside your logseq folder (or just in the folder with the indexes) and help with your workflows.


## Features

- **Fast**: Scans 100+ files in <1 second
- **Task Extraction**: Finds all NOW/LATER/TODO/DOING/DONE tasks with context
- **Time Tracking**: Parses Logseq's `:LOGBOOK:` CLOCK entries
- **Reference Graph**: Builds a network of `[[page links]]`
- **AI-Optimized**: Generates markdown indexes perfect for Claude Code
- **Git Integration**: Ready for git post-commit hooks

## Installation

### From Source

```bash
git clone https://github.com/dyluth/logseq-claude-indexer.git
cd logseq-claude-indexer
make install
```

### Via Go Install

```bash
go install github.com/dyluth/logseq-claude-indexer/cmd/logseq-claude-indexer@latest
```

### Pre-built Binaries

Download from the [releases page](https://github.com/dyluth/logseq-claude-indexer/releases).

## Quick Start

```bash
# Generate indexes for your Logseq repo
logseq-claude-indexer generate --repo /path/to/logseq

# View the generated indexes
cat /path/to/logseq/.claude/indexes/task-index.md
cat /path/to/logseq/.claude/indexes/reference-graph.md
```

## Usage

### Basic Commands

```bash
# Generate all indexes
logseq-claude-indexer generate --repo /path/to/logseq

# Show verbose output
logseq-claude-indexer generate --repo /path/to/logseq --verbose

# Dry run (show what would be generated)
logseq-claude-indexer generate --repo /path/to/logseq --dry-run

# Custom output directory
logseq-claude-indexer generate --repo /path/to/logseq --output /custom/path

# Silent mode (for git hooks)
logseq-claude-indexer generate --repo /path/to/logseq --quiet

# Show version
logseq-claude-indexer version
```

### Flags

- `--repo` - Path to Logseq repository (default: current directory)
- `--output` - Output directory for indexes (default: `.claude/indexes`)
- `--verbose` - Show detailed logging
- `--quiet` - Suppress output (useful for git hooks)
- `--dry-run` - Preview without writing files

## Generated Indexes

### Task Index (`task-index.md`)

Organized view of all tasks with:
- Tasks grouped by status (NOW, DOING, TODO, LATER, DONE)
- Time tracking from `:LOGBOOK:` entries
- Page references and file locations
- Project summaries

Example:
```markdown
## NOW Tasks (3)

### [[Hearth Insights]] - setup actions
- **File**: `journals/2025_04_06.md:1`
- **References**: [[Hearth Insights]]
- **Time Logged**: 260h 46m
- **Last Activity**: 2025-04-17 09:17
```

### Reference Graph (`reference-graph.md`)

Network view of page connections:
- Hub pages (most referenced)
- Inbound and outbound references for each page
- Identifies non-existent pages
- Sorted by connection strength

Example:
```markdown
## Hub Pages (Most Referenced)

1. **[[Hearth Insights]]** - 12 inbound references
   - File: `pages/Hearth Insights.md`

### [[Hearth Insights]]
- **Outbound References** (4):
  - [[HIPPA]]
  - [[FHIR]]
  ...
- **Inbound References** (12):
  - [[2025_04_06]]
  - [[business plan]]
  ...
```

## Integration with Claude Code

These indexes help Claude Code understand your Logseq knowledge base by:

1. **Task Context**: Claude can see all active tasks and their relationships
2. **Knowledge Graph**: Claude understands how your notes connect
3. **Time Tracking**: Claude can see effort invested in different areas
4. **Quick Navigation**: Claude can reference specific files and line numbers

## Git Hook Integration

### Manual Setup

Add to `.git/hooks/post-commit`:

```bash
#!/bin/bash
# Auto-generate indexes after each commit

REPO_PATH="$(git rev-parse --show-toplevel)"

if command -v logseq-claude-indexer &> /dev/null; then
    logseq-claude-indexer generate --repo "$REPO_PATH" --quiet
fi
```

Make it executable:
```bash
chmod +x .git/hooks/post-commit
```

## Development

### Prerequisites

- Go 1.21 or later
- golangci-lint (optional, for linting)

### Build

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Run with coverage
make test-coverage

# Run linter
make lint

# Try it on test fixtures
make demo
```

### Project Structure

```
logseq-claude-indexer/
├── cmd/
│   └── logseq-claude-indexer/    # CLI entry point
├── internal/
│   ├── scanner/                   # File system walker
│   ├── parser/                    # Markdown parsers
│   ├── indexer/                   # Index builders
│   └── writer/                    # Output generators
├── pkg/
│   └── models/                    # Shared data structures
├── testdata/
│   └── fixtures/                  # Sample Logseq files
└── Makefile
```

### Running Tests

```bash
# All tests
go test ./...

# With coverage
go test -cover ./...

# Specific package
go test ./internal/parser/

# Verbose
go test -v ./...
```

## Performance

- **100 files** (<1MB): <100ms
- **1000 files** (~10MB): <500ms
- **10,000 files** (~100MB): <3s

## Roadmap

See [futures.md](futures.md) for planned features:

- **v2.0**: Incremental indexing (only process changed files)
- **v2.1**: Tag extraction, property parsing, block references
- **v3.0**: JSON/CSV/GraphML export formats
- **v4.0**: Server mode with HTTP API
- **Future**: Web-based visualization UI

## Contributing

Contributions welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass: `make check`
5. Submit a pull request

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Acknowledgments

Built for use with [Claude Code](https://claude.com/claude-code) to enhance AI understanding of Logseq knowledge bases.

## Support

- **Issues**: [GitHub Issues](https://github.com/dyluth/logseq-claude-indexer/issues)
- **Discussions**: [GitHub Discussions](https://github.com/dyluth/logseq-claude-indexer/discussions)

## Examples

See `testdata/fixtures/` for example Logseq files and generated indexes.

---

**Made with ❤️ for the Logseq and Claude Code communities**
