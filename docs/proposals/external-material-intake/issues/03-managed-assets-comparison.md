---
flowforge:
  schema: 1
  role: ticket
  id: managed-assets-comparison
  revision: 1
  consumes:
    requirements:
      external-material-intake-requirements: 1
    design:
      external-material-intake-design: 1
---

# 03: Compare embedded managed assets with a project non-destructively

**Blocked by:** None

**Status:** open

**Delivery:** One internal comparison service reports every embedded managed Skill/agent-rule target as current, missing, or drifted and reports extra project files as project-owned without modifying the project.

## Design context

Use the embedded assets extracted by the running binary as described in [managed-resource verification](../design.md#external-material-intake-design), fulfilling [the requirement authority's resource-verification outcome](../requirements.md#external-material-intake-requirements). The comparator is the common seam for read-only verification and later synchronization reporting.

## Touch points

`internal/command/assets_deploy.go`; `internal/command/assets_deploy_test.go`; new focused asset-comparison tests.

## Changes

1. Model managed asset targets for `assets/skills/` and `assets/agents/` using the configured absolute or relative docs root.
2. Compare source and target content deterministically and enumerate unknown target files as project-owned.
3. Return a stable, testable result that separates comparison from copying and from `AGENTS.md` block merging.

## Constraints

- Embedded assets are the only comparison source of truth; do not introduce an installation manifest.
- The comparison path must never overwrite, delete, rename, or classify project-owned files as drifted.

## Done and verify

- Fixtures cover current, missing, drifted, project-owned, relative docs root, and absolute docs root results without file mutation.
- `go test ./internal/command/...`.
