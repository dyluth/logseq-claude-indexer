# Build Prompt: logseq-index - Golang CLI for Logseq Knowledge Base Indexing

## Project Overview

Build a production-ready Golang CLI tool called `logseq-index` that scans a Logseq knowledge base repository and generates Claude Code-optimized index files. The tool must be fast, reliable, and integrate seamlessly with git workflows.

### Core Purpose
- Parse Logseq markdown files (pages and journals)
- Extract task markers (NOW/LATER/TODO/DOING/DONE)
- Build a reference graph of [[page links]]
- Generate structured index files in `.claude/indexes/` for AI consumption
- Run automatically via git post-commit hooks

### Key Requirements
1. **Zero dependencies on Logseq** - Pure file system operations
2. **Fast execution** - Complete scan in <1 second for 100+ files
3. **Idempotent** - Running multiple times produces consistent output
4. **Git-hook friendly** - Silent mode, non-zero exit on errors
5. **Extensible** - Easy to add new index types in future

---

## Architecture

### Data Flow
```
1. Scanner → Walk pages/ and journals/ directories
2. Parser → Read each .md file, extract tasks and [[references]]
3. Indexer → Aggregate tasks by status, build reference graph
4. Writer → Generate .claude/indexes/*.md files
```

### Components
```
logseq-index/
├── cmd/logseq-index/main.go       # CLI entry point
├── internal/
│   ├── scanner/                   # Walk filesystem, find .md files
│   ├── parser/                    # Parse Logseq markdown
│   ├── indexer/                   # Build indexes from parsed data
│   └── writer/                    # Write index files
├── pkg/models/                    # Shared data structures
└── testdata/fixtures/             # Sample Logseq files for testing
```

---

## Complete File Structure

```
logseq-index/
├── cmd/
│   └── logseq-index/
│       └── main.go                 # CLI entry point with Cobra
├── internal/
│   ├── scanner/
│   │   ├── scanner.go              # Filesystem walker
│   │   └── scanner_test.go
│   ├── parser/
│   │   ├── parser.go               # Main parser interface
│   │   ├── tasks.go                # Extract task markers
│   │   ├── references.go           # Extract [[page]] references
│   │   ├── logbook.go              # Parse :LOGBOOK: time tracking
│   │   └── parser_test.go
│   ├── indexer/
│   │   ├── task_indexer.go         # Build task index
│   │   ├── graph_indexer.go        # Build reference graph
│   │   └── indexer_test.go
│   └── writer/
│       ├── writer.go               # Write index files
│       └── writer_test.go
├── pkg/
│   └── models/
│       ├── task.go                 # Task data structure
│       ├── page.go                 # Page metadata
│       ├── reference.go            # Page reference relationship
│       └── logbook.go              # Time tracking entry
├── testdata/
│   └── fixtures/
│       ├── journals/
│       │   └── 2025_04_06.md      # Sample journal
│       └── pages/
│           └── Hearth Insights.md  # Sample page
├── scripts/
│   ├── install-hook.sh             # Install git post-commit hook
│   └── uninstall-hook.sh           # Remove git hook
├── .golangci.yml                   # Linter config
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

---

## Detailed Implementation Specifications

### 1. CLI Structure (cmd/logseq-index/main.go)

Use `github.com/spf13/cobra` for CLI framework.

#### Commands:
```bash
# Primary command (used by git hook)
logseq-index generate --repo /path/to/logseq-repo

# Specific index generation
logseq-index tasks --repo /path/to/logseq-repo
logseq-index graph --repo /path/to/logseq-repo

# Utility commands
logseq-index validate --repo /path/to/logseq-repo   # Validate repo structure
logseq-index install-hook --repo /path/to/logseq-repo  # Install git hook
logseq-index version                                 # Show version
```

#### Flags:
```bash
--repo string       Path to Logseq repository (default: current directory)
--output string     Output directory (default: .claude/indexes)
--quiet            Suppress output (for git hooks)
--verbose          Show detailed logging
--dry-run          Show what would be generated without writing files
```

#### Example Implementation Structure:
```go
package main

