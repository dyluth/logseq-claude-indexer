# Skills Reference & Creation Guide

**Last Updated:** 2025-11-07
**Purpose:** Catalog of available skills and guide to creating custom automation

---

## What Are Skills?

**Skills** are workflow automation templates that:
- Auto-activate based on keyword patterns in your prompts
- Execute predefined multi-step workflows
- Leverage the index system for efficiency

**Think of them as:** Keyboard shortcuts for complex multi-step workflows.

---

## Available Skills Catalog

### status-report.md
**Trigger patterns:**
- "status report"
- "what's happening"
- "what's my status"
- "executive summary"
- "current state"

**What it does:**
Generates executive status report using progressive disclosure through indexes.

**Workflow:**
1. Read dashboard.md → Quick stats, priorities, top projects
2. Read timeline-recent.md → Last 7 days activity
3. Read tasks-by-status.md (NOW section) → Active tasks
4. Generate formatted executive report

**Output sections:**
- Overview (total tasks, completion rate, active projects)
- Current Focus (NOW tasks)
- Recent Wins (last 7 days DONE tasks)
- Project Health (per project breakdown)
- Blockers & Risks (old NOW tasks)
- Next Priorities

**Token usage:** ~1,850 tokens vs ~8,000+ without indexes

**When to use:** Daily/weekly check-ins, stakeholder updates

---

### project-closure.md
**Trigger patterns:**
- "close project"
- "archive project"
- "project closure"
- "finish project"

**What it does:**
Systematically closes/archives a project with comprehensive summary.

**Workflow:**
1. **Discovery Phase:**
   - reference-graph.md → Find all project references
   - tasks-by-status.md → Count tasks by status
   - time-tracking.md → Total time invested
   - timeline-full.md → Activity timeline
   - Read project page

2. **Closure Report Generation:**
   - Summary (duration, time, completion rate)
   - References across knowledge base
   - Task summary (completed vs incomplete)
   - Activity timeline
   - Key outcomes

3. **Cleanup Actions** (asks before executing):
   - Update project page with closure report
   - Close/migrate incomplete tasks
   - Update INDEX.md (move to Completed Projects)

4. **Validation:**
   - Verify changes completed
   - Check cross-index consistency

**When to use:** Project completion, pivot cleanup, archiving

**Example:** "Close project WorkJobThing"

---

### weekly-review.md
**Trigger patterns:**
- "weekly review"
- "week summary"
- "what did I do this week"

**What it does:**
Synthesizes weekly activity from indexes and generates progress report.

**Workflow:**
1. Read timeline-recent.md → Last 7 days overview
2. Read tasks-by-status.md (DONE) → Completed tasks
3. Read time-tracking.md → Project time allocation
4. Generate weekly summary

**Output sections:**
- Wins (completed tasks)
- Active work (NOW tasks)
- Time allocation by project
- Velocity (tasks completed vs created)
- Next week priorities

**Token usage:** ~2,750 tokens

**When to use:** Weekly planning, progress tracking, team updates

---

### project-deep-dive.md
**Trigger patterns:**
- "deep dive on [project]"
- "full context for [project]"
- "everything about [project]"

**What it does:**
Comprehensive project analysis across all sources (Logseq + GitHub + Google Drive).

**Workflow:**
1. **Phase 1: Index Reconnaissance**
   - dashboard.md → Check project status
   - tasks-by-status.md → Filter for project
   - reference-graph.md → Connected pages

2. **Phase 2: Deep Context**
   - Read specific pages identified
   - Read recent journals mentioning project

3. **Phase 3: Synthesis**
   - Overview, active tasks, time invested
   - Recent activity, external status
   - Connected context, recommendations

**Token usage:** ~3,000 tokens + external reads

**When to use:** Project planning, stakeholder questions, context switching

**Example:** "Deep dive on MoonLaunch"

---

### knowledge-gaps.md
**Trigger patterns:**
- "knowledge gaps"
- "missing pages"
- "what should I document"

**What it does:**
Identifies and prioritizes missing documentation using missing-pages.md and reference-graph.md.

**Workflow:**
1. Read missing-pages.md → Pages with 5+ refs
2. Read reference-graph.md → Hub pages analysis
3. Read tasks-by-status.md → Grep undocumented concepts
4. For each missing page: Grep for context
5. Generate prioritized creation plan

