---
name: flowforge
description: |
  Workflow orchestration skill for FlowForge.

  Use this skill whenever the user wants to explore a new topic, create or refine a proposal, map work into tasks, archive completed work, or understand the current workflow state.
---

# FlowForge

This skill is a thin adapter over the canonical workflow specification in
`workflow/` and the installed project rule bundle when one exists.

## Routing first

Always start from the agent action routing contract:

- `workflow/guides/agent-action-routing.md`

The scenario decides which deeper guide matters next. Do not try to describe the
whole workflow as a fixed load order in this skill.

When editing any guide content, use the guide contract and validator as the
gate:

- `workflow/guides/guide-contract.md`
- `scripts/flowforge-validate-guides.js`

When metadata shape matters, use the schema guides:

- `workflow/schema/proposal.schema.yaml`
- `workflow/schema/exploration.schema.yaml`
- `workflow/schema/task-map.schema.yaml`

## Project rule bundle

If the current project contains `docs/flowforge/_rules/`, treat that bundle as
the project-default working policy on top of the canonical workflow.

- Materialize the bundle with `scripts/flowforge-rules-context.js` when you
  need the actual rule text in context.
- Use the bundle to refine analysis posture, writing posture, and archive
  emphasis.
- Do not let the project bundle override lifecycle, schema, or validation
  requirements from the core workflow.

## Context loading

Use the loader script as the single mechanism for merged rules context.
Combine it with intake or exploration context only when the scenario requires
those inputs.

- `scripts/flowforge-rules-context.js`
- `scripts/flowforge-intake-context.js`
- `scripts/flowforge-explore-context.js`

## Operating rules

- Exploration, proposal, task map, notes, and archive targets are revisitable
  artifacts. Update the artifact that matches the current scenario instead of
  treating the workflow as one-way.
- Proposal metadata remains authoritative for lifecycle state.
- Task maps remain authoritative for backend task decomposition.
- Notes are operational history, not a replacement for proposal or design
  changes.
- Archive targets must be updated before a proposal is marked archived.
- `apply` may promote a proposed proposal inline and then move straight into
  execution when the scenario allows it.

## Command surface

- `/flowforge:explore`, `/flowforge:intake`, `/flowforge:propose`,
  `/flowforge:approve`, `/flowforge:apply`, `/flowforge:archive`,
  `/flowforge:status`, `/flowforge:list`, and `/flowforge:notes` are thin
  command adapters over the scripts in `.flowforge/scripts/`.

## Validation

Before reporting a proposal as ready or archivable, use the matching
`.flowforge/scripts/flowforge-*.js` helpers for creation, context loading,
approval, apply, note updates, validation, status, and archive checks.

When modifying `workflow/guides/*.md`, run the guide validator and keep the
document action-oriented. If the file starts to read like reference prose,
move that material out of `workflow/guides/`.
