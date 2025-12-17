# /aidlc-status - Display AIDLC Status

Display current AIDLC workflow status and progress.

## What This Command Does

1. **Load State**
   - Detect current Git branch
   - Read `aidlc-docs/state/{branch}.md`
   - If not found → Prompt to run `/aidlc-init`

2. **Calculate Progress**
   - Count completed stages
   - Calculate percentage

3. **Display Status**
   - Current phase and stage
   - Completed stages with checkmarks
   - Pending stages
   - Next recommended action

## Usage

Basic status:
```
/aidlc-status
```

Detailed with timestamps:
```
/aidlc-status --detailed
```

## Output Format

```markdown
📊 AIDLC Status

**Project**: {name}
**Type**: Greenfield/Brownfield
**Current Stage**: {stage}
**Progress**: 60%

🔵 INCEPTION PHASE
- [x] Workspace Detection
- [x] Requirements Analysis
- [x] User Stories
- [ ] Workflow Planning ← Current
- [ ] Application Design
- [ ] Units Generation

🟢 CONSTRUCTION PHASE
- [ ] Functional Design
- [ ] Code Generation
- [ ] Build and Test

⏭️ Next Action: Continue with Workflow Planning
```

## Related Commands

- `/aidlc` - Continue workflow
- `/aidlc-init` - Initialize project

