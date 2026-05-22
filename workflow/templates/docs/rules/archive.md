---
doc_type: note
title: Workflow Rules - Archive
status: draft
workspace: default
module_scope: []
system_scope: []
convention_scope: []
ownership: []
information_class: note
topics: []
related_docs: []
archive_target: none
created: <ISO-8601 timestamp>
updated: <ISO-8601 timestamp>
---

# Archive Rules

## Ownership summary

- Primary module: none
- System / architecture targets: none
- Convention targets: none
- Canonical reading path: this archive rule file

## Archive defaults

- Archive behavior is driven by proposal metadata, not by platform adapters.
- Confirm the proposal is implemented and the task backend has no open tasks.
- Update the primary archive target and any secondary archive targets.
- Promote validated `reusable_rules` into `docs/conventions/` when the proposal
  validated them.
- Record superseded decisions when applicable.
- Set the proposal status to `archived` only after the final corpus is updated.

## Maintenance expectations

- Keep the overview and linked subdocs in sync so readers can still navigate
  the full system from the canonical entry point.
- Preserve the older fact in history or changelog sections when a final doc is
  narrowed or replaced.
- Do not archive only the proposal directory and skip the target docs.