import (
    "fmt"
    "os"
    "github.com/spf13/cobra"
    "github.com/yourusername/logseq-index/internal/scanner"
    "github.com/yourusername/logseq-index/internal/parser"
    "github.com/yourusername/logseq-index/internal/indexer"
    "github.com/yourusername/logseq-index/internal/writer"
)

var (
    repoPath   string
    outputDir  string
    quiet      bool
    verbose    bool
)

func main() {
    rootCmd := &cobra.Command{
        Use:   "logseq-index",
        Short: "Generate Claude Code-optimized indexes from Logseq repositories",
        Long:  `Scan Logseq markdown files and generate task indexes and reference graphs.`,
    }

    generateCmd := &cobra.Command{
        Use:   "generate",
        Short: "Generate all indexes",
        RunE:  runGenerate,
    }

    // Add flags
    generateCmd.Flags().StringVar(&repoPath, "repo", ".", "Path to Logseq repository")
    generateCmd.Flags().StringVar(&outputDir, "output", ".claude/indexes", "Output directory")
    generateCmd.Flags().BoolVar(&quiet, "quiet", false, "Suppress output")
    generateCmd.Flags().BoolVar(&verbose, "verbose", false, "Verbose logging")

    rootCmd.AddCommand(generateCmd)
    // Add other commands: tasks, graph, validate, install-hook, version

    if err := rootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}

func runGenerate(cmd *cobra.Command, args []string) error {
    // 1. Validate repo structure
    // 2. Scan files
    // 3. Parse markdown
    // 4. Build indexes
    // 5. Write output files
    return nil
}
```

---

### 2. Scanner (internal/scanner/scanner.go)

Walks the filesystem and identifies Logseq files.

#### Responsibilities:
- Find all `.md` files in `pages/` and `journals/`
- Ignore hidden files, `.recycle/`, `logseq/bak/`
- Return file metadata (path, mod time)

#### Data Structures:
```go
package scanner

type File struct {
    Path         string    // Relative path from repo root
    AbsolutePath string    // Absolute filesystem path
    Type         FileType  // Journal or Page
    ModTime      time.Time // Last modified timestamp
}

type FileType int

const (
    FileTypeJournal FileType = iota
    FileTypePage
)

type Scanner struct {
    repoPath string
}

func New(repoPath string) *Scanner {
    return &Scanner{repoPath: repoPath}
}

// Scan walks the repository and returns all markdown files
func (s *Scanner) Scan() ([]File, error) {
    var files []File

    // Scan journals/
    journalFiles, err := s.scanDirectory("journals", FileTypeJournal)
    if err != nil {
        return nil, err
    }
    files = append(files, journalFiles...)

    // Scan pages/
    pageFiles, err := s.scanDirectory("pages", FileTypePage)
    if err != nil {
        return nil, err
    }
    files = append(files, pageFiles...)

    return files, nil
}

func (s *Scanner) scanDirectory(dir string, fileType FileType) ([]File, error) {
    // Use filepath.Walk to recursively find .md files
    // Filter out .recycle, hidden files
    // Return File structs with metadata
    return nil, nil
}
```

#### Edge Cases:
- Handle symlinks gracefully
- Skip files without `.md` extension
- Handle permission errors without failing entire scan

---

### 3. Parser (internal/parser/)

Parse Logseq markdown and extract structured data.

#### 3.1 Task Parser (parser/tasks.go)

Extract task markers with context.

##### Logseq Task Syntax:
```markdown
- NOW [[Project Name]] - Task description
  :LOGBOOK:
  CLOCK: [2025-04-06 Sun 12:30:42]--[2025-04-17 Thu 09:17:26] =>  260:46:44
  :END:
- LATER Do something else
- DONE [[Page Reference]] - Completed task
```

##### Data Structure:
```go
package models

type Task struct {
    Status      TaskStatus      // NOW, LATER, TODO, DOING, DONE
    Description string          // Full task text
    PageRefs    []string        // [[Page Name]] references in task
    SourceFile  string          // File path where task was found
    LineNumber  int             // Line number in source file
    Logbook     []LogbookEntry  // Time tracking entries (if present)
}

type TaskStatus string

