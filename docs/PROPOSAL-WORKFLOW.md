# Workflow Guide

This document is the high-level guide to `tg-workflow`. The canonical operational details live under [`workflow/guides/`](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/guides/lifecycle.md) and [`workflow/schema/`](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/schema/proposal.schema.yaml).

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

- `/tg:explore`
- `/tg:propose`
- `/tg:apply`
- `/tg:archive`
- `/tg:status`
- `/tg:list`
- `/tg:notes`

Platform commands are wrappers. They should load workflow guidance instead of owning the business rules themselves.
