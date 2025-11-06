# Dashboard Design

## Purpose

The `dashboard.md` file is the **single source of truth** for a quick overview of your entire knowledge base. Claude Code should check this first before diving into specific index files.

---

## Dashboard Structure

```markdown
# Knowledge Base Dashboard

**Last Updated**: 2025-11-06 16:45:00
**Repository**: logseq-data-git
**Scan Coverage**: 116 files (96 journals + 20 pages)

---

## 📊 At a Glance

| Metric | Value | Trend |
|--------|-------|-------|
| Total Tasks | 324 | +12 this week |
| Completion Rate | 27.5% (89 completed) | ↑ from 24% |
| Time Logged | 342h 45m | +42h this week |
| Active Projects | 12 | 3 highly active |
| Knowledge Pages | 226 | +3 this week |
| Missing Pages | 8 (5+ refs) | Priority: create 3 |

---

## 🔥 What Needs Attention Right Now

### High Priority Tasks [#A]
- **12 NOW tasks** - Immediate action required
  - [[Project Phoenix]] - Authentication module (3h logged today)
  - [[Mobile App]] - Beta testing deadline approaching
  - [[Business Plan]] - Due Nov 24th, 2025
- **8 DOING tasks** - In progress, need completion
  - [[User Analytics]] - Dashboard integration
  - [[Sprint 23]] - Rate limiting implementation

👉 **Action**: See `tasks-by-priority.md` for full list

### Overdue / Dated Tasks
- 3 tasks reference past dates
- [[Nov 24th, 2025]] deadline approaching (18 days)

---

## 📅 This Week's Activity

### Today (Nov 6)
- ✅ 8 tasks created (3 NOW, 5 LATER)
- ⏱️ 6h 30m logged
- 🎯 Focus: [[Project Phoenix]], [[Mobile App]]
- 📝 Journal: `journals/2025_11_06.md`

### Yesterday (Nov 5)
- ✅ 5 tasks created
- ⏱️ 4h 15m logged
- 🎯 Focus: [[Sprint 23]] planning

### Last 7 Days Summary
- Tasks: 45 created, 32 completed
- Time: 38h 15m logged
- Most productive: Nov 4 (12 tasks, 8h 45m)
- Focus: [[Project Phoenix]] (18h), [[Mobile App]] (12h)

👉 **Action**: See `timeline-recent.md` for daily breakdown

---

## ⏱️ Time Allocation (Last 30 Days)

### Top Projects by Time
1. **[[Project Phoenix]]** - 58h (40.7%) ⭐ Highest focus
2. **[[Mobile App]]** - 42h (29.5%) 🚀 Active development
3. **[[Sprint 23]]** - 28h (19.6%) 📋 Planning
4. **[[Hearth Insights]]** - 14h (9.8%)

### By Priority
- [#A] High: 42.4% of time (good prioritization ✅)
- [#B] Medium: 22.9%
- [#C] Low: 8.2%
- No priority: 26.5% ⚠️ Consider tagging

### Completion Status
- DONE: 57.9% of logged time
- NOW: 24.0%
- DOING: 13.3%
- Backlog (LATER/TODO): 4.8%

**Insight**: High completion rate (57.9%) indicates good follow-through ✅

👉 **Action**: See `time-tracking.md` for detailed breakdown

---

## 🔗 Knowledge Graph Insights

### Hub Pages (Most Connected)
1. **[[Unitary]]** - 14 references (top hub)
2. **[[Hearth Insights]]** - 12 references
3. **[[Matt Stammers - SETT Lead]]** - 11 references (person)
4. **[[Adrian Braine - mentor]]** - 10 references (person)
5. **[[chris kipps - wessex data]]** - 8 references (person)

### Key Contacts
- 12 people pages with 5+ references
- Top mentors: Adrian Braine, Gemma Snell
- Top collaborators: Matt Stammers, chris kipps

👉 **Action**: See `reference-graph.md` for full network

---

## ⚠️ Missing Pages (Should Create)

High-impact pages referenced 5+ times but don't exist yet:

1. **[[Matt Stammers - SETT: Theme Lead for Data & AI]]** - 11 refs
   - Type: Person (contact)
   - Create: Add contact info, expertise, meeting notes

2. **[[catalyst - marketing stream]]** - 7 refs
   - Type: Project/Program
   - Create: Document what this is, how it relates to work

3. **[[Observability Stack]]** - 8 refs (from synthetic data)
   - Type: Technical infrastructure
   - Create: Architecture, tools, timeline

**Total missing**: 8 pages with 5+ references

👉 **Action**: See `missing-pages.md` for full list and suggestions

---

## 📈 Trends & Insights

### Productivity Trends
- ✅ Consistent 32h/week average (last 8 weeks)
- ✅ 57.9% completion rate (healthy)
- ⚠️ 48% of tasks have time tracking (could improve)
- ✅ Strong prioritization (42% time on [#A] tasks)

### Project Momentum
- 🚀 [[Project Phoenix]]: Accelerating (18h this week vs 14h last week)
- 📈 [[Mobile App]]: Steady progress (12h/week average)
- 📉 [[Hearth Insights]]: Slowing (14h this month vs 30h last month)

### Recommendations
1. **Create missing people pages** - 12 contacts need pages
2. **Tag untagged tasks** - 168 tasks (52%) have no priority marker
3. **Review [[Hearth Insights]]** - Momentum down, needs attention?
4. **Maintain current pace** - 32h/week sustainable

---

## 🗂️ Index Files Reference

| File | Purpose | When to Use |
|------|---------|-------------|
| `tasks-by-status.md` | All 324 tasks grouped by NOW/LATER/TODO/DOING/DONE | "Show me all tasks" |
| `tasks-by-priority.md` | High priority [#A] tasks only (48 tasks) | "What's urgent?" |
| `timeline-recent.md` | Last 7 days detailed + 30 day summary | "What happened this week?" |
| `timeline-full.md` | Complete chronological history | "Deep dive into timeline" |
| `reference-graph.md` | Full page connection network (226 pages) | "How are pages connected?" |
| `missing-pages.md` | 8 pages with 5+ refs that don't exist | "What pages should I create?" |
| `time-tracking.md` | Time allocation analytics | "Where's my time going?" |

---

## 💡 Quick Answers

**"What should I work on now?"**
→ Check "🔥 What Needs Attention" above → 12 NOW [#A] tasks

**"What did I do yesterday?"**
→ Check "📅 This Week's Activity" → Yesterday section

**"How's Project Phoenix going?"**
→ Check "⏱️ Time Allocation" → 58h this month, trending up 🚀

**"Who should I talk to about X?"**
→ Check "🔗 Knowledge Graph" → Hub pages show key contacts

**"Am I being productive?"**
→ Check "📈 Trends & Insights" → 57.9% completion rate ✅

**"What pages am I missing?"**
→ Check "⚠️ Missing Pages" → 8 high-priority pages to create

---

**📖 Start here, then explore specific index files for details.**
```

