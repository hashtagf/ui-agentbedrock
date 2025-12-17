# State Management - Branch-Based System

## Overview

AIDLC uses a **branch-based state system** that tracks workflow progress per Git branch. This aligns with the branch-based audit system and provides independent progress tracking for each feature/requirement.

---

## Directory Structure

```text
aidlc-docs/
├── state/
│   ├── state-index.md              # Master index of all state files
│   ├── main.md                     # State for main/master branch
│   ├── {sanitized-branch-name}.md  # State for each feature branch
│   └── archived/                   # Merged/completed branch states
│       └── {sanitized-branch-name}.md
├── audit/
│   └── ... (branch-based audit files)
└── ... (other aidlc-docs)
```

---

## Branch Detection Logic

Uses the **same logic as audit-management.md**:

### Step 1: Get Current Branch
```bash
git branch --show-current
```

### Step 2: Sanitize Branch Name
- Replace `/` with `-`
- Replace spaces with `-`
- Remove special characters except `-` and `_`
- Convert to lowercase

### Step 3: Determine State File Path
```
aidlc-docs/state/{sanitized-branch-name}.md
```

---

## State File Template

When starting work on a new branch, create the state file if it doesn't exist:

```markdown
# AI-DLC State: {branch-name}

## Branch Info
- **Branch**: {original-branch-name}
- **Base Branch**: {base-branch, e.g., main}
- **Created**: {ISO timestamp}
- **Current Stage**: 🔵 INCEPTION - Workspace Detection

## Project Context
- **Project Type**: {Greenfield/Brownfield - TBD}
- **Request Summary**: {Brief description of the feature/fix}

## Stage Progress

### 🔵 INCEPTION PHASE
- [ ] Workspace Detection
- [ ] Reverse Engineering (Brownfield only)
- [ ] Requirements Analysis
- [ ] User Stories
- [ ] Workflow Planning
- [ ] Application Design
- [ ] Units Generation

### 🟢 CONSTRUCTION PHASE
- [ ] Functional Design
- [ ] NFR Requirements
- [ ] NFR Design
- [ ] Infrastructure Design
- [ ] Code Generation
- [ ] Build and Test

### 🟡 OPERATIONS PHASE
- [ ] Operations (Placeholder)

## Session Notes
<!-- Notes and decisions specific to this branch -->
```

---

## State Index Format

The `state-index.md` file tracks all state files:

```markdown
# AIDLC State Index

## Active Branches
| Branch | State File | Current Stage | Last Updated |
|--------|------------|---------------|--------------|
| main | [main.md](./main.md) | ✅ COMPLETE | {timestamp} |
| feature/user-auth | [feature-user-auth.md](./feature-user-auth.md) | 🔵 Requirements | {timestamp} |

## Recently Merged
| Branch | Merged To | Final Stage | Merge Date |
|--------|-----------|-------------|------------|
| feature/onboarding | main | ✅ COMPLETE | {date} |

## Archived
See `archived/` folder for completed branch states.
```

---

## Workflow Integration

### On AIDLC Session Start
1. Detect current Git branch
2. Sanitize branch name to filename
3. Check if state file exists for this branch
4. If not exists:
   - For feature branches: Create new state file with branch info
   - For main branch: Check if migration is needed from `aidlc-state.md`
5. Update state-index.md with branch entry

### On Each Stage Transition
1. Update branch-specific state file
2. Mark completed stages with `[x]`
3. Update `Current Stage` field
4. Add relevant notes to `Session Notes`

### On Branch Merge (Manual Process)
1. Move state file to `archived/` folder
2. Add summary to main branch state (if applicable)
3. Update state-index.md

---

## Relationship with Audit System

| Aspect | Audit | State |
|--------|-------|-------|
| **Purpose** | Track all inputs/outputs/decisions | Track workflow progress |
| **Content** | Detailed logs with timestamps | Stage completion status |
| **Update Frequency** | Every interaction | Stage transitions |
| **Format** | Append-only log entries | Checkbox-based progress |
| **Location** | `aidlc-docs/audit/{branch}.md` | `aidlc-docs/state/{branch}.md` |

Both systems use **identical branch detection and naming** for consistency.

---

## Special Cases

### No Git Repository
- Default to `main.md` as the state file
- Log warning that branch detection is unavailable

### Detached HEAD State
- Use `detached-{short-commit-hash}.md` as filename

### Legacy Single aidlc-state.md
If project has existing `aidlc-state.md` in root:
1. Move to `state/main.md`
2. Create `state/state-index.md`
3. Remove old `aidlc-state.md` (or keep as symlink for backwards compatibility)

---

## Benefits

| Benefit | Description |
|---------|-------------|
| **Independent Progress** | Each branch tracks its own AIDLC progress |
| **Team Collaboration** | Multiple features can be worked on in parallel |
| **Clear Context** | Resume work knowing exactly where you left off |
| **Merge Clarity** | Know what stages were completed for each feature |
| **Audit Alignment** | Matches branch-based audit system structure |

