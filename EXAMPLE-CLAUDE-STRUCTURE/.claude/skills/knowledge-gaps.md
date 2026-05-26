---
metadata:
  patterns:
    - "knowledge gaps"
    - "missing pages"
    - "what should I document"
    - "documentation gaps"
    - "pages to create"
  description: "Identify and prioritize missing documentation using index analysis"
  auto_load_instructions: true
---

instructions: |
  Identify and prioritize missing documentation using cross-index analysis.

  Workflow:

  Step 1: Read Missing Pages Index
  1. Read .claude/indexes/missing-pages.md (300 tokens)
     Extract:
     - Total count of missing pages (pages with 5+ refs that don't exist)
     - Categorization: People vs Concepts
     - Each missing page with:
       - Name
       - Reference count
       - Pages that reference it

  Step 2: Read Reference Graph for Context
  2. Read .claude/indexes/reference-graph.md (5,000 tokens)
     - Find hub pages (most referenced pages that DO exist)
     - Compare missing pages to hub pages
     - Understand network importance of missing pages
     - Identify clusters (groups of related missing pages)

  Step 3: Read Tasks for Additional Context
  3. Read .claude/indexes/tasks-by-status.md (2,000 tokens)
     - grep for mentions of missing page names
     - Identify tasks that reference undocumented concepts/people
     - Count how many active tasks (NOW/TODO) need this documentation

  Step 4: Analyze Patterns
  4. For each missing page, determine:
     - **Urgency**: How many active tasks reference it?
     - **Impact**: How many total pages reference it?
     - **Type**: Person, Concept, Project, Tool, etc.
     - **Context**: Why is it being referenced? (infer from source pages)

  Step 5: Generate Prioritized Report

  # Knowledge Gaps Analysis

  **Analysis Date**: [Current Date]
  **Total Missing Pages**: [X] pages with 5+ references

  ## Executive Summary

  - **High urgency**: [X] pages referenced in active tasks
  - **High impact**: [X] pages with 10+ references
  - **People**: [X] contacts/collaborators without pages
  - **Concepts**: [X] topics/tools without documentation

  **Recommendation**: Create [top 3] pages in next 7 days

  ## High Priority (Create First)

  **Pages with 10+ references OR mentioned in active tasks**

  ### [Page Name 1] ([X] references)
  - **Type**: [Person / Concept / Project / Tool]
  - **Reference count**: [X] mentions across [Y] pages
  - **Urgency**: [HIGH/MEDIUM/LOW based on active tasks]
  - **Active tasks mentioning this**: [count NOW/TODO tasks]

  **Referenced from**:
  - [Page 1]
  - [Page 2]
  - [Page 3]
  [... show top 5-7 references]

  **Why create this page**:
  - [Infer from context: collaboration contact, key technology, ongoing project, etc.]

  **Suggested structure**:
  ```markdown
  # [Page Name]

  **Type**: [Person/Concept/Project/Tool]
  **Status**: [If person: Active contact / If concept: In use / etc.]

  ## Overview
  [Brief description]

  ## [Section based on type]
  [For Person: Role, Organization, Key Topics]
  [For Concept: What it is, How we use it, Key resources]
  [For Project: Status, Goals, Timeline]

  ## Related Pages
  [List connected pages from reference-graph]

  ## Notes
  [Ongoing notes and updates]
  ```

  **Estimated time to create**: [5-15 minutes depending on complexity]

  ---

  ### [Page Name 2] ([X] references)
  [Similar breakdown]

  ---

  ### [Page Name 3] ([X] references)
  [Similar breakdown]

  ---

  ## Medium Priority (Create Next)

  **Pages with 7-9 references**

  ### [Page Name] ([X] references)
  - **Type**: [Type]
  - **Referenced from**: [Top 3-5 pages]
  - **Context**: [Why this matters]
  - **Suggested structure**: [Brief template]

  [... repeat for all medium priority]

  ## Lower Priority (Create When Relevant)

  **Pages with 5-6 references**

  ### [Page Name] ([X] references)
  - **Type**: [Type]
  - **Context**: [Why referenced]
  - **Create when**: [Condition - e.g., "when working on related project" or "when contact becomes active again"]

  [... repeat for all lower priority]

  ## Patterns & Insights

  **Missing People** ([X] total):
  - [List names with ref counts]
  - **Pattern**: [Are these contacts from specific project? Time period? Domain?]
  - **Action**: Consider creating contact pages for active relationships

  **Missing Concepts** ([X] total):
  - [List concepts with ref counts]
  - **Pattern**: [Are these related to specific project? Technology stack? Research area?]
  - **Action**: Document concepts currently in use first

  **Clustering** (related missing pages):
  - Cluster 1: [Group of related pages, e.g., "NHS contacts"]
    - [Page 1] ([X] refs)
    - [Page 2] ([X] refs)
    - [Page 3] ([X] refs)
  - **Recommendation**: Create all pages in cluster together for coherent documentation

  ## Comparison to Hub Pages

  **Existing hub pages** (for reference):
  1. [Existing hub 1] - [X] references
  2. [Existing hub 2] - [X] references
  3. [Existing hub 3] - [X] references

  **Analysis**:
  - [Missing page with highest refs] ([X] refs) would rank #[Y] if it existed
  - This indicates [High/Medium/Low] network importance
  - Creating this page would [significantly/moderately] improve knowledge graph completeness

  ## Recommendations

  ### Immediate Actions (This Week)

  1. **Create [Page 1]** - Highest priority
     - Impact: [X] references, [Y] active tasks
     - Time: ~[Z] minutes
     - Start with: [Template section from above]

  2. **Create [Page 2]** - Second priority
     - Impact: [X] references
     - Time: ~[Z] minutes

  3. **Create [Page 3]** - Third priority or cluster lead
     - Impact: [X] references
     - Time: ~[Z] minutes

  ### Medium-Term (Next 30 Days)

  4. **Address [cluster name] cluster**
     - Create [X] related pages together
     - Estimated time: [X] hours total
     - Benefit: Complete documentation for [domain/project/etc.]

  5. **Document active contacts**
     - [List people pages with ongoing tasks/references]
     - Benefit: Relationship management and context preservation

  ### Process Improvements

  6. **Prevent future gaps**:
     - When starting new project → Create project page immediately (use "start project" skill)
     - When meeting new contact → Create contact page during/after meeting
     - When learning new concept → Document while fresh

  7. **Regular gap review**:
     - Run "knowledge gaps" skill monthly
     - Track missing pages that appear repeatedly
     - Set goal: Keep missing pages <5 at any time

  ## Quick Wins (5-Minute Pages)

  **Pages that can be created in <5 minutes** (simple structure, clear context):

  1. [Page name] - [Type] - [X] refs
     - Just need: [Minimal info required]
  2. [Page name] - [Type] - [X] refs
     - Just need: [Minimal info required]
  3. [Page name] - [Type] - [X] refs
     - Just need: [Minimal info required]

  **Batch creation recommended**: Spend 20 minutes, knock out 4-5 quick wins

  ## Next Steps

  After creating pages:
  - Run "knowledge gaps" again to see updated list
  - Use "status report" to see if active tasks now have better context
  - Check reference-graph.md to see new page connections

  Related skills:
  - Use "start project" when new project identified in gaps
  - Use "contact mapper" after creating person pages to build relationship timeline
  - Use "deep dive on [topic]" to understand context before creating complex concept pages

  ---

  **Total context used**: ~7,300 tokens (missing-pages + reference-graph + tasks)
  **vs. manual analysis**: Would require reading dozens of journal files to identify gaps
  **Efficiency**: Direct identification with prioritization built-in

  Note: This analysis runs automatically from indexes - no manual searching required.
  Creating the identified pages will improve knowledge graph completeness and reduce future re-explanation.
