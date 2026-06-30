# Card Templates

Use these as the minimum body structure when creating or reviewing FlowForge cards.

Do not put internal card links in these bodies by hand. Use frontmatter links through CLI commands, then run `flowforge card refresh <id>` when a requirement or design needs generated navigation. Hand-written Markdown links are only for external source references.

## Requirement

```markdown
# <requirement title>

## Summary

## Source

## Acceptance

## Scope

## Open Questions

## Dependencies

## See Also
```

Review rules:
- One card = one user-visible behavior, constraint, or acceptance point.
- `Acceptance` must contain at least one testable condition.
- `Open Questions` must not be omitted; write `None` if closed.
- `Dependencies` must list: other REQ cards this card depends on AND the reason; external systems/modules; which REQ cards depend on this one. Write `None` if no dependencies.
- `See Also` must list: related DESIGN or DEC cards. Write `None` if none.

## Analysis Task

```markdown
# <analysis task title>

## Goal

## Inputs

## Investigation Plan

## Expected Outputs

## Done When
```

Review rules:
- Must name the uncertainty being investigated.
- Must say what sources or modules will be checked.
- Must not be title-only.

## Structure

```markdown
# <structure title>

## Purpose

## Synthesis

## Key Decisions

## Entries
```

Review rules:
- `Purpose` must state the core problem this structure addresses (1-2 sentences).
- `Synthesis` is mandatory: explain how indexed cards collaborate, key design constraints, and cross-cutting concerns (3-8 lines). Must not be placeholder text.
- `Key Decisions` records design choices that span multiple indexed cards.
- `Entries` is auto-managed by CLI (`structure add`/`structure refresh`); do not hand-write.
- An STR that contains only `Purpose` + `Entries` (no `Synthesis`, no `Key Decisions`) is incomplete -- use `card update` to add synthesis before the proposal is inspectable.

## Design

```markdown
# <design title>

## Goal

## Decision

## Rationale

## Constraints

## Impact

## Verification

## Follow-up Tasks
```

Review rules:
- Must express one stable design focus.
- `Constraints` must come from confirmed context or library evidence.
- Do not summarize the entire proposal in one design card.

## Implementation Task

```markdown
# <task title>

## Goal

## Inputs

## Deliverables

## Acceptance

## Out of Scope

## Read Before Work
```

Review rules:
- A ready task must have linked requirement, design, constraints, and acceptance.
- If the task depends on assumptions, mark it `not_ready` or blocked.
- Do not create implementation tasks with only a title.

## Log

```markdown
# <log title>

## Kind

## Event

## Context

## Result
```

Review rules:
- One log = one event.
- Use logs for process evidence, not as a replacement for requirement, design, or finding cards.
- Logs should point at the relevant proposal, task, requirement, or design.

## Library Candidate

Use this structure before calling `flowforge library import`:

```markdown
# <candidate title>

## Summary

## Rule or Finding

## Applies When

## Source Evidence
```

Review rules:
- Candidate must be atomic and reusable beyond the current proposal.
- Candidate must have a confirmed type such as `convention`, `module`, `decision`, `design`, or `finding`.
- Candidate must keep traceability with `--source-card` or explicit `--links`.

## Content Density Guidelines

Cards are living objects. Create them when content warrants, not when a source document has a bullet point.

| Density | Effective Content | Action |
|---------|-------------------|--------|
| **too-thin** | < 5 lines of business content | Do not create an independent card. Merge into parent or related card. |
| **suitable** | 5--20 lines | Suitable for an independent card. |
| **too-thick** | > 50 lines or any single section > 15 lines | Consider splitting into sub-cards with cross-references. |

"Effective content" = body text after removing frontmatter, template section headings, and auto-generated navigation sections (FlowForge Navigation, Links, Outgoing).

### Progressive Creation Strategy

1. **Coarse seeding first**: Create 1--3 REQ cards per EPIC with rich content. Write `Synthesis` in the EPIC STR.
2. **Split only when content grows**: When a REQ card exceeds 30 lines of effective content, split it into sub-cards. Each sub-card must reference the parent via `Dependencies`.
3. **Design after seeding**: After creating >= 3 REQ cards, create at least one DESIGN card before creating more REQ cards.
