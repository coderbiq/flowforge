# tg-workflow

`tg-workflow` is a docs-first workflow for AI-assisted software design and delivery.

It is built around one idea: durable artifacts should outlive chat sessions and platform integrations. Exploration, proposals, task mapping, implementation notes, archive targets, and reusable experience are all first-class files.

Core workflow rules live in [`workflow/`](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/README.md). Canonical skill definitions live in [`agents/skills/`](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/agents/skills/tg-workflow/SKILL.md). Platform adapters under `configs/` only expose commands, hooks, plugins, or platform-specific entry guidance.

Codex support uses project `AGENTS.md` plus workflow scripts rather than a repo-local slash-command registry.

## Repository layout

```text
tg-workflow/
├── workflow/                   # canonical lifecycle, schemas, templates
├── agents/                     # canonical agent-facing skill definitions
├── configs/                    # platform adapters
├── docs/                       # tool documentation
└── scripts/                    # installation helpers
```

## Core model

- Workflow lifecycle: `explore -> propose -> approve -> apply -> implement -> archive`
- Primary docs root in a target project: `docs/`
- Local work-restoration state: `.workflow/state/`
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

- `workflow/`
- `agents/`
- requested platform adapters

Codex installation adds:

- root `AGENTS.md` when absent
- `.codex/tg-workflow.md` adapter notes
- workflow scripts under `scripts/`

## Commands

The default adapter surface remains:

```text
/tg:explore "topic"
/tg:propose "proposal title"
/tg:approve CR26052001
/tg:apply CR26052001
/tg:archive CR26052001
```

These commands now load the canonical workflow guides instead of embedding their own business rules.

## Read next

- [Architecture](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/docs/ARCHITECTURE.md)
- [Workflow Guide](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/docs/PROPOSAL-WORKFLOW.md)
- [Getting Started](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/docs/GETTING-STARTED.md)
