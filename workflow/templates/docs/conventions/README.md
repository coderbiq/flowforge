# Conventions

`docs/conventions/` holds reusable consensus rules. Each file describes one rule that applies whenever the matching situation appears in the codebase.

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
