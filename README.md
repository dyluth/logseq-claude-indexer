# logseq-claude-indexer

A production-ready Golang CLI tool that scans Logseq knowledge bases and generates Claude Code-optimized index files.
This allows claude code to be run directly inside your logseq folder (or just in the folder with the indexes) and help with your workflows.


## Features

- **Fast**: Scans 100+ files in <1 second
- **Task Extraction**: Finds all NOW/LATER/TODO/DOING/DONE tasks with context
- **Priority Support**: Extracts and organizes by [#A], [#B], [#C] priority markers
- **Time Tracking**: Parses Logseq's `:LOGBOOK:` CLOCK entries and generates analytics
- **Timeline View**: Recent activity (7 days) + complete history in condensed format
- **Missing Pages**: Identifies frequently referenced pages that don't exist yet (5+ refs)
- **Reference Graph**: Builds a network of `[[page links]]`
- **Dashboard**: Aggregated overview with quick stats, priorities, and recent activity
- **AI-Optimized**: Generates token-efficient markdown indexes perfect for Claude Code
- **Git Integration**: Automatic post-commit hook setup with `make setup-git-hook`

## Installation

> **New to installation?** See [INSTALL.md](INSTALL.md) for a step-by-step guide with troubleshooting.

### From Source (Recommended)

```bash
git clone https://github.com/dyluth/logseq-claude-indexer.git
cd logseq-claude-indexer

# Build and test
make build
make test

# Install to ~/.local/bin (user-local, recommended)
make install-user

# OR install to GOPATH/bin (requires Go setup)
make install
```

**Note**: If using `install-user`, ensure `~/.local/bin` is in your PATH:
```bash
export PATH="$HOME/.local/bin:$PATH"
```

### Via Go Install

```bash
go install github.com/dyluth/logseq-claude-indexer/cmd/logseq-claude-indexer@latest
```

### Pre-built Binaries

Download from the [releases page](https://github.com/dyluth/logseq-claude-indexer/releases).

## Quick Start

```bash
# 1. Generate indexes for your Logseq repo (one-time or manual)
logseq-claude-indexer generate --repo /path/to/logseq

# 2. View the generated dashboard (start here!)
cat /path/to/logseq/.claude/indexes/dashboard.md

# 3. Set up git hook for automatic updates (see Git Hook Integration below)
cd /path/to/logseq
cat > .git/hooks/post-commit << 'EOF'
#!/bin/bash
logseq-claude-indexer generate --repo . --output .claude/indexes --quiet
EOF
chmod +x .git/hooks/post-commit
```

**8 index files are generated**:
1. `dashboard.md` - **Overview** (start here for Claude)
2. `tasks-by-status.md` - All tasks by workflow stage
3. `tasks-by-priority.md` - High-priority tasks ([#A], [#B], [#C])
4. `timeline-recent.md` - Last 7 days detailed activity
5. `timeline-full.md` - Complete history (condensed)
6. `missing-pages.md` - Suggested pages to create (5+ refs)
7. `time-tracking.md` - Time allocation analytics
8. `reference-graph.md` - Page connections

See `.claude/indexes/README.md` for detailed documentation of each file.

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

All indexes are optimized for Claude with token-efficient formatting. See `.claude/indexes/README.md` for detailed documentation.

### Dashboard (`dashboard.md`) 🏠

**Your knowledge base at a glance** - Share this first with Claude!

Contains:
- Quick stats (total tasks, completion rate, time tracking adoption)
- Current high-priority tasks ([#A] items)
- Recent activity (last 3 days)
- Top projects by time invested
- Suggested pages to create
- Links to all detailed reports

### Tasks by Status (`tasks-by-status.md`)

All tasks organized by workflow stage (NOW, TODO, DOING, DONE, LATER).

Contains:
- Task counts and completion statistics
- Tasks grouped by status with file locations
- Time tracking data per task
- Page references and project summaries

Example:
```markdown
## NOW (12 tasks)

- [[Project Alpha]] - Implement user authentication `journals/2025_11_06.md:15`
  - ⏱ 5h 30m logged
  - Refs: [[Project Alpha]], [[Auth]]
```

### Tasks by Priority (`tasks-by-priority.md`)

High-priority tasks marked with [#A], [#B], [#C].

Contains:
- Tasks grouped by priority level (A = highest)
- Only tasks with explicit priority markers
- Completion status per priority

### Timeline Recent (`timeline-recent.md`)

Last 7 days detailed activity.

Contains:
- Daily task summaries
- Time logged per day
- Key highlights (🔥 markers for important items)
- Full task details with file locations

### Timeline Full (`timeline-full.md`)

Complete activity history in condensed format.

Contains:
- All days with activity (token-optimized)
- Task counts by status per day
- Time logged summaries
- Highlighted activities only

### Missing Pages (`missing-pages.md`)

Suggested pages to create based on reference frequency.

Contains:
- Pages referenced 5+ times that don't exist yet
- Categorized by type: person, project, concept, date
- Reference count and source pages (top 10)
- Helps identify knowledge gaps

### Time Tracking (`time-tracking.md`)

Time allocation analytics from LOGBOOK entries.

Contains:
- Total time logged across all tasks
- Time tracking adoption rate
- Top 10 projects by time invested
- Weekly breakdown (last 8 weeks)
- Time by priority and status

### Reference Graph (`reference-graph.md`)

Network view of page connections.

Contains:
- Hub pages (most referenced)
- Inbound and outbound references per page
- Orphan pages (no connections)
- Bi-directional link indicators

## Integration with Claude Code

These indexes help Claude Code understand your Logseq knowledge base by:

1. **Task Context**: Claude can see all active tasks and their relationships
2. **Knowledge Graph**: Claude understands how your notes connect
3. **Time Tracking**: Claude can see effort invested in different areas
4. **Quick Navigation**: Claude can reference specific files and line numbers

## Git Hook Integration

### Setup (Recommended)

After installing the indexer, set up automatic index generation on every commit:

```bash
# Navigate to your Logseq repository
cd /path/to/your/logseq

# Create the post-commit hook
cat > .git/hooks/post-commit << 'EOF'
#!/bin/bash
# Auto-generate Claude indexes after each commit
logseq-claude-indexer generate --repo . --output .claude/indexes --quiet
EOF

# Make it executable
chmod +x .git/hooks/post-commit
```

**That's it!** Indexes will now auto-generate after every `git commit`.

### Alternative: Using Make (from source directory)

If you cloned the source and want to use the Makefile:

```bash
# From the logseq-claude-indexer source directory (where you ran make install-user)
cd /path/to/logseq-claude-indexer
make setup-git-hook LOGSEQ_REPO=/path/to/your/logseq

# Example:
make setup-git-hook LOGSEQ_REPO=~/Documents/logseq
```

### Verify It Works

```bash
# Make a commit in your Logseq repo
cd /path/to/your/logseq
echo "test" >> pages/test.md
git add pages/test.md
git commit -m "test hook"

# Check that indexes were updated (should show recent timestamp)
ls -lt .claude/indexes/ | head
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

# Test on all datasets (fixtures, synthetic, user-provided)
make test-datasets

# Run with coverage
make test-coverage

# Run linter
make lint

# Try it on test fixtures
make demo

# Format code
make fmt

# Run all checks (fmt, lint, test)
make check

# View all available targets
make help
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

### Recently Completed ✅

- **v1.5**: Priority extraction ([#A], [#B], [#C])
- **v1.5**: Timeline indexes (recent + full history)
- **v1.5**: Missing pages report (5+ references)
- **v1.5**: Time tracking analytics
- **v1.5**: Dashboard aggregation
- **v1.5**: Automated git hook setup

### Planned Features

See [futures.md](futures.md) for detailed plans:

- **v2.0**: Incremental indexing (only process changed files)
- **v2.1**: Tag extraction (#tags), property parsing, block references
- **v2.2**: Task dependencies and relationships
- **v3.0**: JSON/CSV/GraphML export formats
- **v3.1**: Query API for programmatic access
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
