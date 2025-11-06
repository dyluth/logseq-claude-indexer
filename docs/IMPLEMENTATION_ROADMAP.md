# Implementation Roadmap - Enhanced Indexes

## Overview

Transform `logseq-claude-indexer` from a file indexer into a **knowledge intelligence system** that answers natural questions.

**Time Estimate**: 14-19 hours
**Phases**: 6
**New Files**: 9 indexes (1 enhanced, 5 new, 1 dashboard, 2 supporting)

---

## What We're Building

### Before (Current - v0.1.0)
```
.claude/indexes/
├── task-index.md           # All tasks by status
└── reference-graph.md      # Page connections
```

**Claude's capabilities**:
- ✅ "Show me all tasks"
- ✅ "What pages are connected?"
- ❌ "What happened yesterday?"
- ❌ "What's urgent?"
- ❌ "Where's my time going?"

### After (Enhanced - v0.2.0)
```
.claude/indexes/
├── README.md                  # System documentation
├── dashboard.md               # 🎯 START HERE - overview
├── tasks-by-status.md         # All tasks (enhanced with stats)
├── tasks-by-priority.md       # High priority [#A] only
├── timeline-recent.md         # Last 7 days + summary
├── timeline-full.md           # Complete chronology
├── reference-graph.md         # Page connections (existing)
├── missing-pages.md           # Pages to create (5+ refs)
└── time-tracking.md           # Time analytics
```

**Claude's new capabilities**:
- ✅ "What happened yesterday?" → `timeline-recent.md`
- ✅ "What's urgent?" → `tasks-by-priority.md`
- ✅ "Where's my time going?" → `time-tracking.md`
- ✅ "What should I create?" → `missing-pages.md`
- ✅ "Am I productive?" → `dashboard.md` → 57.9% completion rate
- ✅ "Quick status check" → `dashboard.md` → instant overview

---

## The 7 Questions We're Answering

| # | Question | Index File | Status |
|---|----------|------------|--------|
| 1 | "What happened last week?" | `timeline-recent.md` | NEW |
| 2 | "What are my key contacts?" | `reference-graph.md` | ✅ Existing |
| 3 | "What's completed vs pending?" | `tasks-by-status.md` | Enhanced |
| 4 | "How much time invested in X?" | `time-tracking.md` | NEW |
| 5 | "Which pages should I create?" | `missing-pages.md` | NEW |
| 6 | "What tasks are urgent?" | `tasks-by-priority.md` | NEW |
| 7 | "What's my time allocation?" | `time-tracking.md` | NEW |

---

## Phase Breakdown

### Phase 1: Priority Extraction (2-3 hours)

**Goal**: Extract `[#A]`, `[#B]`, `[#C]` markers from tasks

**Changes**:
- Add `Priority` field to `Task` model
- Parse priority markers in task parser
- Group tasks by priority in indexer
- Generate priority sections in task writer

**Files Modified**:
- `pkg/models/task.go` - Add Priority field
- `internal/parser/tasks.go` - Extract `[#A]`, `[#B]`, `[#C]`
- `internal/parser/parser_test.go` - Test priority extraction
- `internal/indexer/task_indexer.go` - Add ByPriority map, Statistics struct
- `internal/writer/task_writer.go` - Write priority sections + statistics