const (
    StatusNOW   TaskStatus = "NOW"
    StatusLATER TaskStatus = "LATER"
    StatusTODO  TaskStatus = "TODO"
    StatusDOING TaskStatus = "DOING"
    StatusDONE  TaskStatus = "DONE"
)

type LogbookEntry struct {
    Start    time.Time
    End      time.Time
    Duration time.Duration
}
```

##### Parsing Algorithm:
```go
func ParseTasks(content string, filePath string) ([]Task, error) {
    var tasks []Task
    lines := strings.Split(content, "\n")

    for i, line := range lines {
        // 1. Check if line starts with "- " or "	- "
        // 2. Check if it contains task markers (NOW|LATER|TODO|DOING|DONE)
        // 3. Extract task description
        // 4. Extract [[page references]] using regex: \[\[([^\]]+)\]\]
        // 5. Check following lines for :LOGBOOK: block
        // 6. Parse CLOCK entries if present

        task := Task{
            SourceFile: filePath,
            LineNumber: i + 1,
        }

        // Extract status
        if strings.Contains(line, "NOW ") {
            task.Status = StatusNOW
        }
        // ... handle other statuses

        // Extract page references
        task.PageRefs = extractPageReferences(line)

        // Parse logbook if present
        if i+1 < len(lines) && strings.Contains(lines[i+1], ":LOGBOOK:") {
            task.Logbook = parseLogbook(lines[i+1:])
        }

        tasks = append(tasks, task)
    }

    return tasks, nil
}

func extractPageReferences(line string) []string {
    // Use regex: \[\[([^\]]+)\]\]
    // Return slice of page names without brackets
    re := regexp.MustCompile(`\[\[([^\]]+)\]\]`)
    matches := re.FindAllStringSubmatch(line, -1)

    var refs []string
    for _, match := range matches {
        if len(match) > 1 {
            refs = append(refs, match[1])
        }
    }
    return refs
}
```

#### 3.2 Reference Parser (parser/references.go)

Extract all [[page reference]] links from content.

##### Data Structure:
```go
package models

type PageReference struct {
    SourceFile string   // File containing the reference
    SourcePage string   // Page name (derived from filename)
    TargetPage string   // Referenced page name
    LineNumber int      // Line where reference appears
    Context    string   // Surrounding text for context
}
```

##### Parsing Algorithm:
```go
func ParseReferences(content string, filePath string) ([]PageReference, error) {
    var refs []PageReference
    lines := strings.Split(content, "\n")

    sourcePage := extractPageNameFromPath(filePath)

    for i, line := range lines {
        // Find all [[page]] references in line
        pageRefs := extractPageReferences(line)

        for _, targetPage := range pageRefs {
            refs = append(refs, PageReference{
                SourceFile: filePath,
                SourcePage: sourcePage,
                TargetPage: targetPage,
                LineNumber: i + 1,
                Context:    extractContext(line, 50), // 50 chars context
            })
        }
    }

    return refs, nil
}

func extractPageNameFromPath(filePath string) string {
    // journals/2025_04_06.md -> "2025_04_06"
    // pages/Hearth Insights.md -> "Hearth Insights"
    base := filepath.Base(filePath)
    return strings.TrimSuffix(base, ".md")
}
```

---

### 4. Indexer (internal/indexer/)

Aggregate parsed data into index structures.

#### 4.1 Task Indexer (indexer/task_indexer.go)

Group tasks by status, chronology, and page references.

##### Data Structure:
```go
package indexer

type TaskIndex struct {
    GeneratedAt time.Time
    TotalTasks  int
    ByStatus    map[models.TaskStatus][]models.Task
    ByProject   map[string][]models.Task  // Group by page reference
    Recent      []models.Task              // Last 30 days
}

func BuildTaskIndex(tasks []models.Task) *TaskIndex {
    index := &TaskIndex{
        GeneratedAt: time.Now(),
        ByStatus:    make(map[models.TaskStatus][]models.Task),
        ByProject:   make(map[string][]models.Task),
    }

    for _, task := range tasks {
        // Group by status
        index.ByStatus[task.Status] = append(index.ByStatus[task.Status], task)

        // Group by project (first page reference)
        if len(task.PageRefs) > 0 {
            project := task.PageRefs[0]
            index.ByProject[project] = append(index.ByProject[project], task)
        }
    }

    index.TotalTasks = len(tasks)

    return index
}
```

#### 4.2 Graph Indexer (indexer/graph_indexer.go)

Build a navigable reference graph.

##### Data Structure:
```go
package indexer

