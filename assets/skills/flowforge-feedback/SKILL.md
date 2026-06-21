---
name: flowforge-feedback
description: Use ONLY when the user is reviewing test results, code reviews, or task execution output and reports a bug, finding, gap, unexpected behavior, or design issue. Also use when the user explicitly says "反馈", "feedback", "报告问题", or "handoff feedback". Do NOT use for simple card lookups, proposal design, task execution, or bulk knowledge import.
---

# flowforge-feedback

## Start

Run `flowforge context task --task <id>` (if triggered by a task execution) or
`flowforge context proposal --proposal <id>` to gather current context. If no
task or proposal is active, ask the user which task or proposal the discovery
belongs to.

Do **not** use for:
- Simple card lookup or navigation (use `flowforge card read` / `search`).
- Designing or decomposing requirements (use `flowforge-design`). Feedback may
  create draft requirement cards as tracking items, but detailed design happens
  in flowforge-design.
- Executing an existing ready task (use `flowforge-implement`).
- Bulk importing knowledge from external documents (use `flowforge-curate`).

## Workflow

Follow `references/workflow-rules.md` for the 5-step turn loop
(receive → classify → route → record → verify).  Use
`references/classification-rules.md` for the 5-type decision tree.

## Hard Rules

- CLI is the only read/write path for cards.
- Never read wiki files or `02-library/` directly.
- Every discovery produces a log card with `--kind feedback`.
- Bug / missing-requirement / design-flaw must immediately produce a tracking
  card (task or requirement); do not record a log only.
- Knowledge findings must be routed to `library import` or `library promote`.
- Use `--body -` with heredoc (`<<'EOF' ... EOF`) for multi-line body content.
  Single-quoted heredoc delimiter prevents all shell expansion.
- Run `flowforge validate all` after creating or changing card structure.

## Output

Report cards created/updated, relations added, unresolved gaps, and one next step.