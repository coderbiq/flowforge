---
name: tg-workflow
description: |
  Workflow orchestration skill for tg-workflow.

  Use this skill whenever the user wants to explore a new topic, create or refine a proposal, map work into tasks, archive completed work, or understand the current workflow state.
---

# tg-workflow

This skill is a thin adapter over the canonical workflow specification in `workflow/`.

## Source of truth

Always load and follow:

- `workflow/guides/lifecycle.md`
- `workflow/guides/authoring-rules.md`
- `workflow/guides/archive-rules.md`
- `workflow/guides/adapter-contract.md`

When metadata shape matters, use:

- `workflow/schema/proposal.schema.yaml`
- `workflow/schema/exploration.schema.yaml`
- `workflow/schema/task-map.schema.yaml`

## Responsibilities

- create and maintain exploration artifacts
- convert validated exploration into proposals
- maintain proposal metadata and task maps
- keep implementation notes aligned with execution
- archive to modules, architecture docs, and ADRs

## Workflow rules

- Exploration persists important findings before implementation.
- Proposal metadata is authoritative for lifecycle state.
- Task maps are authoritative for backend task decomposition.
- Notes are operational history, not a replacement for proposal/design changes.
- Archive targets must be updated before a proposal is marked archived.

## Default command intents

- `/tg:explore`: create or extend an exploration
- `/tg:propose`: create or revise a proposal from an exploration
- `/tg:apply`: create backend tasks and switch to active execution
- `/tg:archive`: complete archive updates and close the proposal
- `/tg:status`: summarize proposal and task state
- `/tg:list`: list proposals by lifecycle status
- `/tg:notes`: append implementation history

## Validation hooks

Before reporting a proposal as ready or archivable, use:

- `scripts/tg-create-proposal.js --title ... --source-exploration ... --archive-target ...`
- `scripts/tg-apply-proposal.js <proposal-id|proposal-dir>`
- `scripts/tg-validate-proposal.js <proposal-id|proposal-dir>`
- `scripts/tg-proposal-status.js <proposal-id|proposal-dir>`
- `scripts/tg-check-archive.js <proposal-id|proposal-dir>`
