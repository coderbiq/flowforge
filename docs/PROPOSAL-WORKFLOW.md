# Workflow Guide

This document is the high-level guide to `FlowForge`. The canonical operational details live under [`workflow/guides/`](../workflow/guides/lifecycle.md) and [`workflow/schema/`](../workflow/schema/proposal.schema.yaml).

## Canonical directories

```text
docs/
├── explorations/
├── proposals/
├── modules/
├── architecture/
├── conventions/
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

`index.md` must declare `ownership`, `expected_size_class`, and any `reusable_rules` surfaced during exploration. It must also summarize that ownership in plain language so the reader can see the owning module, any system or architecture target, and any convention target without parsing metadata syntax.

## Template customization

Default workflow templates are reference defaults. If a workspace needs a project-specific variant, copy the whole template or the relevant section files into a workspace-local `_templates/` directory and edit those copies directly.

Example:

```text
docs/flowforge/_templates/
```

The workflow does not do automatic merge or override behavior. Template guidance should live in the copied files themselves so an agent can understand why the local version differs from the default.

### Proposal (small / medium single-file)

```text
docs/proposals/<proposal-id>/
├── meta.yaml
├── proposal.md
├── design.md
├── task-map.md
└── notes.md
```

### Proposal (medium split / large)

```text
docs/proposals/<proposal-id>/
├── meta.yaml
├── proposal.md
├── design/
│   ├── README.md
│   ├── architecture.md
│   ├── model.md
│   ├── lifecycle.md
│   ├── flow.md          (optional)
│   ├── api.md           (optional)
│   ├── constraints.md   (optional)
│   └── tradeoffs.md     (optional)
├── model/
│   ├── README.md
│   └── <ModelName>.md
├── task-map.md
└── notes.md
```

- `medium` proposals may use either layout. When two or more business models are introduced, the `model/` directory is required regardless of which design layout is chosen.
- `large` proposals must use the directory layout and must include one document per core business model under `model/`.
- `meta.yaml` sets `links.design` to `design.md` or `design/README.md` accordingly, and may set `links.model` to `model/README.md`.

`task-map.md` should follow the canonical task-splitting rules in [`workflow/guides/task-splitting.md`](../workflow/guides/task-splitting.md).

## Classification

Every exploration and proposal declares:

- `size_class`: `small`, `medium`, or `large`. See [`workflow/guides/sizing.md`](../workflow/guides/sizing.md).
- `ownership`: one or more entries of type `module`, `system`, `cross-module`, or `convention`. See [`workflow/guides/ownership.md`](../workflow/guides/ownership.md).

These two fields determine the document skeleton and the archive destination. They must be locked before design starts.

Human-readable docs must not leave this information only in machine metadata. Exploration, proposal, design, task-map, and model entry docs should all surface an ownership summary that answers:

- what module this belongs to
- what system or architecture surface it affects
- what reusable conventions it introduces or updates

## Lifecycle summary

1. `explore`: create durable findings, declare ownership and expected size
2. `propose`: convert validated exploration into a scoped proposal, lock size class and ownership
3. `approve`: lock scope, backend, and archive targets
4. `apply`: create executable tasks from `task-map.md`
5. `implement`: execute tasks and keep notes current
6. `archive`: update the declared primary and secondary archive targets, promote validated reusable rules into `docs/conventions/`

## Archive behavior

- Module change: primary target in `docs/modules/`
- Cross-cutting or system design: primary target in `docs/architecture/`
- Reusable rule or shared convention: primary target in `docs/conventions/`
- Stable architectural decision: ADR in `docs/decisions/`
- The archived knowledge base is the default reference corpus for future explorations.
- Archive is also the point where existing final docs are reconciled, not just appended to.
- If a change replaces or narrows an existing fact, the final doc should preserve a trace of the older fact in a history or changelog section.

## Proposal behavior

- Proposals start from the canonical corpus and record deltas against it.
- New proposals should identify the existing modules, architecture docs, conventions, and ADRs that form the baseline for the change.
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
- `.flowforge/scripts/flowforge-validate-exploration.js`
- `.flowforge/scripts/flowforge-approve-proposal.js`
- `.flowforge/scripts/flowforge-apply-proposal.js`
- `.flowforge/scripts/flowforge-add-note.js`
- `.flowforge/scripts/flowforge-list-proposals.js`
- `.flowforge/scripts/flowforge-archive-proposal.js`
- `.flowforge/scripts/flowforge-validate-proposal.js`
- `.flowforge/scripts/flowforge-proposal-status.js`
- `.flowforge/scripts/flowforge-check-archive.js`