**Output:**
- High Priority (10+ refs) - People/concepts
- Medium Priority (5-9 refs)
- Quick Wins (clear single-purpose pages)
- For each: Suggested structure, context

**When to use:** Documentation planning, knowledge base maintenance

---

### time-roi-analysis.md
**Trigger patterns:**
- "time analysis"
- "where did I spend time"
- "ROI analysis"
- "time investment"

**What it does:**
Analyzes time allocation across projects and calculates ROI.

**Workflow:**
1. Read time-tracking.md → Time per project
2. Read timeline-recent.md → Recent activity
3. Read tasks-by-status.md (DONE) → Completed tasks
4. Calculate metrics:
   - Total time per project
   - Avg time per task
   - Tasks completed per project
   - Time ROI (completed / invested)
   - Trend (this week vs previous)

**Output:**
- Time Allocation table (hours, tasks, avg, completed)
- This Week Focus
- Recommendations (over/under-invested, quick wins)

**When to use:** Weekly planning, project prioritization, resource allocation

---

### contact-mapper.md
**Trigger patterns:**
- "relationship with [person]"
- "interactions with [person]"
- "contact history [person]"

**What it does:**
Maps relationship timeline with contact using cross-index analysis.

**Workflow:**
1. missing-pages.md → Does person have page?
2. reference-graph.md → Count references
3. task-index.md → Grep for person, find tasks
4. timeline-full.md → When did interactions occur?
5. Read referenced journal entries
6. Generate relationship map

**Output:**
- Person Name
- References (X mentions across Y entries)
- Timeline (first interaction to most recent)
- Associated Projects
- Open Actions (tasks involving person)
- Recommended Next Steps

**When to use:** Outreach prep, relationship management, networking

**Example:** "Relationship with John smith"

---

## How Skills Work

### Activation Flow

```
1. You submit prompt: "status report"
   ↓
2. Claude Code detects keyword pattern match
   ↓
3. Matching skill auto-loads
   ↓
4. Skill provides Claude with:
   - Which indexes to read (and in what order)
   - What to extract from each
   - How to format output
   ↓
5. Claude executes workflow
   ↓
6. Result delivered
```

**No manual invocation needed** - just use the trigger words naturally.

---

## Creating Custom Skills

### Skill File Structure

Skills are Markdown files in `.claude/skills/` with YAML frontmatter:

```yaml
---
metadata:
  patterns:
    - "trigger phrase 1"
    - "trigger phrase 2"
    - "trigger phrase 3"
  description: "Brief description of what skill does"
  auto_load_instructions: true
---

instructions: |
  When user requests [trigger], follow this workflow:

  Step 1: [First action]
  - Read [specific index]
  - Extract [specific information]

  Step 2: [Second action]
  - [Processing steps]

  Step 3: [Final output]
  - Generate [formatted result]

  Reference: [Link to relevant guide if needed]
```

---

### Example: Simple Skill

**File:** `.claude/skills/priority-check.md`

```yaml
---
metadata:
  patterns:
    - "what's urgent"
    - "priority check"
    - "high priority tasks"
  description: "Quick check of high-priority tasks"
  auto_load_instructions: true
---

instructions: |
  Quick priority check workflow:

  1. Read .claude/indexes/tasks-by-priority.md
     Extract: All [#A] and [#B] tasks

  2. Group by status:
     - NOW (active)
     - TODO (backlog)
     - DONE (reference)

  3. Format output:
     ## Priority Tasks

     **NOW (X tasks):**
     - [List]

     **TODO (Y tasks):**
     - [List]

  Token usage: ~400 tokens
```

---

### Skill Creation Checklist

When creating a new skill:

**1. Define Purpose**
- ☐ What problem does this solve?
- ☐ How often will this be used?
- ☐ Is it worth automating?

**2. Choose Trigger Patterns**
- ☐ What natural phrases would user say?
- ☐ Are patterns specific enough (not too broad)?
- ☐ Are patterns memorable?

**3. Map Workflow**
- ☐ Which indexes are needed?
- ☐ What order to read them in?
- ☐ How to synthesize results?

