# AGENTS.md

This repository uses `FlowForge` as the default design and delivery workflow.

## Expected behavior

- Use durable files, not chat-only conclusions, for important findings and decisions.
- Start with exploration before implementation when the work is not already decision-complete.
- Keep proposal metadata, design, task mapping, and notes aligned while work is in progress.
- Archive completed work into the declared target docs instead of leaving the proposal as the final source of truth.

## Canonical paths

- Workflow rules: `.flowforge/workflow/guides/`
- Schemas: `.flowforge/workflow/schema/`
- Project docs: `docs/`
- Local work-restoration state: `.flowforge/state/`

## Default lifecycle

1. Explore in `docs/explorations/`
2. Propose in `docs/proposals/`
3. Approve the approach and archive targets
4. Apply tasks from `task-map.md`
5. Implement while keeping `notes.md` current
6. Archive into modules, architecture docs, and ADRs as needed

## Memory model

- Local state tracks what is currently being worked on.
- External memory stores reusable experience only.
- Do not mix routine progress updates into the external memory provider.

## Task model

- Default backend: `Beads`
- `task-map.md` is the source of truth for task decomposition
- `task-map.md` must use deliverable-first splitting, milestone boundaries, and explicit checkpoint rules from `workflow/guides/task-splitting.md`
- Do not create unrelated standalone tasks outside the proposal workflow

## Archive model

- Module-scoped changes archive to `docs/modules/`
- Cross-cutting changes archive to `docs/architecture/`
- Stable architectural decisions should also produce an ADR in `docs/decisions/`

## When in doubt

- Prefer updating exploration, proposal, or archive target docs over leaving context only in chat
- Prefer changing the canonical artifact rather than adding side notes that become stale

## Validation commands

- `.flowforge/scripts/flowforge-create-proposal.js --title ... --source-exploration ... --archive-target ...`
- `.flowforge/scripts/flowforge-approve-proposal.js <proposal-id|proposal-dir>`
- `.flowforge/scripts/flowforge-add-note.js <proposal-id|proposal-dir> <note text>`
- `.flowforge/scripts/flowforge-list-proposals.js`
- `.flowforge/scripts/flowforge-archive-proposal.js <proposal-id|proposal-dir>`
- `.flowforge/scripts/flowforge-apply-proposal.js <proposal-id|proposal-dir>`
- `.flowforge/scripts/flowforge-validate-proposal.js <proposal-id|proposal-dir>`
- `.flowforge/scripts/flowforge-proposal-status.js <proposal-id|proposal-dir>`
- `.flowforge/scripts/flowforge-check-archive.js <proposal-id|proposal-dir>`
