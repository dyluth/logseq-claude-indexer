# Query Patterns & Optimization Guide

**Last Updated:** 2025-11-07
**Purpose:** Decision trees for efficient query routing using the index system

---

## Overview

This guide provides **decision trees** for common queries, showing:
- Which indexes to read (and in what order)
- Token usage for each approach
- Good vs bad query patterns
- Optimization strategies

**Core principle:** Start with smallest index, expand only if needed.

---

## Query Routing Decision Trees

### Status & Overview Queries

#### "What's my status?" / "What's happening?" / "Current state?"

**Optimal path:**
```
1. Read: dashboard.md (500 tokens)
   ├─ Extract: Quick stats, current priorities, top projects
   └─ If sufficient → DONE

2. If more detail needed:
   Read: timeline-recent.md (350 tokens)
   ├─ Extract: Last 7 days activity
   └─ If sufficient → DONE

3. If task-level detail needed:
   Read: tasks-by-status.md (NOW section only, ~500 tokens)
   └─ Extract: Currently active tasks

Total: 500-1,350 tokens
vs. Reading task-index.md: 8,000 tokens
Efficiency: 6-16x faster
```

**Example output:**
```
Status Report:
- 324 total tasks, 60.8% completion rate
- Current priorities: Business plan (NOW), Working demo (TODO)
- Top projects: Data shield (3,529h), WorkJobThing (2,440h)
- Recent activity: 7 NOW tasks yesterday, 1 TODO task today
```

---

#### "What are my priorities?" / "What's urgent?"

**Optimal path:**
```
1. Read: tasks-by-priority.md (400 tokens)
   └─ Get all [#A] and [#B] tasks (7 total)

Alternative:
1. Read: dashboard.md (500 tokens)
   └─ "Current Priorities" section shows top 2

Total: 400-500 tokens
vs. Reading tasks-by-status.md: 2,000 tokens
Efficiency: 4-5x faster
```

**When to use which:**
- tasks-by-priority.md → Complete priority list
- dashboard.md → Quick top 2 priorities only

---

### Time-Based Queries

#### "What happened this week?" / "Recent activity?" / "What did I do?"

**Optimal path:**
```
1. Read: timeline-recent.md (350 tokens)
   ├─ Get last 7 days summary
   └─ If sufficient → DONE

2. If completed task detail needed:
   Read: tasks-by-status.md (DONE section only, ~400 tokens)
   └─ Extract completed tasks

Total: 350-750 tokens
vs. Reading 7 journal files: 2,600 tokens
Efficiency: 3-7x faster
```

**❌ BAD:**
```
Read: journals/2025_11_06.md (500 tokens)
Read: journals/2025_11_05.md (800 tokens)
Read: journals/2025_11_04.md (600 tokens)
[... 4 more files]
Total: 2,600 tokens
```

**✅ GOOD:**
```
Read: timeline-recent.md (350 tokens)
Result: Same information, 7x faster
```

---

#### "What happened in [month/specific date]?"

**Optimal path:**
```
1. Read: timeline-full.md (5,000 tokens)
   ├─ Filter for specific month/date
   └─ Extract activity summary

2. If detailed context needed:
   Read specific journal file: journals/YYYY_MM_DD.md
   └─ Only for days requiring deep context

Total: 5,000-6,000 tokens
vs. Reading all journals for month: 20,000+ tokens
Efficiency: 3-4x faster
```

**Use case:** Historical analysis, "When did I last work on X?"

---

### Task-Specific Queries

#### "Show all [NOW/TODO/DONE] tasks"

**Optimal path:**
```
1. Read: tasks-by-status.md (2,000 tokens)
   └─ Navigate to specific status section

Note: This is already optimal - tasks-by-status.md is the filtered view
```

**Alternative for quick counts:**
```
1. Read: dashboard.md (500 tokens)
   └─ Get task counts by status (no descriptions)
```

---

#### "Find tasks about [keyword/project]"

**Optimal path:**
```
1. If known status:
   Read: tasks-by-status.md (2,000 tokens)
   └─ grep for keyword in specific section

2. If comprehensive search needed:
   Read: task-index.md (8,000 tokens)
   └─ grep for keyword across all tasks

Total: 2,000-8,000 tokens
```

**Optimization tip:** If you know the project's primary status (e.g., "MoonLaunch is mostly NOW"), start with tasks-by-status.md to save 6,000 tokens.

---

### Project-Focused Queries

#### "Deep dive on [project]" / "Full context for [project]"

**Optimal multi-index workflow:**
```
1. dashboard.md (500 tokens)
   └─ Check if project in "Top Projects"

2. tasks-by-status.md (2,000 tokens)
   └─ grep for project name → Find all tasks

3. time-tracking.md (1,500 tokens)
   └─ Project time investment

4. reference-graph.md (5,000 tokens)
   └─ Find connected pages

5. timeline-recent.md (350 tokens)
   └─ Recent project activity

6. Read specific pages identified in step 4

Total: ~3,000 tokens (indexes) + external sources
vs. Reading all journals mentioning project: 10,000+ tokens
Efficiency: 3x faster
```

