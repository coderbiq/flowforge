---
name: flowforge-feedback
description: Use ONLY when the user is reviewing test results, code reviews, or task execution output and reports a bug, finding, gap, unexpected behavior, or design issue. Also use when the user explicitly says "反馈", "feedback", "报告问题", or "handoff feedback". Do NOT use for simple card lookups, proposal design, task execution, or bulk knowledge import.
---

# flowforge-feedback

## Start

Run `context feature --feature <id>` (if triggered by a task execution) or
`proposal inspect <id>` to gather current context. If no proposal is active,
ask the user which proposal the discovery belongs to.

## Workflow

Follow `references/workflow-rules.md` for the 5-step turn loop.
Use `references/classification-rules.md` for the 5-type decision tree.

## Hard Rules

- Create tracking cards via `card init --type feature`; then edit directly.
- Bug / missing-requirement / design-flaw → create a FEATURE card (draft) or annotate existing.
- Knowledge findings → route to `library import` or `library promote`.
- Use `card log` (not `log create`) for progress recording.
- Use `card link`/`card unlink` for all link operations.
- Run `flowforge validate all` after any changes.
- `task`, `structure`, `log create` are DEPRECATED.

## Output

Report cards created/updated, relations added, unresolved gaps, and one next step.
