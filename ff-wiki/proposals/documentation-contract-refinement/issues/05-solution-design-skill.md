# 05: Add the solution-design owner

**Blocked by:** 02, 04

**Status:** closed

**What to build:** `flowforge-solution-design` turns approved requirements into authoritative module, responsibility, seam, flow, migration, and verification decisions, returning scoped open items when affected planning cannot proceed.

See [solution-design invocation and process](../solution-design-interface-v0.2.md) and [production Skill shape](../skill-coordination-design-v0.1.md#flowforge-solution-design-new).

## Touch points

New Skill package, shared contract reference, router-visible description, scenario fixtures.

## Changes

1. Implement the approved invocation triggers, required inputs, decision-frontier loop, and caller-visible completion report.
2. Orchestrate research, prototype, domain modeling, and codebase-design only when their branch condition fires.
3. Maintain adaptive design authority and area revisions while recording scoped `open_items` for unresolved facts.
4. Apply information-value compression before handing resolved areas to planning.

## Constraints

- The Skill must not persist `design-ready`, publish tickets, implement production code, or silently resolve requirement ambiguity.
- Human design prose is authority; metadata supports deterministic consumption.

## Done and verify

- Scenario tests cover simple-feature rejection, cross-module design, requirement return, research/prototype branches, partial planning, revision scope, and compression.
- Every affected area receives a resolved, gap, warning, or blocker result with an authority link.

## Completion evidence

- Added and validated the production `flowforge-solution-design` Skill with discriminating triggers, decision-frontier loop, Research/Prototype/Domain/Codebase-design/Grilling branches, incremental authority updates, scoped open items, partial planning, adaptive packaging, implementation-return classification, and information-value compression.
- Shared schema v1 is progressively disclosed only when metadata is authored; deployment tests verify every required contract anchor and relative reference.
- One independent complex forward test drove six instruction corrections. A second independent ten-input evaluation passed 10/10 routes and all forbidden-behavior audits; its prompts, results, and reproduction protocol are recorded in `solution-design-production-validation-v0.1.md`.
- `quick_validate.py`, `make dev`, `go test ./internal/...`, `git diff --check`, and repeated Standards/Spec review passed.
