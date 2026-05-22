---
doc_type: "note"
title: "Getting Started"
status: "draft"
workspace: "default"
module_scope: []
system_scope: []
convention_scope: []
ownership: []
information_class: "note"
topics: []
related_docs: []
archive_target: "default:GETTING-STARTED.md"
created: "2026-05-22T08:16:57.269Z"
updated: "2026-05-22T08:16:57.269Z"
---

# Getting Started

## Ownership summary

- Primary module: none
- System / architecture targets: none
- Convention targets: none
- Canonical reading path: GETTING-STARTED.md

## Requirements

- a source repository
- at least one supported AI coding environment
- optional: `Beads`
- optional: a memory provider such as `Memory MCP`

## 1. Install the workflow core and adapters

```bash
cd your-project
/path/to/flowforge/scripts/install.sh all
```

This installs:

- `.flowforge/`
- `.claude/` when requested
- `.opencode/` when requested
- `AGENTS.md` and `.codex/flowforge.md` when using `codex` or `all`

If the project already has FlowForge installed, use upgrade mode instead of reinstalling from scratch:

```bash
cd your-project
/path/to/flowforge/scripts/install.sh upgrade
```

Upgrade mode refreshes the managed workflow payload and platform command surfaces while preserving:

- `.flowforge/config.json`
- `.flowforge/state/`

## 2. Initialize project configuration

Create `.flowforge/config.json`:

```json
{
  "project": {
    "id": "your-project",
    "name": "Your Project",
    "slug": "your-project"
  },
  "paths": {
    "tool_root": ".flowforge",
    "state_root": ".flowforge/state"
  },
  "docs": {
    "default_workspace": "default",
    "workspaces": {
      "default": {
        "root": "docs",
        "scope": ".",
        "kind": "repository"
      }
    }
  },
  "task_backend": {
    "type": "beads"
  },
  "memory_provider": {
    "type": "memory-mcp",
    "enabled": false,
    "endpoint": "http://127.0.0.1:8000",
    "tags": ["project:your-project"]
  }
}
```

## 3. Create the docs layout

Recommended approach: start from the minimal project template under `workflow/templates/project/`.

At minimum, create:

```text
AGENTS.md
docs/explorations
docs/proposals
docs/modules
docs/architecture
docs/conventions
docs/decisions
.flowforge/state
```

If a workspace needs its own template variants, add a workspace-local template copy area such as:

```text
docs/flowforge/_templates
```

Use:

- `workflow/templates/project/AGENTS.md`
- `.flowforge/config.json` (bootstrapped from the project template)
- `workflow/templates/docs/`

## 4. Configure task management

Recommended backend: `Beads`

```bash
bd init
```

`task-map.md` is the bridge from proposals to the task backend. The workflow core does not store task logic in proposal prose.

Task maps should follow [`workflow/guides/task-splitting.md`](../workflow/guides/task-splitting.md): split by deliverable, not by file list, and insert explicit checkpoints for long-running proposals.

Template customization is copy-and-edit only. If a workspace needs a specialized model field layout, copy the default model template or the relevant part file into the workspace-local `_templates` area and adjust the copy there.

## 5. Configure memory

Local work-restoration state requires no external service. It is stored in `.flowforge/state/`.

Reusable experience memory is optional. If enabled, configure the provider in `.flowforge/config.json`.

Per-user authentication can be supplied through environment variables:

- `FLOWFORGE_MEMORY_ENDPOINT`
- `FLOWFORGE_MEMORY_API_KEY`

Platform-specific aliases are also supported for compatibility:

- Claude: `CLAUDE_FLOWFORGE_MEMORY_ENDPOINT`, `CLAUDE_FLOWFORGE_MEMORY_API_KEY`
- OpenCode: `OPENCODE_FLOWFORGE_MEMORY_ENDPOINT`, `OPENCODE_FLOWFORGE_MEMORY_API_KEY`
- Legacy alias: `OPENCODE_MEMORY_ENDPOINT`, `OPENCODE_MEMORY_API_KEY`

If one user needs different tokens for different projects, use the user-level config at `~/.config/flowforge/memory.json`:

```json
{
  "memory_provider": {
    "endpoint": "http://127.0.0.1:8000"
  },
  "projects": {
    "project-a": {
      "memory_provider": {
        "apiKey": "token-a"
      }
    },
    "project-b": {
      "memory_provider": {
        "endpoint": "https://memory.example.com",
        "apiKey": "token-b"
      }
    }
  }
}
```

## 6. Start the lifecycle

```text
/flowforge:explore "topic"
/flowforge:propose "proposal title"
/flowforge:upgrade
/flowforge:approve CR26052001
/flowforge:apply CR26052001
/flowforge:archive CR26052001
```

When opening an exploration, declare `ownership`, `expected_size_class`, and any `reusable_rules` in `index.md` frontmatter. Every exploration, proposal, design, model, and support doc carries its own frontmatter. `meta.yaml` remains the proposal bundle manifest; the generated proposal docs mirror the routing fields that matter to their own document-level indexing. See [`workflow/guides/doc-properties.md`](../workflow/guides/doc-properties.md) for the canonical property contract.

## 7. Create a proposal skeleton

Use the generator instead of creating proposal files by hand:

```bash
.flowforge/scripts/flowforge-create-proposal.js \
  --title "Example Proposal" \
  --source-exploration explorations/example-topic \
  --size-class large \
  --ownership module:modules/example-module:primary \
  --archive-target module:modules/example-module:primary \
  --archive-target convention:conventions/example-rule.md:secondary
```

This creates `docs/proposals/CRYYMMDDNN-<slug>/` with a valid initial skeleton.

## 8. Approve a proposal

Once the proposal content and archive targets are ready:

```bash
.flowforge/scripts/flowforge-approve-proposal.js CR26052001
```

This moves a valid proposal into `approved` state.

## 9. Apply an approved proposal

Once `meta.yaml` status is `approved`, apply it:

```bash
.flowforge/scripts/flowforge-apply-proposal.js CR26052001
```

For `Beads`, this creates one epic plus tasks from `task-map.md`, links dependencies, and moves the proposal to `active`.

## 10. Add AGENTS.md

`AGENTS.md` is the project-local contract that makes the workflow visible to any agent working in the repo.

Start from:

- `workflow/templates/project/AGENTS.md`

## 11. Operate and archive proposals

```bash
.flowforge/scripts/flowforge-add-note.js CR26052001 "Implemented API adapter and updated validation"
.flowforge/scripts/flowforge-list-proposals.js
.flowforge/scripts/flowforge-archive-proposal.js CR26052001
```

## 12. Validate proposals

Use the built-in checks while operating the workflow:

```bash
.flowforge/scripts/flowforge-validate-proposal.js CR26052001
.flowforge/scripts/flowforge-proposal-status.js CR26052001
.flowforge/scripts/flowforge-check-archive.js CR26052001
```

## Read next

- [Architecture](ARCHITECTURE.md)
- [Workflow Guide](PROPOSAL-WORKFLOW.md)
- [Lifecycle guide](../workflow/guides/lifecycle.md)
- [Sizing guide](../workflow/guides/sizing.md)
- [Ownership guide](../workflow/guides/ownership.md)
