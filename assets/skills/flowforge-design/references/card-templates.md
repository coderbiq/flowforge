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
