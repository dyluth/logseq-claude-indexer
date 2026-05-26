---
metadata:
  patterns:
    - "deep dive on"
    - "deep dive"
    - "full context for"
    - "everything about"
    - "comprehensive analysis of"
  description: "Comprehensive project analysis across all sources (Logseq + GitHub + Google Drive)"
  auto_load_instructions: true
---

instructions: |
  Perform comprehensive project analysis combining internal (Logseq) + external (GitHub/Drive) sources.

  When user requests deep dive on [Project Name]:

  Phase 1: Index Reconnaissance

  1. Read .claude/indexes/dashboard.md (500 tokens)
     - Check if [Project Name] in "Top Projects" section
     - Note time investment if present
     - Check if mentioned in current priorities

  2. Read .claude/indexes/tasks-by-status.md (2,000 tokens)
     - grep for [Project Name] across all status sections
     - Extract ALL tasks mentioning project:
       - NOW: [list]
       - DOING: [list]
       - TODO: [list]
       - LATER: [list]
       - DONE: [list]
     - Count total tasks per status
     - Calculate completion rate: DONE / Total

  3. Read .claude/indexes/time-tracking.md (1,500 tokens)
     - Find [Project Name] in "Top Projects"
     - Extract:
       - Total time invested
       - Number of tasks tracked
       - Average time per task
       - Weekly breakdown if available
     - If not in top projects: "Minimal time tracking"

  4. Read .claude/indexes/reference-graph.md (5,000 tokens)
     - Find [Project Name] in hub pages or page details
     - Count inbound references (pages linking to project)
     - Count outbound references (pages project links to)
     - Identify connected pages for additional context

  5. Read .claude/indexes/timeline-recent.md (350 tokens)
     - Filter for [Project Name] mentions in last 7 days
     - Note recent activity pattern
     - Identify last activity date

  Phase 2: Deep Context Reading

  6. Read project main page: pages/[Project Name].md
     - Full project overview
     - Stated goals and objectives
     - Current status notes
     - Key links and references

  7. Read connected pages (from reference-graph, top 5):
     - [Connected Page 1].md
     - [Connected Page 2].md
     - [Connected Page 3].md
     - Extract relevant context about project

  8. Read recent journal entries mentioning project:
     - From timeline-recent, get dates with [Project Name]
     - Read those specific journal files
     - Extract: Decision notes, progress updates, challenges

  Phase 3: Synthesis & Comprehensive Report

  Generate detailed analysis report:

  # Deep Dive: [Project Name]

  **Analysis Date**: [Current Date]
  **Last Updated**: [Date from timeline or project page]

  ## Executive Summary

  [2-3 sentence overview of project current state]

  - **Status**: [Active/On Hold/Completed based on task distribution]
  - **Health**: [Healthy/Needs Attention/Critical based on analysis]
  - **Momentum**: [Increasing/Steady/Declining based on timeline]

  ## Project Overview

  [From project page and dashboard]

  **Purpose**: [Project description]
  **Type**: [Work/Research/Side Project/etc.]
  **Duration**: [From timeline-full if available, or "Ongoing"]

  ## Current Status

  ### Task Distribution

  | Status | Count | % of Total |
  |--------|-------|------------|
  | NOW | [X] | [%] |
  | DOING | [X] | [%] |
  | TODO | [X] | [%] |
  | LATER | [X] | [%] |
  | DONE | [X] | [%] |
  | **Total** | **[X]** | **100%** |

  **Completion Rate**: [X]% ([DONE] / [Total])

  **Active Tasks** (NOW):
  - [Task 1]
  - [Task 2]
  - [Task 3]
  [... list all NOW tasks]

  **Backlog** (TODO):
  - [Count] tasks in backlog
  - [Mention any high-priority items]

  ### Time Investment

  - **Total time logged**: [X] hours across [Y] tasks
  - **Average per task**: [Z] hours
  - **Time tracking adoption**: [%] ([Tracked] / [Total tasks])

  **Weekly time allocation** (if available from time-tracking):
  - This week: [X]h
  - Last 4 weeks average: [Y]h
  - Trend: [Increasing/Decreasing/Steady]

  ### Recent Activity

  [From timeline-recent]

  **Last 7 days**:
  - [Date]: [Activity description]
  - [Date]: [Activity description]
  - [Date]: [Activity description]

  **Last activity**: [Date] - [What happened]
  **Activity frequency**: [X] days with activity in last 7 days

  ## Network Analysis

  **References across knowledge base**:
  - **Inbound references**: [X] pages link to this project
  - **Outbound references**: This project links to [Y] pages
  - **Hub rank**: [If in top 10 hub pages, mention rank]

  **Connected pages** (top 5):
  1. [Page 1] - [Context: why connected]
  2. [Page 2] - [Context: why connected]
  3. [Page 3] - [Context: why connected]
  4. [Page 4] - [Context: why connected]
  5. [Page 5] - [Context: why connected]

  ## External Context

  [If external sources available]

  ### GitHub Repository

  **Repository**: [URL]
  **Last commit**: [Date from recent commits]
  **Recent activity**: [X] commits in last 7 days

  **Technical Overview**:
  - [From README: architecture, tech stack, features]

  **Development Status**:
  - [From commits: what's being worked on]

  **Alignment Check**:
  - Internal tasks ↔ External commits: [Aligned/Misaligned]
  - [Note any discrepancies]

  ### Strategy Documents (Google Drive)

  **Documents found**:
  - [Document 1 name]
  - [Document 2 name]

  **Strategic Direction**:
  - [From strategy docs: target market, positioning, goals]

  **Alignment Check**:
  - Internal tasks ↔ Strategy docs: [Aligned/Misaligned]
  - [Note any discrepancies]

  ## Key Insights

  [Synthesize from all sources]

  **Strengths**:
  1. [Insight from analysis - high completion rate, good momentum, etc.]
  2. [Insight from external alignment]
  3. [Insight from network analysis]

  **Challenges**:
  1. [Insight from blockers, stalled tasks, etc.]
  2. [Insight from misalignment between sources]
  3. [Insight from time allocation or activity patterns]

  **Opportunities**:
  1. [Based on connected pages, external context, or gaps]
  2. [Based on TODO backlog or LATER tasks]

  ## Decision History

  [From journal entries and project page]

  **Key decisions made**:
  - [Date]: [Decision from journal or project notes]
  - [Date]: [Decision]
  - [Date]: [Decision]

  **Open questions**:
  - [Any unresolved items from tasks or notes]

  ## Recommendations

  [Provide 3-5 specific, actionable recommendations]

  ### Immediate Actions (Next 7 Days)

  1. **[Recommendation 1]**: [Based on NOW tasks or blockers]
     - Why: [Rationale from analysis]
     - Action: [Specific next step]

  2. **[Recommendation 2]**: [Based on external alignment or strategy]
     - Why: [Rationale from analysis]
     - Action: [Specific next step]

  ### Medium-Term Focus (Next 30 Days)

  3. **[Recommendation 3]**: [Based on backlog or opportunities]
     - Why: [Rationale from analysis]
     - Action: [Specific next step]

  4. **[Recommendation 4]**: [Based on time allocation or velocity]
     - Why: [Rationale from analysis]
     - Action: [Specific next step]

  ### Strategic Considerations

  5. **[Recommendation 5]**: [Based on network analysis or external context]
     - Why: [Rationale from analysis]
     - Action: [Specific next step]

  ## Next Steps

  **Immediate focus**:
  - [Top 3 NOW tasks to prioritize]

  **Resources needed**:
  - [Any blockers requiring external resources]

  **Follow-up queries**:
  - Use "time ROI analysis" to assess time investment effectiveness
  - Use "weekly review" to track progress against these recommendations
  - Use "status report" for high-level project summary

  ---

  **Total context used**: ~3,000 tokens (indexes) + external reads + connected pages
  **vs. reading all journals**: ~10,000+ tokens
  **Efficiency**: ~3x faster with comprehensive cross-source intelligence

  Example Next Steps:
  - If misalignment found → Update strategy docs or adjust tasks
  - If blockers identified → Schedule time to resolve
  - If momentum low → Re-evaluate priorities or resources
  - If completion rate low → Break large tasks into smaller chunks

  This deep dive provides complete picture from:
  - Internal status (Logseq indexes + pages)
  - External technical (GitHub)
  - External strategic (Google Drive)
  - Network context (reference graph)
  - Historical timeline (activity patterns)