type ReferenceGraph struct {
    GeneratedAt time.Time
    Nodes       map[string]*GraphNode  // Page name -> Node
    HubPages    []string               // Most referenced pages
}

type GraphNode struct {
    PageName       string
    FilePath       string
    OutboundRefs   []string  // Pages this page references
    InboundRefs    []string  // Pages that reference this page
    ReferenceCount int       // Total inbound references
}

func BuildReferenceGraph(refs []models.PageReference, files []scanner.File) *ReferenceGraph {
    graph := &ReferenceGraph{
        GeneratedAt: time.Now(),
        Nodes:       make(map[string]*GraphNode),
    }

    // Build nodes
    for _, file := range files {
        pageName := extractPageNameFromPath(file.Path)
        graph.Nodes[pageName] = &GraphNode{
            PageName: pageName,
            FilePath: file.Path,
        }
    }

    // Add references
    for _, ref := range refs {
        // Add outbound reference
        if node, exists := graph.Nodes[ref.SourcePage]; exists {
            node.OutboundRefs = append(node.OutboundRefs, ref.TargetPage)
        }

        // Add inbound reference
        if node, exists := graph.Nodes[ref.TargetPage]; exists {
            node.InboundRefs = append(node.InboundRefs, ref.SourcePage)
            node.ReferenceCount++
        }
    }

    // Identify hub pages (most referenced)
    graph.HubPages = findHubPages(graph.Nodes, 5) // Top 5

    return graph
}
```

---

### 5. Writer (internal/writer/)

Write index files in Claude Code-optimized markdown format.

#### 5.1 Task Index Output (writer/task_writer.go)

##### Output File: `.claude/indexes/task-index.md`

```markdown
# Logseq Task Index

Generated: 2025-11-06 14:30:00 UTC
Total Tasks: 47

---

## NOW Tasks (3)

### [[Hearth Insights]] - setup actions
- **File**: `journals/2025_04_06.md:1`
- **Time Logged**: 260h 46m 44s
- **Started**: 2025-04-06 12:30:42

### [[HI-prototype]] - Build MVP
- **File**: `pages/HI-prototype.md:15`
- **Status**: In Progress

---

## LATER Tasks (12)

### Research FHIR implementation patterns
- **File**: `pages/FHIR.md:8`
- **Context**: Need to evaluate different approaches

---

## DONE Tasks (32)

### [[Hearth Insights - Action Streams]] - Initial setup
- **File**: `journals/2025_04_06.md:5`
- **Completed**: 2025-04-07
- **Time Logged**: 1 second

---

## Tasks by Project

### Hearth Insights (18 tasks)
- 2 NOW
- 8 LATER
- 8 DONE

### FHIR Implementation (6 tasks)
- 0 NOW
- 4 LATER
- 2 DONE
```

##### Writer Implementation:
```go
package writer

import (
    "fmt"
    "os"
    "path/filepath"
    "time"
)

func WriteTaskIndex(index *indexer.TaskIndex, outputDir string) error {
    filePath := filepath.Join(outputDir, "task-index.md")

    f, err := os.Create(filePath)
    if err != nil {
        return err
    }
    defer f.Close()

    // Write header
    fmt.Fprintf(f, "# Logseq Task Index\n\n")
    fmt.Fprintf(f, "Generated: %s\n", index.GeneratedAt.Format(time.RFC3339))
    fmt.Fprintf(f, "Total Tasks: %d\n\n", index.TotalTasks)
    fmt.Fprintf(f, "---\n\n")

    // Write NOW tasks
    if tasks, exists := index.ByStatus[models.StatusNOW]; exists && len(tasks) > 0 {
        fmt.Fprintf(f, "## NOW Tasks (%d)\n\n", len(tasks))
        for _, task := range tasks {
            writeTask(f, task)
        }
    }

    // Write LATER, TODO, DOING, DONE tasks...

    // Write by project section
    fmt.Fprintf(f, "---\n\n## Tasks by Project\n\n")
    for project, tasks := range index.ByProject {
        writeProjectSummary(f, project, tasks)
    }

    return nil
}

