# Enhanced Indexes Design

## Executive Summary

This design extends `logseq-claude-indexer` to generate additional indexes that enable Claude Code to answer natural questions about your knowledge base without external context.

**Goals**:
1. Enable time-based queries ("What happened last week?")
2. Surface missing pages that need creation
3. Provide time allocation analytics
4. Highlight urgent/priority tasks
5. Maintain <3s performance for 116 files

---

## User Questions & Solutions

| Question | Solution | Index File |
|----------|----------|------------|
| "What happened last week?" | Timeline view of journal entries | `timeline-index.md` |
| "What are my key contacts?" | Already in reference-graph.md | *(enhance existing)* |
| "What's completed vs pending?" | Already in task-index.md | *(add statistics)* |
| "How much time invested in X?" | Aggregate LOGBOOK by project | `time-tracking-report.md` |
| "Which pages are most referenced but don't exist?" | List non-existent pages by reference count | `missing-pages-report.md` |
| "What tasks are urgent?" | Extract priority markers ([#A], [#B], [#C]) | *(enhance task-index.md)* |
| "What's my time allocation by project?" | LOGBOOK aggregation by PageRefs | `time-tracking-report.md` |

---

## Architecture Overview

### New Components

```
Parser (Enhanced)
  ├─ Extract priority markers [#A], [#B], [#C]
  └─ Existing: tasks, references, LOGBOOK

Indexers (New)
  ├─ TimelineIndexer → chronological journal view
  ├─ MissingPagesIndexer → referenced but non-existent pages
  └─ TimeTrackingIndexer → LOGBOOK aggregations

Writers (New)
  ├─ TimelineWriter → timeline-index.md
  ├─ MissingPagesWriter → missing-pages-report.md
  └─ TimeTrackingWriter → time-tracking-report.md

Indexers (Enhanced)
  └─ TaskIndexer → add priority grouping, statistics

Writers (Enhanced)
  └─ TaskWriter → add priority sections, completion stats
```

### Data Flow

```
Scanner
  ↓
Parser (extract priority markers)
  ↓
[Task, Reference, File, LOGBOOK] data
  ↓
  ├─→ TaskIndexer (with priority) → task-index.md (enhanced)
  ├─→ GraphIndexer → reference-graph.md
  ├─→ TimelineIndexer → timeline-index.md (NEW)
  ├─→ MissingPagesIndexer → missing-pages-report.md (NEW)
  └─→ TimeTrackingIndexer → time-tracking-report.md (NEW)
```

---

## Detailed Component Design

### 1. Priority Extraction (Parser Enhancement)

**Goal**: Extract `[#A]`, `[#B]`, `[#C]` priority markers from tasks

**Implementation**: `internal/parser/tasks.go`

```go
type Priority string

const (
    PriorityHigh   Priority = "A"
    PriorityMedium Priority = "B"
    PriorityLow    Priority = "C"
    PriorityNone   Priority = ""
)

// Add to Task model (pkg/models/task.go)
type Task struct {
    Status      TaskStatus
    Priority    Priority        // NEW
    Description string
    PageRefs    []string
    SourceFile  string
    LineNumber  int
    Logbook     []LogbookEntry
}

// In extractTaskDescription, also extract priority
func extractPriority(line string) Priority {
    // Regex: \[#([ABC])\]
    re := regexp.MustCompile(`\[#([ABC])\]`)
    matches := re.FindStringSubmatch(line)
    if len(matches) > 1 {
        return Priority(matches[1])
    }
    return PriorityNone
}
```

**Example Input**:
```markdown
- NOW [#A] [[Project Phoenix]] - Complete authentication
- LATER [#B] Research alternatives
- TODO Low priority task (no marker)
```

**Parsed Output**:
```go
Task{Status: "NOW", Priority: "A", Description: "[[Project Phoenix]] - Complete authentication"}
Task{Status: "LATER", Priority: "B", Description: "Research alternatives"}
Task{Status: "TODO", Priority: "", Description: "Low priority task"}
```

---

### 2. Timeline Index (NEW)

**Goal**: Chronological view of journal activity

**File**: `timeline-index.md`

**Indexer**: `internal/indexer/timeline_indexer.go`

```go
type TimelineIndexer struct {
    GeneratedAt time.Time
    Entries     []TimelineDay
}

type TimelineDay struct {
    Date           time.Time
    JournalPath    string
    TasksCreated   []models.Task
    TimeLogged     time.Duration
    ReferencesAdded []models.PageReference
}

func BuildTimelineIndex(tasks []models.Task, refs []models.PageReference, files []models.File) *TimelineIndexer {
    // Group tasks by source journal date
    // Group LOGBOOK entries by date
    // Group references by source journal date
    // Sort by date descending (newest first)
}
```

**Output Format** (`timeline-index.md`):

```markdown
# Logseq Timeline

Generated: 2025-11-06T16:00:00Z

---

## Last 7 Days

### Nov 6, 2025 (Today)
- **Journal**: `journals/2025_11_06.md`
- **Tasks Created**: 8
  - 3 NOW tasks
  - 5 LATER tasks
- **Time Logged**: 6h 30m
- **New References**: 12 pages referenced

**Key Activity**:
- NOW [#A] [[Project Phoenix]] - Authentication module (3h 15m logged)
- NOW [[Mobile App]] - Beta testing
- LATER [[Kubernetes]] research

**Pages Referenced**:
- [[Sarah Chen - Tech Lead]] (5 times)
- [[Sprint 24]] (3 times)

---

### Nov 5, 2025
- **Journal**: `journals/2025_11_05.md`
- **Tasks Created**: 5
- **Time Logged**: 4h 15m
- **Completed Tasks**: 3

---

## Last 30 Days (Summary)

- Total journal entries: 30
- Tasks created: 124
- Tasks completed: 89
- Time logged: 142h 30m
- Most active day: Nov 1 (12 tasks, 8h logged)

---

## By Week

### Week of Oct 28 - Nov 3
- Tasks: 45 created, 32 completed
- Time: 38h 15m
- Focus: [[Project Phoenix]] (18h), [[Mobile App]] (12h)

### Week of Oct 21 - Oct 27
- Tasks: 38 created, 28 completed
- Time: 31h 45m
- Focus: [[Sprint 23]] planning
```

**Grouping Logic**:
- Last 7 days: Show daily detail
- Last 30 days: Show summary + weekly rollup
- Older: Monthly summary

**Benefits for Claude**:
- "What did I work on yesterday?" → Check yesterday's entry
- "What happened last week?" → Check weekly summary
- "Show my November activity" → Scan monthly sections

---

### 3. Missing Pages Report (NEW)

**Goal**: Identify non-existent pages with 5+ references

**File**: `missing-pages.md`
**Threshold**: Minimum 5 references (reduces noise)

**Indexer**: `internal/indexer/missing_pages_indexer.go`

```go
type MissingPagesIndex struct {
    GeneratedAt   time.Time
    TotalMissing  int
    MissingPages  []MissingPage
}

type MissingPage struct {
    PageName       string
    ReferenceCount int
    ReferencedFrom []PageReference
    SuggestedType  string // "person", "project", "concept", "date"
}

func BuildMissingPagesIndex(graph *ReferenceGraph) *MissingPagesIndex {
    var missing []MissingPage

    for pageName, node := range graph.Nodes {
        // Only include pages with 5+ references (confirmed threshold)
        if node.FilePath == "" && node.ReferenceCount >= 5 {
            // Determine type heuristically
            pageType := classifyPageType(pageName)

            missing = append(missing, MissingPage{
                PageName:       pageName,
                ReferenceCount: node.ReferenceCount,
                ReferencedFrom: getReferences(graph, pageName),
                SuggestedType:  pageType,
            })
        }
    }

    // Sort by reference count descending
    sort.Slice(missing, func(i, j int) bool {
        return missing[i].ReferenceCount > missing[j].ReferenceCount
    })

    return &MissingPagesIndex{
        GeneratedAt:  time.Now(),
        TotalMissing: len(missing),
        MissingPages: missing,
    }
}

func classifyPageType(pageName string) string {
    // Heuristics:
    if strings.Contains(pageName, " - ") {
        return "person" // e.g., "Sarah Chen - Tech Lead"
    }
    if strings.HasSuffix(pageName, "th, 2025") {
        return "date" // e.g., "Nov 24th, 2025"
    }
    if strings.Contains(strings.ToLower(pageName), "sprint") {
        return "project" // e.g., "Sprint 24"
    }
    return "concept"
}
```

**Output Format** (`missing-pages.md`):

```markdown
# Missing Pages Report

Generated: 2025-11-06T16:00:00Z
Total Missing: 8 pages (showing only pages with 5+ references)

These pages are frequently referenced but don't exist yet. Creating them could improve your knowledge base connectivity.

---

## High Priority (5+ references)

### [[Matt Stammers - SETT: Theme Lead for Data & AI]]
- **Type**: Person
- **References**: 11
- **Referenced from**:
  - `journals/2025_04_15.md:12` - Meeting notes
  - `journals/2025_05_03.md:8` - Discussion about [[Data Strategy]]
  - `journals/2025_06_12.md:23` - Follow-up on [[SDE Access]]
  - `pages/SETT Team.md:5` - Team structure
  - *... and 7 more*

**Suggested Content**:
- Contact information
- Role and expertise
- Key projects
- Meeting history

---

### [[Observability Stack]]
- **Type**: Project/Concept
- **References**: 8
- **Referenced from**:
  - `journals/2025_11_01.md:25` - Research task
  - `pages/Project Phoenix.md:34` - Infrastructure dependency
  - `pages/Microservices Architecture.md:56` - Monitoring solution
  - *... and 5 more*

**Suggested Content**:
- Architecture overview
- Tools and technologies
- Implementation timeline
- Team ownership

---

## Medium Priority (2-4 references)

### [[Sprint 24]]
- **Type**: Project
- **References**: 4
- Referenced from 3 journals, 1 page

### [[Database Migration]]
- **Type**: Concept
- **References**: 3
- Referenced from 2 journals, 1 page

---

## Low Priority (1 reference)

- [[Technical RFC]] (1 ref) - Concept
- [[Junior Engineers Program]] (1 ref) - Project
- [[Conference Presentations]] (1 ref) - Concept

---

## Summary by Type

- **People**: 12 missing contact pages
- **Projects**: 8 missing project pages
- **Concepts**: 10 missing concept pages
- **Dates**: 2 missing date pages

**Action**: Consider creating pages for high-priority items first to improve knowledge connectivity.
```

**Benefits for Claude**:
- "What pages should I create?" → Top of missing pages report
- "Is there a page about X?" → Check if listed here
- "Who are my key contacts without pages?" → Filter by type=person

---

### 4. Time Tracking Report (NEW)

**Goal**: Aggregate LOGBOOK data to show time investment

**File**: `time-tracking-report.md`

**Indexer**: `internal/indexer/time_tracking_indexer.go`

```go
type TimeTrackingIndex struct {
    GeneratedAt      time.Time
    TotalTimeLogged  time.Duration
    ByProject        map[string]time.Duration
    ByWeek           map[string]time.Duration // Week of YYYY-MM-DD
    ByPriority       map[Priority]time.Duration
    ByStatus         map[TaskStatus]time.Duration
    TopProjects      []ProjectTime
}

type ProjectTime struct {
    ProjectName  string
    TotalTime    time.Duration
    Percentage   float64
    TaskCount    int
    RecentActivity time.Time
}

func BuildTimeTrackingIndex(tasks []models.Task) *TimeTrackingIndex {
    index := &TimeTrackingIndex{
        GeneratedAt: time.Now(),
        ByProject:   make(map[string]time.Duration),
        ByWeek:      make(map[string]time.Duration),
        ByPriority:  make(map[Priority]time.Duration),
        ByStatus:    make(map[TaskStatus]time.Duration),
    }

    for _, task := range tasks {
        totalTime := task.TotalDuration()
        index.TotalTimeLogged += totalTime

        // Aggregate by project (from PageRefs)
        if len(task.PageRefs) > 0 {
            project := task.PageRefs[0]
            index.ByProject[project] += totalTime
        }

        // Aggregate by week (from LOGBOOK timestamps)
        for _, entry := range task.Logbook {
            week := getWeekStart(entry.Start)
            index.ByWeek[week] += entry.Duration
        }

        // Aggregate by priority and status
        index.ByPriority[task.Priority] += totalTime
        index.ByStatus[task.Status] += totalTime
    }

    // Calculate top projects
    index.TopProjects = calculateTopProjects(index.ByProject, tasks)

    return index
}

func getWeekStart(t time.Time) string {
    // Get Monday of the week containing t
    weekday := int(t.Weekday())
    if weekday == 0 {
        weekday = 7 // Sunday = 7
    }
    monday := t.AddDate(0, 0, -(weekday - 1))
    return monday.Format("2006-01-02")
}
```

**Output Format** (`time-tracking.md`):

```markdown
# Time Tracking Report

Generated: 2025-11-06T16:00:00Z
Total Time Logged: 342h 45m

---

## Time by Project

### Top Projects (by time invested)

1. **[[Project Phoenix]]** - 127h 30m (37.2%)
   - Tasks: 28 (8 NOW, 12 DOING, 8 DONE)
   - Recent activity: 2025-11-01
   - Key focus: Authentication module (45h), API design (32h)

2. **[[Mobile App]]** - 68h 15m (19.9%)
   - Tasks: 15 (3 NOW, 5 DOING, 7 DONE)
   - Recent activity: 2025-11-03
   - Key focus: UI implementation (28h), Performance (18h)

3. **[[Hearth Insights]]** - 52h 45m (15.4%)
   - Tasks: 22 (4 NOW, 3 DOING, 15 DONE)
   - Recent activity: 2025-10-30
   - Key focus: Research (30h), Planning (15h)

4. **[[Unitary]]** - 38h 20m (11.2%)
   - Tasks: 18 (2 NOW, 4 DOING, 12 DONE)
   - Recent activity: 2025-10-28

5. **Other projects** - 55h 55m (16.3%)

---

## Time by Week

### Last 8 Weeks

| Week Starting | Hours | Top Project | Notes |
|---------------|-------|-------------|-------|
| Nov 4, 2025   | 42h 15m | [[Project Phoenix]] (28h) | Sprint 23 |
| Oct 28, 2025  | 38h 30m | [[Project Phoenix]] (22h) | Authentication push |
| Oct 21, 2025  | 35h 45m | [[Mobile App]] (20h) | Beta prep |
| Oct 14, 2025  | 31h 20m | [[Hearth Insights]] (18h) | Planning phase |
| Oct 7, 2025   | 28h 15m | [[Unitary]] (15h) | Research |
| Sep 30, 2025  | 25h 40m | [[Mobile App]] (12h) | UI work |
| Sep 23, 2025  | 22h 30m | [[Project Phoenix]] (14h) | Initial setup |
| Sep 16, 2025  | 18h 45m | Various | Exploration |

**Trend**: Averaging 32h/week over last 8 weeks

---

## Time by Priority

- **High Priority [#A]**: 145h 30m (42.4%)
- **Medium Priority [#B]**: 78h 20m (22.9%)
- **Low Priority [#C]**: 28h 15m (8.2%)
- **No Priority**: 90h 40m (26.5%)

**Insight**: High priority tasks getting good attention, but 26.5% of time on untagged tasks.

---

## Time by Status

- **DONE** (completed): 198h 30m (57.9%)
- **NOW** (active): 82h 15m (24.0%)
- **DOING** (in progress): 45h 30m (13.3%)
- **LATER** (backlog): 12h 20m (3.6%)
- **TODO** (planned): 4h 10m (1.2%)

**Completion Rate**: 57.9% of logged time is on completed tasks ✅

---

## Recent Activity (Last 30 Days)

- Total time: 142h 30m
- Most productive day: Nov 1 (8h 45m)
- Average per day: 4h 45m
- Active days: 24 out of 30

**Top focus areas**:
1. [[Project Phoenix]] - 58h (40.7%)
2. [[Mobile App]] - 42h (29.5%)
3. [[Sprint 23]] - 28h (19.6%)

---

## Insights & Recommendations

- ✅ Good time distribution across active projects
- ⚠️  Consider prioritizing [[Unitary]] (only 11.2% of time, but 14 references)
- 💡 High completion rate (57.9%) indicates good follow-through
- 📊 Averaging 32h/week - consistent productivity
```

**Benefits for Claude**:
- "How much time on Project Phoenix?" → Check "Time by Project" section
- "What's my time allocation?" → See percentage breakdown
- "Am I completing tasks?" → Check completion rate (57.9%)
- "What did I focus on last week?" → Check "Time by Week" table

---

### 5. Enhanced Task Index

**Goal**: Add priority sections and statistics to existing task-index.md

**Changes to** `internal/indexer/task_indexer.go`:

```go
type TaskIndex struct {
    GeneratedAt   time.Time
    TotalTasks    int
    ByStatus      map[TaskStatus][]Task
    ByPriority    map[Priority][]Task      // NEW
    ByProject     map[string][]Task
    Recent        []Task
    Statistics    TaskStatistics           // NEW
}

type TaskStatistics struct {
    CompletionRate float64 // DONE / Total
    PriorityBreakdown map[Priority]int
    StatusBreakdown   map[TaskStatus]int
    WithTimeTracking  int
    TotalTimeLogged   time.Duration
}
```

**Enhanced Output** (`task-index.md`):

Add at the top:
```markdown
# Logseq Task Index

Generated: 2025-11-06T16:00:00Z
Total Tasks: 324

---

## Summary Statistics

- **Completion Rate**: 27.5% (89 of 324 tasks completed)
- **High Priority Tasks**: 48 ([#A])
- **Active Tasks**: 145 (NOW + DOING)
- **Time Tracked**: 342h 45m across 156 tasks

### By Priority
- [#A] High: 48 tasks (14.8%)
- [#B] Medium: 32 tasks (9.9%)
- [#C] Low: 12 tasks (3.7%)
- No priority: 232 tasks (71.6%)

### By Status
- NOW: 48 (14.8%)
- DOING: 42 (13.0%)
- TODO: 89 (27.5%)
- LATER: 56 (17.3%)
- DONE: 89 (27.5%)

---

## 🔥 High Priority Tasks

### NOW [#A] Tasks (12)

(All high priority NOW tasks here...)

### DOING [#A] Tasks (8)

(All high priority DOING tasks here...)

### TODO [#A] Tasks (18)

(All high priority TODO tasks here...)

---

## NOW Tasks (48)

(Rest of NOW tasks, organized as before...)
```

---

## Implementation Plan

### Phase 1: Priority Extraction (2-3 hours)
- [ ] Update Task model with Priority field
- [ ] Enhance task parser to extract `[#A]`, `[#B]`, `[#C]`
- [ ] Add tests for priority extraction
- [ ] Update TaskIndexer to group by priority
- [ ] Enhance TaskWriter to show priority sections

**Files to modify**:
- `pkg/models/task.go` - Add Priority field
- `internal/parser/tasks.go` - Extract priority markers
- `internal/parser/parser_test.go` - Test priority extraction
- `internal/indexer/task_indexer.go` - Add ByPriority map, Statistics
- `internal/writer/task_writer.go` - Write priority sections

### Phase 2: Timeline Index (3-4 hours)
- [ ] Create TimelineIndexer
- [ ] Create TimelineWriter (generates both recent + full)
- [ ] Add tests for timeline grouping
- [ ] Integrate into generate command

**New files**:
- `internal/indexer/timeline_indexer.go`
- `internal/indexer/timeline_indexer_test.go`
- `internal/writer/timeline_writer.go` (writes both timeline-recent.md and timeline-full.md)
- `internal/writer/timeline_writer_test.go`

### Phase 3: Missing Pages Report (2-3 hours)
- [ ] Create MissingPagesIndexer (uses existing ReferenceGraph)
- [ ] Create MissingPagesWriter
- [ ] Add page type classification heuristics
- [ ] Add tests

**New files**:
- `internal/indexer/missing_pages_indexer.go`
- `internal/indexer/missing_pages_indexer_test.go`
- `internal/writer/missing_pages_writer.go`
- `internal/writer/missing_pages_writer_test.go`

### Phase 4: Time Tracking Report (3-4 hours)
- [ ] Create TimeTrackingIndexer
- [ ] Create TimeTrackingWriter
- [ ] Implement aggregation by project/week/priority/status
- [ ] Add trend calculations
- [ ] Add tests

**New files**:
- `internal/indexer/time_tracking_indexer.go`
- `internal/indexer/time_tracking_indexer_test.go`
- `internal/writer/time_tracking_writer.go`
- `internal/writer/time_tracking_writer_test.go`

### Phase 5: Dashboard Generator (2-3 hours)
- [ ] Create DashboardWriter (aggregates all indexes)
- [ ] Pull statistics from all indexers
- [ ] Format dashboard sections
- [ ] Add tests

**New files**:
- `internal/writer/dashboard_writer.go`
- `internal/writer/dashboard_writer_test.go`

### Phase 6: Integration & Testing (2 hours)
- [ ] Update main.go to generate all indexes
- [ ] Update README with new indexes
- [ ] Test with all three datasets (fixtures, synthetic, user-provided)
- [ ] Update TEST_RESULTS.md
- [ ] Performance testing (should still be <3s)

**Total estimated time**: 14-19 hours (added dashboard phase)

---

## Output Files

After implementation, `.claude/indexes/` will contain:

```
.claude/indexes/
├── README.md                  # Documentation of index system
├── dashboard.md               # 🎯 START HERE - overview (NEW)
├── tasks-by-status.md         # All tasks by NOW/LATER/TODO/DOING/DONE (enhanced)
├── tasks-by-priority.md       # High priority [#A] tasks only (NEW)
├── timeline-recent.md         # Last 7 days + 30 day summary (NEW)
├── timeline-full.md           # Complete chronological detail (NEW)
├── reference-graph.md         # Page connection network (existing)
├── missing-pages.md           # Pages with 5+ refs to create (NEW)
└── time-tracking.md           # Time allocation analytics (NEW)
```

**Reading order for Claude**:
1. Start with `dashboard.md` - quick overview
2. Then dive into specific files as needed

**File naming rationale**:
- Descriptive names (self-documenting)
- Paired files for different detail levels (recent + full)
- Consistent format (noun-by-attribute.md)

---

## Performance Considerations

**Target**: Maintain <3s for 116 files

**Optimizations**:
1. **Timeline**: Group during single pass of journals
2. **Missing Pages**: Use existing ReferenceGraph (no extra parsing)
3. **Time Tracking**: Aggregate during task indexing
4. **Memory**: Stream processing, no data duplication

**Expected impact**: +0.5s total (well under 3s target)

---

## Configuration (Future)

Could add `.claude/config.yml`:

```yaml
indexes:
  task-index: true
  reference-graph: true
  timeline-index: true
  missing-pages-report: true
  time-tracking-report: true

timeline:
  recent_days: 7
  summary_days: 30

missing-pages:
  min_references: 1
  max_results: 50

time-tracking:
  recent_days: 30
  weeks_to_show: 8
```

But for MVP, use sensible defaults.

---

## Testing Strategy

### Unit Tests
- Priority extraction from various formats
- Timeline grouping by date
- Missing page classification
- Time aggregation accuracy
- Week calculation logic

### Integration Tests
- Generate all indexes on synthetic data
- Verify output format
- Check cross-references between indexes
- Performance testing

### Validation
- Run on user-provided data (116 files)
- Verify statistics accuracy
- Check timeline chronology
- Validate time tracking totals

---

## Benefits for Claude Code

With these enhanced indexes, Claude can answer:

| Question | Index to Check | Time Saved |
|----------|----------------|------------|
| "What happened yesterday?" | timeline-index.md → Nov 5 section | Instant vs manual journal reading |
| "Show urgent tasks" | task-index.md → High Priority [#A] section | Instant vs scanning 324 tasks |
| "What pages should I create?" | missing-pages-report.md → High Priority | Instant vs reference analysis |
| "Time on Project Phoenix?" | time-tracking-report.md → By Project | Instant vs LOGBOOK aggregation |
| "Am I completing tasks?" | task-index.md → Statistics → Completion Rate | Instant vs manual counting |
| "Who are my key contacts?" | reference-graph.md → Hub Pages (people) | Already instant ✅ |
| "What's my time allocation?" | time-tracking-report.md → By Project % | Instant vs complex calculation |
| "Most productive week?" | time-tracking-report.md → By Week | Instant vs date arithmetic |

**Estimated time savings**: 5-10 minutes per query → 1-2 seconds

---

## Success Criteria

✅ **Functional**:
- [ ] All 5 indexes generate correctly
- [ ] Priority extraction works for all variants
- [ ] Timeline shows last 7 days + weekly summaries
- [ ] Missing pages sorted by reference count
- [ ] Time tracking totals match LOGBOOK data
- [ ] Statistics accurate (completion rate, distributions)

✅ **Performance**:
- [ ] <3s for 116 files (all indexes)
- [ ] <5s for 500 files (scalability test)

✅ **Quality**:
- [ ] 80%+ test coverage maintained
- [ ] All tests pass
- [ ] Documentation updated
- [ ] Works on all test datasets

✅ **Usability**:
- [ ] Claude can answer all 7 questions from indexes
- [ ] Output format readable and actionable
- [ ] Cross-references between indexes work

---

## Future Enhancements (v3)

Once this is stable, could add:

1. **Query-specific indexes**:
   - "People I haven't contacted in 30 days"
   - "Stale projects" (no activity in X weeks)
   - "Overdue tasks" (referenced dates in past)

2. **Visualizations**:
   - Mermaid gantt charts in timeline
   - Burndown charts from time tracking
   - Network graphs for references

3. **Smart suggestions**:
   - "You spend 40% time on X, consider delegating"
   - "These 5 tasks are blocking 12 others"
   - "Create [[Person Name]] page - referenced 10 times"

4. **Trends**:
   - Week-over-week velocity
   - Project momentum (increasing/decreasing time)
   - Completion rate trends

---

## ✅ Design Decisions (Confirmed)

1. **Priority format**: `[#A]`, `[#B]`, `[#C]` only ✅
2. **Timeline depth**: Two files:
   - `timeline-recent.md`: Last 7 days detailed + 30 days summary ✅
   - `timeline-full.md`: Complete chronological detail ✅
3. **Missing pages threshold**: Only show pages with 5+ references ✅
4. **Time tracking**: Show all tasks in totals, flag tracking adoption % ✅
   - Example: "324 total tasks, 156 with time tracking (48.1%), 342h 45m logged"
5. **Directory structure**: Flat structure in `.claude/indexes/` ✅
   - Added: `dashboard.md` as the "home page" / quick overview
   - Use descriptive names: `tasks-by-status.md` not `task-index.md`
6. **Index files**: ✅
   - `dashboard.md` (NEW - overview of everything)
   - `tasks-by-status.md` (was task-index.md)
   - `tasks-by-priority.md` (NEW - [#A] tasks only)
   - `timeline-recent.md` (last 7 days + summary)
   - `timeline-full.md` (complete detail)
   - `reference-graph.md` (unchanged)
   - `missing-pages.md` (5+ refs only)
   - `time-tracking.md` (analytics)

---

## Conclusion

This design enhances `logseq-claude-indexer` to be a comprehensive knowledge base analytics tool. By generating 5 complementary indexes, it enables Claude Code to answer natural questions instantly without external context.

**Key Innovation**: From "file indexer" to "knowledge intelligence system"

**Next Step**: Review this design, confirm decisions, then implement in phases.
