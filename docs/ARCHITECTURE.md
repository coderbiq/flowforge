---
doc_type: "note"
title: "Architecture"
status: "draft"
workspace: "default"
module_scope: []
system_scope: []
convention_scope: []
ownership: []
information_class: "note"
topics: []
related_docs: []
archive_target: "default:ARCHITECTURE.md"
created: "2026-05-22T08:16:57.269Z"
updated: "2026-05-22T08:16:57.269Z"
---

# Architecture

## Ownership summary

- Primary module: none
- System / architecture targets: none
- Convention targets: none
- Canonical reading path: ARCHITECTURE.md

## Why this split exists

Earlier iterations mixed workflow design with Claude/OpenCode implementation details. That caused drift in proposal states, exploration layouts, and memory behavior. The current architecture fixes that by making the workflow spec platform-agnostic and forcing adapters to load it instead of re-declaring it.

## Layer model

### 1. Canonical workflow spec

Lives in [`workflow/`](../workflow/README.md).

Contains:

- lifecycle rules
- sizing rules
- ownership rules
- task splitting and checkpoint rules
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
- install and upgrade surfaces that keep managed payloads in sync

Adapter note:

- Claude Code and OpenCode use repo-local command surfaces, including upgrade wrappers
- Codex uses project `AGENTS.md` plus workflow scripts

Adapters must not invent states, directory structures, or archive semantics.

### 4. Project artifacts

Lives in the target project after installation.

Contains:

- `docs/explorations/`
- `docs/proposals/`
- `docs/modules/`
- `docs/architecture/`
- `docs/conventions/`
- `docs/decisions/`
- `.flowforge/state/`

## Data model

### Local state memory

Purpose:

- restore active work quickly
- remember current focus, touched files, and next steps

Storage:

- `.flowforge/state/active-session.json`
- `.flowforge/state/sessions/*.json`
- `.flowforge/state/workstreams/*.json`

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

## Classification model

Every exploration, proposal, and durable subdocument carries two classification axes in its own YAML frontmatter:

- `size_class`: `small | medium | large` — controls the document skeleton (see `workflow/guides/sizing.md`).
- `ownership`: one or more of `module | system | cross-module | convention` — controls the archive destination (see `workflow/guides/ownership.md`).

These axes are independent. A `small` proposal can still introduce a `convention` archive target, and a `large` module proposal can still carry zero conventions.

These axes also need a human-readable mirror in the document bodies. Readers should not have to reconstruct module or architecture ownership by inspecting only `meta.yaml`.
`meta.yaml` remains the proposal bundle manifest, but document frontmatter is the document-level contract for Obsidian indexing and doc-local routing.

## Archive model

The archive view has four first-class destinations:

- Module-scoped change: archive primarily to `docs/modules/<module>/`
- Cross-cutting or system design: archive primarily to `docs/architecture/<topic>.md`
- Reusable rule or consensus standard: archive primarily to `docs/conventions/<topic>.md`
- Stable high-cost decision: additionally record an ADR in `docs/decisions/`

This resolves the earlier problem where every change was forced through a module lens even when the real artifact was architectural or a shared convention.

## External references

The current model borrows good ideas from:

- OpenSpec: staged explore/propose/apply flow
- ADR/MADR: durable decision records
- C4/documentation-as-code: architecture views as maintained text artifacts
- Beads: dependency-aware task execution for AI agents
