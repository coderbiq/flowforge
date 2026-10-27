---
flowforge:
  schema: 1
  role: ticket
  id: derive-execution-contract-gap
  revision: 1
  consumes:
    requirements:
      ticket-refinement-contract-requirements: 1
    design:
      ticket-refinement-contract-design: 1
---

# 01: Derive execution-contract completeness as a frontier gap

**Blocked by:** None

**Status:** closed

## Delivery

An open code ticket without the five-part machine execution contract is reported as `execution-contract-incomplete` and excluded from the normal frontier, while remaining observable through `--include-gaps`.

## Design context

This ticket owns only the derived catalog/frontier signal. It does not create a persisted readiness state or teach an agent to author the contract. It realizes the [frontier eligibility outcome](../requirements.md#ticket-refinement-contract-requirements) through the [execution contract design](../design.md#ticket-refinement-contract-design).

## Touch points

- `internal/tracker/catalog.go` — executable ticket diagnostics
- `internal/tracker/catalog_test.go` — catalog fixtures
- `internal/command/frontier.go` and `internal/command/frontier_test.go` — gap projection

## Changes

- [x] 1. Add a structural detector for the five required execution-detail sections and placeholder-only content on nonterminal code tickets.
- [x] 2. Emit `execution-contract-incomplete` as a scoped gap without changing legacy-ticket compatibility or closed-ticket evidence checks.
- [x] 3. Add catalog and frontier fixtures proving default exclusion and `--include-gaps` visibility.
- [x] 4. Fix: Bound every required-section lookup to the `Execution detail` block so headings under `Implementation note` or later sections cannot satisfy the contract.
- [x] 5. Fix: Treat empty Markdown task-list entries (for example `- [ ]`) as placeholder-only contract content.
- [x] 6. Fix: Add catalog fixtures for each missing required section, an empty task-list placeholder, and a `Verified contracts`-only execution block followed by later headings.

## Constraints

- Do not introduce `ready` or `execution-ready` persisted status.
- Do not diagnose authority/design gaps as execution-contract gaps.
- Write set: `internal/tracker/`, `internal/command/`, `docs/proposals/ticket-refinement-contract/issues/01-derive-execution-contract-gap.md`
- standards: none found

## Done and verify

- Catalog cases cover each missing section and a complete contract: `go test ./internal/tracker/...` — all pass.
- Frontier cases distinguish normal and include-gaps projections: `go test ./internal/command/...` — all pass.

---

## Execution detail

### Verified contracts

- Ticket status is parsed from the human-visible header by `internal/tracker/parser.go`; catalog emits closed-ticket evidence diagnostics in `internal/tracker/catalog.go`.
- `internal/command/frontier.go` already excludes diagnostics classified as gaps from its default effective frontier.

### Execution scenarios

- Success: a ticket containing all five non-placeholder sections has no execution-contract diagnostic.
- Failure: a ticket missing `Verified contracts` remains visible only with `--include-gaps`.

### Expected tests

- `go test ./internal/tracker/...` — detector and compatibility fixtures pass.
- `go test ./internal/command/...` — JSON and human frontier projections preserve gap semantics.

### Generated artifacts

- Not applicable — this ticket changes tracker behavior, not generated proposal payloads.

### Conventions

- Reuse existing catalog diagnostic constructors and fixture style; do not add a second ticket parser.

## Implementation note

- Completed Changes: 1, 2, 3.
- Added `execution-contract-incomplete` as a catalog gap for schema-managed executable tickets missing a complete five-section execution contract. Legacy tickets and closed tickets remain outside this gate.
- Added catalog and frontier fixtures for complete, missing-section, placeholder, normal-frontier and `--include-gaps` cases.
- Commands: `go test ./internal/tracker/... ./internal/command/...` — 68 passed.
- Files modified: `internal/tracker/catalog.go`, `internal/tracker/catalog_test.go`, `internal/command/frontier_test.go`, this ticket.
- Write set compliance: All modifications within write set.
- Round-1 repair: bounded section lookup to the `Execution detail` H2 block, reject unchecked task-list placeholders, and added fixtures for every required section, boundary leakage, and `- [ ]`.
- Round-1 verification: `go test ./internal/tracker/... ./internal/command/...` — all pass.

## Review rounds

### Round 1

- Fixed point: current working tree
- Standards: [Low] `internal/tracker/catalog.go` is not gofmt-formatted.
- Spec: [High] required headings can be satisfied outside `Execution detail`; [High] `- [ ]` is accepted as content; [Medium] fixtures omit several required-section and boundary cases.
- Fix changes: 4, 5, 6
- Design returns: none

### Round 2

- Fixed point: current working tree
- Standards: none
- Spec: none
- Fix changes: none
- Design returns: none

## Completion evidence

- Delivered: `execution-contract-incomplete` gap diagnostic for schema-managed executable tickets missing any of the five required execution-detail sections (`Verified contracts`, `Execution scenarios`, `Expected tests`, `Generated artifacts`, `Conventions`). Legacy tickets (schema 0) and closed tickets are excluded. Placeholder-only content (empty task-list entries, TODO/TBD/N/A tokens, angle-bracket placeholders) is rejected. Section lookup is bounded to the `## Execution detail` H2 block to prevent headings under later sections from satisfying the contract.
- Commands: `go test ./internal/tracker/... ./internal/command/...` — 68 passed in 2 packages; `flowforge check --dir docs/proposals/ticket-refinement-contract` — dependency graph healthy, no cycles.
- Review axes: Round 1 — Standards [Low] gofmt (fixed); Spec [High] boundary leakage (fixed), [High] empty task-list (fixed), [Medium] missing fixtures (fixed). Round 2 — Standards none; Spec none.
- Deviations: none.
- Implementation reference: working-tree diff — `internal/tracker/catalog.go` (+82 -17), `internal/tracker/catalog_test.go` (+82), `internal/command/frontier_test.go` (+42).
