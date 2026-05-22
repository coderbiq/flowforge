# AGENTS.md

This repository uses `FlowForge` as the default design and delivery workflow.

## Expected behavior

- Use durable files, not chat-only conclusions, for important findings and decisions.
- Start with exploration before implementation when the work is not already decision-complete.
- Classify every exploration and proposal before writing the design surface.
- Keep proposal metadata, design, task mapping, and notes aligned while work is in progress.
- Archive completed work into the declared target docs instead of leaving the proposal as the final source of truth.

## Canonical paths

- Workflow rules: `.flowforge/workflow/guides/`
- Schemas: `.flowforge/workflow/schema/`
- Project docs: `docs/`
- Local work-restoration state: `.flowforge/state/`
- Workspace-local template copies, when needed: `docs/<workspace>/_templates/`

## Default lifecycle

1. Explore in `docs/explorations/`
2. Propose in `docs/proposals/`
3. Approve the approach and archive targets
4. Apply tasks from `task-map.md`
5. Implement while keeping `notes.md` current
6. Archive into modules, architecture docs, conventions, and ADRs as needed

## Classification model

Every exploration and proposal must declare:

- `size_class`: `small`, `medium`, or `large`. See `.flowforge/workflow/guides/sizing.md`.
- `ownership`: one or more entries of type `module`, `system`, `cross-module`, or `convention`. See `.flowforge/workflow/guides/ownership.md`.
- Template customization is copy-and-edit only. If a workspace needs project-specific template variants, place them in `docs/<workspace>/_templates/` and make the customization explicit in the copied files.

Document layout follows the size class:

- `small`: single-file `design.md`
- `medium`: single-file `design.md`, or `design/` plus `model/` when split is justified
- `large`: `design/` plus `model/` directories, with one document per business model

## Memory model

- Local state tracks what is currently being worked on.
- External memory stores reusable experience only.
- Do not mix routine progress updates into the external memory provider.

## Task model

- Default backend: `Beads`
- `task-map.md` is the source of truth for task decomposition
- `task-map.md` must use deliverable-first splitting, milestone boundaries, and explicit checkpoint rules from `.flowforge/workflow/guides/task-splitting.md`
- Reference `model/<Model>.md` and `docs/conventions/<topic>.md` through `model_refs` and `convention_refs` when applicable
- Do not create unrelated standalone tasks outside the proposal workflow

## Archive model

- Module-scoped changes archive to `docs/modules/`
- Cross-cutting or system-level changes archive to `docs/architecture/`
- Reusable rules and consensus standards archive to `docs/conventions/`
- Stable architectural decisions should also produce an ADR in `docs/decisions/`
- The archived docs corpus is the default baseline for future exploration; prefer reading it before opening new exploratory threads.

## When in doubt

- Prefer updating exploration, proposal, or archive target docs over leaving context only in chat
- Prefer changing the canonical artifact rather than adding side notes that become stale

## Validation commands

- `.flowforge/scripts/flowforge-create-proposal.js --title ... --source-exploration ... --size-class ... --ownership ... --archive-target ...`
- `.flowforge/scripts/flowforge-approve-proposal.js <proposal-id|proposal-dir>`
- `.flowforge/scripts/flowforge-add-note.js <proposal-id|proposal-dir> <note text>`
- `.flowforge/scripts/flowforge-list-proposals.js`
- `.flowforge/scripts/flowforge-archive-proposal.js <proposal-id|proposal-dir>`
- `.flowforge/scripts/flowforge-apply-proposal.js <proposal-id|proposal-dir>`
- `.flowforge/scripts/flowforge-validate-proposal.js <proposal-id|proposal-dir>`
- `.flowforge/scripts/flowforge-proposal-status.js <proposal-id|proposal-dir>`
- `.flowforge/scripts/flowforge-check-archive.js <proposal-id|proposal-dir>`
