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

## Command surface

- `/flowforge:explore`
- `/flowforge:propose`
- `/flowforge:approve`
- `/flowforge:apply`
- `/flowforge:archive`
- `/flowforge:status`
- `/flowforge:list`
- `/flowforge:notes`

Platform commands are wrappers. They should load workflow guidance instead of owning the business rules themselves.

## Script surface

- `scripts/flowforge-create-proposal.js`
- `scripts/flowforge-approve-proposal.js`
- `scripts/flowforge-apply-proposal.js`
- `scripts/flowforge-add-note.js`
- `scripts/flowforge-list-proposals.js`
- `scripts/flowforge-archive-proposal.js`
- `scripts/flowforge-validate-proposal.js`
- `scripts/flowforge-proposal-status.js`
- `scripts/flowforge-check-archive.js`
