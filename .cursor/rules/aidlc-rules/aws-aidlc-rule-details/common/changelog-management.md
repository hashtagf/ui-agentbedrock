# CHANGELOG Management

## Overview

AIDLC manages the project's `CHANGELOG.md` file following [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) standard.

**Use `/aidlc-changelog` command to update CHANGELOG.md when ready.**

### Types of Changes
| Type | Usage |
|------|-------|
| `Added` | New features |
| `Changed` | Changes in existing functionality |
| `Deprecated` | Soon-to-be removed features |
| `Removed` | Now removed features |
| `Fixed` | Bug fixes |
| `Security` | Vulnerability fixes |

---

## CHANGELOG Location

The CHANGELOG is created at the project root:

```
your-project/
├── CHANGELOG.md          # ← Project changelog
├── .cursor/
├── aidlc-docs/
└── [source code]
```

---

## When to Create/Update CHANGELOG

| Trigger | Action |
|---------|--------|
| **aidlc-init** | Create initial CHANGELOG.md if not exists |
| **`/aidlc-changelog`** | User triggers to add entries for completed work |
| **Release** | Convert [Unreleased] to versioned entry |

---

## CHANGELOG Template

When creating a new CHANGELOG.md:

```markdown
# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

### Added
<!-- New features will be added here by AIDLC -->

### Changed
<!-- Changes to existing functionality -->

### Fixed
<!-- Bug fixes -->

### Removed
<!-- Removed features -->

---

<!-- Previous versions will be added below -->
```

---

## How to Update CHANGELOG

### During Code Generation

After each unit's code generation is complete, add entries under `[Unreleased]`:

```markdown
## [Unreleased]

### Added
- **[Unit Name]**: [Brief description of what was added]
  - [Specific feature 1]
  - [Specific feature 2]
  - Related stories: [Story IDs]
```

### Example Entry

```markdown
## [Unreleased]

### Added
- **User Authentication Module**: Implemented complete user authentication system
  - JWT-based authentication with refresh tokens
  - Login/logout endpoints
  - Password reset functionality
  - Related stories: US-001, US-002, US-003

- **User Management API**: RESTful API for user CRUD operations
  - Create, read, update, delete user endpoints
  - Role-based access control
  - Input validation and error handling
  - Related stories: US-004, US-005
```

---

## During Build and Test (Version Finalization)

When the Build and Test stage completes successfully:

1. **Determine Version Number**:
   - If user specifies version → Use that version
   - If not specified → Suggest based on changes:
     - Major: Breaking changes
     - Minor: New features
     - Patch: Bug fixes only

2. **Convert [Unreleased] to Version**:

```markdown
## [1.0.0] - 2025-12-15

### Added
- **User Authentication Module**: Implemented complete user authentication system
  - JWT-based authentication with refresh tokens
  - Login/logout endpoints
  - Password reset functionality

### Changed
- Updated database schema for user roles

---

## [Unreleased]

<!-- Ready for next development cycle -->
```

3. **Ask User for Version Confirmation**:

```markdown
📋 **CHANGELOG Update**

The following changes are ready to be released:
- [Summary of changes]

**Suggested Version**: 1.0.0 (Major release - new project)

Would you like to:
1. Use suggested version (1.0.0)
2. Specify a different version
3. Keep as [Unreleased]
```

---

## CHANGELOG Entry Guidelines

### Good Entry Examples

```markdown
### Added
- **Payment Gateway Integration**: Stripe payment processing for subscriptions
  - Monthly/yearly billing cycles
  - Webhook handlers for payment events
  - Invoice generation

### Changed
- **User Model**: Added `subscription_tier` and `billing_date` fields
- **API Response Format**: Standardized error responses across all endpoints

### Fixed
- **Login Flow**: Fixed session expiration not redirecting to login page
- **Data Validation**: Fixed email validation allowing invalid formats
```

### Bad Entry Examples

```markdown
### Added
- Added some stuff          ❌ Too vague
- Fixed bugs               ❌ Not descriptive
- Updated code             ❌ Not meaningful
```

---

## Integration with AIDLC Stages

### 1. aidlc-init Stage

```markdown
**Action**: Check if CHANGELOG.md exists
- If NOT exists → Create from template
- If exists → Leave unchanged

**Log in branch audit file** (`aidlc-docs/audit/{branch}.md`):
- "CHANGELOG.md created" or "CHANGELOG.md already exists"
```

### 2. Code Generation Stage (Per Unit)

```markdown
**Action**: Update CHANGELOG.md after each unit completes
- Add entries under [Unreleased] section
- Include unit name and key features
- Reference related user stories

**Update Format**:
### Added
- **{Unit Name}**: {Brief description}
  - {Feature 1}
  - {Feature 2}
  - Related stories: {Story IDs}
```

### 3. Build and Test Stage

```markdown
**Action**: Finalize CHANGELOG version
- Summarize all [Unreleased] entries
- Ask user for version number
- Convert [Unreleased] to versioned entry
- Add release date
- Create new empty [Unreleased] section
```

---

## Rules

1. **Never Delete Existing Entries**: Only append or modify [Unreleased]
2. **Keep Format Consistent**: Follow Keep a Changelog format
3. **Be Descriptive**: Each entry should explain WHAT and WHY
4. **Link to Stories**: Reference related user stories when applicable
5. **Date Format**: Use ISO 8601 (YYYY-MM-DD)
6. **Version Format**: Follow Semantic Versioning (X.Y.Z)

