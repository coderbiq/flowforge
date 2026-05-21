# Monorepo Document Workspaces

- Status: active
- Primary proposal: CR26052001

## Scope

`FlowForge` uses `document workspace` as the stable unit for document ownership, proposal routing, and archive targeting.

This document defines how workspace configuration, canonical corpus selection, command routing, and archive maintenance work together in single-workspace projects and monorepos.

## Workspace model

Each workspace carries four core responsibilities:

- document root
- code scope
- lifecycle routing identity
- canonical corpus boundary

### Configuration fields

- `paths.tool_root`
  - project-local tool root, recommended as `.flowforge`
- `docs.default_workspace`
  - default workspace name when routing is ambiguous
- `docs.workspaces.<name>.root`
  - documents root for the workspace
- `docs.workspaces.<name>.scope`
  - code scope associated with the workspace
- `docs.workspaces.<name>.kind`
  - `repository` or `project`
- `docs.workspaces.<name>.label`
  - optional display name
- `docs.workspaces.<name>.owners`
  - optional ownership metadata

### Compatibility rule

- If a project does not define `docs.workspaces`, `FlowForge` constructs a default `default` workspace rooted at `docs/`
- Simple projects continue to work without extra configuration
- Monorepos must declare workspaces explicitly

### Constraints

- Only one `repository` workspace should exist
- workspace roots must not nest inside each other
- nested code scopes are allowed
- scope resolution prefers the deepest matching workspace

## Canonical corpus

The final knowledge base is not a terminal artifact. It is the baseline for later exploration and later proposals.

### Canonical sources

- modules
- architecture docs
- ADRs

### Maintenance rules

- proposals record deltas against the canonical corpus instead of rewriting it from scratch
- when existing facts change, update the owning doc in place
- when new facts are added, place them in the nearest stable canonical doc
- when facts are superseded, preserve the old fact in history or changelog sections
- for large modules, keep one canonical entry point and split details into linked subdocs

## Command routing

### Resolution order

1. explicit `--workspace`
2. proposal or exploration metadata `workspace`
3. `cwd` scope match
4. `docs.default_workspace`

### Runtime behavior

- `create-proposal` should infer a workspace when not explicitly provided
- `list-proposals` should default to the current workspace and support `--all-workspaces`
- `proposal-status` should resolve both proposal id and proposal dir in a workspace-aware way
- `archive-proposal` should update the declared archive targets under their owning workspace

## Archive behavior

Archive is a knowledge-base maintenance pass.

- primary archive targets define the main reading path
- secondary archive targets preserve alternate reading paths
- archive should update linked final docs together so the corpus stays coherent
- archive should preserve historical facts in `history.md`, changelog sections, or ADR context where applicable

## Views to maintain

- context view for workspace placement
- container view for workspace-to-doc relationships
- sequence or flow view for lifecycle and archive transitions

