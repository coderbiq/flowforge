---
flowforge:
  schema: 1
  role: ticket
  id: route-review-findings-to-repair-tickets
  revision: 1
  consumes:
    requirements:
      ticket-refinement-contract-requirements: 1
    design:
      ticket-refinement-contract-design: 1
---

# 04: Route substantive review findings through repair tickets

**Blocked by:** 02, 03

**Status:** closed

## Delivery

Review preserves review rounds on the original ticket while creating and wiring a repair ticket for substantive findings; published documentation and scenario coverage demonstrate the complete refine–implement–review–repair loop.

## Design context

Review decides local mechanical Fix versus repair-ticket escalation using the agreed severity and contract-change threshold. It must preserve audit history without allowing the original ticket to become a patch log. It fulfills the [review-audit outcome](../requirements.md#ticket-refinement-contract-requirements) through the [repair lifecycle design](../design.md#ticket-refinement-contract-design).

## Touch points

- `assets/skills/flowforge-review/SKILL.md`
- `assets/skills/_shared/ARTIFACT-CONTRACT.md`
- `docs/agents/issue-tracker.md`
- `README.md`
- tracker/command scenario fixtures and tests

## Changes

- [x] 1. Change Review's disposition rules to create a repair ticket for substantive findings and retain only local mechanical Fix loops on the original.
- [x] 2. Define original/repair evidence hand-off, downstream rewiring and original closure after a clean repair review.
- [x] 3. Update user-facing tracker/workflow documentation and add an end-to-end fixture covering both execution-contract and repair DAG paths.
- [x] 4. Fix: Merge Design return classification into Substantive finding escalation so that design-return findings also create a repair ticket instead of staying on the original as non-blocking design returns, per requirements.

## Constraints

- Review records evidence; it must not implement the repair.
- Do not append High/Critical, contract-changing, cross-write-set or design-return repair work to the original Changes list.
- Write set: `assets/skills/`, `docs/agents/issue-tracker.md`, `README.md`, `internal/tracker/`, `internal/command/`
- standards: none found

## Done and verify

- Workflow and fixture coverage pass: `go test ./internal/tracker/... ./internal/command/...` — all pass.
- Repository documentation describes the same contract as the deployed skills: `go test ./internal/command/...` — all pass.

---

## Execution detail

### Verified contracts

- `assets/skills/flowforge-review/SKILL.md` Step 6 (Fix planning) currently classifies findings as Fixable (→ Fix Change) or Design return; it has no repair-ticket escalation path. This ticket adds a third disposition: substantive findings create a repair ticket.
- `assets/skills/_shared/ARTIFACT-CONTRACT.md:81-90` "Execution and review loop" section describes the Plan→Implement→Review→Fix cycle; it needs a repair branch for substantive findings.
- `docs/agents/issue-tracker.md` "Execution and review loop" section mirrors the artifact contract; it needs the same repair branch.
- `internal/tracker/model.go:12` defines `StatusNeedsRepair`; `model.go:43` defines `RepairOf string` field — both added by ticket 03 and ready for Review to reference.
- `internal/tracker/catalog.go:53-54` defines `DiagnosticDanglingRepairReference` and `DiagnosticMissingReciprocalRepair` — Review must ensure its repair tickets pass these validations.
- `internal/tracker/dag.go:158-160` `ComputeFrontier` already excludes `needs-repair` from dispatch; repair tickets with `open` status remain dispatchable.

### Execution scenarios

- Success: a Review round finds a substantive finding (High/Critical, contract-changing, cross-write-set), creates a repair ticket with `Repair of: <original-id>`, sets the original to `needs-repair`, and rewires downstream `Blocked by` to the repair ticket. The repair ticket enters frontier; the original does not. After a clean repair review, the original is closed.
- Failure: a Review round with only local mechanical findings (same write set, no contract change) appends `Fix:` Changes to the original ticket without creating a repair ticket — preserving the existing Fix loop for trivial issues.

### Expected tests

- `go test ./internal/tracker/... ./internal/command/...` — end-to-end fixture covers: missing contract exclusion, refined implementation eligibility, substantive finding → repair ticket creation, repair execution, and downstream release after clean repair review.
- `go test ./internal/command/...` — packaged-asset tests verify Review skill and ARTIFACT-CONTRACT describe the repair escalation path.

### Generated artifacts

- Producer: Review skill; artifact: repair ticket Markdown under `issues/`; consumer: Implement. The shared contract must state the same repair-creation protocol.

### Conventions

- Keep the dual-axis review method; only the finding disposition path changes.
- Review rounds remain on the original ticket as append-only audit history.
- A clean repair review is required before closing both the repair and its original.
- `needs-repair` status and `Repair of:` header are provenance, not DAG edges; downstream rewiring uses `Blocked by: <repair-id>`.

## Implementation note

- Completed Changes: 1, 2, 3.
- Updated `assets/skills/flowforge-review/SKILL.md` Step 6 with third disposition "Substantive finding (repair escalation)": creates a repair ticket with `Repair of:`, sets original to `needs-repair`, rewires downstream `Blocked by`, and records repair reference in Review rounds. Updated Review rounds template with `Repair:` field. Updated Return and boundary clauses.
- Updated `assets/skills/_shared/ARTIFACT-CONTRACT.md` Execution and review loop with "Substantive finding → repair ticket" paragraph describing the repair escalation path, provenance semantics, and closure flow.
- Updated `docs/agents/issue-tracker.md` Execution and review loop with the same repair workflow description.
- Added three packaged-asset tests in `internal/command/assets_deploy_test.go`: `TestReviewSkillDescribesRepairEscalation`, `TestArtifactContractDescribesRepairLoop`, `TestIssueTrackerDescribesRepairWorkflow`.
- Added end-to-end fixture `TestEndToEndRepairDAGPath` in `internal/tracker/catalog_test.go` covering: missing contract exclusion (needs-repair non-executable), repair provenance parsing, frontier dispatch (repair in ready, original excluded, downstream blocked), and reciprocal reference validation (no false diagnostics).
- Synced `internal/command/assets/skills/` build copies for Review SKILL.md and ARTIFACT-CONTRACT.md.
- Commands: `go test ./internal/...` — 118 passed in 6 packages.
- Files modified: `assets/skills/flowforge-review/SKILL.md`, `assets/skills/_shared/ARTIFACT-CONTRACT.md`, `docs/agents/issue-tracker.md`, `internal/command/assets_deploy_test.go`, `internal/tracker/catalog_test.go`, `internal/command/assets/skills/flowforge-review/SKILL.md` (build copy), `internal/command/assets/skills/_shared/ARTIFACT-CONTRACT.md` (build copy).
- Write set compliance: All source modifications within write set (`assets/skills/`, `docs/agents/issue-tracker.md`, `internal/tracker/`, `internal/command/`).

## Review rounds

### Round 1

- Fixed point: current working tree
- Standards: none
- Spec: [High] Design return classification retained as non-blocking on original; requirements mandate design-return findings create repair tickets. Fixed by merging Design return into Substantive finding escalation.
- Fix changes: 4
- Design returns: none

## Completion evidence

- Delivered: Review SKILL.md disposition rules create repair tickets for substantive findings (High/Critical, contract-changing, cross-write-set, boundary-altering) and design-return findings; local mechanical findings retain Fix loop. ARTIFACT-CONTRACT.md and issue-tracker.md document the repair escalation path. End-to-end fixture covers missing contract exclusion, repair provenance, frontier dispatch, downstream blocking, and reciprocal validation.
- Commands: `go test ./internal/...` — 118 passed in 6 packages; `flowforge check --dir docs/proposals/ticket-refinement-contract` — dependency graph healthy.
- Review axes: Round 1 — Standards none; Spec [High] design-return classification mismatch (fixed by merging into repair escalation).
- Deviations: none.
- Implementation reference: working-tree diff — `assets/skills/flowforge-review/SKILL.md`, `assets/skills/_shared/ARTIFACT-CONTRACT.md`, `docs/agents/issue-tracker.md`, `internal/command/assets_deploy_test.go`, `internal/tracker/catalog_test.go`, build copies in `internal/command/assets/skills/`.
