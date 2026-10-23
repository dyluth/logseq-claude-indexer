# Future Features (v2+)

This document tracks features deferred from the MVP to future versions.

## v2.0 - Incremental Indexing

### Motivation
Currently, the tool performs a full scan and reindex on every run. For large Logseq repositories (1000+ files), this can take several seconds. Incremental mode would only reprocess changed files.

### Implementation Approach
1. **Change detection**:
   - Store last index timestamp in `.claude/indexes/.last-index`
   - Use `git diff --name-only` to find changed files since last commit
   - Alternatively, compare file ModTime against last index time

2. **Partial reindexing**:
   - Load previous index files (task-index.md, reference-graph.md)
   - Parse only changed files
   - Remove old entries for changed files
   - Add new entries for changed files
   - Rebuild affected portions of graph

3. **CLI flag**:
   ```bash
   logseq-claude-indexer generate --incremental
   ```

### Trade-offs
- **Complexity**: Must handle edge cases (deleted files, renamed files)
- **Correctness**: Risk of stale data if change detection fails
- **Performance**: 10-100x faster for small changesets

### Decision
Defer to v2.0. MVP should prioritize correctness and simplicity. Full reindex is acceptable for repos <1000 files (target <3s execution time).

---

## v2.1 - Advanced Features

### Tag Extraction
- Parse `#tags` from markdown content
- Generate tag index showing which pages use each tag
- Group tasks by tag

### Property Extraction
- Parse Logseq properties (`:property: value` syntax)
- Index by property type
- Useful for custom metadata (status, priority, etc.)

### Block References
- Parse `((block-id))` references
- Build block-level reference graph
- More granular than page-level references

---

## v3.0 - Alternative Outputs

### JSON Export
- Machine-readable format for tooling integration
- Easier to consume in other programs

### CSV Export
- Simple format for spreadsheet analysis
- Useful for task time tracking reports

### GraphML Export
- Standard graph format
- Import into visualization tools (Gephi, Cytoscape)

---

## v4.0 - Server Mode

### Long-Running Process
- Watch filesystem for changes
- Auto-regenerate indexes without git hook
- Serve indexes over HTTP API

### Use Cases
- Real-time index updates during Logseq editing
- Integration with other tools (VSCode extension, etc.)

---

## Validated Write Path

### Motivation
The indexes give Claude a reliable, deterministic way to *read* the knowledge base,
but writing back (marking a task DONE, rescheduling, adding a block) still means
hand-editing raw Logseq markdown, where a wrong indent or marker fails silently.

### Approach
- The human keeps editing freehand; the agent writes only through a checked command,
  e.g. `logseq-claude-indexer task done <file>:<line>` or `task schedule <file>:<line> <date>`
- Each command parses the target block, applies one structured change and re-parses
  the file, refusing the write if anything other than the intended block changed
- Round-trip tests: parse → serialise must reproduce every fixture byte for byte
- Regenerate the affected indexes after a successful write

### Trade-offs
- Logseq's format is loosely specified, so the round-trip guarantee is only as good as
  the fixtures; start with task status and SCHEDULED only

---

## Far Future - Web UI

### Browser-Based Viewer
- Interactive graph visualization
- Search interface
- Task management view
- Alternative to Claude Code for non-technical users

---

## Notes
- Each version should maintain backward compatibility with v1 index formats
- Consider breaking changes only for major versions
- Prioritize stability and correctness over features
