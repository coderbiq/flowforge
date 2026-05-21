# Workflow Guide

This document is the high-level guide to `FlowForge`. The canonical operational details live under [`workflow/guides/`](../workflow/guides/lifecycle.md) and [`workflow/schema/`](../workflow/schema/proposal.schema.yaml).

## Canonical directories

```text
docs/
├── explorations/
├── proposals/
├── modules/
├── architecture/
└── decisions/
```

## Canonical artifacts

### Exploration

```text
docs/explorations/<slug>/
├── index.md
├── journal/
├── findings/
├── decisions/
└── artifacts/
```

### Proposal

```text
docs/proposals/<proposal-id>/
├── meta.yaml
├── proposal.md
├── design.md
├── task-map.md
└── notes.md
```

`task-map.md` should follow the canonical task-splitting rules in [`workflow/guides/task-splitting.md`](../workflow/guides/task-splitting.md).

## Lifecycle summary

1. `explore`: create durable findings before implementation
2. `propose`: convert validated exploration into a scoped proposal
3. `approve`: lock scope, backend, and archive targets
4. `apply`: create executable tasks from `task-map.md`
5. `implement`: execute tasks and keep notes current
6. `archive`: update the declared primary and secondary archive targets

## Archive behavior

- Module change: primary target in `docs/modules/`
- Cross-cutting or system design: primary target in `docs/architecture/`
- Stable architectural decision: ADR in `docs/decisions/`
- The archived knowledge base is the default reference corpus for future explorations.
- Archive is also the point where existing final docs are reconciled, not just appended to.
- If a change replaces or narrows an existing fact, the final doc should preserve a trace of the older fact in a history or changelog section.

## Proposal behavior

- Proposals start from the canonical corpus and record deltas against it.
- New proposals should identify the existing modules, architecture docs, and ADRs that form the baseline for the change.
- Proposal creation may infer the baseline corpus from archive targets and same-type final docs in the workspace, with explicit overrides for broader review sets.
- For complex modules, proposals should say which canonical doc remains the entry point and which subdocs carry the detailed knowledge.

## Command surface

- `/flowforge:explore`
- `/flowforge:propose`
- `/flowforge:upgrade`
- `/flowforge:approve`
- `/flowforge:apply`
- `/flowforge:archive`
- `/flowforge:status`
- `/flowforge:list`
- `/flowforge:notes`

Platform commands are wrappers. They should load workflow guidance instead of owning the business rules themselves.

## Script surface

- `.flowforge/scripts/flowforge-create-proposal.js`
- `.flowforge/scripts/flowforge-approve-proposal.js`
- `.flowforge/scripts/flowforge-apply-proposal.js`
- `.flowforge/scripts/flowforge-add-note.js`
- `.flowforge/scripts/flowforge-list-proposals.js`
- `.flowforge/scripts/flowforge-archive-proposal.js`
- `.flowforge/scripts/flowforge-validate-proposal.js`
- `.flowforge/scripts/flowforge-proposal-status.js`
- `.flowforge/scripts/flowforge-check-archive.js`
