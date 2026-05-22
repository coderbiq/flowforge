# Proposal Sizing

`FlowForge` requires every proposal to declare a `size_class` before design work starts. Sizing controls how the design surface is structured, not how important the change is. The same `size_class` is mirrored in proposal frontmatter so the design surface can be indexed independently from `meta.yaml`.

## Size classes

- `small`: a localized change inside a single module. No new model family, no new lifecycle, no new convention.
- `medium`: a non-trivial change inside one module, or a focused change that touches a small set of related modules. May introduce a few new fields, flows, or partial conventions, but does not redefine the module.
- `large`: a new module, a cross-module redesign, a new business model family, a new lifecycle, or a new architecture-level convention.

A proposal is `large` if any one of these is true:

- introduces a new module
- introduces a new business model family or restructures an existing one
- introduces or replaces a lifecycle (state machine, audit model, validation gate)
- establishes or supersedes an architecture-level convention
- requires coordinated changes across more than one module

A proposal is `medium` if it materially changes a module without meeting any `large` trigger, for example adding a new flow, a new bounded capability, or a new persistence surface.

A proposal is `small` if it only adjusts existing behavior, fields, or wording inside one module without changing model boundaries, lifecycle, or conventions.

## Design surface per size

### small

- Single-file `design.md`.
- Optional sections may be omitted: `architecture`, `lifecycle`, `flow`, `tradeoffs`.
- Models, if mentioned, may be described inline in `design.md`.

### medium

- Default is single-file `design.md`.
- A proposal may opt into the directory layout (`design/` plus `model/`) when the change spans several concerns or introduces two or more new models.
- When more than one business model is introduced, a `model/` directory becomes mandatory regardless of whether `design.md` stays single-file.

### large

- `design/` directory is mandatory. `design.md` at the proposal root is not used.
- `model/` directory is mandatory. Every core business model gets its own document.
- The `design/` directory must contain at minimum `README.md`, `architecture.md`, `model.md`, and `lifecycle.md`. Other subdocs (`flow.md`, `api.md`, `constraints.md`, `tradeoffs.md`) are added when the proposal needs them.

## Choosing the size class

The size class is declared in `meta.yaml` (`size_class`) and mirrored in proposal frontmatter and `proposal.md`.

Explorations should predict the size class via `expected_size_class` in `index.md`. The proposal author may override it when the scope changes between exploration and proposal.

A size class may only be revised by re-running the `propose` phase, not silently during `implement`. Down-classing requires removing the larger surface; up-classing requires creating the larger surface and migrating existing content.

## Anti-patterns

- Treating size as a quality signal. Size only describes the design surface, not the change's importance.
- Forcing every proposal to `large` because the module is large. Size reflects this change, not the surrounding module.
- Keeping a proposal at `small` to avoid writing `model/` docs when the change introduces multiple models. Add the docs or split the proposal.
