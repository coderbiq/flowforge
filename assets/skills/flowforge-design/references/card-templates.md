# Card Templates

Use these as the minimum body structure when creating or reviewing FlowForge cards.

Cards are created via `card init --type <type>`, then edited directly. Use `card link`/`card unlink` for link operations. Run `validate all` after edits.

## Feature

```markdown
# <feature title>

## Summary

## Motivation

## Design

### Key Decisions

### Architecture

### Alternatives Considered

## Constraints

## Implementation Plan

### Step N: <step goal>

<!-- step-status: not_started -->

- **Goal**:
  One observable deliverable.
- **Files**:
  - Exact paths relative to project root.
- **Symbols**:
  - Exact classes, functions, commands, or sections; use `None (documentation-only)` only when accurate.
- **Actions**:
  1. Ordered mechanical action with input and output.
  2. Next action; do not hide design choices behind "as needed".
- **Constraints**:
  - Forbidden changes, compatibility rules, and boundary behavior.
- **Done When**:
  - Observable completion conditions and assertions.
- **Dependencies**: FEATURE IDs and wait strategy
- **Parallel**: yes or no
- **Verification**:
  - Exact commands and expected results.

## Verification

## History

## Open Questions

## Dependencies
```

Review rules:
- Each step must include Goal, Files, Symbols, Actions, Constraints, Done When, and Verification.
- Put each action, constraint, completion condition, and command on its own line.
- Reject TBD, TODO, "as needed", "as appropriate", “必要时”, “视情况”, and other unresolved execution choices.
- No cross-card references (no "参考 DES-xxx").
- Step status via HTML comments, CLI-managed.

### Complex analysis opt-in

Simple FEATURE cards keep the minimum structure above. Only add the marker below when the requirement needs a recoverable investigation loop; validation then requires every listed section to contain a value (`None` is valid, `TBD` is not).

```markdown
<!-- analysis-mode: complex -->

## Objective

## Current Understanding

## Evidence

Use FIND card IDs for accepted investigation evidence, or `None` before evidence exists.

## Working Design

## Rejected or Revised Assumptions

## Open Questions

## Next Investigation
```

Each investigation FIND uses the same opt-in marker and a stable work ID from the published plan:

```markdown
<!-- analysis-mode: complex -->
<!-- analysis-work-id: <revision-and-work-id> -->

## Evidence

## Source

## Impact

## Open Questions
```

Investigators edit only those four FIND sections. Sources must be independently recoverable; do not rely on a live Agent session.

## Content Density Guidelines

| Density | Effective Content | Action |
|---------|-------------------|--------|
| **too-thin** | < 5 lines | Do not create independently |
| **suitable** | 5--20 lines | Suitable for independent card |
| **too-thick** | > 50 lines or section > 15 lines | Consider splitting |

### Progressive Creation Strategy

1. **Coarse seeding first**: Create 1--3 FEATURE cards per proposal.
2. **Split when needed**: When a FEATURE exceeds ~10 steps, use `card split`.
3. **Design after seeding**: Fill Design and evolve to designed before more FEATUREs.

---

<!-- DEPRECATED below -->

## Requirement (DEPRECATED)

## Design (DEPRECATED)

## Implementation Task (DEPRECATED)

## Log (DEPRECATED)

## Structure (DEPRECATED)
