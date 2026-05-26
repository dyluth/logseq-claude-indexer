# Guides Directory - README

**Last Updated:** 2025-11-07
**Purpose:** Guide to the specialized documentation guides

---

## Overview

This directory contains **specialized reference guides** that keep CLAUDE.md concise while providing comprehensive guidance for the auto-generated index system, query optimization, cross-source integration, and automation.

**Navigation principle:** CLAUDE.md → Guides (you are here) → Indexes → Raw files

---

## Available Guides

### 📊 [index-system.md](./index-system.md)
**Complete reference for the 9 auto-generated indexes**

**When to read:**
- First time using this repository
- Need to understand which index for which query
- Want to learn cross-index patterns
- Troubleshooting index issues

**What's inside:**
- Complete inventory of all 9 indexes
- Progressive disclosure hierarchy (Levels 1-4)
- Detailed specs for each index (size, purpose, structure)
- Cross-index query patterns
- Token efficiency comparisons
- Troubleshooting FAQ

**Read this first if:** You're new to the index system

**Key sections:**
- Progressive Disclosure Hierarchy → How to navigate levels
- Complete Index Inventory → Detailed specs for each index
- Cross-Index Query Patterns → Combining indexes for comprehensive answers
- Token Efficiency Comparisons → See 5-30x efficiency gains

---

### 🎯 [query-patterns.md](./query-patterns.md)
**Decision trees for efficient query routing**

**When to read:**
- Want to optimize token usage
- Need query examples for common questions
- Learning which index for which query type
- Want to avoid common mistakes

**What's inside:**
- Decision trees for 20+ common queries
- Token usage comparisons (good vs bad)
- Optimization strategies
- Cross-index combination patterns
- Common mistakes and fixes

**Read this if:** You want to understand optimal query workflows

**Key sections:**
- Query Routing Decision Trees → Step-by-step workflows
- Status & Overview Queries → Most common use cases
- Token Efficiency Quick Reference → At-a-glance comparison table
- Common Mistakes & Fixes → Learn from bad patterns

---

### ⚡ [skills-reference.md](./skills-reference.md)
**Catalog of automation skills and creation guide**

**When to read:**
- Want to automate repetitive queries
- Need to understand available skills
- Planning to create custom skills
- Want workflow templates

**What's inside:**
- Complete catalog of 7 available skills
- Trigger patterns for each skill
- Token usage for each skill
- How to create custom skills
- Skill best practices
- Troubleshooting guide

**Read this if:** You want to automate common workflows

**Key sections:**
- Available Skills Catalog → What exists and how to use it
- How Skills Work → Activation flow explanation
- Creating Custom Skills → Step-by-step guide
- Skill Best Practices → Do's and don'ts

---

## Quick Navigation Flowchart

```
WHERE DO I START?

New to this system?
├─ Yes → Read index-system.md (Progressive Disclosure Hierarchy)
└─ No → Continue below

What do you need?

"How do I query efficiently?"
└─ query-patterns.md → Decision trees

"What automation exists?"
└─ skills-reference.md → Skills catalog

"Which index for [specific query]?"
└─ query-patterns.md → Query routing section

"How does [specific index] work?"
└─ index-system.md → Index inventory section

"I want to create automation"
└─ skills-reference.md → Creating custom skills

"Index not working as expected"
└─ index-system.md → Troubleshooting section
```

---

## Guide Interdependencies

**How guides reference each other:**

```
index-system.md
├─ References query-patterns.md for optimization examples
└─ References skills-reference.md for automation

query-patterns.md
├─ Relies on index-system.md for index specs
└─ References skills-reference.md for automation templates

skills-reference.md
├─ Built on foundation of index-system.md
├─ Uses patterns from query-patterns.md
```

**Reading order recommendation:**
1. **index-system.md** - Foundation
2. **query-patterns.md** - Optimization
4. **skills-reference.md** - Automation

---

## Common Use Cases

### Use Case: "I'm new, where do I start?"

**Path:**
1. Read [index-system.md](./index-system.md) - Progressive Disclosure Hierarchy section
2. Skim [query-patterns.md](./query-patterns.md) - Status & Overview Queries section
3. Try a query: "What's my status?"
4. Observe which indexes Claude reads
5. Return to guides as needed for deeper understanding

---

### Use Case: "My queries are using too many tokens"

**Path:**
1. Read [query-patterns.md](./query-patterns.md) - Common Mistakes & Fixes
2. Read [index-system.md](./index-system.md) - Progressive Disclosure Hierarchy
3. Review [query-patterns.md](./query-patterns.md) - Token Efficiency Quick Reference
4. Apply progressive disclosure (Level 1 → Level 2 → Level 3)

---

### Use Case: "I repeat the same query weekly"

**Path:**
1. Read [skills-reference.md](./skills-reference.md) - Available Skills Catalog
2. Check if skill exists (e.g., "weekly-review.md")
3. If exists: Use trigger pattern
4. If doesn't exist: Read [skills-reference.md](./skills-reference.md) - Creating Custom Skills
5. Build custom skill following template

---

### Use Case: "I want to close a project cleanly"

**Path:**
1. Read [skills-reference.md](./skills-reference.md) - project-closure.md section
2. Use trigger: "Close project [name]"
3. Review generated closure report
4. Execute cleanup actions

---

## Maintenance Notes

### When to Update Guides

**index-system.md:**
- When logseq-claude-indexer adds/removes indexes
- When index structure changes significantly
- When new cross-index patterns discovered

**query-patterns.md:**
- When new query patterns emerge
- When token efficiency metrics change
- When new common mistakes identified

**skills-reference.md:**
- When new skills created
- When skill trigger patterns change
- When skill best practices evolve

### Last Updated Dates

Check "Last Updated" at top of each guide:
- index-system.md: 2025-11-07
- query-patterns.md: 2025-11-07
- skills-reference.md: 2025-11-07

---

## Related Documentation

- **Main Entry Point:** [CLAUDE.md](../CLAUDE.md) - Repository overview
- **Navigation Hub:** [INDEX.md](../../INDEX.md) - Central navigation
- **Design Document:** [claude.md-optimisation-design.md](../../claude.md-optimisation-design.md) - Implementation plan
- **Auto-Generated Indexes:** [.claude/indexes/](../indexes/) - 9 index files
- **Skills:** [.claude/skills/](../skills/) - Automation workflows

---

## Quick Reference Table

| Guide | Size | Purpose | Read When |
|-------|------|---------|-----------|
| [index-system.md](./index-system.md) | ~400 lines | Complete index reference | First time, need specs |
| [query-patterns.md](./query-patterns.md) | ~330 lines | Query optimization | Need efficiency, examples |
| [skills-reference.md](./skills-reference.md) | ~210 lines | Skills catalog & creation | Want automation |
| **This file (README.md)** | ~100 lines | Guide navigation | Finding right guide |

**Total documentation:** ~1,300 lines across 5 files

**CLAUDE.md size:** ~280 lines (references these guides)

**Benefit:** Comprehensive documentation with concise entry point

---

## Support

**Issues with guides:**
- Outdated information → Check "Last Updated" date
- Broken internal links → Verify file paths
- Missing information → Suggest additions

**Issues with indexes:**
- See [index-system.md](./index-system.md) - Troubleshooting section

**Issues with skills:**
- See [skills-reference.md](./skills-reference.md) - Troubleshooting Skills section

---

**Guide to Guides Complete**

Start with [index-system.md](./index-system.md) if you're new, or jump to the specific guide you need using the navigation above.
