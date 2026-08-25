# 03: Project Catalog facts through check and frontier

**Blocked by:** 02

**Status:** open

**What to build:** `check` and `frontier` combine Catalog diagnostics with DAG facts so humans and automation can distinguish ready work, warnings, gaps, claimed work, and blockers without a persisted readiness state.

See [check projection](../artifact-catalog-interface-v0.1.md#check-projection) and [frontier projection](../artifact-catalog-interface-v0.1.md#frontier-projection).

## Touch points

Check/frontier command policies, text/quiet/JSON rendering, exit behavior, existing DAG computation.

## Changes

1. Feed both commands from Artifact Catalog while keeping DAG calculation focused on tickets.
2. Render stable diagnostic fields and all frontier groups consistently across text and JSON.
3. Make strict policy affect filtering and exit status without mutating diagnostic severity.
4. Provide an explicit ephemeral gap override that keeps explanations visible and never includes true blockers.

## Constraints

- Existing command defaults and automation-oriented output require compatibility tests before any public flag/signature change.
- Quiet output must keep executable paths on stdout and diagnostics on stderr.
- Commands must not write readiness or waiver state.

## Done and verify

- Command tests cover default, strict, override, quiet, and JSON behavior plus unchanged legacy DAG results.
- `go test ./internal/command/... ./internal/tracker/...`
- `go run ./cmd/flowforge check --dir docs/proposals`
- `go run ./cmd/flowforge frontier --dir docs/proposals`
