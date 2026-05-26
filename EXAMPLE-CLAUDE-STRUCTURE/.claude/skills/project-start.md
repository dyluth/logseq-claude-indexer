---
metadata:
  patterns:
    - "start project"
    - "new project"
    - "create project"
    - "initialize project"
    - "begin project"
  description: "Initialize new project with structure and auto-update INDEX.md"
  auto_load_instructions: true
---

instructions: |
  Initialize a new project systematically with proper structure and documentation.

  When user requests to start a new project:

  Phase 1: Gather Project Details

  Ask user for required information:

  1. **Project name**: [Required]
     - Use for page name and references
     - Example: "Customer Analytics Platform" or "Research: AI Ethics"

  2. **Project purpose**: [Required]
     - 1-2 sentence description
     - Example: "Build analytics dashboard for customer behavior insights"

  3. **Project type**: [Required]
     - Options: Work, Research, Side Hustle, Learning, Client Project
     - Used for categorization

  4. **Expected duration**: [Optional]
     - Options: Days, Weeks, Months, Ongoing
     - Helps determine if dev-docs needed

  5. **Target completion**: [Optional]
     - Specific date if known
     - Example: "End of Q1 2026" or "December 2025"

  Phase 2: Create Project Structure

  1. Create main project page: `pages/[Project Name].md`

  Content template:
  ```markdown
  # [Project Name]

  **Status**: Active - Just Started
  **Created**: [Current Date]
  **Type**: [Project Type]
  **Expected Duration**: [Duration]
  **Target Completion**: [Date if provided]

  ## Overview

  [Project purpose from user input]

  ## Goals

  - [Ask user or leave as TODO]

  ## Current Status

  - Just initialized
  - Setting up initial structure

  ## Key Links

  - [Add as project progresses]

  ## Notes

  [Space for ongoing notes]
  ```

  2. If project is multi-week (expected duration = Weeks or Months):
     Create dev-docs structure in project page:

     Add section to project page:
     ```markdown
     ## Development Documentation

     For multi-week project tracking:
     - **Plan**: [High-level approach and decisions]
     - **Context**: [Current state, what's happening now]
     - **Tasks**: [Checklist of work items]

     See: .claude/guides/index-system.md (Dev Docs section)
     ```

     Note: Dev docs can be added to project page or separate files as needed

  Phase 3: Create Initial Tasks (Optional)

  Ask user: "Would you like to create initial tasks for this project?"

  If yes:
  - Ask for 3-5 initial tasks
  - Create tasks in today's journal: `journals/[YYYY_MM_DD].md`
  - Format:
    ```markdown
    - TODO [[Project Name]] - [Task description]
    ```

  Phase 4: Update INDEX.md

  1. Read current INDEX.md

  2. Find "## Active Projects" section

  3. Add new project entry (after existing projects, before "## Completed Projects"):

  ```markdown
  ### [Project Name] - [Project Type]
  **Overview:** [Project purpose]

  - **Main page**: [Link to pages/[Project Name].md]
  - **Status**: Active - Just started ([Current Date])
  - **Expected completion**: [Date if provided, or "TBD"]
  - **Recent activity**: Check [timeline-recent.md](.claude/indexes/timeline-recent.md) and filter for "[Project Name]"
  - **Time invested**: See [time-tracking.md](.claude/indexes/time-tracking.md)
  - **External context**: [Add if GitHub/Drive relevant, or "TBD"]

  ---
  ```

  4. Write updated INDEX.md

  Phase 5: Confirmation & Next Steps

  Generate summary for user:

  ```
  ✅ Project "[Project Name]" initialized successfully!

  Created:
  - Project page: pages/[Project Name].md
  - [If applicable] Initial tasks in journals/[Date].md
  - Updated INDEX.md (added to Active Projects)

  Next steps:
  1. Review project page: pages/[Project Name].md
  2. Add specific goals and key links as needed
  3. [If tasks created] Start working on initial tasks
  4. [If multi-week] Consider creating detailed dev-docs as project progresses

  Track progress:
  - Use "status report" to see overall project health
  - Use "deep dive on [Project Name]" for detailed analysis
  - Use "weekly review" to track across all projects

  Helpful commands:
  - "Add tasks for [Project Name]" - Create new tasks
  - "What's the status of [Project Name]?" - Project-specific status
  - "Time spent on [Project Name]" - Time investment analysis
  ```

  Example Usage:

  User: "start project Customer Analytics Platform"

  Claude:
  1. Asks for project details
  2. Creates pages/Customer Analytics Platform.md
  3. Creates initial tasks if requested
  4. Updates INDEX.md
  5. Confirms creation and provides next steps

  Notes:
  - Project name used in [[Page References]] for linking
  - INDEX.md automatically updated (no manual editing needed)
  - Dev-docs structure can be expanded later for complex projects
  - Initial tasks optional - can add them later via journals

  Total time: ~2 minutes for complete project initialization
  vs. manual setup: ~10-15 minutes + risk of forgetting INDEX.md update
