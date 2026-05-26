---
metadata:
  patterns:
    - "status report"
    - "what's happening"
    - "what's my status"
    - "executive summary"
    - "current state"
    - "project status"
  description: "Generate executive status report using index system"
  auto_load_instructions: true
---

instructions: |
  Generate executive status report using progressive disclosure.

  Reference: .claude/guides/index-system.md for complete index usage
  Reference: .claude/guides/query-patterns.md for query optimization

  Workflow:

  Step 1: Read Dashboard (Master Overview)
  1. Read .claude/indexes/dashboard.md (500 tokens)
     Extract:
     - Quick stats (total tasks, completion rate, time tracking adoption)
     - Current priorities ([#A] tasks, top 2 shown)
     - Top projects (by time investment, top 5)
     - Pages to create (knowledge gaps, top 5)

  Step 2: Read Recent Activity
  2. Read .claude/indexes/timeline-recent.md (350 tokens)
     Extract:
     - Last 7 days activity summary
     - Task counts by day and status
     - Completed tasks (DONE)
     - Active work days

  Step 3: Read Active Tasks (NOW section only)
  3. Read .claude/indexes/tasks-by-status.md (NOW section only, ~500 tokens)
     Extract:
     - Currently active tasks (48 NOW tasks total)
     - Identify tasks >7 days old (potential blockers)
     - Group by project if possible

  Step 4: Generate Formatted Report

  # Status Report - [Current Date]

  ## Overview
  - **Total tasks**: [from dashboard]
  - **Completion rate**: [from dashboard]
  - **Active projects**: [count from top projects]
  - **Time tracking adoption**: [from dashboard]

  ## Current Focus (NOW Tasks)

  [Organize by project if identifiable, otherwise list top 5-10 NOW tasks]

  ### [Project Name]
  - [NOW task 1]
  - [NOW task 2]

  ### [Project Name]
  - [NOW task 1]

  ### General
  - [Other NOW tasks]

  **Blockers** (tasks >7 days in NOW):
  - [List any tasks that have been NOW for >7 days]

  ## Recent Wins (Last 7 Days)

  [List DONE tasks from timeline-recent, organized by date or project]

  - [Date or Project]: [Completed task 1]
  - [Date or Project]: [Completed task 2]

  ## Project Health

  ### [Project 1 Name] (e.g., Test Company)
  - **Time invested**: [from time-tracking via dashboard]
  - **Active tasks**: [count NOW/TODO tasks mentioning project]
  - **Recent activity**: [from timeline-recent]
  - **Status**: [synthesize from tasks + context]
  - **External alignment**: [if Google Drive/GitHub checked]

  ### [Project 2 Name] (e.g., MoonLaunch)
  - **Time invested**: [from time-tracking via dashboard]
  - **Active tasks**: [count NOW/TODO tasks mentioning project]
  - **Recent activity**: [from timeline-recent]
  - **Status**: [synthesize all sources]

  ### [Project 3 Name]
  - [Similar breakdown]

  ## Next Priorities

  [From tasks-by-priority or dashboard current priorities]

  **High priority items** ([#A] tasks):
  - [List from dashboard or tasks-by-priority]

  ## Knowledge Gaps

  [If dashboard shows significant missing pages, mention top 2-3]

  Pages to create (5+ references):
  - [Page name] ([X] refs) - [Type: Person/Concept]

  ## Recommendations

  [Based on analysis, provide 2-3 actionable recommendations]

  - [Recommendation based on blockers, priorities, or project health]
  - [Recommendation for knowledge gaps if significant]

  ---

  **Context used**: ~1,850 tokens (dashboard + timeline + NOW tasks)
  **vs. without indexes**: ~8,000+ tokens reading all task sources
  **Efficiency**: ~4-5x faster

  Note: This report combines internal (Logseq) + external (GitHub/Drive) sources for comprehensive intelligence.
