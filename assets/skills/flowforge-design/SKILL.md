# flowforge-design

Use when the user wants to analyze, clarify, design, or decompose a FlowForge proposal before implementation. Grow cards through `flowforge`; do not write long proposal docs.

## Start

Run `flowforge project current`, `flowforge proposal current`, `flowforge proposal inspect <id>`, then `flowforge context proposal --proposal <id>`. If project/proposal is missing, ask the user to create or select it.

## Workflow

Follow `references/workflow-rules.md`. Use `structure add/remove`, atomic requirement cards, analysis tasks for uncertainty, `library suggest` / `card search --scope library` / `card read --summary/--section` for library discovery, then focused design cards.

Use `references/card-templates.md` whenever creating or reviewing card bodies. Use `references/library-discovery.md` before reading or linking library cards. Record each real design turn with `log create --kind <kind>`.

## Hard Rules

- CLI is the only read/write path for cards.
- Never read wiki files or `02-library/` directly.
- Never load the whole proposal or library.
- Never create title-only tasks.
- Do not execute implementation work here.

## Output

Report updated cards, relations, gaps, and one next step.
