---
name: flowforge-curate
description: Use ONLY when the user wants to import knowledge from external documents into the library, or archive a completed proposal to the library. Do NOT use for single-card creation, proposal design, or task execution.
---

# flowforge-curate

## Start

Check for an in-progress plan card (tag `curation-plan`, status `active`). If one exists, resume batch execution. Otherwise determine the mode:

- **Mode A (external import)**: user provided a file path — read the source file.
- **Mode B (proposal archive)**: user provided a proposal ID — run `flowforge proposal inspect <id>` then `flowforge context proposal --proposal <id>`.

## Workflow

Follow `references/workflow-rules.md` for mode-specific extraction, clustering, and batch execution. Use `references/extraction-guide.md` for knowledge unit criteria and card type mapping.

## Hard Rules

- Stop and wait for user review before writing any cards.
- Always create cards with `status: draft`; promote to `active` only after user confirms.
- Batch size: 5-10 items per activation. The plan card tracks progress.
- CLI is the only read/write path for cards.
- Never read wiki files directly (except source files for Mode A).
- Never hand-write card files, frontmatter, or wikilinks.
- Always read the plan card first on each activation to resume state.
- Use `--body -` with heredoc (`<<'EOF' ... EOF`) for multi-line body content. Single-quoted heredoc delimiter prevents all shell expansion.
- For batch card creation, generate a YAML manifest and use `card create --batch <file>` instead of per-card CLI calls.
- Use `card update --section "<name>" --body -` with heredoc to update the plan card's batch progress.
- Use `-o json` to capture created card IDs for scripting.

## Output

Report batch number, completed/total items, created card IDs, and next step ("continue" or "done").