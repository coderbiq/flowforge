# Incremental Adoption Guide v0.1

Existing v5 proposals remain executable without repository-wide migration. Adopt the refined contract at the next meaningful edit, role by role.

## Existing proposal

Keep current `issues/*.md`, human `Blocked by`, and existing statuses. Catalog recognizes them as legacy tickets, emits `legacy-metadata` warnings, and preserves their DAG. Do not add empty requirement/design/spec/evidence files merely to match the new vocabulary.

## Next feature or semantic change

- Keep a small existing-seam change in one compact ticket.
- Promote requirements or design only when they gain independent consumers, review, lifecycle, or would obscure execution.
- Add schema v1 metadata when creating or independently revising an authority. Record semantic revisions in downstream `consumes` while retaining meaningful Markdown links.
- Replace superseded authority meaning; do not migrate chronological discussion or template filler.
- Add Completion evidence after actual verification and review, inline unless independent proof earns promotion.

## Diagnostics and policy

Default `check` and `frontier` preserve legacy compatibility and expose diagnostics. Teams may adopt `--strict` in CI after reviewing existing warnings; strictness is policy, not stored ticket state. `--include-gaps` is an explicit, visible exception for scoped gaps. It never hides the diagnostic, cannot override blockers, and does not mutate artifacts.

Persistent waivers match one exact diagnostic and target with a reason. Recheck them after authority changes: a waiver whose diagnostic disappears becomes stale and should be removed.

## Safe migration sequence

1. Deploy the shared Skills and references; run packaged-pointer tests.
2. Run default `check/frontier` and inventory warnings without blocking current work.
3. Upgrade authority metadata only where a new consumer needs revision diagnostics.
4. Upgrade a ticket when it is next planned or implemented, preserving its real DAG edges.
5. Introduce strict CI only for the chosen scope after warnings have dispositions.

At every step, Markdown files remain the content interface and the CLI remains graph/catalog projection. No readiness phase or bulk status transition is introduced.
