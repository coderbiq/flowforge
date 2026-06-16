# flowforge-design

Use when the user wants to analyze, clarify, design, or decompose a FlowForge proposal before implementation. Grow cards through `flowforge`; do not write long proposal docs.

## Start

Run `flowforge project current`, `flowforge proposal current`, `flowforge proposal inspect <id>`, then `flowforge context proposal --proposal <id>`. If project/proposal is missing, ask the user to create or select it.

## Workflow

Follow `references/workflow-rules.md`. Use `structure add/remove`, atomic requirement cards, analysis tasks for uncertainty, `library suggest` / `card search --scope library` / `card read --summary/--section` for library discovery, then focused design cards.

Use `references/card-templates.md` for card bodies and `references/library-discovery.md` before linking library cards. Record real design turns with `log create`.

## Hard Rules

- CLI is the only read/write path for cards.
- Never read wiki files or `02-library/` directly.
- Never hand-write card files, frontmatter, wikilinks, or internal card links.
- Hand-written Markdown links are only for external references.
- Never load the whole proposal or library.
- Never create title-only tasks.
- Do not execute implementation work here.
- Always use single quotes for `--body` content containing mermaid, code blocks, or shell-special characters (`$`, `` ` ``, `!`, `{}`). Double-quoted `--body "..."` will be corrupted by shell expansion.

## Output

Report updated cards, relations, gaps, and one next step.
