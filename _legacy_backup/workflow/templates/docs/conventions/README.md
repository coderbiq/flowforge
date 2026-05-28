---
doc_type: note
title: Conventions
status: draft
workspace: default
module_scope: []
system_scope: []
convention_scope: []
ownership: []
information_class: convention
topics: []
related_docs: []
archive_target: none
created: <ISO-8601 timestamp>
updated: <ISO-8601 timestamp>
---

# Conventions

`docs/conventions/` holds reusable consensus rules. Each file describes one rule that applies whenever the matching situation appears in the codebase.

## Ownership summary

- Primary module: none
- System / architecture targets: none
- Convention targets: this directory
- Canonical reading path: this conventions overview

Use this directory for:

- "this class of problem is solved with this standard approach"
- "this kind of field uses this storage shape"
- "this layer must depend only on these modules"
- "this artifact must use this naming pattern"

Do not use this directory for:

- module-internal behavior (use `docs/modules/`)
- system or cross-module structural views (use `docs/architecture/`)
- one-off architectural decisions (use `docs/decisions/`)

New convention files should follow `workflow/templates/docs/conventions/convention.md`.