---

#### "How much time on [project]?"

**Optimal path:**
```
1. Read: time-tracking.md (1,500 tokens)
   └─ Find project in "Top Projects" section

2. If task-level breakdown needed:
   Read: task-index.md (8,000 tokens)
   └─ grep for project → Find tasks with LOGBOOK entries
   └─ Sum time across tasks

Total: 1,500-9,500 tokens
vs. grepping all journals for LOGBOOK: 10,000+ tokens
Efficiency: 1-7x faster
```

**Quick answer path (Level 1 only):**
```
Read: dashboard.md (500 tokens)
└─ "Top Projects" shows time for top 5 projects

If project not in top 5, use time-tracking.md
```

---

### Documentation & Knowledge Queries

#### "What should I document?" / "Knowledge gaps?" / "Missing pages?"

**Optimal path:**
```
1. Read: missing-pages.md (300 tokens)
   └─ Get pages with 5+ refs that don't exist (9 pages)
   └─ Categorized by People vs Concepts
   └─ Sorted by reference count

2. If context about why page is needed:
   Read: reference-graph.md (5,000 tokens)
   └─ See who references the missing page

Total: 300-5,300 tokens
Result: Prioritized list with context
```

**Example output:**
```
High Priority (10+ refs):
- John smith (11 refs) - Person, Lead
  → Create contact page with role, projects, interactions

Medium Priority (7-9 refs):
- alistair bastian (7 refs) - Person, InfoSum contact
- Gemma Snell (7 refs) - Person, mentor
- Data shield (5 refs) - Concept, research project
```

---

#### "What's connected to [page]?" / "Show relationships"

**Optimal path:**
```
1. Read: reference-graph.md (5,000 tokens)
   ├─ Find page in hub list or page details
   ├─ See inbound references (who links to this)
   └─ See outbound references (what this links to)

Total: 5,000 tokens
Result: Complete network view of page relationships
```

**Use case:** Understanding context, finding related work, network analysis

---

### Relationship & Contact Queries

#### "Relationship with [person]" / "Interactions with [person]"

**Optimal multi-index workflow:**
```
1. missing-pages.md (300 tokens)
   └─ Check if person has a page (or should have one)

2. reference-graph.md (5,000 tokens)
   └─ Count references to person
   └─ See where person is mentioned

3. task-index.md (8,000 tokens)
   └─ grep for person name → Find all tasks

4. timeline-full.md (5,000 tokens)
   └─ When did interactions occur?

5. Read 2-3 journal entries for conversation notes

Total: ~18,000 tokens
Result: Comprehensive relationship timeline
```

**Optimization:** If person in missing-pages.md with high ref count (e.g., John smith - 11 refs), you know immediately it's a significant relationship requiring documentation.

---

## Cross-Index Combination Patterns

### Pattern 1: Weekly Progress Report

**Query:** "What did I accomplish this week?"

**Optimal combination:**
```
1. timeline-recent.md (350 tokens)
   └─ 7-day activity overview with task counts

2. tasks-by-status.md (DONE section, ~400 tokens)
   └─ List of completed tasks this week

3. time-tracking.md (1,500 tokens)
   └─ Time allocation by project this week

4. dashboard.md (500 tokens)
   └─ Overall completion rate and trends

Total: 2,750 tokens
Result: Complete weekly summary
vs. Reading all journals: 10,000+ tokens
Efficiency: 3.6x faster
```

---

### Pattern 2: Project Health Check

**Query:** "Is [project] on track?"

**Optimal combination:**
```
1. dashboard.md (500 tokens)
   └─ Check "Top Projects" for time investment

2. tasks-by-status.md (2,000 tokens)
   └─ grep for project → Count NOW vs DONE tasks

3. timeline-recent.md (350 tokens)
   └─ Recent activity on project

4. time-tracking.md (1,500 tokens)
   └─ Time trend (increasing or decreasing?)

Total: 4,350 tokens
Result: Project health metrics
```

---

### Pattern 3: ROI Analysis

**Query:** "Which projects have best ROI?"

**Optimal combination:**
```
1. time-tracking.md (1,500 tokens)
   └─ Time invested per project

2. tasks-by-status.md (DONE section, ~400 tokens)
   └─ Completed tasks per project

3. Calculate: Completed tasks / Hours invested

Total: 1,900 tokens
Result: ROI ranking
```

---

## Token Efficiency Quick Reference

