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

- **Goal**: Verifiable deliverable
- **Files**: Paths relative to project root
- **Approach**: Method signatures, pseudocode, algorithms
- **Edge Cases**: At least 1 boundary condition
- **Dependencies**: FEATURE IDs and wait strategy
- **Parallel**: yes or no
- **Verification**: Test scenarios, key assertions

## Verification

## History

## Open Questions

## Dependencies
```

Review rules:
- Each step must include Files, Approach, and Edge Cases.
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
