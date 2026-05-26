# Skills README - Making the Most of Your Optimized System

**Last Updated:** 2025-11-07
**Purpose:** Practical guide to using your fully optimized Claude Code + Logseq system

---

## What You've Built

You now have a **fully optimized AI-powered knowledge management system** with:

✅ **9 auto-generated indexes** (refreshes on every commit)
✅ **5 specialized guides** (~1,300 lines of documentation)
✅ **8 automation skills** (workflow templates)
✅ **CLAUDE.md** (concise 205-line entry point)
✅ **INDEX.md** (310-line navigation hub)
✅ **Cross-source integration** (Logseq + GitHub + Google Drive)

**Result:** 5-30x more efficient queries, automated workflows, and comprehensive cross-source intelligence.

---

## Quick Start: Your First 5 Minutes

### 1. Try a Status Check
```
You: "What's my status?"
```

**What happens:**
- Claude reads dashboard.md (500 tokens)
- Shows: Current priorities, top projects, recent activity
- **vs. before:** Would read 8,000+ tokens

### 2. Check Recent Activity
```
You: "What happened this week?"
```

**What happens:**
- Claude reads timeline-recent.md (350 tokens)
- Shows: Last 7 days activity summary
- **vs. before:** Would read 7 journal files (2,600 tokens) = 7x faster

### 3. See Your Priorities
```
You: "What are my priorities?"
```

