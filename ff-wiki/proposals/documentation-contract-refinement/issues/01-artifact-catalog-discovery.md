# 01: Discover proposal artifacts through one Catalog

**Blocked by:** None

**Status:** closed

**What to build:** Proposal discovery returns one Artifact Catalog whose executable projection contains only compatible tickets under `issues/`; legacy repositories continue to work and non-ticket Markdown never enters the DAG.

See [Artifact Catalog build pipeline](../artifact-catalog-interface-v0.1.md#build-pipeline) and [executable discovery safety](../physical-markdown-schema-v0.1.md#executable-discovery).

## Touch points

Tracker discovery/parsing, artifact and ticket domain models, compatibility adapter, filesystem fixtures.

## Changes

1. Introduce the Catalog, artifact-role, source-location, and executable-ticket projection types behind one discovery entry point.
2. Split YAML frontmatter from the Markdown body, then preserve legacy human-field parsing for ticket state and blockers.
3. Apply the role/location safety matrix and keep `DiscoverIssues` as a temporary compatibility projection.
4. Cover current FlowForge and Tangram proposal layouts with filesystem-level tests.

## Constraints

- Only `issues/*.md` may become executable.
- Missing or invalid new metadata must degrade visibly and must not break legacy tickets by default.
- The CLI must not rewrite proposal files.

## Done and verify

- Catalog tests prove every role/location and legacy fallback case, including `spec.md` exclusion.
- `go test ./internal/tracker/...`
- The filesystem fixture reports 11 Tangram tickets and no executable `spec.md`.

## Completion evidence

- Delivered `DiscoverArtifacts` with explicit artifact roles, diagnostics, source locations, and stable executable-ticket projection; `DiscoverIssues` now acts as the legacy adapter.
- Frontmatter and human ticket fields are parsed from one file snapshot. Invalid delimiters, schema values, roles, and role/location combinations degrade safely with diagnostics.
- `go test ./internal/tracker/...` passed, including the complete role/location matrix and the real FlowForge/Tangram proposal regressions.
- `GOPROXY=https://goproxy.cn,direct go test ./internal/...` and `git diff --check` passed.
- Dual-axis review found two Standards risks and four Spec gaps on the first pass; all were corrected. Standards re-review passed, and the final Spec matrix gap was covered by `TestArtifactRoleLocationMatrix` plus an explicit Tangram `spec.md` exclusion assertion.
