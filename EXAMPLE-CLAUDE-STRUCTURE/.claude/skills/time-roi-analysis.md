---
metadata:
  patterns:
    - "time analysis"
    - "where did I spend time"
    - "ROI analysis"
    - "time investment"
    - "time allocation"
    - "productivity analysis"
  description: "Analyze time allocation across projects and calculate ROI"
  auto_load_instructions: true
---

instructions: |
  Analyze time investment across projects and calculate return on investment (ROI).

  Reference: .claude/guides/index-system.md for time-tracking details

  Workflow:

  Step 1: Read Time Tracking Data
  1. Read .claude/indexes/time-tracking.md (1,500 tokens)
     Extract:
     - Total time logged (all time)
     - Total tasks tracked
     - Average time per task (overall)
     - Most productive week
     - Top projects (by time):
       - Project name
       - Total time
       - Task count
       - Average time per task
     - Weekly breakdown (last 8 weeks shown)
     - Time by priority ([#A] vs [#None])
     - Time by status (DONE vs TODO vs DOING vs LATER)

  Step 2: Read Recent Activity
  2. Read .claude/indexes/timeline-recent.md (350 tokens)
     Extract:
     - Last 7 days activity
     - Which projects active recently
     - Time logged per day (if shown)

  Step 3: Read Task Completion Data
  3. Read .claude/indexes/tasks-by-status.md (DONE section, ~400 tokens)
     Extract:
     - Total DONE tasks
     - DONE tasks by project (count)
     - DONE tasks with time logged vs without

  4. Read .claude/indexes/tasks-by-status.md (ALL sections, 2,000 tokens)
     - Count total tasks per project (all statuses)
     - Calculate completion rate per project: DONE / Total

  Step 4: Calculate ROI Metrics

  For each project:
  - **Time invested**: [from time-tracking]
  - **Tasks completed**: [DONE count for project]
  - **Tasks total**: [all statuses count]
  - **Completion rate**: [DONE / Total]%
  - **Time per completed task**: [Time invested / Tasks completed]
  - **ROI score**: [Tasks completed / Time invested] = tasks per hour

  Sort projects by:
  - Highest time investment (most hours)
  - Highest ROI (most tasks/hour)
  - Lowest ROI (least tasks/hour)
  - Highest completion rate

  Step 5: Generate Analysis Report

  # Time & ROI Analysis

  **Analysis Date**: [Current Date]
  **Analysis Period**: [All time / or specify if focusing on recent weeks]

  ## Executive Summary

  - **Total time tracked**: [X] hours across [Y] tasks
  - **Time tracking adoption**: [Z]% of all tasks
  - **Average time/task**: [A] hours
  - **Most productive week**: [Date] with [X] hours logged
  - **Top ROI project**: [Project name] with [X] tasks/hour
  - **Highest time investment**: [Project name] with [X] hours

  ## Time Allocation by Project

  | Project | Time | Tasks | Completed | Completion % | Avg/Task | ROI (tasks/h) | Rank |
  |---------|------|-------|-----------|--------------|----------|---------------|------|
  | [Project 1] | [X]h | [Y] | [Z] | [%] | [A]h | [B] | [#] |
  | [Project 2] | [X]h | [Y] | [Z] | [%] | [A]h | [B] | [#] |
  | [Project 3] | [X]h | [Y] | [Z] | [%] | [A]h | [B] | [#] |
  | [Project 4] | [X]h | [Y] | [Z] | [%] | [A]h | [B] | [#] |
  | **Total** | **[X]h** | **[Y]** | **[Z]** | **[%]** | **[A]h** | **[B]** | - |

  **Key metrics**:
  - Highest time: [Project] - [X]h ([%] of total time)
  - Lowest time: [Project] - [X]h ([%] of total time)
  - Best ROI: [Project] - [X] tasks/h
  - Worst ROI: [Project] - [X] tasks/h

  ## Time Distribution

  **By project category**:
  [If projects can be categorized: Work, Research, Side Projects, etc.]

  - Work projects: [X]h ([%] of total)
  - Research: [X]h ([%] of total)
  - Side projects: [X]h ([%] of total)
  - Other: [X]h ([%] of total)

  **By priority**:
  - High priority [#A]: [X]h ([%] of total)
  - No priority: [X]h ([%] of total)

  **By status**:
  - Completed (DONE): [X]h ([%] of total)
  - In progress (DOING): [X]h ([%] of total)
  - Planned (TODO): [X]h ([%] of total)
  - Deferred (LATER): [X]h ([%] of total)

  ## Recent Trends (Last 8 Weeks)

  **Weekly time allocation**:
  [From time-tracking weekly breakdown]

  | Week | Hours | Tasks | Trend |
  |------|-------|-------|-------|
  | [Most recent] | [X]h | [Y] | ↑/↓/→ |
  | [Week -1] | [X]h | [Y] | ↑/↓/→ |
  | [Week -2] | [X]h | [Y] | ↑/↓/→ |
  | [Week -3] | [X]h | [Y] | ↑/↓/→ |
  | [Week -4] | [X]h | [Y] | ↑/↓/→ |

  **Trend analysis**:
  - Time investment: [Increasing / Decreasing / Stable]
  - Task volume: [Increasing / Decreasing / Stable]
  - Productivity (tasks/week): [Improving / Declining / Stable]

  **Most active week**: [Date] - [X]h logged
  **Least active week**: [Date] - [X]h logged
  **Average weekly time**: [X]h

  ## ROI Rankings

  ### Top Performers (Best ROI)

  **Projects with highest tasks/hour** (efficient time use):

  1. **[Project 1]** - [X] tasks/hour
     - Time: [Y]h
     - Completed: [Z] tasks
     - Why efficient: [Small tasks / Clear scope / etc. - infer from data]

  2. **[Project 2]** - [X] tasks/hour
     - Time: [Y]h
     - Completed: [Z] tasks
     - Why efficient: [Reason]

  3. **[Project 3]** - [X] tasks/hour
     - Time: [Y]h
     - Completed: [Z] tasks
     - Why efficient: [Reason]

  ### Needs Attention (Low ROI)

  **Projects with lowest tasks/hour** (time-intensive):

  1. **[Project 1]** - [X] tasks/hour
     - Time: [Y]h (significant investment)
     - Completed: [Z] tasks (low output)
     - Completion rate: [%]
     - **Concern**: [High time, low output / Many incomplete tasks / etc.]
     - **Recommendation**: [Break into smaller tasks / Re-evaluate scope / etc.]

  2. **[Project 2]** - [X] tasks/hour
     - [Similar breakdown]

  ### Incomplete Projects (Low Completion Rate)

  **Projects with <50% completion**:

  1. **[Project 1]** - [X]% complete
     - Total tasks: [Y]
     - Completed: [Z]
     - Remaining: [A]
     - **Concern**: Many started tasks, few finished
     - **Recommendation**: [Finish existing before starting new / Abandon if no longer relevant]

  ## Time vs. Completion Analysis

  **Quadrant analysis**:

  **High Time + High Completion** ✅ (Productive investment):
  - [List projects: Significant effort, good results]
  - **Action**: Continue current approach

  **High Time + Low Completion** ⚠️ (Concerning):
  - [List projects: Significant effort, poor results]
  - **Action**: Re-evaluate scope, break into smaller tasks, or consider sunsetting

  **Low Time + High Completion** ⭐ (Efficient):
  - [List projects: Minimal effort, good results]
  - **Action**: Model this efficiency in other projects

  **Low Time + Low Completion** 📊 (Early stage or abandoned):
  - [List projects: Minimal effort, minimal results]
  - **Action**: Either commit more time or formally close

  ## Recommendations

  ### Optimize High-ROI Projects

  1. **[Project with best ROI]**: Continue current approach
     - Current: [X] tasks/hour
     - Why it works: [Efficient task breakdown / Clear scope / etc.]
     - Apply to: [List other projects that could benefit from same approach]

  ### Improve Low-ROI Projects

  2. **[Project with worst ROI]**: Break down large tasks
     - Current: [X] tasks/hour ([Y]h per task average)
     - Target: Reduce avg task time to [Z]h
     - How: [Specific actions based on project]

  3. **[Another low-ROI project]**: Re-evaluate scope
     - Current completion: [X]%
     - Action: [Specific recommendation]

  ### Time Reallocation

  4. **Shift [X] hours from [Low-ROI project] to [High-ROI project]**
     - Expected impact: [Y] more tasks completed per month
     - Rationale: [Based on ROI differential]

  ### Process Improvements

  5. **Improve time tracking adoption** (currently [X]%):
     - [Y]% of tasks have no time data
     - Better tracking enables better analysis
     - Action: Add LOGBOOK entries to ongoing tasks

  6. **Set time budgets per project**:
     - [Project 1]: Max [X]h/week
     - [Project 2]: Max [Y]h/week
     - Monitor via weekly reviews

  ### Project Health Actions

  7. **Close or revive stalled projects**:
     - [List projects with <20% completion and low recent activity]
     - Decision needed: Close formally or commit resources

  ## This Week Focus

  **Based on ROI analysis, recommended time allocation for next week**:

  | Project | Recommended Hours | Why |
  |---------|-------------------|-----|
  | [High-ROI project] | [X]h | Maximize efficient work |
  | [Strategic project] | [Y]h | Important despite low ROI |
  | [Problem project] | [Z]h | Fix or close decision |

  **Total recommended**: [X]h for focused project work

  ## Next Steps

  1. Review low-ROI projects: Use "deep dive on [project]" for detailed analysis
  2. Track this week: Compare actual time vs recommendations
  3. Next review: Run this analysis again in 1 month, compare trends

  Related skills:
  - "weekly review" - Track ongoing time allocation
  - "deep dive on [project]" - Analyze specific project in depth
  - "close project [name]" - Formally close low-value projects
  - "status report" - See current priorities alignment

  ---

  **Total context used**: ~4,250 tokens (time-tracking + timeline + tasks)
  **vs. manual analysis**: Would require reading all LOGBOOK entries across journals (~15,000+ tokens)
  **Efficiency**: ~3.5x faster with automated calculations

  Note: ROI calculations based on task completion as proxy for value.
  Some projects may have high value despite low task count (research, strategy, etc.).
  Use qualitative judgment alongside quantitative metrics.
