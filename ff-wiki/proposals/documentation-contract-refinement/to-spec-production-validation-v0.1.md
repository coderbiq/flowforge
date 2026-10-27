# To-Spec Production Validation v0.1

**Date:** 2026-08-25

## Generated overview

The exact [optional navigation spec](scenario-fixtures/v4/plan-simple/spec.example) consumes the local requirement and design revisions, links rather than copies them, and points to the compact execution/verification contract.

It contains no extensive user-story inventory, invented implementation decision, copied authority section, or executable status.

## Reproduction

1. Copy `scenario-fixtures/v4/plan-simple` to an isolated proposal root.
2. Rename `issues-example` to `issues` and copy each `.example` ticket to the same basename with a `.md` suffix.
3. Copy `spec.example` to `spec.md`.
4. Resolve every Markdown link in `spec.md`, then run `flowforge check` and `flowforge frontier`; record one executable ticket and no spec issue.
5. Delete only `spec.md`, then repeat both commands.

All overview links resolved, including the production execution path `issues/01-warning-count.md`. The recorded before/after result is identical: `Checked 1`, a healthy graph, and ready ticket `#01`. Requirement/design authorities and the ticket's direct semantic consumption remain intact after deletion.

## Gap behavior

When a required decision is absent rather than represented by a scoped authority open item, production To-Spec returns the unresolved owner and writes no overview. A recorded scoped gap is linked with its affected area; synthesis never fills it.
