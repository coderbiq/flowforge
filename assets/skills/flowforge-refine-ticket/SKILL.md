---
name: flowforge-refine-ticket
description: Fill the machine execution contract of one candidate ticket with verified repository evidence. Use when Plan has published a ticket skeleton whose only blocker is the execution-contract gap.
disable-model-invocation: true
---

# Refine one ticket execution contract

Use the shared contract's [hand-offs](../_shared/ARTIFACT-CONTRACT.md#hand-offs) and [information-value test](../_shared/ARTIFACT-CONTRACT.md#information-value).

## Process

### 1. Select the candidate

Identify a ticket whose only blocker is `execution-contract-incomplete`. The ticket must have been published by Plan with a titled `## Execution detail` skeleton and no DAG blockers. If the ticket has unresolved design questions or open authority gaps, return it to `flowforge-solution-design` or `flowforge-align` — do not invent values.

### 2. Gather verified repository evidence

Read the ticket's linked requirement and design authorities, their consumed revisions, and applicable waivers. Inspect the codebase at the Touch points to collect concrete facts: file paths, symbol names, API shapes, test commands, and generation/consumption relationships. Every fact written into the execution contract must be traceable to a repository location or an authority artifact.

If a fact cannot be verified — conflicting API shape, missing schema, unclear seam, or unknown ordering — stop and return a scoped design/research finding. Do not fill the section with placeholders or guessed values.

### 3. Fill the five execution-contract sections

Write or complete all five sections under the ticket's existing `## Execution detail` heading. Each section must carry verified content, not template text:

- **Verified contracts** — precise contract facts with file/symbol/anchor evidence (e.g. `internal/tracker/catalog.go:discoverArtifact` owns ticket diagnostics).
- **Execution scenarios** — at least one success path and one failure path with observable results (e.g. "Success: a complete ticket is eligible" / "Failure: a missing section is diagnosed").
- **Expected tests** — runnable commands or named test cases with assertions (e.g. `go test ./internal/tracker/...` — all pass).
- **Generated artifacts** — producer → artifact → consumer sync assertions, or an explicit `Not applicable` reason stating why no artifact is produced.
- **Conventions** — local non-obvious conventions and transcribed `must`/`must not` standards clauses tagged `[Conventions]` that the implementer needs.

Do not duplicate the `Write set` — verify it exists in Constraints, is narrow enough, and is consistent with the execution scenarios.

### 4. Verify readiness

Run `flowforge check --dir <proposal-dir>` on the refined ticket. The `execution-contract-incomplete` gap must clear. If it persists, the contract is still incomplete — re-read the diagnostic and fix the specific missing section or placeholder.

Verify the ticket carries a `**Mode:** lightweight` line (write it when Plan omitted it and the ticket is mechanical): the mode is declared here for the dispatcher, never self-selected by the implementer.

Do not introduce a persisted `ready` or `execution-ready` status. Readiness is derived from the current Markdown and diagnostics.

### 5. Return

Return the refined ticket path, the verified facts and their sources, and the check result. The ticket is now eligible for `flowforge-implement` lightweight mode.

## Boundaries

- Refine Ticket does not implement, review, create readiness states, or invent unresolved values.
- It fills exactly one ticket's execution contract; it does not slice, re-plan, or re-design.
- If a fact conflict requires a design decision, return to `flowforge-solution-design` with the affected area and evidence.
