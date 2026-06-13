# flowforge-design

Use when the user wants to analyze, clarify, design, or decompose a FlowForge proposal before implementation. Grow cards through `flowforge`; do not write long proposal docs.

## Start

Run `flowforge project current`, `flowforge proposal current`, `flowforge proposal inspect <proposal-id>`, then `flowforge context proposal --proposal <proposal-id>`. If project or proposal is missing, ask the user to create or select it.

## Workflow

Update STR with `structure add/remove`, create atomic requirements, turn uncertainty into analysis tasks, discover library context with `library suggest`, `card search --scope library`, and `card read --summary/--section`, then create focused design cards.

Create implementation tasks only when requirement, design, constraints, and acceptance are present; otherwise mark `not_ready` or blocked. Record each turn with `log create --kind <kind>`.

## Hard Rules

- CLI is the only read/write path for cards.
- Never read wiki files or `02-library/` directly.
- Never load the whole proposal or library.
- Never create title-only tasks.
- Do not execute implementation work here.

## Output

Report updated cards, relations, gaps, and one next step.
