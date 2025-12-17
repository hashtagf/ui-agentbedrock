# Audit Management - Branch-Based System

## Overview

AIDLC uses a **branch-based audit system** that organizes audit logs by Git branch. This provides natural alignment with Git workflow and makes it easy to trace decisions for each feature/requirement.

---

## Directory Structure

```text
aidlc-docs/
├── audit/
│   ├── audit-index.md              # Master index of all audit files
│   ├── main.md                     # Audit for main/master branch
│   ├── {sanitized-branch-name}.md  # Audit for each feature branch
│   └── archived/                   # Merged/completed branch audits
│       └── {sanitized-branch-name}.md
├── state/
│   ├── state-index.md              # Master index of all state files
│   ├── main.md                     # State for main/master branch
│   ├── {sanitized-branch-name}.md  # State for each feature branch
│   └── archived/                   # Merged/completed branch states
│       └── {sanitized-branch-name}.md
└── ... (other aidlc-docs)
```

> **Note**: See `common/state-management.md` for branch-based state tracking rules.

---

## Branch Detection Logic

### Step 1: Get Current Branch
```bash
git branch --show-current
```

If git is not available or not in a git repo, default to `main`.

### Step 2: Sanitize Branch Name
Convert branch name to valid filename:
- Replace `/` with `-`
- Replace spaces with `-`
- Remove special characters except `-` and `_`
- Convert to lowercase

**Examples**:
| Git Branch | Sanitized Filename |
|------------|-------------------|
| `main` | `main.md` |
| `feature/user-auth` | `feature-user-auth.md` |
| `feature/REQ-001-payment` | `feature-req-001-payment.md` |
| `bugfix/login-issue` | `bugfix-login-issue.md` |
| `hotfix/critical-fix` | `hotfix-critical-fix.md` |

### Step 3: Determine Audit File Path
```
aidlc-docs/audit/{sanitized-branch-name}.md
```

---

## Audit File Initialization

When starting work on a new branch, create the audit file if it doesn't exist:

```markdown
# Audit Trail: {branch-name}

## Branch Info
- **Branch**: {original-branch-name}
- **Base Branch**: {base-branch, e.g., main}
- **Created**: {ISO timestamp}
- **Status**: 🟡 In Progress
- **Related Issues**: {if available from branch name or commit}

---

## Sessions

<!-- Audit entries will be appended below -->
```

---

## Audit Entry Format

Each audit entry follows this format:

```markdown
## {Stage Name or Interaction Type}
**Timestamp**: {ISO 8601 timestamp}
**Branch**: {current-branch}
**User Input**: "{Complete raw user input - never summarized}"
**AI Response**: "{AI's response or action taken}"
**Context**: {Stage, action, or decision made}

---
```

---

## Audit Index Format

The `audit-index.md` file tracks all audit files:

```markdown
# AIDLC Audit Index

## Active Branches
| Branch | Audit File | Last Updated | Status |
|--------|------------|--------------|--------|
| main | [main.md](./main.md) | {timestamp} | ✅ Stable |
| feature/user-auth | [feature-user-auth.md](./feature-user-auth.md) | {timestamp} | 🟡 In Progress |

## Recently Merged
| Branch | Merged To | Merge Date | PR/MR |
|--------|-----------|------------|-------|
| feature/onboarding | main | {date} | #{number} |

## Archived
See `archived/` folder for completed branch audits.
```

---

## Workflow Integration

### On AIDLC Session Start
1. Detect current Git branch
2. Sanitize branch name to filename
3. Check if audit file exists for this branch
4. If not exists, create new audit file with branch info
5. Update audit-index.md with branch entry

### On Each Interaction
1. Append entry to branch-specific audit file
2. Include branch name in each entry for cross-reference

### On Branch Merge (Manual Process)
1. Move audit file to `archived/` folder
2. Add summary entry to target branch audit (e.g., `main.md`)
3. Update audit-index.md

---

## Special Cases

### No Git Repository
- Default to `main.md` as the audit file
- Log warning that branch detection is unavailable

### Detached HEAD State
- Use `detached-{short-commit-hash}.md` as filename
- Example: `detached-a1b2c3d.md`

### Very Long Branch Names
- Truncate sanitized name to 100 characters
- Append hash suffix if needed for uniqueness

---

## Migration from Single audit.md

If project has existing `audit.md`:
1. Rename to `audit/main.md`
2. Create `audit/audit-index.md`
3. Update all references in rules

---

## Benefits

| Benefit | Description |
|---------|-------------|
| **Git Alignment** | Natural fit with feature branch workflow |
| **Easy Review** | Audit history stays with feature in PRs |
| **Team Clarity** | Each developer sees their branch's history |
| **Cleaner History** | Main branch audit stays focused |
| **Auto-Archive** | Natural archival when branches merge |

