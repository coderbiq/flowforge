# Architecture

`tg-workflow` is intentionally split into four layers:

1. canonical workflow spec
2. canonical agent definitions
3. platform adapters
4. project artifacts

## Why this split exists

Earlier iterations mixed workflow design with Claude/OpenCode implementation details. That caused drift in proposal states, exploration layouts, and memory behavior. The current architecture fixes that by making the workflow spec platform-agnostic and forcing adapters to load it instead of re-declaring it.

## Layer model

### 1. Canonical workflow spec

Lives in [`workflow/`](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/README.md).

Contains:

- lifecycle rules
- metadata schemas
- archive rules
- adapter contracts
- canonical templates

This is the source of truth for all business semantics.

### 2. Agent definitions

Lives in `agents/skills/`.

Contains:

- canonical skill definitions
- agent-facing workflow semantics
- references to workflow guides and schemas

These files are the portable source for agent behavior and should not be duplicated by platform.

### 3. Platform adapters

Lives in `configs/`.

Contains:

- command entrypoints
- skill wrappers
- hooks and plugins
- adapter-local configuration resolution

Adapter note:

- Claude Code and OpenCode use repo-local command surfaces
- Codex uses project `AGENTS.md` plus workflow scripts

Adapters must not invent states, directory structures, or archive semantics.

### 4. Project artifacts

Lives in the target project after installation.

Contains:

- `docs/explorations/`
- `docs/proposals/`
- `docs/modules/`
- `docs/architecture/`
- `docs/decisions/`
- `.workflow/state/`

## Data model

### Local state memory

Purpose:

- restore active work quickly
- remember current focus, touched files, and next steps

Storage:

- `.workflow/state/active-session.json`
- `.workflow/state/sessions/*.json`
- `.workflow/state/workstreams/*.json`

This layer is operational, not semantic.

### External memory provider

Purpose:

- store reusable decisions, architecture insights, debugging knowledge, and workflow preferences

Contract:

- `store`
- `search`
- `list_due_reviews`
- `supersede`

`Memory MCP` is a default provider, not a hardcoded design assumption.

### Task backend

Purpose:

- manage dependency-aware execution
- expose agent-friendly ready work

Default backend:

- `Beads`

Contract:

- `create_epic`
- `create_tasks`
- `query_by_proposal`
- `close_epic`

## Archive model

The default archive view is module-first, but not module-only.

- Module-scoped change: archive primarily to `docs/modules/<module>/`
- Cross-cutting or system design: archive primarily to `docs/architecture/`
- Stable high-cost decision: additionally record an ADR in `docs/decisions/`

This resolves the earlier problem where every change was forced through a module lens even when the real artifact was architectural.

## External references

The current model borrows good ideas from:

- OpenSpec: staged explore/propose/apply flow
- ADR/MADR: durable decision records
- C4/documentation-as-code: architecture views as maintained text artifacts
- Beads: dependency-aware task execution for AI agents