---

## Design Rationale

### Why Dashboard First?

1. **Cognitive Load**: Don't make Claude read 8 files to answer "how's it going?"
2. **Actionability**: Highlights what needs attention NOW
3. **Context**: Provides overview before details
4. **Speed**: Single file check vs multiple file reads
5. **Trends**: Shows movement over time, not just snapshots

### Information Hierarchy

```
Dashboard (1 minute read)
    ↓
Quick answer found → Done
    ↓
Need more detail → Check specific index
    ↓
tasks-by-priority.md → Full list of urgent items
timeline-recent.md → Detailed daily activity
time-tracking.md → Full time analytics
etc.
```

### Dashboard Sections Explained

1. **At a Glance** - Metrics table for instant status
2. **What Needs Attention** - Action items (urgent tasks, deadlines)
3. **This Week's Activity** - Temporal context (what's happening now)
4. **Time Allocation** - Resource distribution (where effort goes)
5. **Knowledge Graph** - Connection insights (what's important)
6. **Missing Pages** - Gaps to fill (improve knowledge base)
7. **Trends & Insights** - Analysis (am I productive? what's changing?)
8. **Index Files Reference** - Navigation (where to find more)
9. **Quick Answers** - Common questions mapped to sections

---

## Implementation Details

### Dashboard Generator

New file: `internal/writer/dashboard_writer.go`

```go
type DashboardWriter struct {
    TaskIndex     *indexer.TaskIndex
    GraphIndex    *indexer.ReferenceGraph
    TimelineIndex *indexer.TimelineIndex
    TimeIndex     *indexer.TimeTrackingIndex
    MissingPages  *indexer.MissingPagesIndex
}

func WriteDashboard(data *DashboardWriter, outputDir string) error {
    // Aggregate statistics from all indexes
    // Format into dashboard sections
    // Write to dashboard.md
}
```

### Data Sources

Dashboard pulls from:
- TaskIndex → Completion rate, priority breakdown, urgent tasks
- TimelineIndex → Recent activity, weekly summary
- TimeTrackingIndex → Time allocation, trends
- ReferenceGraph → Hub pages, connection insights
- MissingPagesIndex → High-priority missing pages

### Update Frequency

Dashboard is regenerated every time `logseq-claude-indexer generate` runs:
- Post-commit hook → Auto-update
- Manual run → Fresh dashboard
- Always shows "Last Updated" timestamp

---

## Example Usage by Claude

**User**: "What should I focus on today?"

**Claude** (internally):
1. Reads `.claude/indexes/dashboard.md`
2. Sees "🔥 What Needs Attention" section
3. Finds 12 NOW [#A] tasks
4. Top item: "[[Project Phoenix]] - Authentication module"

**Claude** (response):
"Focus on [[Project Phoenix]] authentication module - it's marked as NOW [#A] priority and you've already logged 3h on it today. You have 12 high-priority NOW tasks total. See `.claude/indexes/tasks-by-priority.md` for the complete list."

---

## Dashboard vs Other Indexes

| File | Scope | Depth | Use Case |
|------|-------|-------|----------|
| `dashboard.md` | Everything | Summary | Quick overview, what's important |
| `tasks-by-status.md` | Tasks only | Complete | Find specific task |
| `timeline-recent.md` | Recent time | Detailed | Daily activity review |
| `time-tracking.md` | Time data | Analytical | Understand time allocation |
| `reference-graph.md` | Connections | Complete | Explore knowledge network |

**Dashboard = Your knowledge base's "home page"**

---

## Validation

Test the dashboard with these questions:

✅ "Am I productive?" → Check completion rate (57.9%)
✅ "What's urgent?" → Check High Priority section (12 NOW [#A])
✅ "What happened this week?" → Check This Week's Activity
✅ "Where's my time going?" → Check Time Allocation (40% on Project Phoenix)
✅ "What should I create?" → Check Missing Pages (8 pages)
✅ "How's Project X?" → Check Time Allocation or Trends

If dashboard can answer these 6 questions without reading other files → Success ✅

---

## Future Enhancements

Dashboard v2 could add:

1. **Alerts**:
   ```markdown
   ## ⚠️ Alerts
   - [[Business Plan]] due in 18 days - only 12% complete
   - [[Project Phoenix]] time dropped 40% this week
   - 5 tasks marked DOING for 14+ days (stalled?)
   ```

2. **Predictive insights**:
   ```markdown
   ## 🔮 Predictions
   - At current pace, [[Sprint 23]] will complete on Nov 10 (2 days late)
   - [[Project Phoenix]] needs 40h to finish - 2 weeks at current rate
   ```

3. **Network analysis**:
   ```markdown
   ## 🕸️ Knowledge Clusters
   - Healthcare cluster: 8 pages, 42 references
   - Tech stack cluster: 12 pages, 38 references
   - People network: 15 contacts, 28 interactions
   ```

But for MVP: Keep it simple, focused on answering the 7 key questions.

---

## Conclusion

The dashboard transforms the index system from "files you query" to "intelligence you scan". It's the difference between:

**Before**: "Let me read 8 files to figure out what's happening"
**After**: "Let me check the dashboard - ah, I see everything at a glance"

This makes Claude Code exponentially more useful for knowledge base queries.
