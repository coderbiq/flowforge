---
name: flowforge-implement
description: Use ONLY when the user asks to execute a ready FlowForge implementation task, or provides a task id and wants code changes for that task. Do NOT use for design, analysis, feedback, archive, or general card lookup.
---

# flowforge-implement

## Start

Resolve the task with `flowforge task ready --type i` or an explicit task id. Run `flowforge context task --task <id>`. Confirm linked requirement, design, and constraints are present.

## Workflow

Follow `references/workflow-rules.md` for execution and completion steps.

## Hard Rules

- Stop immediately if the task is `not_ready`, `blocked`, `design`, or `analysis`.
- CLI is the only read/write path for cards.
- Never edit card files, frontmatter, or wikilinks manually.
- Only make changes within the ready task's defined scope.
- Run tests and `flowforge validate all` when card state changed.
- For multi-line body content: use inline `--body` with `\n` for newlines. Example: `--body "## Goal\n\ncontent"`. Never use shell heredoc or redirects with flowforge CLI — redirects trigger agent permission prompts.

## Output

Report changed files, tests run, validation result, gaps or blockers, and one next step.