func writeTask(f *os.File, task models.Task) {
    fmt.Fprintf(f, "### %s\n", task.Description)
    fmt.Fprintf(f, "- **File**: `%s:%d`\n", task.SourceFile, task.LineNumber)

    if len(task.Logbook) > 0 {
        totalDuration := calculateTotalDuration(task.Logbook)
        fmt.Fprintf(f, "- **Time Logged**: %s\n", formatDuration(totalDuration))
    }

    fmt.Fprintf(f, "\n")
}
```

#### 5.2 Reference Graph Output (writer/graph_writer.go)

##### Output File: `.claude/indexes/reference-graph.md`

```markdown
# Logseq Reference Graph

Generated: 2025-11-06 14:30:00 UTC
Total Pages: 22
Total References: 58

---

## Hub Pages (Most Referenced)

1. **[[Hearth Insights]]** - 12 inbound references
   - File: `pages/Hearth Insights.md`

2. **[[FHIR]]** - 8 inbound references
   - File: `pages/FHIR.md`

3. **[[HI-pitch deck]]** - 4 inbound references
   - File: `pages/HI-pitch deck.md`

---

## Page Clusters

### Hearth Insights Cluster

#### [[Hearth Insights]]
- **File**: `pages/Hearth Insights.md`
- **Outbound References** (4):
  - [[Hearth Insights - principles]]
  - [[HIPPA]]
  - [[FHIR]]
  - [[Healthcare research empowerment - company idea]]
- **Inbound References** (12):
  - `journals/2025_04_06.md:1`
  - `journals/2025_04_07.md:3`
  - `pages/business plan.md:8`
  - `pages/HI-pitch deck.md:15`
  - ... (8 more)

#### [[Hearth Insights - Action Streams]]
- **File**: `pages/Hearth Insights - Action Streams.md`
- **Outbound References** (2):
  - [[Hearth Insights]]
  - [[HI-prototype]]
- **Inbound References** (3):
  - `journals/2025_04_06.md:5`
  - `pages/Hearth Insights.md:20`

---

## Orphan Pages (No References)

- `pages/contents.md` (empty page)

---

## Journal References

Recent journal entries mentioning key topics:

### [[Hearth Insights]]
- 2025-04-06: Initial setup and action items
- 2025-04-07: Action streams defined
- 2025-04-17: Architecture decisions

### [[FHIR]]
- 2025-04-06: Research public test servers
- 2025-04-15: Implementation planning
```

##### Writer Implementation:
```go
func WriteReferenceGraph(graph *indexer.ReferenceGraph, outputDir string) error {
    filePath := filepath.Join(outputDir, "reference-graph.md")

    f, err := os.Create(filePath)
    if err != nil {
        return err
    }
    defer f.Close()

    // Write header
    fmt.Fprintf(f, "# Logseq Reference Graph\n\n")
    fmt.Fprintf(f, "Generated: %s\n", graph.GeneratedAt.Format(time.RFC3339))
    fmt.Fprintf(f, "Total Pages: %d\n\n", len(graph.Nodes))
    fmt.Fprintf(f, "---\n\n")

    // Write hub pages
    fmt.Fprintf(f, "## Hub Pages (Most Referenced)\n\n")
    for i, pageName := range graph.HubPages {
        node := graph.Nodes[pageName]
        fmt.Fprintf(f, "%d. **[[%s]]** - %d inbound references\n",
            i+1, node.PageName, node.ReferenceCount)
        fmt.Fprintf(f, "   - File: `%s`\n\n", node.FilePath)
    }

    // Write detailed page clusters
    // Group by identifying clusters (pages that reference each other heavily)

    return nil
}
```

---

## 6. Git Hook Integration

### Install Hook Command

```go
// cmd/logseq-index/install_hook.go

