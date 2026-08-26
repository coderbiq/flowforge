---
flowforge:
  schema: 1
  role: ticket
  id: strict-authority-publication-check
  revision: 1
  consumes:
    requirements:
      external-material-intake-requirements: 1
    design:
      external-material-intake-design: 1
---

# 02: Make strict check a complete authority publication check

**Blocked by:** None

**Status:** open

**Delivery:** `flowforge check --dir <feature-dir> --strict` rejects every unwaived metadata, anchor, open-item, consumed-revision, and semantic-link diagnostic created while publishing authority, even when the feature has no tickets.

## Design context

The [authority publication check](../design.md#external-material-intake-design) implements [the requirement authority's publication outcome](../requirements.md#external-material-intake-requirements) through the current Catalog projection; it must remain diagnostic-only and must not write workflow state or require Plan to have created an issue.

## Touch points

`internal/tracker/catalog.go`; `internal/tracker/catalog_semantics.go`; `internal/tracker/catalog_test.go`; `internal/command/check.go`; `internal/command/frontier_test.go`.

## Changes

1. Establish fixtures for requirement/design-only features, including valid whole-document authority and invalid area/open-item/consumption/link cases.
2. Preserve the Catalog's semantic diagnostics for all artifact roles and make the command projection clearly report authority validation when zero executable tickets exist.
3. Verify strict/non-strict policy behavior and JSON output without changing frontier eligibility or ticket status.

## Constraints

- Whole-document revision stays the default; require an explicit anchor only for declared areas and open items.
- A strict failure is a publication result, never a persisted readiness phase or global cross-feature block.

## Done and verify

- A valid feature containing only requirement/design returns success in strict mode; each malformed authority relationship returns the expected diagnostic and strict non-zero result.
- `GOPROXY=https://goproxy.cn,direct go test ./internal/tracker/... ./internal/command/...`.