**4. Define Output Format**
- ☐ What sections?
- ☐ How detailed?
- ☐ Structured or narrative?

**5. Estimate Token Usage**
- ☐ Calculate: Sum of index reads + processing
- ☐ Compare to alternative approach
- ☐ Document efficiency gain

**6. Test**
- ☐ Does it activate correctly?
- ☐ Does it produce expected output?
- ☐ Is token usage as estimated?

---

## Skill Best Practices

### DO:
✅ **Use progressive disclosure** - Start with smallest indexes
✅ **Reference guides** - Link to index-system.md, query-patterns.md, etc.
✅ **Specify order** - Numbered steps for clarity
✅ **Include examples** - Show expected output format
✅ **Document tokens** - Help users understand efficiency

### DON'T:
❌ **Overly broad patterns** - "write" triggers on everything
❌ **Skip index optimization** - Reading all large indexes wastes tokens
❌ **Make skills too long** - Break complex workflows into multiple skills

---

## Skill Composition Patterns

### Pattern 1: Chaining Skills

**Instead of:** One mega-skill for entire workflow

**Do:** Multiple focused skills that can be combined

```
"status report" → Generates overview
"deep dive on MoonLaunch" → Detailed MoonLaunch analysis
"time ROI analysis" → Time allocation review

User can chain: "status report" followed by "deep dive on MoonLaunch"
```

**Benefit:** Modular, reusable, easier to maintain

---

### Pattern 2: Conditional Branching

**Inside a skill, branch based on context:**

```yaml
instructions: |
  When analyzing project:

  1. Read dashboard.md → Is project in top 5?

  2. If YES:
     - Read time-tracking.md for details
     - Project is well-tracked

  3. If NO:
     - Read reference-graph.md to find pages
     - Project may need more documentation
```

---

## Advanced Skill Examples

### Multi-Project Comparison

```yaml
---
metadata:
  patterns: ["compare projects", "project comparison"]
  description: "Compare metrics across projects"
---

instructions: |
  Compare all active projects:

  1. dashboard.md → Top projects list
  2. For each project:
     - time-tracking.md → Time invested
     - tasks-by-status.md → Active vs complete tasks
     - timeline-recent.md → Recent activity

  3. Generate comparison table:
     | Project | Time | Active | Done | Velocity |
     |---------|------|--------|------|----------|

  4. Rank by velocity (Done / Time)
  5. Recommend: Focus area based on metrics
```

---

### Contact Network Map

```yaml
---
metadata:
  patterns: ["network map", "contact network", "who do I know"]
  description: "Map contact network by category"
---

instructions: |
  Map contact network:

  1. missing-pages.md → People with high ref counts
  2. reference-graph.md → Find all person pages
  3. Categorize by:
     - Role (mentor, client, collaborator)
     - Project (which projects mentioned)
     - Status (active vs past)

  4. Generate network diagram:
     ## Active Contacts
     ### Mentors
     - [Name] - [X] interactions - [Last contact]

     ### Clients
     - [Name] - [X] interactions - [Last contact]
```

---

## Troubleshooting Skills

### Issue: Skill Not Activating

**Possible causes:**
1. Trigger pattern too specific
2. Typo in pattern
3. Skill file not in `.claude/skills/`
4. YAML frontmatter syntax error

**Fix:**
- Check pattern matches your query exactly
- Verify file location
- Validate YAML syntax (use YAML validator)

---

### Issue: Skill Loads Wrong Indexes

**Possible cause:** Workflow specifies wrong files

**Fix:** Review instructions section, ensure correct index names:
- `.claude/indexes/dashboard.md` (not `dashboard.yaml`)
- `.claude/indexes/tasks-by-status.md` (not `task-status.md`)

---

### Issue: Skill Uses Too Many Tokens

**Possible cause:** Loading large indexes unnecessarily

**Fix:**
- Review workflow, ensure progressive disclosure
- Start with Level 1 indexes (dashboard, timeline-recent)
- Only load Level 3 (task-index, timeline-full) when necessary

---

## Related Documentation

- **Index System:** See [index-system.md](./index-system.md) for index reference
- **Query Patterns:** See [query-patterns.md](./query-patterns.md) for optimization


---

**Last Updated:** 2025-11-07
