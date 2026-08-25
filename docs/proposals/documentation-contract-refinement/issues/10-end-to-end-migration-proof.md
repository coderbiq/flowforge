# 10: Prove the refined workflow end to end

**Blocked by:** 03, 08, 09

**Status:** open

**What to build:** Maintainers can run documented simple and Tangram-style complex scenarios showing adaptive artifacts, deterministic execution order, scoped diagnostics/overrides, legacy compatibility, and evidence-backed completion.

See [verification scenarios](../skill-coordination-design-v0.1.md#verification-scenarios) and [artifact scenario validation](../artifact-scenario-validation-v0.3.md).

## Touch points

End-to-end fixtures, Tangram compatibility sample, migration guidance, CLI and Skill integration tests.

## Changes

1. Run the simple feature through compact routing, execution, review, and inline evidence.
2. Run the complex feature through requirement authority, solution design, tracer tickets, frontier execution, and promoted evidence.
3. Demonstrate warnings, a scoped gap override, a non-overridable blocker, an exact persistent waiver, and stale-waiver detection.
4. Document incremental adoption for existing v5 proposals without repository-wide migration.

## Constraints

- Representative artifacts, not only parser fixtures, must prove the workflow.
- Migration guidance must preserve local Markdown and existing v5 executable tickets.
- No scenario may depend on a manually advanced readiness state.

## Done and verify

- Both scenarios satisfy every coordination verification case and retain human-readable authority links.
- Tangram still projects 11 executable legacy tickets and excludes `spec.md`.
- `go test ./internal/...`
- `git diff --check`
