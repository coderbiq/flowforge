# Codex Adapter

Codex uses `tg-workflow` through:

- project-root `AGENTS.md`
- workflow scripts in `scripts/`
- canonical project artifacts under `docs/` and `.workflow/`

Unlike Claude Code and OpenCode, Codex does not use a repository-local slash-command registry in this toolkit.

Recommended entrypoints in Codex:

- `scripts/tg-create-proposal.js`
- `scripts/tg-approve-proposal.js`
- `scripts/tg-apply-proposal.js`
- `scripts/tg-add-note.js`
- `scripts/tg-list-proposals.js`
- `scripts/tg-archive-proposal.js`
- `scripts/tg-proposal-status.js`
- `scripts/tg-check-archive.js`
