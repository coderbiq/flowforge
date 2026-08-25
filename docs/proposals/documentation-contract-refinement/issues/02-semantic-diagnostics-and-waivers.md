# 02: Explain semantic dependencies, open items, and waivers

**Blocked by:** 01

**Status:** open

**What to build:** The Catalog deterministically explains stale or missing authorities, scoped gaps/blockers, and exact reasoned waivers from current Markdown files without inferring prose quality or persisting readiness.

See [diagnostic scope projection](../artifact-catalog-interface-v0.1.md#diagnostic-scope-projection) and [physical open-item and waiver metadata](../physical-markdown-schema-v0.1.md#optional-common-fields).

## Touch points

Catalog identity index, authority areas and revisions, Markdown link resolution, diagnostic projection, schema fixtures.

## Changes

1. Resolve feature-local and explicit cross-feature semantic identities, anchors, and consumed revisions.
2. Compare machine dependencies with marked human semantic links and project findings only to affected consumers.
3. Parse scoped `open_items` and emit their declared gap or blocker facts.
4. Match persistent waivers by exact diagnostic and target while retaining the original diagnostic and reason.

## Constraints

- Diagnostic severity describes a fact and must not change under strict policy.
- Invalid scope must produce a structure diagnostic, never an implicit feature-wide gate.
- Wildcard or blanket waivers must not suppress diagnostics.

## Done and verify

- Tests cover duplicate/missing/future identities, stale revisions, link mismatches, scoped open items, and valid/invalid/stale waivers.
- `go test ./internal/tracker/...`
