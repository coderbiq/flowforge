---
name: flowforge
description: |
  Workflow orchestration skill for FlowForge.

  Use this skill whenever the user wants to explore a new topic, create or refine a proposal, map work into tasks, archive completed work, or understand the current workflow state.
---

# FlowForge

This skill is a thin adapter over the canonical workflow specification in
`workflow/` and the installed project rule bundle when one exists.

## Source of truth

Always load and follow:

- `workflow/guides/rule-loading.md`
- `workflow/guides/intake-packages.md`
- `workflow/guides/lifecycle.md`
- `workflow/guides/authoring-rules.md`
- `workflow/guides/archive-rules.md`
- `workflow/guides/adapter-contract.md`

When metadata shape matters, use:

- `workflow/schema/proposal.schema.yaml`
- `workflow/schema/exploration.schema.yaml`
- `workflow/schema/task-map.schema.yaml`

## Responsibilities

- create and maintain intake packages
- create and maintain exploration artifacts
- convert validated exploration into proposals
- maintain proposal metadata and task maps
- keep implementation notes aligned with execution
- archive to modules, architecture docs, and ADRs

## Project rule bundle

If the current project contains `docs/flowforge/_rules/`, treat that bundle as
the project-default working policy on top of the canonical workflow.

- Materialize the bundle with `scripts/flowforge-rules-context.js` when you
  need the actual rule text in context, following
  `workflow/guides/rule-loading.md`.
- Use the bundle to refine analysis posture, writing posture, and archive
  emphasis.
- Do not let the project bundle override lifecycle, schema, or validation
  requirements from the core workflow.

## Workflow rules

- Exploration persists important findings before execution.
- Proposal metadata is authoritative for lifecycle state.
- Task maps are authoritative for backend task decomposition.
- Notes are operational history, not a replacement for proposal/design changes.
- Archive targets must be updated before a proposal is marked archived.
- Approve is an optional review gate; apply can promote a proposed proposal inline and then move straight into execution.

## Default command intents

- `/flowforge:explore`: create or extend an exploration
- `/flowforge:intake`: create or revise an intake package before exploration
- `/flowforge:propose`: create or revise a proposal from an exploration
- `/flowforge:approve`: optional review gate that can move a valid proposal into approved state
- `/flowforge:apply`: create backend tasks, auto-promote a proposed proposal if needed, and switch straight into active execution
- `/flowforge:archive`: complete archive updates and close the proposal
- `/flowforge:status`: summarize proposal and task state
- `/flowforge:list`: list proposals by lifecycle status
- `/flowforge:notes`: append implementation history

## Validation hooks

Before reporting a proposal as ready or archivable, use:

- `.flowforge/scripts/flowforge-create-proposal.js --title ... --source-exploration ... --archive-target ...`
- `.flowforge/scripts/flowforge-create-intake.js --title ...`
- `.flowforge/scripts/flowforge-intake-context.js <intake-slug|intake-dir>`
- `.flowforge/scripts/flowforge-explore-context.js <intake-slug|intake-dir>`
- `.flowforge/scripts/flowforge-approve-proposal.js <proposal-id|proposal-dir>`
- `.flowforge/scripts/flowforge-apply-proposal.js <proposal-id|proposal-dir>`
- `.flowforge/scripts/flowforge-add-note.js <proposal-id|proposal-dir> <note text>`
- `.flowforge/scripts/flowforge-list-proposals.js`
- `.flowforge/scripts/flowforge-archive-proposal.js <proposal-id|proposal-dir>`
- `.flowforge/scripts/flowforge-validate-proposal.js <proposal-id|proposal-dir>`
- `.flowforge/scripts/flowforge-proposal-status.js <proposal-id|proposal-dir>`
- `.flowforge/scripts/flowforge-check-archive.js <proposal-id|proposal-dir>`
