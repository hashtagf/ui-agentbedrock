# /aidlc-multi-repo - Multi-Repository Setup

Configure AIDLC for multi-repository projects (frontend, backend, jobs, etc.)

## What This Command Does

1. **Create Related Projects Config**
   - Create `aidlc-docs/related-projects.md`
   - Define project ecosystem

2. **Scan for Related Projects**
   - Check parent directory for sibling projects
   - Identify project types (frontend, backend, etc.)

3. **Generate Integration Map**
   - Document how projects interact
   - Identify shared dependencies

## Usage

```
/aidlc-multi-repo
```

With specific projects:
```
/aidlc-multi-repo ../frontend ../backend ../shared
```

## Output Format

```markdown
# Related Projects Configuration

## Project Ecosystem

| Project | Type | Path | Description |
|---------|------|------|-------------|
| my-frontend | Frontend | ../my-frontend | React SPA |
| my-backend | Backend | ../my-backend | Node.js API |

## Integration Points

### API Contracts
- Frontend → Backend: REST API at /api

## Cross-Repo Notes
[Notes about integration]
```

## Use Cases

### Microservices
```
/aidlc-multi-repo
```
→ Detects all services, maps communication

### Frontend + Backend
```
/aidlc-multi-repo ../frontend ../backend
```
→ Maps API contracts between them

### Monorepo with packages
```
/aidlc-multi-repo ./packages/web ./packages/api ./packages/shared
```
→ Maps internal package dependencies

## During AIDLC Workflow

When `related-projects.md` exists:

- **Requirements**: Shows impact on all projects
- **Code Generation**: Generates cross-repo change notes
- **Build & Test**: Includes integration test instructions

## Related Commands

- `/aidlc` - Main workflow
- `/aidlc-status` - Shows related projects status
- `/aidlc-reverse` - Analyzes cross-repo integration

## See Also

- Rule: `common/multi-repo-context.md`

