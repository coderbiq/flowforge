# Proposal Ownership

`FlowForge` requires every exploration and proposal to declare ownership tags. Ownership controls where knowledge is archived, not which team owns the work. The same ownership graph is also surfaced in each document's YAML frontmatter so Obsidian and other indexers can route the document without parsing the proposal bundle first.

## Ownership types

- `module`: the change belongs to a specific module. Final knowledge archives under `docs/modules/<module>/`.
- `system`: the change belongs to a system-wide architectural concern. Final knowledge archives under `docs/architecture/<topic>.md`.
- `cross-module`: the change spans multiple modules and produces shared behavior. Final knowledge archives under `docs/architecture/<topic>.md` and each affected module records the impact in its history.
- `convention`: the change establishes or revises a reusable rule, pattern, or policy. Final knowledge archives under `docs/conventions/<topic>.md`.

A single proposal may declare multiple ownership entries. Most non-trivial proposals will. A new module proposal, for example, often has:

- one `module` ownership entry for the new module
- one or more `convention` ownership entries for new shared rules it introduces
- one `system` ownership entry if it modifies architecture-level boundaries

## Ownership entry shape

Each ownership entry has three fields:

- `type`: one of `module`, `system`, `cross-module`, `convention`
- `target`: the canonical archive reference, relative to the workspace docs root
- `role`: `primary` or `secondary`

`primary` is where a future reader should start. `secondary` is preserved for cross-referencing, traceability, and parallel reading paths.

There must be exactly one ownership entry with `role: primary`.

## Human-readable ownership summary

Machine-readable `ownership` in metadata is not enough on its own. The exploration and proposal reading surfaces must also summarize ownership in plain language so a reader can immediately answer:

- which module this work belongs to
- whether it also belongs to a system or architecture surface
- whether it introduces or updates reusable conventions

At minimum, the human-readable docs should surface:

- owning modules
- system or cross-module targets
- reusable convention targets
- the primary reading path

## Conventions as a first-class category

Conventions are reusable consensus rules that survive beyond the proposal that introduced them. Typical examples:

- "this class of problem is solved with this standard approach"
- "this kind of field uses this storage shape"
- "this layer must depend only on these modules"
- "this artifact must use this naming pattern"

A convention is not a module behavior and not a one-off architecture diagram. It is a rule that applies whenever the matching situation appears in the codebase.

When a proposal introduces a convention, the convention must be archived under `docs/conventions/<topic>.md` and not embedded only inside a module or architecture document.

## Relationship to archive targets

`ownership` and `archive_targets` are aligned. Every ownership entry should resolve to a matching archive target. Document frontmatter mirrors the same graph at the document level, but `meta.yaml` remains the proposal bundle contract.

- ownership `module` maps to archive target `module`
- ownership `system` maps to archive target `architecture`
- ownership `cross-module` maps to archive target `architecture` plus per-module `history` updates
- ownership `convention` maps to archive target `convention`

A proposal must declare at least one `primary` archive target that corresponds to its `primary` ownership.

## Exploration ownership

Explorations declare ownership too. This lets the proposal phase inherit the ownership graph instead of rediscovering it, while each exploration file still carries its own frontmatter for Obsidian indexing.

When an exploration spans multiple ownership types, the resulting proposals may either:

- declare the same ownership graph, or
- split into multiple proposals that each carry a subset of the ownership graph.

Use the second approach when the change scope is naturally separable.

## Anti-patterns

- Tagging every change as `module` because that is the default archive path.
- Hiding convention-grade rules inside a module design doc and never extracting them.
- Declaring `cross-module` ownership without identifying which modules are affected.
- Declaring `primary` ownership on more than one entry. There must be exactly one primary entry per ownership graph.
