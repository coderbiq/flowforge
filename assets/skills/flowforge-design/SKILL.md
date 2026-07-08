---
name: flowforge-design
description: Use ONLY when the user wants to design or decompose a FlowForge proposal before implementation. Do NOT use for implementation, feedback, archiving, or single-card creation.
---

# flowforge-design

## Start

Run `flowforge project current`, `flowforge proposal current`, `flowforge proposal inspect <id>`, then `flowforge context proposal --proposal <id>`. If project/proposal is missing, ask the user.

## Workflow

Follow `references/workflow-rules.md` for the 7-mode turn loop. Use `references/card-templates.md` for card bodies. Use `references/library-discovery.md` for library context discovery.

## Hard Rules

- CLI is the only read/write path for cards.
- Never read wiki files or `02-library/` directly.
- Never load the whole proposal or library at once.
- Never hand-write card files, frontmatter, wikilinks, or internal card links.
- Never create title-only tasks.
- Do not execute implementation work here.
- Run `flowforge validate all` after creating or changing proposal structure.
- For multi-line body content: use inline `--body` with `\n` for newlines. Example: `--body '## Goal\n\ncontent'`. Never use shell heredoc or redirects with flowforge CLI — redirects trigger agent permission prompts.
- Never create > 10 REQ cards in a single index pass without creating at least 1 DESIGN card.
- Never create a REQ card with < 5 lines of effective business content; merge into parent instead.
- After index mode, the next recommended step must be design or clarify mode; never recommend another index pass.
- STR cards must contain `## Synthesis` section; propose `card update` when synthesis is missing.
- Run `flowforge proposal inspect <id>` and address all health issues before proposing implementation readiness.

## Output

Report cards created/updated, relations added, unresolved gaps, and one next step.