---
flowforge:
  schema: 1
  role: ticket
  id: add-refine-ticket-and-implement-preflight
  revision: 1
  consumes:
    requirements:
      ticket-refinement-contract-requirements: 1
    design:
      ticket-refinement-contract-design: 1
---

# 02: Add Refine Ticket and reject incomplete lightweight work

**Blocked by:** 01

**Status:** closed

## Delivery

FlowForge authors can use `flowforge-refine-ticket` to turn one candidate ticket into a fact-backed machine work package, and Implement refuses an incomplete package before either implementation mode starts.

## Design context

Refine Ticket validates one candidate ticket only and fills its lower execution contract; it returns unresolved facts to Explore/Solution Design. Plan publishes the skeleton, while Implement consumes the diagnostic produced by ticket 01. It fulfills the [lightweight-implementation boundary](../requirements.md#ticket-refinement-contract-requirements) through the [responsibility design](../design.md#ticket-refinement-contract-design).

## Touch points

- `assets/skills/flowforge-refine-ticket/SKILL.md` — new skill
- `assets/skills/flowforge-plan/SKILL.md` — ticket skeleton publication
- `assets/skills/flowforge-implement/SKILL.md` — preflight
- `assets/skills/_shared/ARTIFACT-CONTRACT.md` — shared contract
- `internal/command/assets_deploy_test.go` — packaged-asset references

## Changes

- [x] 1. Add the Refine Ticket input, fact-checking, return and five-section writing contract.
- [x] 2. Make Plan publish only the titled Execution detail skeleton and delegate fact-backed contents to Refine Ticket.
- [x] 3. Require Implement preflight to reject `execution-contract-incomplete`, including lightweight mode and include-gaps invocation.
- [x] 4. Update the shared artifact contract and packaged-asset tests for the new skill and hand-off.
- [x] 5. Fix: Add `flowforge-refine-ticket` to `assertRequiredArtifactContractPointers` required map with `hand-offs` and `information-value` anchors.
- [x] 6. Fix: Extract a `readSkillBody` helper for the repeated read-skill-file-and-check-substring test pattern in `internal/command/assets_deploy_test.go`.

## Constraints

- Refine Ticket must not implement, review, create readiness states, or invent unresolved values.
- Implement must not treat `--include-gaps` as permission to bypass the preflight.
- Write set: `assets/skills/`, `internal/command/assets_deploy_test.go`
- standards: none found

## Done and verify

- Every packaged skill reference resolves: `go test ./internal/command/...` — all pass.
- Skill text contains all five execution-contract sections and the refusal path: `go test ./internal/command/...` — all pass.

---

## Execution detail

### Verified contracts

- `assets/skills/flowforge-plan/SKILL.md` currently owns ticket packaging, and `assets/skills/flowforge-implement/SKILL.md` selects lightweight mode from ticket structure.
- Ticket 01 supplies the canonical diagnostic; this ticket must consume it rather than reimplement section parsing in a skill.

### Execution scenarios

- Success: a Refine Ticket run with verified repository evidence fills all five sections and leaves the ticket eligible for implementation.
- Failure: a conflicting API/Schema fact returns a scoped design/research return with no code or invented execution step.

### Expected tests

- `go test ./internal/command/...` — deployment/reference coverage passes.

### Generated artifacts

- Producer: Plan/Refine Ticket skills; artifact: issue Markdown; consumer: Implement. The shared contract must state the same five-section protocol.

### Conventions

- Keep the human-priority upper ticket tiers free of implementation paths and source coordinates.

## Implementation note

- Completed Changes: 1, 2, 3, 4.
- Created `assets/skills/flowforge-refine-ticket/SKILL.md` with input selection, fact-checking, five-section writing contract, readiness verification, and boundaries. Names all five execution-contract sections and links `../_shared/ARTIFACT-CONTRACT.md#hand-offs` and `#information-value`.
- Updated `assets/skills/flowforge-plan/SKILL.md` Tier 3 to publish only the `## Execution detail` skeleton (five empty `###` headings) and delegate fact-backed content to `flowforge-refine-ticket`. Updated ticket template accordingly.
- Added execution-contract pre-flight to `assets/skills/flowforge-implement/SKILL.md`: rejects `execution-contract-incomplete` before any implementation mode, including `--include-gaps` invocations.
- Updated `assets/skills/_shared/ARTIFACT-CONTRACT.md` Tier 3 to document the five required execution-contract sections and the Plan→Refine→Implement flow.
- Added four packaged-asset tests in `internal/command/assets_deploy_test.go`: refine-ticket packaging and links, Plan skeleton delegation, Implement preflight rejection, and ARTIFACT-CONTRACT section documentation.
- Synced `internal/command/assets/skills/flowforge-refine-ticket/` from `assets/skills/` (build-time copy maintained by Makefile/scripts/build.sh).
- Commands: `go test ./internal/...` — 111 passed in 6 packages.
- Files modified: `assets/skills/flowforge-refine-ticket/SKILL.md` (new), `assets/skills/flowforge-plan/SKILL.md`, `assets/skills/flowforge-implement/SKILL.md`, `assets/skills/_shared/ARTIFACT-CONTRACT.md`, `internal/command/assets_deploy_test.go`, `internal/command/assets/skills/flowforge-refine-ticket/SKILL.md` (build copy).
- Write set compliance: All source modifications within write set (`assets/skills/`, `internal/command/assets_deploy_test.go`). `internal/command/assets/` is a build-time copy maintained by Makefile.

## Review rounds

### Round 1

- Fixed point: current working tree
- Standards: [Low] `assertRequiredArtifactContractPointers` map omits `flowforge-refine-ticket`; [Low] duplicated read-skill-and-check-substring pattern in three new tests.
- Spec: none
- Fix changes: 5, 6
- Design returns: none

### Round 2

- Fixed point: current working tree
- Standards: none
- Spec: none
- Fix changes: none
- Design returns: none

## Completion evidence

- Delivered: `flowforge-refine-ticket` skill with input selection, fact-checking, five-section writing contract, readiness verification, and boundaries. Plan publishes only the `## Execution detail` skeleton (five empty `###` headings) and delegates fact-backed content to Refine Ticket. Implement preflight rejects `execution-contract-incomplete` before any implementation mode, including `--include-gaps` invocations. ARTIFACT-CONTRACT.md documents the five required sections and the Plan→Refine→Implement flow.
- Commands: `go test ./internal/...` — 111 passed in 6 packages; `flowforge check --dir docs/proposals/ticket-refinement-contract` — dependency graph healthy.
- Review axes: Round 1 — Standards [Low] missing required-pointer map entry (fixed), [Low] duplicated test pattern (fixed); Spec none. Round 2 — Standards none; Spec none.
- Deviations: none.
- Implementation reference: working-tree diff — `assets/skills/flowforge-refine-ticket/SKILL.md` (new), `assets/skills/flowforge-plan/SKILL.md` (+18 -6), `assets/skills/flowforge-implement/SKILL.md` (+2), `assets/skills/_shared/ARTIFACT-CONTRACT.md` (+9 -3), `internal/command/assets_deploy_test.go` (+75), `internal/command/assets/skills/flowforge-refine-ticket/SKILL.md` (build copy).
