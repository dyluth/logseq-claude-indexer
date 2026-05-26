# Index System Complete Reference

**Last Updated:** 2025-11-07
**Purpose:** Complete guide to the 9 auto-generated index files in `.claude/indexes/`

---

## Overview

This repository uses [logseq-claude-indexer](https://github.com/dyluth/logseq-claude-indexer) to automatically generate 9 optimized index files on **every git commit**. These indexes are materialized views providing:

- 📊 **Dashboard KPIs** - Instant project health metrics
- 🎯 **Task Management** - 324 tasks across workflow stages
- ⏰ **Time Tracking** - 33,690 hours logged across projects
- 🔗 **Knowledge Graph** - 226 pages, 295 bidirectional references
- 🕐 **Timeline Analysis** - 96 days of complete activity history
- 🔍 **Gap Detection** - Missing pages identified by reference count

**Context efficiency:** Indexes are **5-30x more token-efficient** than reading raw journal files.

---

## Progressive Disclosure Hierarchy

**CRITICAL: Always follow this order (don't skip levels):**

```
Level 1: Start Here (400-500 tokens)
├── dashboard.md - Master KPI overview
├── timeline-recent.md - Last 7 days activity
└── tasks-by-priority.md - High-priority tasks only

    ↓ [If more detail needed]

Level 2: Filtered Views (300-2,000 tokens)
├── tasks-by-status.md - All tasks by workflow stage
├── time-tracking.md - Project time allocation
└── missing-pages.md - Documentation gaps

    ↓ [If comprehensive view needed]

Level 3: Complete Datasets (5,000-8,000 tokens)
├── task-index.md - All 324 tasks with metadata
├── timeline-full.md - 96-day complete history
└── reference-graph.md - Full knowledge graph

    ↓ [Only if indexes insufficient]

Level 4: Raw Source Files
└── journals/*.md, pages/*.md - Last resort
```

**Rule:** Start Level 1 → Move to Level 2 only if insufficient → Use Level 3 sparingly → Avoid Level 4 when possible.

---

## Complete Index Inventory

### Level 1: Start Here (Aggregated Views)

#### dashboard.md
- **Size:** ~500 tokens
- **Purpose:** Master KPI dashboard with quick stats and current priorities
- **When to use:** ALWAYS START HERE for any status query
- **Refresh:** Every commit (auto-generated)

**Contents:**
- Quick stats (total tasks, completion rate, time tracking adoption)
- Current priorities ([#A] high-priority tasks, top 2 shown)
- Recent activity (last few days with task counts)
- Top projects (by time investment, top 5 shown)
- Pages to create (high-reference missing pages, top 5 shown)
- Links to detailed reports

**Structure example:**
```markdown
## 📊 Quick Stats
- Total Tasks: 324
- Completion Rate: 60.8% (197 DONE)
- Time Tracking: 36.1% adoption, 33690h 10m logged

## 🎯 Current Priorities [#A]
- [NOW] [[HI-setup]] start [[business plan]]
- [TODO] Working demo [[HI-prototype]]

## 📅 Recent Activity
### Thursday, Nov 6
- 1 NOW task, 1 TODO task

## 📁 Top Projects
- Data shield: 3529h 19m
- WorkJobThing: 2440h 53m (5 active tasks)
```

**Use for:**
- "What's my status?"
- "What are my priorities?"
- "Show me current project health"
- Quick health check queries

**Token efficiency:** 500 tokens vs 8,000 (task-index.md) = **16x faster**

---

#### timeline-recent.md
- **Size:** ~350 tokens
- **Purpose:** Last 7 days activity summary
- **When to use:** "What happened recently?" queries
- **Refresh:** Every commit (auto-generated)

**Contents:**
- Last 7 days with activity (date, day of week)
- Task counts by status (NOW/TODO/DONE per day)
- Time logged per day (if applicable)
- Actual task descriptions for recent days

**Structure example:**
```markdown
## Thursday, November 6, 2025
- 1 NOW task
- 1 TODO task

**Tasks** (2):
- [NOW] FOLLOWUP Tertiary referrals...
- [TODO] [[linkedin-post-idea]] - AI as a black box...

## Wednesday, November 5, 2025
- 7 NOW tasks
- 1 TODO task
```

**Use for:**
- "What happened this week?"
- "Recent activity summary"
- "What did I work on yesterday?"

**Token efficiency:** 350 tokens vs 2,600 (7 journal files) = **7x faster**

---

#### tasks-by-priority.md
- **Size:** ~400 tokens
- **Purpose:** High-priority tasks only ([#A], [#B])
- **When to use:** "What's urgent?" or "Show priorities"
- **Refresh:** Every commit (auto-generated)

**Contents:**
- Total high-priority count
- Tasks grouped by status (NOW/TODO/DONE)
- Each task shows: description, file location, references, time logged
- Only includes [#A] and [#B] priority tasks

**Structure example:**
```markdown
**Total High Priority**: 7 tasks

## NOW (1)
### [[HI-setup]] start [[business plan]] - with Gemini
- File: journals/2025_04_24.md:6
- References: [[HI-setup]], [[business plan]]

## TODO (1)
### Working demo [[HI-prototype]]
- File: journals/2025_04_22.md:2
```

**Use for:**
- "What are my priorities?"
- "Show high-priority tasks"
- "What's urgent?"

**Token efficiency:** 400 tokens vs 2,000 (tasks-by-status.md) = **5x faster**

---

### Level 2: Filtered Views (Detailed)

#### tasks-by-status.md
- **Size:** ~2,000 tokens
- **Purpose:** All 324 tasks organized by workflow stage
- **When to use:** Need task details by specific status
- **Refresh:** Every commit (auto-generated)

**Contents:**
- Statistics (total, completion rate, time tracking, priority breakdown, status breakdown)
- Tasks grouped by status: NOW (48), DOING (3), TODO (73), LATER (3), DONE (197)
- Each task shows: description, priority, file location, time logged (if any)

**Structure example:**
```markdown
## Statistics
- Total Tasks: 324
- Completion Rate: 60.8%
- By Status:
  - NOW: 48
  - DOING: 3
  - TODO: 73
  - LATER: 3
  - DONE: 197

## NOW (48)
- [[HI-setup]] start [[business plan]]
- Traction - BOOK - READ
- [... 46 more]

## DOING (3)
- work through accelerator advice
- investigate [[CQL]] further
- work on [[HI-Financial model]]
```

**Use for:**
- "Show all NOW tasks"
- "What's in TODO backlog?"
- "List completed tasks (DONE)"
- Task status filtering

**Token efficiency:** Use this only when you need ALL tasks in a specific status, not for general queries.

---

#### time-tracking.md
- **Size:** ~1,500 tokens
- **Purpose:** Project time allocation analytics
- **When to use:** "How much time on [project]?" queries
- **Refresh:** Every commit (auto-generated)

**Contents:**
- Summary (total time, tasks tracked, avg time/task, most productive week)
- Top projects (ranked by time, shows task count and avg)
- Weekly breakdown (last 8 weeks shown)
- Breakdown by priority ([#A] vs no priority)
- Breakdown by status (DONE vs TODO vs others)

**Structure example:**
```markdown
## Summary
- Total Time Logged: 33690h 10m
- Tasks Tracked: 117 / 324 (36.1%)
- Avg Time/Task: 287h 57m

## Top Projects
### WorkJobThing
- Time: 2440h 53m (3 tasks, avg 813h 37m/task)
### Data shield
- Time: 3529h 19m (4 tasks, avg 882h 19m/task)

## Weekly Breakdown
- Week of 2025-10-27: 210h 16m (2 tasks)
- Week of 2025-10-20: 1000h 28m (7 tasks)
```

**Use for:**
- "How much time on WorkJobThing?"
- "Time investment by project"
- "Which projects consume most time?"
- ROI analysis (time vs output)

**Token efficiency:** 1,500 tokens vs 10,000+ (grepping all journals) = **7x faster**

**Note:** Time entries can span long periods - hours represent cumulative project time, not session duration.

---

#### missing-pages.md
- **Size:** ~300 tokens
- **Purpose:** Pages with 5+ references that don't exist yet
- **When to use:** "What should I document?" or gap analysis
- **Refresh:** Every commit (auto-generated)

**Contents:**
- Total missing page count
- Categorized by type (People, Concepts)
- Each entry shows: reference count, pages that reference it
- Sorted by reference count (highest first)

**Structure example:**
```markdown
**Pages with 5+ references that don't exist yet**: 9

## People (4)
### [[John smith]]
- References: 11
- Referenced from: [[2025_09_09]], [[2025_10_13]], ...

### [[Gemma Snell - mentor]]
- References: 7
- Referenced from: [[2025_10_05]], [[2025_10_08]], ...

## Concepts (5)
### [[Data shield]]
- References: 5
- Referenced from: [[2025_06_04]], [[2025_06_30]], ...
```

**Use for:**
- "What pages should I create?"
- "Knowledge gaps in my system"
- "Who/what am I referencing but haven't documented?"
- Prioritizing documentation work

**Token efficiency:** Direct answer to "what to create next" without reading all references.

---

### Level 3: Complete Datasets (Use Sparingly)

#### task-index.md
- **Size:** ~8,000 tokens
- **Purpose:** Complete task inventory with full metadata
- **When to use:** Deep task search by keyword or comprehensive analysis
- **Refresh:** Every commit (auto-generated)

**Contents:**
- All 324 tasks with complete details
- Grouped by status (NOW/DOING/TODO/LATER/DONE)
- Each task shows: full description, file location (path:line), references ([[links]]), time logged

**Structure example:**
```markdown
## NOW Tasks (48)

### [#A] [[HI-setup]] start [[business plan]] - with Gemini
- File: journals/2025_04_24.md:6
- References: [[HI-setup]], [[business plan]]

### Traction - BOOK - READ - written by duckduckgo creator
- File: journals/2025_05_06.md:77

[... all 324 tasks follow]
```

**Use for:**
- grep for specific keywords across all tasks
- Finding tasks by reference ([[specific-page]])
- Comprehensive task audit
- When you need complete task context

**Token efficiency:** Only use when tasks-by-status.md or tasks-by-priority.md insufficient. This is large.

**Search pattern example:**
```
User: "Find all tasks mentioning 'MoonLaunch'"
Claude: Read task-index.md → grep "MoonLaunch" → return matching tasks
```

---

#### timeline-full.md
- **Size:** ~5,000 tokens
- **Purpose:** Complete 96-day activity history
- **When to use:** Historical analysis or specific date queries
- **Refresh:** Every commit (auto-generated)

**Contents:**
- 96 days with activity (reverse chronological)
- Each day shows: task counts by status, time logged
- Includes days with no tasks (marked as "No tasks")

**Structure example:**
```markdown
**Total Days**: 96 days with activity

## 2025-11-06 (Thu)
- 1 NOW task
- 1 TODO task

## 2025-11-05 (Wed)
- 7 NOW tasks
- 1 TODO task

## 2025-11-04 (Tue)
- *No tasks*

## 2025-10-29 (Wed)
- 5 NOW tasks
- 1 DONE task
- ⏱ 169h 7m logged
```

**Use for:**
- "What happened in October?"
- "Activity on specific date"
- Long-term pattern analysis
- "When did I last work on X?"

**Token efficiency:** Use timeline-recent.md for recent queries. Only use this for historical analysis.

---

#### reference-graph.md
- **Size:** ~5,000 tokens
- **Purpose:** Complete knowledge graph (226 pages, 295 references)
- **When to use:** Network analysis, discovering connections
- **Refresh:** Every commit (auto-generated)

**Contents:**
- Total page and reference counts
- Hub pages (most referenced, top 10)
- Detailed page analysis showing:
  - Outbound references (what this page links to)
  - Inbound references (what links to this page)
  - File location (if exists) or "Page not yet created"

**Structure example:**
```markdown
Total Pages: 226
Total References: 295

## Hub Pages (Most Referenced)
1. [[WorkJobThing]] - 14 inbound references
2. [[Test Company]] - 12 inbound references
3. [[John smith]] - 11 inbound (Page not yet created)

## Page Details
### [[Test Company]]
- File: pages/Test Company.md
- Outbound References (5):
  - [[HIPPA]]
  - [[FHIR]]
  - [...]
- Inbound References (12):
  - [[2025_04_06]]
  - [[2025_04_22]]
  - [...]
```

**Use for:**
- "What's connected to WorkJobThing?"
- "Show me all pages linking to X"
- Hub page identification (most central topics)
- Network analysis
- Discovering related context

**Token efficiency:** Only use when you need to understand page relationships. Very large file.

---

## Cross-Index Query Patterns

**Indexes are designed to be combined. Common patterns:**

### Pattern 1: Project Health Report
```
Query: "Status of WorkJobThing project?"

Optimal workflow:
1. dashboard.md → Check "Top Projects" → WorkJobThing: 2,440h
2. tasks-by-status.md → grep "WorkJobThing" → Find active tasks
3. time-tracking.md → WorkJobThing: 3 tasks, avg 813h/task
4. timeline-recent.md → Recent WorkJobThing activity
5. reference-graph.md → WorkJobThing: 14 inbound refs → find pages

Result: Comprehensive project context (~3,000 tokens)
vs. reading all WorkJobThing journals (~10,000 tokens) = 3x efficiency
```

### Pattern 2: Time ROI Analysis
```
Query: "Which projects had best ROI this month?"

Optimal workflow:
1. timeline-recent.md → Recent tasks by project
2. time-tracking.md → Time invested per project
3. tasks-by-status.md (DONE) → Completed tasks per project
4. Calculate: Completed tasks / Hours invested
5. Generate ranked report

All data in indexes, no journal reads needed
```

### Pattern 3: Knowledge Gap Discovery
```
Query: "What am I forgetting to document?"

Optimal workflow:
1. missing-pages.md → 9 pages with 5+ refs:
   - John smith (11 refs) - Person
   - alistair bastian (7 refs) - Concept
   - Gemma Snell (7 refs) - Person
2. reference-graph.md → Hub pages with high refs
3. tasks-by-status.md → Tasks mentioning undocumented concepts
4. Generate prioritized creation list

Direct answer without extensive journal reading
```

### Pattern 4: Contact Relationship Mapping
```
Query: "Map my relationship with John smith"

Optimal workflow:
1. missing-pages.md → John smith: 11 refs, no page
2. reference-graph.md → Find all 11 reference sources
3. task-index.md → grep "John smith" → tasks
4. timeline-full.md → When did interactions occur?
5. Read 2-3 journal entries for conversation notes
6. Generate relationship timeline + action items

Comprehensive relationship view with targeted journal reads
```

### Pattern 5: Weekly Progress Synthesis
```
Query: "What did I accomplish this week?"

Optimal workflow:
1. timeline-recent.md → 7-day activity overview
2. tasks-by-status.md (DONE section) → Completed tasks
3. time-tracking.md → Time allocation this week
4. dashboard.md → Overall completion rate trend

Weekly summary without reading all journals
Token usage: ~2,500 vs ~10,000+ = 4x efficiency
```

---



### Pattern: Weekly Cross-Project Review
```
Query: "Weekly review across all projects"

Workflow:
1. timeline-recent.md → 7-day activity
2. time-tracking.md → Hours per project this week
3. tasks-by-status.md (DONE) → Completed tasks
4. For each active project: Check Logseq tasks
5. Generate comparative weekly report

Cross-source synthesis showing complete picture
```

---

## Token Efficiency Comparisons

| Query Type | Index Approach | Raw File Approach | Efficiency Gain |
|------------|---------------|-------------------|-----------------|
| Status check | dashboard.md (500 tokens) | task-index.md (8,000 tokens) | **16x** |
| Recent activity | timeline-recent.md (350 tokens) | 7 journals (2,600 tokens) | **7x** |
| Priorities | tasks-by-priority.md (400 tokens) | tasks-by-status.md (2,000 tokens) | **5x** |
| Time on project | time-tracking.md filtered (300 tokens) | grep all journals (10,000+ tokens) | **33x** |
| Project deep dive | Multi-index (3,000 tokens) | All related journals (10,000+ tokens) | **3x** |

**Average efficiency gain: 5-30x faster queries**

---

## Troubleshooting

### Q: Indexes seem stale?
**A:** Run: `git add . && git commit -m "refresh indexes"`

Indexes regenerate automatically on every commit via `.git/hooks/post-commit` hook.

### Q: Can't find information in indexes?
**A:** Check these possibilities:
1. **missing-pages.md** - May need to create the page first
2. **Information in raw journals** - Indexes show tasks, not all journal content. Read specific journal file.
3. **Spelling variation** - Try alternate search terms

### Q: Which index for task search?
**A:**
- **By status** → tasks-by-status.md
- **By priority** → tasks-by-priority.md
- **By keyword/project** → task-index.md (grep within file)
- **By reference** → task-index.md (grep for [[page-name]])

### Q: Dashboard shows wrong current date?
**A:** Indexes use UTC timestamps - "Generated: 2025-11-06T23:27:44Z" = when last committed. Commit to refresh.

### Q: How do I update index generation?
**A:**
- See: `.git/hooks/post-commit` for hook configuration
- Configuration: logseq-claude-indexer settings
- **Don't manually edit** `.claude/indexes/*` - they're auto-generated

### Q: Can I disable specific indexes?
**A:** Yes - configure logseq-claude-indexer. But all 9 are lightweight and provide different views. Recommend keeping all.

---

## Best Practices

### DO:
✅ **Start with Level 1** - dashboard.md, timeline-recent.md, tasks-by-priority.md
✅ **Use progressive disclosure** - Only move to Level 2/3 if needed
✅ **Combine indexes** - Cross-reference for comprehensive analysis
✅ **Commit frequently** - Keeps indexes fresh
✅ **Use for navigation** - Indexes show file:line locations, jump directly

### DON'T:
❌ **Skip levels** - Don't jump straight to task-index.md or timeline-full.md
❌ **Load multiple large indexes** - Use filtered views (Level 2) instead
❌ **Read journals first** - Check indexes before reading raw files
❌ **Manually edit indexes** - They regenerate on every commit
❌ **Ignore missing-pages.md** - It shows your documentation gaps

---

## Index Data Quality

### Time Tracking Statistics
- **Total logged:** 33,690 hours across 117 tasks (36% adoption)
- **Top projects:** Data shield (3,529h), WorkJobThing (2,440h), DiscoverNow (2,322h)
- **Note:** Some entries span long periods - hours = cumulative project time, not session duration

### Task Distribution
- **Total tasks:** 324
  - NOW: 48 (active focus)
  - DOING: 3 (in progress)
  - TODO: 73 (backlog)
  - LATER: 3 (deferred)
  - DONE: 197 (60.8% completion rate)

### Knowledge Graph Health
- **Pages:** 226 total
- **References:** 295 bidirectional links
- **Hub pages:** WorkJobThing (14 refs), Test Company (12 refs), John smith (11 refs)
- **Missing pages:** 9 with 5+ references (should create)

---

## Related Documentation

- **Query Patterns:** See [query-patterns.md](./query-patterns.md) for decision trees
- **Skills Reference:** See [skills-reference.md](./skills-reference.md) for automation

---

**Last Updated:** 2025-11-07
