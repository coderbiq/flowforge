# flowforge-feedback

## Start

Use when the user is reviewing test results, code reviews, or task execution
output and reports a bug, finding, gap, unexpected behavior, or design issue.
Also use when the user explicitly says "反馈", "feedback", "报告问题", or
"handoff feedback".

Do **not** use for:
- Simple card lookup or navigation (use `flowforge card read` / `search`).
- Creating or editing proposal requirements or designs (use `flowforge-design`).
- Executing an existing ready task (use `flowforge-implement`).
- Bulk importing knowledge from external documents (use `flowforge-curate`).

## Workflow

Follow `references/workflow-rules.md` for the 5-step turn loop
(receive → classify → route → record → verify).  Use
`references/classification-rules.md` for the 5-type decision tree.

Key commands:

| Command | Purpose |
|---|---|
| `flowforge log create --kind feedback --for <card-id>` | Record each discovery |
| `flowforge card create --type task/finding/requirement` | Create output cards |
| `flowforge library import` / `library promote` | Route knowledge to library |
| `flowforge structure add` | Index new requirements in STR |
| `flowforge validate all` | Verify card state before closing turn |

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

Report:

- Cards created or updated (type, ID, title, status).
- Relations added (which card links to which, relation type).
- Unresolved gaps (open questions, missing requirements, pending library
  candidates).
- One next step (the most important card or task to address next).
