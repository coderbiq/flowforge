# Archive Rules

Archive behavior is driven by proposal metadata, not by platform adapters.

## Target types

- `module`: default for bounded changes in one area of the system
- `architecture`: required for cross-module or system-wide design changes
- `decision`: required when the proposal introduces or supersedes stable architectural decisions

## Required archive steps

1. Confirm proposal status is `implemented`.
2. Confirm task backend has no open tasks for the proposal.
3. Update the primary archive target.
4. Update any secondary archive targets.
5. Record superseded decisions if applicable.
6. Set proposal status to `archived`.

## Typical mappings

- New module or major module change:
  - primary target: `docs/modules/<module>/`
- Cross-cutting architecture work:
  - primary target: `docs/architecture/<topic>.md`
  - secondary target(s): impacted module docs
- Stable technical decision:
  - secondary target: `docs/decisions/ADR-*.md`

## Anti-patterns

- Archiving only the proposal directory and skipping target docs
- Writing architecture decisions only in implementation notes
- Treating modules as the only valid final documentation view
