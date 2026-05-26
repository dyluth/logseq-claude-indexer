---
metadata:
  patterns:
    - "relationship with"
    - "interactions with"
    - "contact history"
    - "who is"
    - "tell me about"
  description: "Map relationship timeline with contact using cross-index analysis"
  auto_load_instructions: true
---

instructions: |
  Map comprehensive relationship timeline with a contact using cross-index and external source analysis.

  When user requests relationship mapping for [Person Name]:

  Phase 1: Check if Page Exists

  1. Read .claude/indexes/missing-pages.md (300 tokens)
     - Check if [Person Name] in missing pages list
     - If yes: Note reference count and categorization (Person/Concept)
     - This indicates: Page should be created (high reference count = important relationship)

  2. Read .claude/indexes/reference-graph.md (5,000 tokens)
     - Search for [Person Name] in page details or hub pages
     - If page exists:
       - Count inbound references (mentions of this person)
       - Count outbound references (topics person is linked to)
       - Identify connected pages
     - If page doesn't exist:
       - Count how many pages reference this person
       - Identify all source pages

  Phase 2: Task Analysis

  3. Read .claude/indexes/task-index.md (8,000 tokens)
     - grep for [Person Name] across all tasks
     - Extract tasks mentioning person:
       - NOW: [list]
       - DOING: [list]
       - TODO: [list]
       - LATER: [list]
       - DONE: [list]
     - For each task, note:
       - Task description
       - File location (journal date)
       - Status
       - Context (what was discussed/planned)

  Phase 3: Timeline Construction

  4. Read .claude/indexes/timeline-full.md (5,000 tokens)
     - Search for all dates mentioning [Person Name]
     - Create chronological list:
       - First mention (relationship start)
       - All interaction dates
       - Most recent mention (last interaction)
     - Calculate:
       - Total days of relationship
       - Frequency (mentions per month)
       - Activity pattern (regular vs sporadic)

  Phase 4: Context from Journals

  5. Based on dates from timeline, read specific journal entries:
     - Read journals/[date of first mention].md
     - Read journals/[date of most recent mention].md
     - Read journals/[2-3 key interaction dates].md (dates with most tasks/activity)

     Extract:
     - How met / context of relationship
     - Key topics discussed
     - Decisions made together
     - Commitments or follow-ups
     - Person's role, organization, interests

  Phase 5: Generate Relationship Map

  # Relationship Map: [Person Name]

  **Analysis Date**: [Current Date]
  **Page Status**: [Exists at pages/[Person Name].md / **Should be created** - [X] references]

  ## Overview

  - **First interaction**: [Date from timeline]
  - **Most recent interaction**: [Date from timeline]
  - **Relationship duration**: [X] days / [Y] months
  - **Total mentions**: [X] references across [Y] pages/journal entries
  - **Interaction frequency**: ~[X] mentions per month

  ## Contact Information

  [If available from journal entries or existing page]

  - **Role**: [Title/Position if mentioned]
  - **Organization**: [Company/Institution if mentioned]
  - **Email**: [If mentioned]
  - **LinkedIn**: [If mentioned]
  - **Other**: [Phone, website, etc. if mentioned]

  ## Relationship Timeline

  **Key Interactions** (chronological):

  ### [Date - First Interaction]
  - **Context**: [How met / why connected]
  - **Topics**: [What was discussed]
  - **Source**: [Journal entry or page]

  ### [Date - Significant Interaction]
  - **Context**: [What happened]
  - **Topics**: [Discussion points]
  - **Outcomes**: [Decisions, commitments]
  - **Source**: [Journal entry]

  ### [Date - Another Interaction]
  - **Context**: [What happened]
  - **Topics**: [Discussion points]
  - **Source**: [Journal entry]

  [... list all significant interactions]

  ### [Date - Most Recent Interaction]
  - **Context**: [What happened]
  - **Topics**: [Discussion points]
  - **Follow-up**: [Any pending actions]
  - **Source**: [Journal entry]

  ## Associated Projects & Topics

  **Projects involving [Person Name]**:
  - [Project 1] - [Context: role in project, collaboration type]
  - [Project 2] - [Context]

  **Topics of discussion** (from tasks and journals):
  - [Topic 1] - [Frequency: how often discussed]
  - [Topic 2] - [Frequency]
  - [Topic 3] - [Frequency]
  - [Topic 4] - [Frequency]

  **Common themes**:
  - [Pattern 1: e.g., "Technical consultations on healthcare AI"]
  - [Pattern 2: e.g., "Networking and introductions"]

  ## Open Actions & Commitments

  **Active tasks** (NOW/TODO involving [Person Name]):
  - [Task 1] - Status: [NOW/TODO]
    - Context: [What needs to be done]
    - Deadline: [If mentioned]

  - [Task 2] - Status: [NOW/TODO]
    - Context: [What needs to be done]

  **Pending follow-ups** (from recent interactions):
  - [Follow-up 1] - [Based on last interaction]
  - [Follow-up 2] - [Based on commitments]

  **Waiting on** (if any tasks blocked by this person):
  - [Item 1] - Since: [Date]
  - [Item 2] - Since: [Date]

  ## Completed Collaborations

  **DONE tasks** (involving [Person Name]):
  - [Date]: [Completed task 1]
  - [Date]: [Completed task 2]
  - [Date]: [Completed task 3]

  **Achievements**:
  - [Outcome 1 from collaboration]
  - [Outcome 2 from collaboration]

  ## Relationship Context

  **Person's interests** (inferred from interactions):
  - [Interest 1]
  - [Interest 2]
  - [Interest 3]

  **Person's expertise** (based on discussions):
  - [Area 1]
  - [Area 2]

  **Value of relationship**:
  - [Why this relationship matters: mentor, client, collaborator, expert in X, etc.]
  - [Mutual benefits or exchange]

  **Relationship stage**:
  - [Initial contact / Active collaboration / Occasional touch-base / Past relationship]

  ## External Context

  [If found in Google Drive]

  **Documents involving [Person Name]**:
  - [Document 1 name] - [Type: meeting notes, email, etc.]
  - [Document 2 name] - [Type]

  **Additional context from documents**:
  - [Key insight 1 from external docs]
  - [Key insight 2]

  ## Network Connections

  **Connected pages** (from reference-graph):
  - [Page 1] - [Relationship: both mention project X]
  - [Page 2] - [Relationship: introduced by, works with, etc.]
  - [Page 3] - [Relationship]

  **Potential introductions** (people this person might know):
  - [Person 2] - [Common context]
  - [Person 3] - [Common context]

  ## Recommended Next Steps

  **Immediate actions**:
  1. [Based on open tasks or pending follow-ups]
  2. [Based on relationship maintenance: reach out if no recent contact]
  3. [Based on opportunities: introduction, collaboration, etc.]

  **Medium-term**:
  - [Action based on project involvement]
  - [Action based on relationship development]

  **Documentation**:
  [If page doesn't exist]
  - **Create page**: pages/[Person Name].md
  - **Use template**: See structure below
  - **Populate with**: Information from this analysis

  [If page exists]
  - **Update page**: Add recent interactions, update status
  - **Review**: Ensure page reflects current relationship state

  ## Suggested Page Structure

  [If page needs to be created]

  ```markdown
  # [Person Name]

  **Role**: [Title/Position]
  **Organization**: [Company/Institution]
  **Relationship**: [Mentor / Client / Collaborator / etc.]
  **Status**: [Active / Occasional contact / Past]

  ## Contact Information
  - Email: [Email]
  - LinkedIn: [URL]
  - [Other contact methods]

  ## Overview
  [Brief description: who they are, why connected, value of relationship]

  ## Key Topics
  - [Topic 1]
  - [Topic 2]
  - [Topic 3]

  ## Interaction History
  **First met**: [Date] - [Context]
  **Recent interactions**:
  - [Date]: [What happened]
  - [Date]: [What happened]

  ## Projects Together
  - [[Project 1]] - [Role/involvement]
  - [[Project 2]] - [Role/involvement]

  ## Open Items
  - [Pending action 1]
  - [Pending action 2]

  ## Notes
  [Ongoing notes and updates from interactions]
  ```

  ## Activity Pattern

  **Interaction frequency over time**:
  - [Period 1, e.g., "First 3 months"]: [X] mentions ([Y] per month) - [Pattern: ramping up, intensive collaboration, etc.]
  - [Period 2]: [X] mentions ([Y] per month) - [Pattern]
  - [Recent period]: [X] mentions ([Y] per month) - [Pattern]

  **Current status**:
  - [Active: regular interaction / Dormant: no recent activity / Occasional: sporadic touch-base]
  - [If dormant]: Last contact [X] days ago - Consider: [Reach out / Let lapse / etc.]

  ---

  **Total context used**: ~18,000 tokens (missing-pages + reference-graph + task-index + timeline + 3-5 journals + Drive)
  **vs. manual**: Would require extensive journal searching and mental timeline reconstruction
  **Efficiency**: Comprehensive relationship view with precise timeline

  Next actions after mapping:
  - If page missing: Create using suggested structure
  - If pending actions: Add to task list
  - If relationship dormant: Decide on re-engagement or archive
  - Use "status report" to see this person in context of all work
  - Use "deep dive on [project]" if person involved in specific project

  Example usage:
  User: "Relationship with John smith"
  → Discovers: 11 references, no page, involved in MoonLaunch discussions, last contact 10 days ago
  → Generates: Complete timeline, open actions, recommendation to create page + follow up

  This skill provides relationship intelligence for:
  - Networking and outreach prep
  - Context before meetings
  - Relationship maintenance decisions
  - Identifying collaboration opportunities
  - Tracking commitments and follow-ups
