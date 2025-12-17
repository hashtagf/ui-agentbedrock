# /aidlc-changelog

Generate CHANGELOG entry for completed work.

## When to Use

Run this command after completing your work to update CHANGELOG.md.

## Instructions

1. Read the current branch state file (`aidlc-docs/state/{branch}.md`)
2. Read recent audit entries (`aidlc-docs/audit/{branch}.md`)
3. Identify all changes made in current session/feature
4. Read existing `CHANGELOG.md` from project root
5. Add entries under `[Unreleased]` section using Keep a Changelog format:
   - `Added` - New features
   - `Changed` - Changes in existing functionality
   - `Deprecated` - Soon-to-be removed features
   - `Removed` - Now removed features
   - `Fixed` - Bug fixes
   - `Security` - Vulnerability fixes
6. Present summary to user for confirmation

## Output Format

```markdown
## CHANGELOG Update

### Changes to add:

**Added**
- [Feature 1]
- [Feature 2]

**Changed**
- [Change 1]

**Fixed**
- [Fix 1]

---

**Ready to update CHANGELOG.md?**
- ✅ Confirm - Add these entries
- ✏️ Edit - Modify before adding
```

## Reference

- Format: [Keep a Changelog](https://keepachangelog.com/en/1.1.0/)
- See: `common/changelog-management.md`

