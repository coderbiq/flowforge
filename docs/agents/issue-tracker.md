# Issue tracker: local Markdown + FlowForge DAG

Feature artifacts live under `<docs_dir>/proposals/<feature-slug>/` (default `docs/proposals/`).

## Authority and layout

- Requirement owns observable behavior and scope; design owns responsibilities, interfaces, seams, migration, and verification strategy.
- Human semantic prose is authoritative. Machine IDs and revisions support traceability.
- Only `issues/*.md` is executable. A ticket keeps `Blocked by` and `Status` directly below its title.
- Use `requirements.md`, `design.md`, `spec.md`, `map.md`, or `evidence/*.md` only when independent consumption, review, lifecycle, or readability requires them.

## Publishing

Write and edit Markdown files directly. A ticket states one Delivery, locally sufficient Design context, stable Touch points, ordered Changes, scoped Constraints, and paired Done/verify conditions. Keep only real DAG dependencies and omit empty template sections.

Do not store readiness phases. Record scoped gaps and blockers on the authority they affect. A closed ticket includes non-empty `Completion evidence` with delivered behavior, observed verification, both review dispositions, deviations, and an implementation reference.

## Deterministic checks

- Run `flowforge check` after changing artifacts or dependencies.
- Run `flowforge frontier` to obtain the next unblocked executable tickets.
- Use `--strict` when warnings and gaps must fail policy; use `--include-gaps` only as an explicit risk acceptance. Diagnostics remain visible.
