# Branch Artifacts Management

## Overview

AIDLC uses a **branch-based artifact system** where each Git branch has its own `inception/` and `construction/` folders. This enables parallel development without artifact conflicts.

---

## Directory Structure

```text
aidlc-docs/
├── audit/                      # Branch-based audit logs
│   └── {branch}.md
├── state/                      # Branch-based state tracking
│   └── {branch}.md
└── branches/                   # Branch-based artifacts
    ├── branches-index.md       # Master index
    ├── main/
    │   ├── inception/
    │   └── construction/
    ├── feature-{name}/
    │   ├── inception/
    │   └── construction/
    └── archived/
        └── {branch-name}/
```

---

## Branch Detection Logic

Uses the **same logic as audit-management.md and state-management.md**:

### Step 1: Get Current Branch
```bash
git branch --show-current
```

### Step 2: Sanitize Branch Name
- Replace `/` with `-`
- Replace spaces with `-`
- Remove special characters except `-` and `_`
- Convert to lowercase

### Step 3: Determine Artifacts Path
```
aidlc-docs/branches/{sanitized-branch-name}/
```

---

## Artifact Paths by Phase

### INCEPTION Phase
```text
branches/{branch}/inception/
├── plans/
│   ├── workspace-detection.md
│   ├── workflow-planning.md
│   ├── story-generation-plan.md
│   └── execution-plan.md
├── reverse-engineering/        # Brownfield only
│   ├── architecture.md
│   ├── code-structure.md
│   ├── api-documentation.md
│   └── ...
├── requirements/
│   ├── requirements.md
│   └── requirement-verification-questions.md
├── user-stories/
│   ├── stories.md
│   └── personas.md
└── application-design/
    ├── components.md
    ├── component-methods.md
    ├── services.md
    ├── component-dependency.md
    ├── unit-of-work.md
    ├── unit-of-work-dependency.md
    └── unit-of-work-story-map.md
```

### CONSTRUCTION Phase
```text
branches/{branch}/construction/
├── plans/
│   └── {unit-name}-*.md
├── {unit-name}/
│   ├── functional-design/
│   ├── nfr-requirements/
│   ├── nfr-design/
│   ├── infrastructure-design/
│   └── code/
└── build-and-test/
    ├── build-instructions.md
    ├── unit-test-instructions.md
    ├── integration-test-instructions.md
    └── build-and-test-summary.md
```

---

## Workflow Integration

### On AIDLC Session Start
1. Detect current Git branch
2. Sanitize branch name
3. Check if branch directory exists: `branches/{branch}/`
4. If not exists, create directory structure:
   ```bash
   mkdir -p branches/{branch}/inception
   mkdir -p branches/{branch}/construction
   ```
5. Update branches-index.md

### On Artifact Creation
1. Always create artifacts in branch-specific path
2. Example: `branches/feature-auth/inception/requirements/requirements.md`
3. Never use root-level `inception/` or `construction/` paths

### On Branch Merge
1. Move branch folder to `archived/{branch-name}/`
2. Update branches-index.md
3. Optionally clean up after retention period

---

## Relationship with Other Systems

| System | Location | Purpose |
|--------|----------|---------|
| **Audit** | `audit/{branch}.md` | Log all interactions |
| **State** | `state/{branch}.md` | Track progress |
| **Artifacts** | `branches/{branch}/` | Store phase outputs |

All three systems use **identical branch detection and naming**.

---

## Special Cases

### No Git Repository
- Default to `main` as branch name
- All artifacts go to `branches/main/`

### Detached HEAD State
- Use `detached-{short-commit-hash}` as branch name

### Main Branch
- Main branch artifacts are in `branches/main/`
- This is the "baseline" for reference

---

## Benefits

| Benefit | Description |
|---------|-------------|
| **Parallel Development** | Multiple features can be developed simultaneously |
| **No Conflicts** | Artifacts don't overwrite each other |
| **Clear Traceability** | Each feature has complete history |
| **Easy PR Review** | All docs included in feature PRs |
| **Clean Organization** | Main branch stays uncluttered |

