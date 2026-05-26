# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Overview

This is a **Logseq knowledge base** repository, not a traditional software project. Logseq is a privacy-first, open-source knowledge management and collaboration platform. This repository contains personal notes, journal entries, and project documentation

## Quick Navigation

For detailed guidance, see specialized guides in `.claude/guides/`:
- **[Index System Guide](.claude/guides/index-system.md)** - Complete reference for 9 auto-generated indexes
- **[Query Patterns](.claude/guides/query-patterns.md)** - Decision trees for efficient queries (5-30x token savings)
- **[Skills Reference](.claude/guides/skills-reference.md)** - Automation workflows and creation guide
- **[Guides README](.claude/guides/README.md)** - Navigation flowchart for all guides

## Repository Structure

- **journals/** - Daily journal entries in the format `YYYY_MM_DD.md`
- **pages/** - Topic-specific pages (e.g., project documentation, meeting notes, concepts)
- **logseq/** - Logseq application configuration and metadata
  - `config.edn` - Main Logseq configuration file (Clojure EDN format)
  - `custom.css` - Custom CSS styling for the Logseq interface
  - `.recycle/` - Deleted pages (gitignored)
- **.claude/** - Claude Code configuration and optimization
  - `guides/` - Specialized documentation (index system, query patterns, skills)
  - `skills/` - Automation workflow templates
  - `indexes/` - 9 auto-generated indexes (regenerate on every commit)

## Logseq-Specific Conventions

### File Format
- All content files are in **Markdown format**
- Journal files use the naming convention: `YYYY_MM_DD.md`
- Page files use human-readable names with spaces (e.g., `Phoenix projext.md`)

### Markdown Syntax
Logseq uses an outliner-based approach with specific markdown conventions:

- **Blocks**: Each line starting with `-` or a number is a block (the fundamental unit in Logseq)
- **Page references**: `[[Page Name]]` creates a bidirectional link to another page
- **Block references**: `((block-id))` references a specific block
- **Tags**: Can be created using `#tag` or `[[tag]]`
- **Task markers**: `NOW`, `LATER`, `TODO`, `DOING`, `DONE` for task management
- **Properties**: Block properties are defined using `:property-name: value` syntax
- **Logbook**: Time tracking data stored in `:LOGBOOK:` blocks with `CLOCK:` entries

### Configuration (config.edn)
- Written in **Clojure EDN** (Extensible Data Notation)
- Key settings:
  - Preferred workflow: `:now` style (NOW/LATER instead of TODO/DOING)
  - Journal file format: `yyyy_MM_dd`
  - Default directories: `journals/` and `pages/`

## Git Workflow

This repository uses **automatic commits** via Logseq's git integration:
- Commits are made automatically with messages like "Auto saved by Logseq"
- The repository is always on the `main` branch
- There is no build process, testing, or CI/CD - it's pure content
- **Index regeneration**: `.git/hooks/post-commit` auto-generates 9 indexes on every commit

## Auto-Generated Index System

This repository uses [logseq-claude-indexer](https://github.com/dyluth/logseq-claude-indexer) to auto-generate 9 optimized indexes on every commit. These indexes provide **5-30x token efficiency** vs reading raw journal files.

### Progressive Disclosure Strategy (ALWAYS follow this order)

**Level 1: Start Here** (400-500 tokens)
- `dashboard.md` - Master KPI overview, current priorities, top projects
- `timeline-recent.md` - Last 7 days activity summary
- `tasks-by-priority.md` - High-priority tasks only ([#A], [#B])

**Level 2: Filtered Views** (300-2,000 tokens)
- `tasks-by-status.md` - All 324 tasks by workflow stage (NOW/DOING/TODO/LATER/DONE)
- `time-tracking.md` - Project time allocation (33,690h logged across 117 tasks)
- `missing-pages.md` - Pages with 5+ refs that don't exist (documentation gaps)

**Level 3: Complete Datasets** (5,000-8,000 tokens, use sparingly)
- `task-index.md` - All tasks with full metadata (file locations, references, time logs)
- `timeline-full.md` - 96-day complete activity history
- `reference-graph.md` - Full knowledge graph (including pages, references & hub identification)

**Critical Rule:** Always read Level 1 first. Only move to Level 2/3 if insufficient.

**For complete reference:** See [Index System Guide](.claude/guides/index-system.md)

## Available Skills (Automation)

Skills are workflow automation templates that auto-activate on keyword patterns:

**Common Skills:**
- **"status report"** → Generates executive summary using indexes (dashboard + timeline + tasks)
- **"weekly review"** → Cross-project activity synthesis (timeline + tasks + time tracking)
- **"deep dive on [project]"** → Comprehensive project analysis (indexes + files)
- **"close project [name]"** → Systematic project closure with cleanup (cross-index analysis)
- **"knowledge gaps"** → Identifies missing documentation (missing-pages + reference-graph)
- **"time ROI"** → Analyzes time investment vs output (time-tracking + tasks)
- **"relationship with [person]"** → Maps contact interactions (references + tasks + timeline)

**For complete catalog and usage:** See [Skills Reference](.claude/guides/skills-reference.md)

## Query Optimization Quick Reference

**Efficient query patterns (use smallest index first):**

| Query | Optimal Index | Tokens | ❌ Avoid | Tokens | Efficiency Gain |
|-------|--------------|--------|----------|--------|-----------------|
| "What's my status?" | dashboard.md | 500 | task-index.md | 8,000 | **16x faster** |
| "Recent activity?" | timeline-recent.md | 350 | 7 journal files | 2,600 | **7x faster** |
| "Priorities?" | tasks-by-priority.md | 400 | tasks-by-status.md | 2,000 | **5x faster** |
| "Time on project?" | time-tracking.md | 300 | grep journals | 10,000+ | **33x faster** |

**Progressive disclosure example:**
```
"What's my status?"
├─ Read dashboard.md (500 tokens) → Sufficient? YES → DONE
└─ If more detail → timeline-recent.md (350 tokens) → Sufficient? YES → DONE
   └─ If task detail → tasks-by-status.md (NOW section, ~500 tokens)
```

**For complete decision trees:** See [Query Patterns Guide](.claude/guides/query-patterns.md)

## Common Operations

### Searching Content
Use grep to search across all markdown files:
```bash
grep -r "search term" pages/ journals/
```

**Better:** Use indexes first before grepping raw files.

### Finding Page References
To find all references to a specific page:
```bash
grep -r "\[\[Page Name\]\]" .
```

**Better:** Use `.claude/indexes/reference-graph.md` to see all references.

### Viewing Recent Journals
```bash
ls -lt journals/ | head -10
```

**Better:** Use `.claude/indexes/timeline-recent.md` for last 7 days summary.

### Checking Logseq Configuration
```bash
cat logseq/config.edn
```

## Content Focus

The primary content in this repository relates to:
<TO BE FILLED IN>

## Important Notes

1. **No Code Execution**: This is a documentation repository with no executable code, tests, or build process
2. **Manual Editing**: When editing files, preserve the Logseq markdown structure (bullet points, indentation)
3. **Git Commits**: If creating commits, follow the existing pattern of simple commit messages
4. **File References**: Maintain bidirectional links by updating all references when renaming pages
5. **Logseq Config**: The `config.edn` file uses Clojure syntax - be careful with formatting when editing
6. **Index Optimization**: Always check indexes before reading raw journal files (5-30x efficiency gain)
7. **Progressive Disclosure**: Start with Level 1 indexes (dashboard, timeline-recent, tasks-by-priority)
