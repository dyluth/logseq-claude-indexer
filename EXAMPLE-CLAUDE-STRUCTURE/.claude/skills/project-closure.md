---
metadata:
  patterns:
    - "close project"
    - "archive project"
    - "project closure"
    - "finish project"
    - "complete project"
  description: "Close/archive project with comprehensive summary and auto-update INDEX.md"
  auto_load_instructions: true
---

instructions: |
  Close a project systematically using cross-index analysis with comprehensive closure report.

  When user requests closure of [Project Name]:

  Phase 1: Discovery (Cross-Index Analysis)

  1. Read .claude/indexes/reference-graph.md (5,000 tokens)
     - Search for [Project Name] in hub pages or page details
     - Count inbound references (how many pages link to this project)
     - Identify all pages referencing this project
     - Note if project page exists

  2. Read .claude/indexes/tasks-by-status.md (2,000 tokens)
     - grep for [Project Name] across all status sections
     - Count tasks by status:
       - NOW: [count]
       - DOING: [count]
       - TODO: [count]
       - LATER: [count]
       - DONE: [count]
     - List any non-DONE tasks (these need decisions)

  3. Read .claude/indexes/time-tracking.md (1,500 tokens)
     - Find [Project Name] in "Top Projects" section
     - Extract:
       - Total time invested
       - Number of tasks tracked
       - Average time per task
     - If not in top projects, note: "Time tracking not significant"

  4. Read .claude/indexes/timeline-full.md (5,000 tokens)
     - Search for first mention of [Project Name] (project start date)
     - Search for last activity with [Project Name] (project end date)
     - Calculate project duration (days between start and end)
     - Note: Most productive periods if visible

  5. Read project page (if exists): pages/[Project Name].md
     - Understand project overview
     - Check for any existing closure notes
     - Extract key outcomes/deliverables

  Phase 2: Closure Report Generation

  Generate comprehensive closure report:

  ```markdown
  # Project Closure: [Project Name]

  **Closed**: [Current Date]
  **Status**: Archived

  ## Summary

  - **Duration**: [Start date] to [End date] ([X] days)
  - **Total time invested**: [X] hours across [Y] tasks
  - **Average time/task**: [X] hours
  - **Completion rate**: [Z]% ([DONE count] / [Total tasks])
  - **Time tracking adoption**: [X]% ([Tracked tasks] / [Total tasks])

  ## References Across Knowledge Base

  - **Total mentions**: [X] references across [Y] pages
  - **Hub strength**: [If in top 10 hub pages, mention rank; otherwise "Not a hub page"]
  - **Connected pages**: [List from reference-graph]
    - [Page 1]
    - [Page 2]
    - [Page 3]
    - [... top 5-10 pages]

  ## Task Summary

  ### Completed Tasks ([X] total)
  - [DONE count] tasks completed successfully

  ### Incomplete Tasks ([Y] total)

  #### NOW Tasks ([count])
  [List each NOW task with recommendation]
  - [Task description] → **Recommended action**: CLOSE (project ended) / MIGRATE to [Other Project]
  - [Task description] → **Recommended action**: [...]

  #### TODO Tasks ([count])
  [List each TODO task with recommendation]
  - [Task description] → **Recommended action**: CLOSE (not critical) / MIGRATE to [Other Project]
  - [Task description] → **Recommended action**: [...]

  #### DOING Tasks ([count])
  [If any DOING tasks, these are concerning - project shouldn't close with in-progress work]
  - [Task description] → **Recommended action**: COMPLETE or MIGRATE

  #### LATER Tasks ([count])
  [List if any]
  - [Task description] → **Recommended action**: CLOSE (deferred indefinitely)

  ## Activity Timeline

  - **First activity**: [Date] - [First task or mention from timeline]
  - **Peak activity**: [If identifiable from timeline, mention period with most tasks]
  - **Last activity**: [Date] - [Last task or mention from timeline]
  - **Project lifespan**: [X] days from first to last activity

  ## Key Outcomes

  [Scan DONE tasks and project page for major accomplishments]

  Completed deliverables:
  - [Outcome 1 from DONE tasks]
  - [Outcome 2 from DONE tasks]
  - [Outcome 3 from project page if mentioned]

  [If project page has specific outcomes listed, include those]

  ## Learnings & Context

  [Extract from project page if available, otherwise note from tasks]

  [Any notes from project page about lessons learned, challenges overcome, etc.]

  [If no explicit learnings documented, infer from task patterns:
   - High completion rate → Well-scoped project
   - Long duration → Complex or evolving scope
   - High time investment → Resource-intensive work]

  ## External Context

  [If project has GitHub or Google Drive references, note them]

  - GitHub: [If HOLT or other repo mentioned]
  - Google Drive: [If strategy docs or other files mentioned]

  ```

  Phase 3: Cleanup Actions (Ask User Before Executing)

  Present proposed cleanup actions to user:

  ```
  Proposed cleanup actions for [Project Name]:

  1. Update project page (pages/[Project Name].md):
     - Add closure date to header: "**Status**: Archived ([Current Date])"
     - Append this closure report to bottom of page

  2. Update incomplete tasks ([X] tasks):
     [For each incomplete task, show recommendation]
     - [Task in file:line] → Action: CLOSE / MIGRATE to [Project]
     - [Task in file:line] → Action: CLOSE / MIGRATE to [Project]

     Would you like to:
     a) CLOSE all incomplete tasks (mark as DONE with note "Project closed")
     b) Review each task individually
     c) Skip task cleanup (leave as-is)

  3. Update INDEX.md:
     - Remove [Project Name] from "Active Projects" section
     - Add [Project Name] to "Completed Projects" section with summary

  4. Optional actions:
     - Create archive file: pages/archive/[Project Name]-closure-report.md
     - Add "Closure" tag to all related journal entries (can be done later)

  Proceed with cleanup? (Options: all / selective / skip)
  ```

  If user chooses "all":
  - Execute all cleanup actions automatically

  If user chooses "selective":
  - Ask which actions to perform (1-4, comma-separated)

  If user chooses "skip":
  - Provide closure report only, no cleanup

  Phase 4: Execute Cleanup (if approved)

  Action 1: Update Project Page

  1. Read pages/[Project Name].md
  2. Find status line (usually in header)
  3. Update status: "**Status**: Archived ([Current Date])"
  4. Append closure report to end of file
  5. Write updated file

  Action 2: Handle Incomplete Tasks (based on user choice)

  For each incomplete task:
  - If CLOSE: Update task marker to DONE, add note "(Project closed: [Date])"
  - If MIGRATE: Update task to reference new project, add note "(Migrated from [Old Project])"

  Action 3: Update INDEX.md

  1. Read INDEX.md
  2. Find "## Active Projects" section
  3. Locate [Project Name] entry
  4. Cut entire entry (from ### to next --- or next ###)
  5. Find "## Completed Projects" section
  6. Paste entry with modifications:

  ```markdown
  ### [Project Name] - [Type]
  **Overview:** [Original description]

  - **Main page**: [Link]
  - **Status**: ⚠️ Completed - Archived on [Current Date]
  - **Duration**: [Start] to [End] ([X] days)
  - **Time invested**: [X] hours across [Y] tasks (from time-tracking.md)
  - **Completion rate**: [Z]%
  - **References**: [X] mentions across knowledge base (from reference-graph.md)
  - **Closure report**: See project page for full report

  ---
  ```

  7. Write updated INDEX.md

  Action 4: Create Archive (if requested)

  - Create pages/archive/ directory if doesn't exist
  - Write closure report to pages/archive/[Project Name]-closure-report-[Date].md

  Phase 5: Validation

  After cleanup (if performed):

  1. Verify project moved in INDEX.md:
     - Check "Active Projects" section → [Project Name] should be gone
     - Check "Completed Projects" section → [Project Name] should be present

  2. Verify closure report added to project page:
     - Read pages/[Project Name].md
     - Confirm status = "Archived"
     - Confirm closure report appended

  3. Verify tasks updated (if cleanup chosen):
     - Check that incomplete tasks are CLOSED or MIGRATED as specified

  4. Check cross-index consistency:
     - New tasks created today should NOT reference closed project (unless migration)

  Generate final summary:

  ```
  ✅ Project "[Project Name]" closed successfully!

  Summary:
  - Duration: [X] days
  - Time invested: [X] hours
  - Completion rate: [X]%
  - References: [X] across [Y] pages

  Changes made:
  ✓ Project page updated with closure report
  ✓ [X] incomplete tasks handled ([CLOSED/MIGRATED])
  ✓ INDEX.md updated (moved to Completed Projects)
  [✓ Archive created] (if applicable)

  Next steps:
  - Project is now archived and tracked in INDEX.md
  - Closure report available in project page
  - Knowledge base references preserved for future reference

  Related projects still active:
  [List any projects mentioned in task migrations or related pages]
  ```

  Example Usage:

  User: "close project Unitary"

  Claude:
  1. Discovers: 14 refs, 5 tasks (2 NOW, 3 TODO), 2,440h invested
  2. Generates closure report with recommendations
  3. Proposes cleanup actions
  4. User approves
  5. Updates project page, tasks, INDEX.md
  6. Validates all changes
  7. Confirms completion

  Total context: ~13,500 tokens for discovery + ~1,000 for updates
  vs. manual: Would require reading dozens of journal files + manual INDEX.md editing
  Efficiency: Complete automation with comprehensive documentation

  Note: This skill combines ALL index sources for thorough project closure