**What happens:**
- Claude reads tasks-by-priority.md (400 tokens)
- Shows: 7 high-priority tasks ([#A] and [#B])
- **vs. before:** Would read all tasks (2,000 tokens) = 5x faster

---

## The 8 Skills: When to Use Each

### Daily Use Skills

#### 1. status-report
**Trigger:** "status report" / "what's happening" / "current state"

**Use when:**
- Starting your day (quick overview)
- Preparing for standup/meeting
- Weekly check-in with yourself
- Need executive summary

**What you get:**
- Overview (total tasks, completion rate)
- Current Focus (NOW tasks by project)
- Recent Wins (DONE tasks last 7 days)
- Project Health (per-project breakdown)
- Next Priorities
- Recommendations

**Time saved:** 4-5x faster than manual status check

---

#### 2. weekly-review
**Trigger:** "weekly review" / "what did I do this week"

**Use when:**
- Monday morning planning
- Friday afternoon wrap-up
- Monthly reflection
- Progress tracking

**What you get:**
- Wins (completed tasks)
- Active work (current NOW tasks)
- Time allocation table (hours per project)
- Velocity metrics (tasks/day)
- Project health per project
- Blockers & challenges
- Next week priorities
- Recommendations

**Best practice:** Run every Monday to plan week and every Friday to reflect

---

### Project Management Skills

#### 3. project-start
**Trigger:** "start project [name]" / "new project"

**Use when:**
- Beginning any new project
- New client engagement
- Starting research initiative
- New side hustle

**What it does:**
- Asks for project details (name, purpose, type, duration)
- Creates project page (pages/[Project Name].md)
- Creates initial tasks (optional)
- **Auto-updates INDEX.md** (adds to Active Projects)
- Provides next steps

**Time saved:** 10-15 minutes of setup + ensures consistent structure

**Example:**
```
You: "start project Customer Analytics Platform"
Claude: [Asks questions] → Creates structure → Updates INDEX.md → Confirms
```

---

#### 4. project-closure
**Trigger:** "close project [name]" / "archive project"

**Use when:**
- Project completed
- Pivot away from initiative
- Contract ends
- Research concluded

**What it does:**
- Discovers all project references (cross-index analysis)
- Counts tasks by status (NOW/TODO/DONE)
- Calculates time invested and completion rate
- Generates comprehensive closure report
- Handles incomplete tasks (close or migrate)
- **Auto-updates INDEX.md** (moves to Completed Projects)
- Validates all changes

**Time saved:** Complete automation vs 30+ minutes manual archiving

**Example:**
```
You: "close project WorkJobThing"
Claude: [Analyzes 14 refs, 5 tasks, 2,440h] → Generates report → Proposes cleanup → Updates INDEX.md
```

---

#### 5. project-deep-dive
**Trigger:** "deep dive on [project]" / "everything about [project]"

**Use when:**
- Context switching to different project
- Preparing for project meeting
- Onboarding someone to project
- Need comprehensive project view
- Strategic planning

**What you get:**
- Executive summary (status, health, momentum)
- Task distribution table
- Time investment analytics
- Recent activity timeline
- Network analysis (connected pages)
- **External context** (GitHub commits + Google Drive docs)
- Key insights (strengths, challenges, opportunities)
- Decision history
- Recommendations (immediate + medium-term + strategic)

**Unique power:** Combines Logseq + GitHub + Google Drive for complete picture

**Example:**
```
You: "deep dive on MoonLaunch"
Claude: [Reads indexes] → [Checks GitHub] → [Checks Google Drive] → Comprehensive analysis
```

---

### Optimization Skills

#### 6. knowledge-gaps
**Trigger:** "knowledge gaps" / "missing pages" / "what should I document"

**Use when:**
- Monthly maintenance
- Noticing repeated re-explanation
- Onboarding prep
- Documentation sprint planning

**What you get:**
- Missing pages prioritized (high/medium/low)
- Each with: reference count, type, context, suggested structure
- Pattern analysis (clusters of related gaps)
- Comparison to hub pages
- Quick wins (5-minute pages)
- Recommendations with time estimates

**Best practice:** Run monthly, keep missing pages <5 at any time

**Example:**
```
You: "knowledge gaps"
Claude: [Identifies 9 missing pages] → John smith (11 refs) → Creates prioritized list
```

---

#### 7. time-roi-analysis
**Trigger:** "time analysis" / "ROI analysis" / "where did I spend time"

**Use when:**
- Monthly planning
- Feeling unproductive
- Need to prioritize projects
- Resource allocation decisions
- Quarterly review

**What you get:**
- Time allocation table (by project)
- ROI rankings (tasks/hour)
- Recent trends (last 8 weeks)
- Top performers (best ROI)
- Needs attention (low ROI)
- Quadrant analysis (time vs completion)
- Recommendations (optimize/improve/reallocate)
- This week focus (time budget per project)

**Key metric:** Tasks completed / Hours invested = ROI score

**Example:**
```
You: "time ROI analysis"
Claude: [Calculates ROI for all projects] → WorkJobThing: 0.82 tasks/h → MoonLaunch: 2.1 tasks/h → Recommendations
```

---

### Relationship Management

#### 8. contact-mapper
**Trigger:** "relationship with [person]" / "contact history"

**Use when:**
- Preparing for meeting
- Haven't talked to someone in a while
- Planning outreach
- Need context before call
- Relationship maintenance

**What you get:**
- Relationship timeline (all interactions chronologically)
- Contact information (if available)
- Associated projects & topics
- Open actions & commitments
- Completed collaborations
- Relationship context (interests, expertise, value)
- **External context** (Google Drive meeting notes)
- Network connections
- Recommended next steps
- Suggested page structure (if page missing)

**Example:**
```
You: "relationship with John smith"
Claude: [Finds 11 refs, no page] → Complete timeline → Open actions → Recommends page creation
```

---

## Cross-Source Power Moves

### Combining Logseq + GitHub + Google Drive

**Your system uniquely combines 3 sources:**

1. **Logseq** (internal knowledge)
   - Tasks, time tracking, journal notes
   - 9 auto-generated indexes

2. **GitHub** (external technical)
   - MoonLaunch repository
   - Code, commits, README

3. **Google Drive** (external strategic)
   - Business plans, strategy docs
   - LinkedIn posts, meeting notes

### Power Query Examples

**Example 1: Complete MoonLaunch Picture**
```
You: "deep dive on MoonLaunch"

Claude combines:
- Logseq tasks → 3 MoonLaunch tasks (2 NOW, 1 TODO)
- Logseq time → 45h invested
- GitHub → 12 commits this week, active development
- Google Drive → Pivot analysis doc, target market research
→ Synthesizes complete picture across all sources
```

**Example 2: LinkedIn Content with Past Style**
```
You: "Draft LinkedIn post about MoonLaunch"

Claude combines:
- Google Drive → Reads past LinkedIn posts for style
- Logseq tasks → Current MoonLaunch focus
- GitHub → MoonLaunch technical features
- Google Drive → MoonLaunch strategy for positioning
→ Generates on-brand post informed by all sources
```

**Example 3: Weekly Review with External Activity**
```
You: "weekly review"

Claude combines:
- Logseq timeline → Internal activity
- Logseq time-tracking → Hours per project
- GitHub → MoonLaunch commits this week
- Google Drive → Strategy doc updates
→ Complete weekly picture: internal + external
```

---

## Best Practices: Getting the Most Value

### Daily Routine

**Morning (5 minutes):**
```
"What's my status?"
```
→ See priorities, plan day

**End of day (2 minutes):**
- Update tasks in Logseq (mark DONE, add new)
- Commit changes (regenerates indexes automatically)

### Weekly Routine

**Monday (10 minutes):**
```
"weekly review"
```
→ Plan week based on last week's data

**Friday (5 minutes):**
```
"weekly review"
```
→ Celebrate wins, identify blockers for next week

### Monthly Routine

**First Monday (30 minutes):**
```
1. "weekly review" (understand trends)
2. "time ROI analysis" (evaluate time allocation)
3. "knowledge gaps" (plan documentation work)
4. "deep dive on [top priority project]" (strategic focus)
```

### Project Lifecycle

**Starting:**
```
"start project [name]"
```
→ Automated setup + INDEX.md update

**Working:**
- Add tasks to daily journals
- Update project page as needed
- Use "deep dive on [project]" when context switching

**Closing:**
```
"close project [name]"
```
→ Automated closure + INDEX.md update

---

## Understanding the Index System

### Progressive Disclosure (Always Follow This Order)

**Level 1: Start Here** (400-500 tokens)
- dashboard.md → Overview
- timeline-recent.md → Recent activity
- tasks-by-priority.md → Urgent items

**Use for:** 80% of queries

**Level 2: Filtered Views** (300-2,000 tokens)
- tasks-by-status.md → All tasks by stage
- time-tracking.md → Project time
- missing-pages.md → Documentation gaps

**Use when:** Level 1 insufficient

**Level 3: Complete Datasets** (5,000-8,000 tokens)
- task-index.md → All tasks with metadata
- timeline-full.md → Complete history
- reference-graph.md → Full knowledge graph

**Use sparingly:** Only when Level 2 insufficient

### Quick Index Reference

| Need | Index | Tokens |
|------|-------|--------|
| Quick status | dashboard.md | 500 |
| This week | timeline-recent.md | 350 |
| Urgent items | tasks-by-priority.md | 400 |
| All tasks | tasks-by-status.md | 2,000 |
| Time invested | time-tracking.md | 1,500 |
| Doc gaps | missing-pages.md | 300 |
| Search tasks | task-index.md | 8,000 |
| History | timeline-full.md | 5,000 |
| Connections | reference-graph.md | 5,000 |

---

## Common Workflows

### Morning Planning
```
1. "What's my status?" → See overview
2. Review NOW tasks
3. Pick top 3 to focus on today
```

### Context Switching
```
You: "deep dive on [project I haven't worked on this week]"
→ Get complete context before diving in
```

### End of Week Reflection
```
You: "weekly review"
→ See wins, time allocation, velocity, blockers
→ Plan next week
```

### Project Handoff
```
You: "deep dive on [project]"
→ Share comprehensive report with team member
→ Complete context transfer in one document
```

### Monthly Planning
```
1. "time ROI analysis" → See where time is going
2. "knowledge gaps" → Plan documentation work
3. For each active project: "deep dive on [project]"
4. Make strategic decisions on time allocation
```

### Preparing for Meeting
```
You: "relationship with [person you're meeting]"
→ See complete interaction history
→ Review open commitments
→ Prepare talking points
```

---

## Troubleshooting

### "Indexes seem out of date"
**Fix:** Indexes regenerate on every commit
```
git add . && git commit -m "refresh indexes"
```

### "Skill not activating"
**Check:** Are you using the exact trigger pattern?
- "status report" ✓
- "show me status" ✗ (use "what's my status?" instead)

**See:** [Skills Reference](./../guides/skills-reference.md) for exact patterns

### "Too much information"
**Fix:** Use progressive disclosure
- Start with Level 1 indexes
- Only expand if needed
- Use specific skills instead of manual queries

### "Can't find something"
**Try:**
1. "knowledge gaps" → See if page is missing
2. Use reference-graph.md → Find connections
3. grep for keyword in task-index.md

---

## Advanced Tips

### Skill Chaining
```
You: "weekly review"
[Reviews shows low ROI on Project X]
You: "deep dive on Project X"
[Deep dive reveals blockers]
You: [Take action on blockers]
```

### Combining Skills
```
Morning: "status report" → Overview
Midday: "deep dive on [priority project]" → Deep focus
Evening: Mark tasks DONE, commit
Friday: "weekly review" → Reflection
Monthly: "time ROI analysis" + "knowledge gaps" → Optimization
```

---

## Efficiency Gains: By the Numbers

| Operation | Before | After | Gain |
|-----------|--------|-------|------|
| Status check | 8,000 tokens | 500 tokens | **16x** |
| Recent activity | 2,600 tokens | 350 tokens | **7x** |
| Priority check | 2,000 tokens | 400 tokens | **5x** |
| Time on project | 10,000+ tokens | 300 tokens | **33x** |
| Weekly review | 10,000+ tokens | 2,750 tokens | **3.6x** |
| Project closure | Manual 30+ min | Automated 2 min | **15x** |

**Average:** 5-30x more efficient queries

---

## Next-Level Usage

### Multi-Project Management
```
You: "status report"
→ See all projects
You: "deep dive on [Project A]"
→ Deep focus on one
You: "time ROI analysis"
→ Evaluate allocation across all
```

### Team Collaboration
```
Share skills-generated reports:
- "status report" → Team standup
- "deep dive on [project]" → Project handoff
- "weekly review" → Manager update
```

### Strategic Planning
```
Monthly: "time ROI analysis" + all project deep dives
→ Data-driven decisions on where to invest time
→ Identify projects to close, start, or double-down on
```

---

## Getting Help

### Documentation
- **Quick reference:** [CLAUDE.md](../../CLAUDE.md)
- **Navigation:** [INDEX.md](../../INDEX.md)
- **Detailed guides:** [.claude/guides/](./../guides/)
- **Design doc:** [claude.md-optimisation-design.md](../../claude.md-optimisation-design.md)

### Skills Reference
- **This file:** Overview and practical usage
- **Full catalog:** [.claude/guides/skills-reference.md](./../guides/skills-reference.md)
- **Creating custom skills:** See guides/skills-reference.md

### Query Optimization
- **Decision trees:** [.claude/guides/query-patterns.md](./../guides/query-patterns.md)
- **Index reference:** [.claude/guides/index-system.md](./../guides/index-system.md)

---

## Your Optimized Workflow

**You now have:**
- ✅ 5-30x more efficient queries
- ✅ Automated project lifecycle (start → close)
- ✅ Cross-source intelligence (Logseq + GitHub + Drive)
- ✅ 8 automation skills for common workflows
- ✅ Complete documentation with progressive disclosure

**Result:** Work smarter, not harder. Let the system handle the busywork while you focus on high-value work.

---

**Last Updated:** 2025-11-07

**Quick Start Reminder:**
1. Try: "What's my status?"
2. Try: "What happened this week?"
3. Try: "knowledge gaps"
4. Explore: Use other skills as needed
5. Optimize: Run "time ROI analysis" monthly

Welcome to your optimized knowledge management system!
