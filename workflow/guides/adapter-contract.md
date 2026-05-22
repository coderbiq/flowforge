# Adapter Contract

Adapters and agent wrappers must preserve the canonical task-splitting rules defined in `task-splitting.md`.

Platform integrations must adapt `FlowForge` without redefining workflow behavior.

## Adapter responsibilities

- expose platform-native commands or prompts
- load workflow guides and templates
- load project rules using `scripts/flowforge-rules-context.js` and
  `workflow/guides/rule-loading.md`
- load project configuration
- update local state snapshots
- call task and memory providers
- keep install and upgrade entrypoints aligned with the same managed payload rules

Adapters may expose different entry surfaces:

- Claude Code and OpenCode: repo-local commands plus hooks/plugins
- Codex: project `AGENTS.md` plus workflow scripts

## Adapter must not

- fork lifecycle definitions
- invent platform-specific proposal states
- hardcode project tags or archive rules
- redefine document directory structures

## Required capabilities

### Workflow adapter

- `explore(topic)`
- `propose(exploration_path, title)`
- `upgrade(target_project)`
- `apply(proposal_id)`
- `archive(proposal_id)`
- `status(proposal_id)`
- `list()`
- `notes(proposal_id, entry)`

### Task backend adapter

- `create_epic(proposal)`
- `create_tasks(task_map)`
- `query_by_proposal(proposal_id)`
- `close_epic(epic_id)`

### Memory provider adapter

- `store(memory)`
- `search(query, filters)`
- `list_due_reviews()`
- `supersede(memory_id)`

## Configuration lookup

Adapters should resolve configuration in this order:

1. environment overrides
2. project-local `.flowforge/config.json`
3. user-level adapter config
4. built-in defaults

## Upgrade behavior

- Upgrade must preserve project-owned `.flowforge/config.json`
- Upgrade must preserve `.flowforge/state/`
- Platform command surfaces should describe the same upgrade boundary as the underlying scripts
