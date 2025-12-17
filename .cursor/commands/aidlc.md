# /aidlc - Main Entry Command

Main entry point for AIDLC (AI Development Life Cycle) workflow.

## What This Command Does

When you use `/aidlc`, the AI will automatically:

1. **Detect Workspace State**
   - Check for existing `state/{branch}.md`
   - Scan for existing source code
   - Determine if Greenfield (new) or Brownfield (existing code)

2. **For New Projects**
   - Create `aidlc-docs/` folder structure
   - Initialize `state/{branch}.md` and `audit/{branch}.md`
   - Progress through AIDLC stages

3. **For Resume**
   - Load existing state
   - Continue from last stage
   - Show current progress

4. **For Specific Requests**
   - Analyze your request context
   - Skip stages that aren't needed
   - Execute relevant stages only

## Why Only One Main Command?

The **AIDLC core-workflow** handles everything automatically:
- ✅ Auto-progress through stages
- ✅ Auto-skip unnecessary stages
- ✅ Context-aware execution
- ✅ Branch-based state management

You don't need separate commands for each stage!

## Usage Examples

### Start New Project
```
/aidlc สร้าง REST API สำหรับ user authentication
```

### Resume Work
```
/aidlc
```

### Jump to Specific Stage
```
/aidlc skip to code generation
```

### Re-run a Stage
```
/aidlc re-run requirements analysis
```

## AIDLC Workflow Overview

```
🔵 INCEPTION PHASE
├── Workspace Detection → Requirements → Stories → Planning → Design → Units

🟢 CONSTRUCTION PHASE  
├── Functional Design → NFR → Infrastructure → Code Generation → Build & Test

🟡 OPERATIONS PHASE
└── (Placeholder for future)
```

## Related Commands

| Command | Description |
|---------|-------------|
| `/aidlc-status` | View current status |
| `/aidlc-multi-repo` | Configure related projects |

## Branch-Based System

AIDLC uses **branch-based tracking** for team collaboration:

```
aidlc-docs/
├── state/
│   └── {branch}.md      # State per branch
├── audit/
│   └── {branch}.md      # Audit per branch
└── branches/
    └── {branch}/        # Artifacts per branch
        ├── inception/
        └── construction/
```

### How It Works
- Detects current Git branch automatically
- Creates separate files per branch
- Archives when branch is merged
