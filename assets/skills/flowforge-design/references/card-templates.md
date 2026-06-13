# Card Templates

Use these as the minimum body structure when creating or reviewing FlowForge cards.

## Requirement

```markdown
# <requirement title>

## Summary

## Source

## Acceptance

## Scope

## Open Questions
```

Review rules:
- One card = one user-visible behavior, constraint, or acceptance point.
- `Acceptance` must contain at least one testable condition.
- `Open Questions` must not be omitted; write `None` if closed.

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
