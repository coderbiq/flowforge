# EXAMPLE — Show warning counts in human-readable frontier output

> Scenario fixture, not an executable FlowForge ticket.

## Why

When a ready ticket has content warnings, the human-readable `frontier` output should show the warning count beside that ticket so a user can see execution risk without switching to JSON output.

JSON output remains unchanged. A ready ticket with no warnings keeps the current text rendering.

## Design increment

Use the existing human-readable frontier rendering seam. The warning count is derived from the diagnostics already associated with each frontier result; this change introduces no new diagnostic state or persistence.

## Change

1. Add a failing text-output case for a ready ticket with two warnings.
2. Render the count beside that ticket without changing its DAG classification.
3. Preserve the zero-warning text and existing JSON representation.

## Constraints

- Text rendering MUST NOT change whether a ticket is ready, blocked, or included by override.
- JSON fields and values MUST NOT change.
- The implementation MUST consume existing diagnostics rather than recompute content quality in the renderer.

## Done and verify

- A ready ticket with two warnings displays a warning count of two.
- A ready ticket with no warnings retains its existing text.
- Existing JSON snapshots remain unchanged.
- Run the focused frontier rendering tests, then `go test ./internal/...`.

## Completion evidence

Written after execution: commands, observed results, review disposition, deviations, and implementation reference. Raw logs do not belong here.