| Query Type | Optimal Index | Tokens | Bad Alternative | Tokens | Gain |
|------------|--------------|--------|-----------------|---------|------|
| "What's my status?" | dashboard.md | 500 | task-index.md | 8,000 | 16x |
| "Recent activity?" | timeline-recent.md | 350 | 7 journals | 2,600 | 7x |
| "Priorities?" | tasks-by-priority.md | 400 | tasks-by-status.md | 2,000 | 5x |
| "Time on project?" | time-tracking.md | 300 | grep journals | 10,000 | 33x |
| "Deep dive project?" | Multi-index | 3,000 | All journals | 10,000 | 3x |
| "Knowledge gaps?" | missing-pages.md | 300 | Manual analysis | N/A | Direct |
| "Weekly review?" | Multi-index | 2,750 | All journals | 10,000+ | 3.6x |

**Average efficiency: 5-30x faster queries**

---

## Common Mistakes & Fixes

### ❌ Mistake 1: Jumping to Complete Datasets

**Bad:**
```
User: "What's my status?"
Claude: Reads task-index.md (8,000 tokens)
```

**Good:**
```
User: "What's my status?"
Claude: Reads dashboard.md (500 tokens) → Sufficient
```

**Fix:** Always start Level 1 (dashboard, timeline-recent, tasks-by-priority)

---

### ❌ Mistake 2: Reading Journals First

**Bad:**
```
User: "What happened this week?"
Claude: Reads 7 journal files (2,600 tokens)
```

**Good:**
```
User: "What happened this week?"
Claude: Reads timeline-recent.md (350 tokens)
```

**Fix:** Check indexes before reading raw journals

---

### ❌ Mistake 3: Loading Multiple Large Indexes

**Bad:**
```
User: "Status of MoonLaunch?"
Claude: Reads task-index.md (8,000) + timeline-full.md (5,000) + reference-graph.md (5,000)
Total: 18,000 tokens
```

**Good:**
```
User: "Status of MoonLaunch?"
Claude: Reads dashboard.md (500) + tasks-by-status.md filtered (200) + key-pages
Total: 700 tokens + pages
```

**Fix:** Use filtered views (Level 2) instead of complete datasets (Level 3)

---

## Advanced Optimization Strategies

### Strategy 1: Keyword Filtering Before Full Read

**Instead of:**
```
Read entire task-index.md (8,000 tokens)
Then filter for keyword
```

**Do:**
```
Read tasks-by-status.md (2,000 tokens)
grep for keyword in relevant section only
If insufficient, then read task-index.md
```

**Savings:** 6,000 tokens if keyword found in tasks-by-status.md

---

### Strategy 2: Progressive Detail Expansion

**Instead of:**
```
Read all indexes for comprehensive answer
```

**Do:**
```
Start: dashboard.md (500 tokens) → Does this answer the question?
If no → timeline-recent.md (350 tokens) → Does this answer?
If no → tasks-by-status.md (2,000 tokens) → Does this answer?
Only if still insufficient → Complete datasets
```

**Savings:** Stop as soon as question is answered

---

### Strategy 3: Cross-Reference Optimization

**Instead of:**
```
Read reference-graph.md to find all pages mentioning X
Then read each page individually
```

**Do:**
```
Read reference-graph.md to find pages
Check if question can be answered from references alone
Only read pages if detailed content needed
```

---

## Query Templates

### Template: Status Report
```
Input: "What's my status?"

Workflow:
1. Read dashboard.md
2. Extract:
   - Quick stats
   - Current priorities
   - Top projects
3. If user wants task detail:
   Read tasks-by-status.md (NOW section)
4. Generate formatted report

Token usage: 500-1,000 tokens
```

---

### Template: Weekly Review
```
Input: "Weekly review"

Workflow:
1. Read timeline-recent.md → 7-day overview
2. Read tasks-by-status.md (DONE) → Completions
3. Read time-tracking.md → Time allocation
4. Generate:
   - Wins (completed tasks)
   - Active work (NOW tasks)
   - Time by project
   - Velocity trend

Token usage: 2,750 tokens
```

---

### Template: Project Deep Dive
```
Input: "Deep dive on [project]"

Workflow:
1. dashboard.md → Project in top 5?
2. tasks-by-status.md → grep project
3. time-tracking.md → Time investment
4. reference-graph.md → Connected pages
5. timeline-recent.md → Recent activity
6. Generate comprehensive report

Token usage: 3,000-4,000 tokens
```

---

## Best Practices Summary

### DO:
✅ **Start smallest** - dashboard.md, timeline-recent.md, tasks-by-priority.md
✅ **Stop when answered** - Don't load more indexes than needed
✅ **Combine strategically** - Use multiple Level 1/2 indexes instead of one Level 3
✅ **Filter first** - grep in smaller indexes before reading large ones

### DON'T:
❌ **Skip levels** - Don't jump to Level 3 without checking Level 1/2
❌ **Read journals first** - Always check indexes first
❌ **Load everything** - Progressive disclosure prevents waste
❌ **Ignore efficiency metrics** - 5-30x gains matter at scale

---

## Related Documentation

- **Index System:** See [index-system.md](./index-system.md) for complete index reference
- **Skills Reference:** See [skills-reference.md](./skills-reference.md) for automation

---

**Last Updated:** 2025-11-07
