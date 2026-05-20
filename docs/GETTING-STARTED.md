# Getting Started

This guide installs `tg-workflow` into a project with both the canonical workflow core and the optional platform adapters.

## Requirements

- a source repository
- at least one supported AI coding environment
- optional: `Beads`
- optional: a memory provider such as `Memory MCP`

## 1. Install the workflow core and adapters

```bash
cd your-project
/path/to/tg-workflow/scripts/install.sh all
```

This installs:

- `workflow/`
- `agents/`
- `.claude/` when requested
- `.opencode/` when requested
- `AGENTS.md` and `.codex/tg-workflow.md` when using `codex` or `all`

## 2. Initialize project configuration

Create `workflow/config.json`:

```json
{
  "project": {
    "id": "your-project",
    "name": "Your Project",
    "slug": "your-project"
  },
  "paths": {
    "docs_root": "docs",
    "state_root": ".workflow/state"
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
docs/decisions
.workflow/state
```

Use:

- `workflow/templates/project/AGENTS.md`
- `workflow/templates/project/workflow/config.json`
- `workflow/templates/docs/`

## 4. Configure task management

Recommended backend: `Beads`

```bash
bd init
```

`task-map.md` is the bridge from proposals to the task backend. The workflow core does not store task logic in proposal prose.

## 5. Configure memory

Local work-restoration state requires no external service. It is stored in `.workflow/state/`.

Reusable experience memory is optional. If enabled, configure the provider in `workflow/config.json`.

## 6. Start the lifecycle

```text
/tg:explore "topic"
/tg:propose "proposal title"
/tg:approve CR26052001
/tg:apply CR26052001
/tg:archive CR26052001
```

## 7. Create a proposal skeleton

Use the generator instead of creating proposal files by hand:

```bash
scripts/tg-create-proposal.js \
  --title "Example Proposal" \
  --source-exploration docs/explorations/example-topic \
  --archive-target module:docs/modules/example-module:primary \
  --archive-target architecture:docs/architecture/system-overview.md:secondary
```

This creates `docs/proposals/CRYYMMDDNN-<slug>/` with a valid initial skeleton.

## 8. Approve a proposal

Once the proposal content and archive targets are ready:

```bash
scripts/tg-approve-proposal.js CR26052001
```

This moves a valid proposal into `approved` state.

## 9. Apply an approved proposal

Once `meta.yaml` status is `approved`, apply it:

```bash
scripts/tg-apply-proposal.js CR26052001
```

For `Beads`, this creates one epic plus tasks from `task-map.md`, links dependencies, and moves the proposal to `active`.

## 10. Add AGENTS.md

`AGENTS.md` is the project-local contract that makes the workflow visible to any agent working in the repo.

Start from:

- `workflow/templates/project/AGENTS.md`

## 11. Operate and archive proposals

```bash
scripts/tg-add-note.js CR26052001 "Implemented API adapter and updated validation"
scripts/tg-list-proposals.js
scripts/tg-archive-proposal.js CR26052001
```

## 12. Validate proposals

Use the built-in checks while operating the workflow:

```bash
scripts/tg-validate-proposal.js CR26052001
scripts/tg-proposal-status.js CR26052001
scripts/tg-check-archive.js CR26052001
```

## Read next

- [Architecture](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/docs/ARCHITECTURE.md)
- [Workflow Guide](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/docs/PROPOSAL-WORKFLOW.md)
- [Lifecycle guide](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/guides/lifecycle.md)