func runInstallHook(cmd *cobra.Command, args []string) error {
    repoPath, err := getRepoPath()
    if err != nil {
        return err
    }

    // Check if .git exists
    gitDir := filepath.Join(repoPath, ".git")
    if _, err := os.Stat(gitDir); os.IsNotExist(err) {
        return fmt.Errorf("not a git repository: %s", repoPath)
    }

    // Create hook file
    hookPath := filepath.Join(gitDir, "hooks", "post-commit")

    hookContent := `#!/bin/bash
# Auto-generated by logseq-index
# Regenerate indexes after each commit

REPO_PATH="$(git rev-parse --show-toplevel)"

if command -v logseq-index &> /dev/null; then
    logseq-index generate --repo "$REPO_PATH" --quiet
else
    echo "Warning: logseq-index not installed"
fi
`

    if err := os.WriteFile(hookPath, []byte(hookContent), 0755); err != nil {
        return err
    }

    fmt.Println("✓ Git post-commit hook installed successfully")
    fmt.Printf("  Location: %s\n", hookPath)

    return nil
}
```

### Verification

After installation, test the hook:
```bash
cd /path/to/logseq-data-git
echo "test" >> test.txt
git add test.txt
git commit -m "Test hook"
# Should see: "Updating Logseq indexes..."
ls .claude/indexes/
# Should show: task-index.md, reference-graph.md
```

---

## 7. Testing Strategy

### Test Data Fixtures (testdata/fixtures/)

Create realistic sample files:

**testdata/fixtures/journals/2025_04_06.md:**
```markdown
- NOW [[Project A]] - Implement feature X
  :LOGBOOK:
  CLOCK: [2025-04-06 Sun 10:00:00]--[2025-04-06 Sun 12:00:00] =>  02:00:00
  :END:
- LATER [[Project B]] - Research options
- DONE [[Project A]] - Setup repository
```

**testdata/fixtures/pages/Project A.md:**
```markdown
- See also: [[Project B]]
- Related: [[FHIR]]
- TODO [[Project A]] - Next milestone
```

### Unit Tests

```go
// internal/parser/tasks_test.go

func TestParseTasksWithLogbook(t *testing.T) {
    content := `- NOW [[Test Project]] - Sample task
  :LOGBOOK:
  CLOCK: [2025-04-06 Sun 10:00:00]--[2025-04-06 Sun 12:00:00] =>  02:00:00
  :END:`

    tasks, err := ParseTasks(content, "test.md")
    if err != nil {
        t.Fatalf("ParseTasks failed: %v", err)
    }

    if len(tasks) != 1 {
        t.Fatalf("Expected 1 task, got %d", len(tasks))
    }

    task := tasks[0]
    if task.Status != models.StatusNOW {
        t.Errorf("Expected status NOW, got %s", task.Status)
    }

    if len(task.PageRefs) != 1 || task.PageRefs[0] != "Test Project" {
        t.Errorf("Expected page ref 'Test Project', got %v", task.PageRefs)
    }

    if len(task.Logbook) != 1 {
        t.Fatalf("Expected 1 logbook entry, got %d", len(task.Logbook))
    }

    if task.Logbook[0].Duration != 2*time.Hour {
        t.Errorf("Expected 2h duration, got %v", task.Logbook[0].Duration)
    }
}
```

### Integration Tests

```go
// cmd/logseq-index/integration_test.go

func TestFullPipeline(t *testing.T) {
    // 1. Create temp directory with test fixtures
    tmpDir := setupTestRepo(t)
    defer os.RemoveAll(tmpDir)

    // 2. Run scanner
    scanner := scanner.New(tmpDir)
    files, err := scanner.Scan()
    if err != nil {
        t.Fatalf("Scan failed: %v", err)
    }

    // 3. Parse all files
    var allTasks []models.Task
    var allRefs []models.PageReference

    for _, file := range files {
        content, _ := os.ReadFile(file.AbsolutePath)
        tasks, _ := parser.ParseTasks(string(content), file.Path)
        refs, _ := parser.ParseReferences(string(content), file.Path)

        allTasks = append(allTasks, tasks...)
        allRefs = append(allRefs, refs...)
    }

    // 4. Build indexes
    taskIndex := indexer.BuildTaskIndex(allTasks)
    graphIndex := indexer.BuildReferenceGraph(allRefs, files)

    // 5. Write outputs
    outputDir := filepath.Join(tmpDir, ".claude", "indexes")
    os.MkdirAll(outputDir, 0755)

    writer.WriteTaskIndex(taskIndex, outputDir)
    writer.WriteReferenceGraph(graphIndex, outputDir)

    // 6. Verify outputs exist and contain expected data
    taskIndexPath := filepath.Join(outputDir, "task-index.md")
    if _, err := os.Stat(taskIndexPath); os.IsNotExist(err) {
        t.Error("task-index.md not created")
    }

    graphIndexPath := filepath.Join(outputDir, "reference-graph.md")
    if _, err := os.Stat(graphIndexPath); os.IsNotExist(err) {
        t.Error("reference-graph.md not created")
    }
}
```

---

## 8. Build and Installation

### Makefile

```makefile
.PHONY: build test install clean lint

