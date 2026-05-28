# Codex Adapter

Codex uses `FlowForge` through:

- project-root `AGENTS.md`
- workflow scripts in `.flowforge/scripts/`
- canonical project artifacts under `docs/` and `.flowforge/`

Unlike Claude Code and OpenCode, Codex does not use a repository-local slash-command registry in this toolkit.

Recommended entrypoints in Codex:

- `.flowforge/scripts/flowforge-create-proposal.js`
- `.flowforge/scripts/flowforge-approve-proposal.js`
- `.flowforge/scripts/flowforge-apply-proposal.js`
- `.flowforge/scripts/flowforge-add-note.js`
- `.flowforge/scripts/flowforge-list-proposals.js`
- `.flowforge/scripts/flowforge-archive-proposal.js`
- `.flowforge/scripts/flowforge-proposal-status.js`
- `.flowforge/scripts/flowforge-check-archive.js`