**Deliverables**:
- `tasks-by-status.md` (enhanced with priority stats)
- `tasks-by-priority.md` (NEW - [#A] tasks only)

**Tests**:
```go
TestExtractPriority("NOW [#A] Task") → Priority = "A"
TestTaskIndexByPriority() → ByPriority["A"] contains high-pri tasks
TestPriorityTaskWriter() → Generates tasks-by-priority.md
```

---

### Phase 2: Timeline Index (3-4 hours)

**Goal**: Generate chronological view of journal activity

**New Components**:
- `TimelineIndexer` - Groups tasks/refs by date
- `TimelineWriter` - Generates two files:
  - `timeline-recent.md` - Last 7 days detailed + 30 day summary
  - `timeline-full.md` - Complete chronological history

**Files Created**:
- `internal/indexer/timeline_indexer.go`
- `internal/indexer/timeline_indexer_test.go`
- `internal/writer/timeline_writer.go`
- `internal/writer/timeline_writer_test.go`

**Data Structures**:
```go
type TimelineIndex struct {
    GeneratedAt time.Time
    Entries     []TimelineDay // Sorted newest first
}

type TimelineDay struct {
    Date           time.Time
    JournalPath    string
    TasksCreated   []models.Task
    TimeLogged     time.Duration
    KeyActivity    []string // Summary bullets
}
```

**Deliverables**:
- `timeline-recent.md` - Quick recent view
- `timeline-full.md` - Complete history

**Tests**:
```go
TestTimelineGrouping() → Tasks grouped by journal date
TestTimelineSorting() → Newest first
TestTimelineRecent() → Last 7 days + 30 day summary
TestTimelineFull() → All days included
```

---

### Phase 3: Missing Pages Report (2-3 hours)

**Goal**: Identify high-value pages to create (5+ references)

**New Components**:
- `MissingPagesIndexer` - Uses existing ReferenceGraph
- `MissingPagesWriter` - Generates report with page type classification

**Files Created**:
- `internal/indexer/missing_pages_indexer.go`
- `internal/indexer/missing_pages_indexer_test.go`
- `internal/writer/missing_pages_writer.go`
- `internal/writer/missing_pages_writer_test.go`

**Page Type Heuristics**:
```go
func classifyPageType(pageName string) string {
    if strings.Contains(pageName, " - ") {
        return "person" // "Sarah Chen - Tech Lead"
    }
    if strings.HasSuffix(pageName, "th, 2025") {
        return "date" // "Nov 24th, 2025"
    }
    if containsKeywords(pageName, ["sprint", "project"]) {
        return "project"
    }
    return "concept"
}
```

**Deliverable**:
- `missing-pages.md` - High-priority missing pages (5+ refs)

**Tests**:
```go
TestMissingPagesThreshold() → Only 5+ refs shown
TestPageTypeClassification() → Correct type detection
TestMissingPagesSorting() → By reference count desc
```

---

### Phase 4: Time Tracking Report (3-4 hours)

**Goal**: Aggregate LOGBOOK data for time analytics

**New Components**:
- `TimeTrackingIndexer` - Aggregates by project/week/priority/status
- `TimeTrackingWriter` - Generates analytical report

**Files Created**:
- `internal/indexer/time_tracking_indexer.go`
- `internal/indexer/time_tracking_indexer_test.go`
- `internal/writer/time_tracking_writer.go`
- `internal/writer/time_tracking_writer_test.go`

**Aggregations**:
```go
type TimeTrackingIndex struct {
    TotalTimeLogged  time.Duration
    ByProject        map[string]time.Duration
    ByWeek           map[string]time.Duration
    ByPriority       map[Priority]time.Duration
    ByStatus         map[TaskStatus]time.Duration
    TopProjects      []ProjectTime
    Statistics       TimeStatistics
}

type TimeStatistics struct {
    TotalTasks       int
    TasksWithTracking int
    TrackingAdoption float64 // TasksWithTracking / TotalTasks
    CompletionRate   float64
}
```

**Key Insight**: Show tracking adoption
- "324 total tasks, 156 with time tracking (48.1%), 342h 45m logged"

**Deliverable**:
- `time-tracking.md` - Time allocation analytics

**Tests**:
```go
TestTimeAggregation() → Correct totals
TestWeekCalculation() → Monday week starts
TestProjectAggregation() → By PageRefs[0]
TestTrackingAdoption() → 156/324 = 48.1%
```

---

### Phase 5: Dashboard Generator (2-3 hours)

**Goal**: Create "home page" that summarizes everything

**New Component**:
- `DashboardWriter` - Pulls statistics from all indexers

**Files Created**:
- `internal/writer/dashboard_writer.go`
- `internal/writer/dashboard_writer_test.go`

**Dashboard Sections**:
1. At a Glance (metrics table)
2. What Needs Attention (urgent tasks)
3. Recent Activity (this week)
4. Time Allocation (project breakdown)
5. Knowledge Graph (hub pages)
6. Missing Pages (top gaps)
7. Trends & Insights (analytics)
8. Index Files Reference (navigation)
9. Quick Answers (question mapping)

**Data Flow**:
```
TaskIndex → Completion rate, urgent tasks
TimelineIndex → Recent activity
TimeTrackingIndex → Time allocation, trends
ReferenceGraph → Hub pages
MissingPages → Top missing pages
    ↓
DashboardWriter → Aggregates & formats
    ↓
dashboard.md
```

**Deliverable**:
- `dashboard.md` - The "start here" file

**Tests**:
```go
TestDashboardGeneration() → All sections present
TestDashboardStatistics() → Correct calculations
TestDashboardCrossReferences() → Links to other indexes
```

---

### Phase 6: Integration & Testing (2 hours)

**Goal**: Wire everything together and validate

**Tasks**:
- [ ] Update `cmd/logseq-claude-indexer/main.go`:
  - Build all new indexes
  - Generate dashboard last (needs all other indexes)
  - Add timing/progress output
- [ ] Create `README.md` in `.claude/indexes/`:
  - Explain each index file
  - Show reading order
  - Document file structure
- [ ] Update main `README.md` with new features
- [ ] Test on all datasets:
  - Fixtures (6 files) - basic validation
  - Synthetic (11 files) - feature coverage
  - User-provided (116 files) - production validation
- [ ] Update `testdata/TEST_RESULTS.md`
- [ ] Performance testing: <3s for 116 files
- [ ] Integration test: Verify all indexes cross-reference correctly

**Validation Checklist**:
- [ ] All 9 files generated
- [ ] Dashboard links to all indexes
- [ ] Priority extraction works
- [ ] Timeline chronology correct
- [ ] Missing pages threshold = 5
- [ ] Time tracking adoption shown
- [ ] Performance <3s
- [ ] Tests pass: 80%+ coverage

---

## Expected Results (User-Provided Data - 116 files)

### Current (v0.1.0)
```
.claude/indexes/
├── task-index.md          # 324 tasks
└── reference-graph.md     # 226 pages, 295 references
```

### After Enhancement (v0.2.0)
```
.claude/indexes/
├── README.md
├── dashboard.md           # Overview of all 324 tasks, 342h time
├── tasks-by-status.md     # 324 tasks + statistics
├── tasks-by-priority.md   # 48 high-priority tasks
├── timeline-recent.md     # Last 7 days activity
├── timeline-full.md       # 96 journal entries chronologically
├── reference-graph.md     # 226 pages (unchanged)
├── missing-pages.md       # 8 pages with 5+ refs
└── time-tracking.md       # 342h across 12 projects
```

**Dashboard Preview**:
```markdown
# Knowledge Base Dashboard

📊 **At a Glance**
- Tasks: 324 (89 completed, 48 urgent)
- Time: 342h 45m (156 tasks tracked)
- Projects: 12 active
- Missing: 8 pages to create

🔥 **What Needs Attention**
- 12 NOW [#A] tasks
- Business Plan due Nov 24th (18 days)

⏱️ **Time Allocation**
- Project Phoenix: 37.2%
- Mobile App: 19.9%
- Hearth Insights: 15.4%
```

---

## Performance Target

**Current**: ~1.5s for 116 files (2 indexes)
**Target**: <3s for 116 files (9 indexes)
**Expected**: ~2.5s (within target ✅)

**Breakdown**:
- Parsing: 1.0s (unchanged)
- Indexing: 0.8s (+0.5s for new indexes)
- Writing: 0.7s (+0.5s for new files)

**Optimizations**:
- Timeline: Group during single pass
- Missing Pages: Use existing graph
- Time Tracking: Aggregate during indexing
- Dashboard: Pull from existing indexes

---

## Testing Strategy

### Unit Tests (40+ new tests)
- Priority extraction: 5 tests
- Timeline grouping: 8 tests
- Missing pages: 5 tests
- Time tracking: 10 tests
- Dashboard: 8 tests
- Integration: 4 tests

### Integration Tests
- Generate all indexes on synthetic data
- Verify cross-references between indexes
- Check dashboard aggregation accuracy
- Validate file structure

### Validation
- Run on user-provided data (116 files)
- Verify all statistics match manual counts
- Check timeline chronology
- Validate time tracking totals against LOGBOOK

---

## Success Criteria

✅ **Functional**:
- [ ] All 9 index files generate
- [ ] Priority extraction: [#A], [#B], [#C]
- [ ] Timeline: Last 7 days + summaries
- [ ] Missing pages: Only 5+ refs
- [ ] Time tracking: Shows adoption %
- [ ] Dashboard: Aggregates all data
- [ ] Statistics: Completion rate, distributions

✅ **Performance**:
- [ ] <3s for 116 files
- [ ] <5s for 500 files (scalability)

✅ **Quality**:
- [ ] 80%+ test coverage maintained
- [ ] All tests pass
- [ ] Documentation updated

✅ **Usability**:
- [ ] Claude can answer all 7 questions
- [ ] Dashboard provides instant overview
- [ ] Cross-references work

---

## Rollout Plan

### Week 1: Core Functionality
- Day 1-2: Phase 1 (Priority extraction)
- Day 3-4: Phase 2 (Timeline index)
- Day 5: Testing & validation

### Week 2: Analytics & Dashboard
- Day 1-2: Phase 3 (Missing pages)
- Day 3-4: Phase 4 (Time tracking)
- Day 5: Phase 5 (Dashboard)

### Week 3: Integration
- Day 1-2: Phase 6 (Integration & testing)
- Day 3: Documentation
- Day 4-5: User testing & refinement

---

## Risk Mitigation

### Risk 1: Performance Degradation
- **Mitigation**: Profile each phase, optimize hot paths
- **Fallback**: Make new indexes optional via flags

### Risk 2: Complex Timeline Logic
- **Mitigation**: Start simple, iterate based on feedback
- **Fallback**: Reduce detail levels if needed

### Risk 3: Dashboard Complexity
- **Mitigation**: Build incrementally, test each section
- **Fallback**: Simplify dashboard to just key metrics

---

## Post-Implementation

### Documentation Updates
- [ ] README.md - New features section
- [ ] CHANGELOG.md - v0.2.0 release notes
- [ ] User guide - How to use new indexes
- [ ] Examples - Dashboard screenshots

### Future Enhancements (v0.3.0+)
- [ ] Configurable thresholds (5+ refs → user choice)
- [ ] Alerts (overdue tasks, stalled projects)
- [ ] Trends (week-over-week velocity)
- [ ] Smart suggestions (create page X, prioritize Y)

---

## Questions Before Starting

1. ✅ Priority format confirmed: [#A], [#B], [#C]
2. ✅ Timeline: recent + full files
3. ✅ Missing pages: 5+ refs threshold
4. ✅ Time tracking: Show tracking adoption %
5. ✅ Structure: Flat `.claude/indexes/` with dashboard
6. ✅ File names confirmed

**All decisions locked in. Ready to implement.** 🚀

---

## Next Action

**Start Phase 1**: Priority Extraction
- Create branch: `feature/enhanced-indexes`
- Begin with: `pkg/models/task.go` - Add Priority field
- TDD: Write tests first, then implement
- Commit often: Small, focused commits

**Command to start**:
```bash
git checkout -b feature/enhanced-indexes
# Begin Phase 1...
```