# Build for current platform
build:
	go build -o bin/logseq-index ./cmd/logseq-index

# Build for multiple platforms
build-all:
	GOOS=darwin GOARCH=amd64 go build -o bin/logseq-index-darwin-amd64 ./cmd/logseq-index
	GOOS=darwin GOARCH=arm64 go build -o bin/logseq-index-darwin-arm64 ./cmd/logseq-index
	GOOS=linux GOARCH=amd64 go build -o bin/logseq-index-linux-amd64 ./cmd/logseq-index
	GOOS=windows GOARCH=amd64 go build -o bin/logseq-index-windows-amd64.exe ./cmd/logseq-index

# Run tests
test:
	go test -v -race -cover ./...

# Install to GOPATH/bin
install:
	go install ./cmd/logseq-index

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Run linter
lint:
	golangci-lint run

# Run on example repo
demo:
	go run ./cmd/logseq-index generate --repo ./testdata/fixtures --verbose
```

### go.mod

```go
module github.com/yourusername/logseq-index

go 1.21

require (
    github.com/spf13/cobra v1.8.0
    github.com/stretchr/testify v1.8.4  // For testing
)
```

### Installation Instructions

```bash
# Install from source
git clone https://github.com/yourusername/logseq-index.git
cd logseq-index
make install

# Or via go install
go install github.com/yourusername/logseq-index/cmd/logseq-index@latest

# Install git hook
cd /path/to/logseq-data-git
logseq-index install-hook
```

---

## 9. Performance Considerations

### Expected Performance
- **100 files** (<1MB total): <100ms
- **1000 files** (~10MB total): <500ms
- **10,000 files** (~100MB total): <3s

### Optimization Strategies
1. **Parallel scanning**: Use goroutines for file parsing
2. **Incremental indexing**: Only reparse changed files (check git diff)
3. **Caching**: Store last-indexed timestamp, skip unchanged files
4. **Memory efficiency**: Stream large files, don't load all into memory

### Future Optimization: Incremental Mode

```go
// Only index files changed since last run
func (s *Scanner) ScanIncremental(since time.Time) ([]File, error) {
    allFiles, err := s.Scan()
    if err != nil {
        return nil, err
    }

    var changedFiles []File
    for _, file := range allFiles {
        if file.ModTime.After(since) {
            changedFiles = append(changedFiles, file)
        }
    }

    return changedFiles, nil
}
```

---

## 10. Error Handling

### Graceful Degradation
- If a single file fails to parse, log error but continue
- If output directory doesn't exist, create it
- If git hook fails, exit with non-zero code (git will show warning)

### Error Types
```go
package errors

type ParserError struct {
    File    string
    Line    int
    Message string
}

func (e *ParserError) Error() string {
    return fmt.Sprintf("%s:%d: %s", e.File, e.Line, e.Message)
}

type ScannerError struct {
    Path    string
    Message string
}

func (e *ScannerError) Error() string {
    return fmt.Sprintf("scanner error for %s: %s", e.Path, e.Message)
}
```

---

## 11. Future Extensions

### Potential Features (Not in MVP)
1. **Search index**: Full-text search capability
2. **Date-based views**: Tasks/refs by date ranges
3. **Tag extraction**: Parse and index `#tags`
4. **Property extraction**: Parse `:property: value` pairs
5. **Block references**: Parse `((block-id))` references
6. **Export formats**: JSON, CSV, GraphML for visualization
7. **Server mode**: Long-running process with file watching
8. **Web UI**: Browser-based index viewer

