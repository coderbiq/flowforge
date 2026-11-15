---
flowforge:
  schema: 1
  role: ticket
  id: add-repair-ticket-lifecycle
  revision: 1
  consumes:
    requirements:
      ticket-refinement-contract-requirements: 1
    design:
      ticket-refinement-contract-design: 1
---

# 03: Represent repair tickets without DAG deadlock

**Blocked by:** 01

**Status:** closed

## Delivery

An original ticket with a substantive review finding becomes non-executable `needs-repair`; a separate repair ticket remains runnable through `Repair of:`, and affected downstream tickets wait for the repair instead of stale delivery.

## Design context

`Repair of:` is provenance, not a blocker. The original cannot block the repair because it remains nonterminal until repair review passes. It delivers the [repair-DAG outcome](../requirements.md#ticket-refinement-contract-requirements) using the [repair lifecycle design](../design.md#ticket-refinement-contract-design).

## Touch points

- `internal/tracker/model.go` — status lifecycle
- `internal/tracker/parser.go` — ticket header parsing
- `internal/tracker/dag.go` — dispatch projection
- `internal/tracker/catalog.go` and tracker tests — reciprocal repair validation

## Changes

- [x] 1. Add nonterminal, non-executable `needs-repair` and make frontier dispatch only executable statuses.
- [x] 2. Parse and expose `Repair of:` provenance while retaining backward-compatible ticket headers.
- [x] 3. Validate original/repair reciprocal references and write fixtures for repair execution, original closure eligibility and downstream blocking.
- [x] 4. Fix: Replace three inline feature-qualified key calculations in `validateRepairReferences` with calls to `issueKey` in `internal/tracker/catalog.go`.
- [x] 5. Fix: Pre-compile the five execution-contract section heading regexes as package-level vars in `internal/tracker/catalog.go` to eliminate per-call compilation in `hasCompleteExecutionContract`.

## Constraints

- `Repair of:` must not create a DAG edge.
- Do not alter terminal semantics for `closed`, `resolved` or `wontfix`.
- Write set: `internal/tracker/`, `internal/command/`
- standards: none found

## Done and verify

- Status and parser tests cover legacy tickets and repair provenance: `go test ./internal/tracker/...` — all pass.
- Frontier tests prove needs-repair is not dispatched and repair unblocks downstream only on completion: `go test ./internal/command/...` — all pass.

---

## Execution detail

### Verified contracts

- `internal/tracker/model.go:6-18` defines `Status` as a `string` type with constants including `StatusOpen` and `StatusReadyForAgent`; `IsExecutable()` at `model.go:26-28` returns `true` only for those two. `IsTerminal()` at `model.go:21-23` returns `true` for `StatusResolved`, `StatusClosed`, `StatusWontfix`.
- `internal/tracker/parser.go:14-22` defines header regexes: `statusRegex` parses `**Status:** <value>`, `blockedByRegex` parses `**Blocked by:** <ids>`. `parseIssueData` at `parser.go:33-146` populates `Issue.Status` and `Issue.BlockedBy` from these headers but has no `Repair of:` parsing.
- `internal/tracker/dag.go:141-196` `ComputeFrontier` skips only `IsTerminal()` and `StatusClaimed`; all other nonterminal statuses (including `needs-triage`, `needs-info`, `ready-for-human`) are dispatched as `Ready` when unblocked — this is the behavior deviation to fix.
- `internal/tracker/model.go:31-44` `Issue` struct has no `RepairOf` field; provenance requires a new field or parser-level extraction.
- `internal/tracker/catalog.go:220-238` `discoverArtifact` already gates `execution-contract-incomplete` on `issue.Status.IsExecutable()`; adding `needs-repair` to the non-executable set requires only changing `IsExecutable()` or `ComputeFrontier`.

### Execution scenarios

- Success: a ticket with `**Status:** needs-repair` is excluded from frontier dispatch; a repair ticket with `**Repair of:** 03` and `**Status:** open` enters the frontier while its original is paused.
- Failure: a repair ticket referencing a non-existent original ID is diagnosed as a dangling reference by catalog; a `needs-repair` ticket that is also `Blocked by:` another ticket is not dispatched, confirming the status gate takes priority over unblocked logic.

### Expected tests

- `go test ./internal/tracker/...` — status lifecycle tests prove `needs-repair` is non-executable; parser tests extract `Repair of:` and retain legacy header compatibility; catalog tests validate reciprocal repair references and downstream reconnection.
- `go test ./internal/command/...` — frontier tests prove `needs-repair` is not dispatched and repair unblocks downstream only on completion.

### Generated artifacts

- Not applicable — this ticket changes tracker status, parser, and DAG projection behavior; no generated proposal payloads are produced.

### Conventions

- Preserve `Blocked by` and `Status` as the first human-visible ticket fields for legacy compatibility.
- Add `Repair of:` as a new human-visible header after `Status` and before `Blocked by`, parsed by a new regex in `parser.go`; do not re-parse the entire ticket body.
- `needs-repair` is a lifecycle state, not a readiness state; do not introduce a persisted `ready` or `execution-ready` phase.

## Implementation note

- Completed Changes: 1, 2, 3.
- Added `StatusNeedsRepair` to `internal/tracker/model.go:12`; `IsExecutable()` already returns `false` for it since it only matches `StatusOpen` and `StatusReadyForAgent`. Added `RepairOf string` field to `Issue` struct at `model.go:43`.
- Modified `ComputeFrontier` in `internal/tracker/dag.go:158-160` to skip any non-terminal, non-claimed status that fails `IsExecutable()`, fixing the "all nonterminal except claimed are ready" deviation.
- Added `repairOfRegex` to `internal/tracker/parser.go:19` and parsing logic after `blockedByRegex` match; `Repair of: None` or empty values leave `RepairOf` empty, preserving legacy header compatibility.
- Added `DiagnosticDanglingRepairReference` and `DiagnosticMissingReciprocalRepair` to `internal/tracker/catalog.go:53-54`. Added `validateRepairReferences()` method on `Catalog` called after `buildSemanticDiagnostics()`; it checks that every `Repair of:` target exists and every `needs-repair` ticket has a reciprocal repair.
- Tracker tests: `TestNeedsRepairIsNonExecutableAndRepairProvenanceIsParsed` and `TestRepairReciprocalReferencesAreValidated` in `catalog_test.go`.
- Command tests: `TestFrontierExcludesNeedsRepairAndDispatchesRepair` in `frontier_test.go` proving needs-repair exclusion, repair dispatch, and downstream blocking.
- Commands: `go test ./internal/...` — 114 passed in 6 packages.
- Files modified: `internal/tracker/model.go`, `internal/tracker/parser.go`, `internal/tracker/dag.go`, `internal/tracker/catalog.go`, `internal/tracker/catalog_test.go`, `internal/command/frontier_test.go`.
- Write set compliance: All modifications within write set (`internal/tracker/`, `internal/command/`).

## Review rounds

### Round 1

- Fixed point: current working tree
- Standards: [Low] Duplicated Code — three inline feature-qualified key calculations in `validateRepairReferences` instead of calling `issueKey`.
- Spec: none
- Fix changes: 4
- Design returns: none

### Round 2

- Fixed point: current working tree
- Standards: [Low] loop-in regex compilation in `hasCompleteExecutionContract` (pre-compiled as `executionContractSectionHeadings`).
- Spec: none
- Fix changes: 5
- Design returns: none

## Completion evidence

- Delivered: `needs-repair` status is nonterminal and non-executable; `ComputeFrontier` dispatches only `IsExecutable()` statuses. `Repair of:` provenance is parsed and does not create a DAG edge. Reciprocal repair references are validated with `dangling-repair-reference` and `missing-reciprocal-repair` diagnostics.
- Commands: `go test ./internal/...` — 114 passed in 6 packages; `flowforge check --dir docs/proposals/ticket-refinement-contract` — dependency graph healthy.
- Review axes: Round 1 — Standards [Low] duplicated key calculation (fixed via `issueKey`); Spec none. Round 2 — Standards [Low] loop-in regex compilation (fixed via pre-compiled `executionContractSectionHeadings`); Spec none.
- Deviations: none.
- Implementation reference: working-tree diff — `internal/tracker/model.go`, `internal/tracker/parser.go`, `internal/tracker/dag.go`, `internal/tracker/catalog.go`, `internal/tracker/catalog_test.go`, `internal/command/frontier_test.go`.
