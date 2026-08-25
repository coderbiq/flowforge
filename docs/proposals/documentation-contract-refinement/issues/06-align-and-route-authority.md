# 06: Make Align and Route preserve requirement authority

**Blocked by:** 04

**Status:** closed

**What to build:** Incoming work is routed by semantic need, and alignment persists concise requirement authority before handing cross-module design questions to solution design.

See [Align responsibility](../skill-coordination-design-v0.1.md#flowforge-align) and [router context-load budget](../skill-coordination-design-v0.1.md#router-context-load-budget).

## Touch points

Align, Route, and domain-modeling pointers; requirement artifact examples; compact-ticket branch.

## Changes

1. Make Align own problem, observable outcomes, scope, scenarios, constraints, terms, and requirement-changing unknowns.
2. Persist accepted requirement changes incrementally instead of waiting for a final synthesis step.
3. Route local work at an existing seam to a compact ticket and responsibility/interface/seam decisions to solution design.
4. Keep router guidance to one trigger sentence per main path and disclose detailed contracts behind pointers.

## Constraints

- Alignment must not invent solution modules or ticket decomposition.
- Route must not inline the artifact schema or create a new workflow state machine.

## Done and verify

- Simple and complex interaction fixtures reach the intended Skill without empty ceremonial artifacts.
- Requirement-changing unknowns return to Align; design-only open items stay with solution design.
- Packaged Skill reference tests pass.

## Completion evidence

- Production Align owns requirement authority, persists accepted facts incrementally, and hands module/interface/seam decisions to Solution Design.
- Production Route selects one semantic owner, supports compact existing-seam work, and creates no readiness or workflow-state artifact.
- Four exact-prompt routing cases and forbidden-behavior checks are recorded in [Align and Route Production Validation v0.1](../align-route-production-validation-v0.1.md).
- Standards and Spec re-reviews both passed after correcting pointer triggers, balanced deployment markers, and evidence scope.
- `GOPROXY=https://goproxy.cn,direct go test ./internal/...` and `git diff --check` pass.