---

## Implementation Checklist

### Phase 1: Core Functionality
- [ ] Set up Go project structure
- [ ] Implement scanner (walk filesystem)
- [ ] Implement task parser (extract NOW/LATER/TODO/DOING/DONE)
- [ ] Implement reference parser (extract [[page]] links)
- [ ] Implement task indexer (group by status)
- [ ] Implement graph indexer (build reference graph)
- [ ] Implement task index writer (markdown output)
- [ ] Implement graph index writer (markdown output)
- [ ] Create CLI with Cobra (generate command)
- [ ] Write unit tests for parsers
- [ ] Write integration tests

### Phase 2: Git Integration
- [ ] Implement install-hook command
- [ ] Create post-commit hook template
- [ ] Add --quiet flag for git hooks
- [ ] Test hook installation and execution
- [ ] Add uninstall-hook command

### Phase 3: Polish
- [ ] Add validate command (check repo structure)
- [ ] Add --dry-run flag
- [ ] Improve error messages
- [ ] Add progress indicators (when not --quiet)
- [ ] Write comprehensive README
- [ ] Add GitHub Actions CI
- [ ] Set up golangci-lint
- [ ] Create release binaries

### Phase 4: Documentation
- [ ] Document CLI usage
- [ ] Document output formats
- [ ] Add examples
- [ ] Create troubleshooting guide
- [ ] Document integration with Claude Code

---

## Success Criteria

The tool is complete when:

1. ✅ Can scan a Logseq repo and find all `.md` files
2. ✅ Extracts all task markers (NOW/LATER/TODO/DOING/DONE) with context
3. ✅ Extracts all [[page references]] and builds graph
4. ✅ Generates `.claude/indexes/task-index.md` with tasks grouped by status
5. ✅ Generates `.claude/indexes/reference-graph.md` with page relationships
6. ✅ Runs in <1 second for 100 files
7. ✅ Installs git post-commit hook that auto-regenerates indexes
8. ✅ Has 80%+ test coverage
9. ✅ Works on macOS, Linux, and Windows
10. ✅ Can be installed via `go install`

---

## Example Usage Workflows

### Initial Setup
```bash
# Install tool
go install github.com/yourusername/logseq-index/cmd/logseq-index@latest

# Navigate to Logseq repo
cd ~/logseq-data-git

# Install git hook
logseq-index install-hook

# Generate initial indexes
logseq-index generate

# View indexes
cat .claude/indexes/task-index.md
cat .claude/indexes/reference-graph.md
```

### Daily Usage (Automatic)
```bash
# Just use Logseq normally
# Indexes auto-update on every commit via git hook
# No manual intervention needed
```

### Manual Updates
```bash
# Force regeneration
logseq-index generate --verbose

# Dry run (see what would be generated)
logseq-index generate --dry-run

# Generate only task index
logseq-index tasks

# Generate only reference graph
logseq-index graph
```

### Validation
```bash
# Check repo structure
logseq-index validate

# Should output:
# ✓ Found pages/ directory
# ✓ Found journals/ directory
# ✓ Found 22 page files
# ✓ Found 96 journal files
# ✓ No issues detected
```

---

## Questions to Resolve Before Implementation

1. **Logbook parsing complexity**: Should we parse all CLOCK entries or just show totals?
2. **Page name normalization**: How to handle spaces, special characters in page names?
3. **Orphan detection**: Should we flag pages with zero references?
4. **Index file naming**: Use `task-index.md` or `tasks.md`? Prefix with timestamp?
5. **Error recovery**: If one file fails, continue or abort?
6. **Incremental mode**: Implement now or defer to v2?

---

## Final Notes

This prompt is designed to be copy-pasted to any AI assistant (Claude Code, ChatGPT, etc.) to build the complete `logseq-index` tool. The implementation should:

- Be production-ready (proper error handling, tests, documentation)
- Follow Go best practices (idioms, project structure, naming)
- Be maintainable (clear code, good comments, modular design)
- Be fast (<1s for 100 files)
- Be reliable (works consistently across platforms)

**Start with Phase 1 (core functionality), get it working end-to-end, then add git integration and polish.**
