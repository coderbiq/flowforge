# Adapter Contract

Platform integrations must adapt `tg-workflow` without redefining workflow behavior.

## Adapter responsibilities

- expose platform-native commands or prompts
- load workflow guides and templates
- load project configuration
- update local state snapshots
- call task and memory providers

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
2. project-local `workflow/config.json`
3. user-level adapter config
4. built-in defaults
