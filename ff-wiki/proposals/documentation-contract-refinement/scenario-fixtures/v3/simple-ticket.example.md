---
prototype: true
artifact_role: ticket
flowforge:
  id: frontier-text-warning-count
  revision: 1
blocked_by: []
---

# EXAMPLE — Show warning counts in human-readable frontier output

> Scenario fixture, not an executable FlowForge ticket or an approved schema.

## Why

Human-readable `frontier` output shows the warning count beside each otherwise-ready ticket, so users can see execution risk without switching to JSON. JSON output remains stable, and tickets with no warnings retain their current text.

## Design increment

Reuse the existing frontier text-rendering seam. Render diagnostics already attached to each frontier result; do not create diagnostic state or recompute content quality in the renderer.

## Touch points

- Human-readable frontier renderer and its focused output tests.
- The existing frontier result diagnostics consumed by that renderer.

## Changes

1. Add a failing text-output case for one ready ticket with two warnings.
2. Render the count without changing the ticket's DAG classification.
3. Preserve zero-warning text and existing JSON behavior.

## Constraints

- Rendering MUST NOT change whether a ticket is ready, blocked, or included by override.
- JSON fields and values MUST NOT change.
- The renderer MUST consume existing diagnostics.

## Done and verify

- The focused text-output test proves a ready ticket with two warnings displays `2`.
- The zero-warning text case remains unchanged.
- Existing JSON tests remain unchanged.
- Run the focused frontier rendering test, then `go test ./internal/...`.

## Completion evidence

Pending implementation.

