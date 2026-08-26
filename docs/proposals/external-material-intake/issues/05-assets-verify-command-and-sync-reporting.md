---
flowforge:
  schema: 1
  role: ticket
  id: assets-verify-command-and-sync-reporting
  revision: 1
  consumes:
    requirements:
      external-material-intake-requirements: 1
    design:
      external-material-intake-design: 1
---

# 05: Expose managed-assets verification and truthful sync reporting

**Blocked by:** 03

**Status:** open

**Delivery:** Users can run `flowforge assets verify [project]` (including `--json`) to inspect managed-resource drift read-only, while `init` and `upgrade` report synchronization success only after the same comparison says every managed target is current.

## Design context

Project-facing projection of the [managed-resource verification design](../design.md#external-material-intake-design), fulfilling [the requirement authority's resource-verification outcome](../requirements.md#external-material-intake-requirements) and built on ticket 03's comparator. `init`/`upgrade` retain their explicit copying behavior; verify never copies.

## Touch points

`internal/command/root.go`; new assets command; `internal/command/init.go`; `internal/command/upgrade.go`; command tests; README and CLI documentation.

## Changes

1. Add the nested `assets verify` command with human and JSON projections and a non-zero result for missing/drifted managed assets.
2. After explicit deployment, invoke the shared comparison result and replace unconditional synchronized messages with current-only success or precise divergence output.
3. Document the command, result categories, and non-destructive semantics.

## Constraints

- Do not use verify to synchronize resources or automatically overwrite drifted content.
- `project-owned` is informational and does not make verification fail.
- Preserve configured docs-root behavior and legacy project configuration.

## Done and verify

- Command-level fixtures prove current, missing, drifted, project-owned, JSON, and non-zero policy behavior; init/upgrade success messages occur only for an all-current post-copy result.
- `GOPROXY=https://goproxy.cn,direct go test ./internal/command/... ./internal/update/...` and `git diff --check`.
