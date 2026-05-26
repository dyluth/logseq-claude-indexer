---
metadata:
  patterns:
    - "weekly review"
    - "week summary"
    - "what did I do this week"
    - "weekly progress"
    - "this week's work"
  description: "Synthesize weekly activity from indexes and generate progress report"
  auto_load_instructions: true
---

instructions: |
  Generate comprehensive weekly review using cross-index analysis.

  Reference: .claude/guides/index-system.md for index details
  Reference: .claude/guides/query-patterns.md for weekly review pattern

  Workflow:

  Step 1: Read Recent Activity Timeline
  1. Read .claude/indexes/timeline-recent.md (350 tokens)
     Extract:
     - Last 7 days activity summary
     - Days with activity vs no activity
     - Task counts by day (NOW/TODO/DOING/DONE)
     - Time logged by day (if applicable)

  Step 2: Read Completed Tasks
  2. Read .claude/indexes/tasks-by-status.md (DONE section only, ~400 tokens)
     Extract:
     - Tasks completed in last 7 days (check file references for recent dates)
     - Group completed tasks by project if identifiable
     - Note any high-priority completions ([#A] or [#B])

  Step 3: Read Time Allocation
  3. Read .claude/indexes/time-tracking.md (1,500 tokens)
     Extract:
     - Weekly breakdown section (last week's data)
     - Time by project for recent week
     - Compare to previous weeks if available
     - Note: Hours logged, tasks tracked

  Step 4: Read Current Active Work
  4. Read .claude/indexes/tasks-by-status.md (NOW section only, ~500 tokens)
     Extract:
     - Currently active tasks (as of end of week)
     - Compare to start of week if possible (what moved from TODO → NOW?)
     - Identify new NOW tasks added this week

  Step 5: Calculate Velocity
  5. Calculate metrics:
     - Tasks completed this week: [count from DONE]
     - Tasks started this week: [new NOW tasks]
     - Net task change: [completed - started]
     - Velocity: [completed / 7 days] = [X] tasks/day

  Step 6: Generate Weekly Review Report

  # Weekly Review - [Date Range]
  *Week of [Start Date] to [End Date]*

  ## Overview

  - **Active days**: [X] of 7 days had activity
  - **Tasks completed**: [X] tasks
  - **Velocity**: [X] tasks/day average
  - **Time logged**: [X] hours (if time tracking used)

  ## Wins (Completed Tasks)

  [Group by project or category, show all DONE tasks from week]

  ### [Project 1 Name]
  - ✅ [Completed task 1]
  - ✅ [Completed task 2]

  ### [Project 2 Name]
  - ✅ [Completed task 1]

  ### General
  - ✅ [Other completed tasks]

  **Highlights**:
  - [Pick 2-3 most significant completions]

  ## Active Work (Current NOW Tasks)

  [Group current NOW tasks by project]

  ### [Project 1 Name]
  - 🔄 [NOW task 1]
  - 🔄 [NOW task 2]

  ### [Project 2 Name]
  - 🔄 [NOW task 1]

  **New this week**:
  - [Tasks that moved to NOW this week]

  ## Time Allocation

  [From time-tracking.md weekly breakdown]

  | Project | Hours | Tasks | Avg/Task | % of Week |
  |---------|-------|-------|----------|-----------|
  | [Project 1] | [X]h | [Y] | [Z]h | [%] |
  | [Project 2] | [X]h | [Y] | [Z]h | [%] |
  | [Project 3] | [X]h | [Y] | [Z]h | [%] |
  | **Total** | **[X]h** | **[Y]** | **[Z]h** | **100%** |

  **Time insights**:
  - Most time: [Project with highest hours]
  - Most productive: [Project with best completion/hour ratio]
  - [Any notable trends: increasing/decreasing time on projects]

  ## Velocity & Progress

  - **Tasks completed**: [X] this week vs [Y] last week → [Trend: ↑/↓/→]
  - **Velocity**: [X] tasks/day
  - **Completion rate**: [Completed / (Completed + Started)]%
  - **Net progress**: [Completed - Started] = [+/- X] tasks

  **Trend analysis**:
  - [Improving/Declining/Steady] velocity compared to previous weeks
  - [Note any patterns: certain days more productive, etc.]

  ## Project Health

  [For each major project, synthesize status]

  ### [Project 1 Name]
  - **Time this week**: [X]h
  - **Tasks completed**: [X]
  - **Tasks active**: [X] NOW tasks
  - **Status**: [On track / Needs attention / Blocked]
  - **External alignment**: [If GitHub/Drive checked]
  - **Next week focus**: [Infer from NOW tasks]

  ### [Project 2 Name]
  - [Similar breakdown]

  ## External Activity

  [If GitHub or Google Drive checked]

  ### GitHub Activity (MoonLaunch)
  - **Commits this week**: [X] commits
  - **Key changes**: [From commit messages if available]
  - **Alignment**: Internal tasks ↔ External code [Aligned/Misaligned]

  ### Google Drive Activity
  - **Documents updated**: [List if found]
  - **Strategy changes**: [Note any significant updates]

  ## Blockers & Challenges

  [Identify from timeline if patterns visible]

  - [Days with no activity → Why?]
  - [Tasks stuck in NOW >7 days → Blockers?]
  - [Projects with low completion rate → Challenges?]

  ## Next Week Priorities

  [Based on current NOW tasks and project status]

  **Must complete**:
  1. [High-priority or blocked tasks to resolve]
  2. [Key deliverables approaching deadlines]

  **Focus areas**:
  - [Project needing most attention based on health]
  - [Knowledge gaps to address if significant]

  **Stretch goals**:
  - [Optional items if time allows]

  ## Recommendations

  [Provide 2-4 actionable recommendations based on data]

  1. **[Recommendation 1]**: [Based on velocity, time allocation, or blockers]
     - Action: [Specific step to take]

  2. **[Recommendation 2]**: [Based on project health or external alignment]
     - Action: [Specific step to take]

  3. **[Recommendation 3]**: [Based on knowledge gaps or process improvements]
     - Action: [Specific step to take]

  ---

  **Total context used**: ~2,750 tokens (timeline + tasks + time-tracking)
  **vs. reading all journals**: ~10,000+ tokens
  **Efficiency**: ~3.6x faster

  Example recommendations based on patterns:
  - Low velocity → "Focus on fewer parallel tasks"
  - Blockers identified → "Schedule time to resolve X blocker"
  - Misaligned external → "Update strategy docs to match actual work"
  - High time, low completion → "Break large tasks into smaller chunks"

  Next Steps After Review:
  - Use "deep dive on [project]" for any project needing attention
  - Use "time ROI analysis" if time allocation concerns arise
  - Use "status report" for high-level update to stakeholders
  - Consider "knowledge gaps" if documentation recommendations made
