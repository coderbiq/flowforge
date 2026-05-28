# FlowForge

`FlowForge` is a docs-first workflow for AI-assisted software design and delivery.

It is built around one idea: durable artifacts should outlive chat sessions and platform integrations. Exploration, proposals, task mapping, implementation notes, archive targets, and reusable experience are all first-class files.

Core workflow rules live in [`workflow/`](workflow/README.md). Canonical skill definitions live in [`agents/skills/`](agents/skills/flowforge/SKILL.md). Platform adapters under `configs/` only expose commands, hooks, plugins, or platform-specific entry guidance.

Codex support uses project `AGENTS.md` plus workflow scripts rather than a repo-local slash-command registry.

## Repository layout

```text
FlowForge/
├── workflow/                   # canonical lifecycle, schemas, templates
├── agents/                     # canonical agent-facing skill definitions
├── configs/                    # platform adapters
├── docs/                       # tool documentation
└── scripts/                    # installation helpers
```

## Core model

- Workflow lifecycle: `explore -> propose -> approve -> apply -> implement -> archive`
- Primary docs root in a target project: `docs/`
- Local work-restoration state: `.flowforge/state/`
- Default task backend: `Beads`
- Default external memory provider: `Memory MCP` through a provider interface

## Install

```bash
./scripts/install.sh all
./scripts/install.sh claude
./scripts/install.sh opencode
./scripts/install.sh codex
./scripts/install.sh global
```

Project installation copies:

- `.flowforge/workflow/`
- `.flowforge/agents/`
- `.flowforge/adapters/`
- requested platform adapters

Codex installation adds:

- root `AGENTS.md` when absent
- `.codex/flowforge.md` adapter notes
- workflow scripts under `.flowforge/scripts/`

## Commands

The default adapter surface remains:

```text
/flowforge:explore "topic"
/flowforge:propose "proposal title"
/flowforge:approve CR26052001
/flowforge:apply CR26052001
/flowforge:archive CR26052001
```

These commands now load the canonical workflow guides instead of embedding their own business rules.

## Read next

- [Architecture](docs/ARCHITECTURE.md)
- [Workflow Guide](docs/PROPOSAL-WORKFLOW.md)
- [Getting Started](docs/GETTING-STARTED.md)
