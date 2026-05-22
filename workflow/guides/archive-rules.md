# Archive Rules

Archive behavior is driven by proposal metadata, not by platform adapters.

## Target types

- `module`: default for bounded changes inside one module
- `architecture`: required for system-wide or cross-module design changes
- `convention`: required when the proposal establishes or revises a reusable rule, pattern, or policy
- `decision`: required when the proposal introduces or supersedes stable architectural decisions

## Required archive steps

1. Confirm proposal status is `implemented`.
2. Confirm task backend has no open tasks for the proposal.
3. Update the primary archive target.
4. Update any secondary archive targets.
5. Promote validated `reusable_rules` from the source exploration into `docs/conventions/` if not already archived.
6. Record superseded decisions if applicable.
7. Set proposal status to `archived`.

## Typical mappings

- New module or major module change:
  - primary target: `docs/modules/<module>/`
- Cross-cutting architecture work:
  - primary target: `docs/architecture/<topic>.md`
  - secondary targets: impacted module docs
- Reusable rule or shared convention:
  - primary target: `docs/conventions/<topic>.md`
  - secondary targets: modules or architecture docs that must reference it
- Stable technical decision:
  - secondary target: `docs/decisions/ADR-*.md`

## Ownership alignment

Each ownership entry on the proposal should map to an archive target:

- ownership `module` → archive `module`
- ownership `system` → archive `architecture`
- ownership `cross-module` → archive `architecture` plus `history` updates in each affected module
- ownership `convention` → archive `convention`

The primary ownership and the primary archive target should describe the same canonical destination.

## Anti-patterns

- Archiving only the proposal directory and skipping target docs
- Writing architecture decisions only in implementation notes
- Treating modules as the only valid final documentation view
- Burying convention-grade rules inside a module design doc instead of promoting them to `docs/conventions/`
