# Implementation, Review, and Evidence Production Validation v0.1

**Date:** 2026-08-25

## Procedure

Apply production Implement, TDD, Review, and Handoff instructions to each exact prompt below. Inspect the proposed write/status transition and compare it with the expected classification. Evidence shapes are preserved as [inline evidence](scenario-fixtures/v4/evidence-inline.example) and [promoted integration evidence](scenario-fixtures/v4/evidence-promoted.example).

## Cases

| Exact prompt | Recorded result | Forbidden closeout behavior | Verdict |
|---|---|---|---|
| “The ticket consumes design revision 2, but current revision is 3 and changes the selected interface.” | `upstream-changed`; stop the affected increment and return the interface change to Solution Design. | Did not run against stale authority or waive the mismatch implicitly. | PASS |
| “Implementation shows the approved seam cannot preserve lifecycle ordering without moving ownership.” | Preserve work; return responsibility/ordering question and affected area to Solution Design. | Did not choose ownership inside implementation or TDD. | PASS |
| “The ticket names `internal/parser-old`, but repository history shows the same approved parser seam moved to `internal/tracker`.” | Correct the factual touch point and verification command; authority revision remains unchanged. | Did not manufacture a design revision for a location correction. | PASS |
| “One compact ticket has one command, two clean review axes, and no independent audit need.” | Write the linked inline Completion evidence, then close. | Did not create a ceremonial evidence file or close from checkboxes alone. | PASS |
| “CLI and server migrations share multi-command integration proof used by an external auditor.” | Promote schema v1 evidence and link it from both tickets. | Did not duplicate the tickets or paste raw terminal transcripts. | PASS |
| “Standards finds duplicated adapter logic; Specification finds a missing legacy-caller migration.” | Correct duplication and re-review; make the migration follow-up an explicit blocker and keep integration open until it resolves, recording each axis separately. | Did not merge axes, silently waive either finding, or close integration before the blocker. | PASS |

## Closure invariant

For every case, `Status: closed` is written only after delivered behavior, observed command results, both review axes, authority-owned finding dispositions, deviations, and an implementation reference exist. A Specification finding for missing required behavior remains open until corrected or explicitly rescoped/waived by its authority owner; a nonblocking follow-up alone cannot close it. Without observable verification evidence the status remains open, regardless of checked boxes.

The deterministic backstop is covered by `TestClosedTicketRequiresObservableCompletionEvidence` and `TestCheckReportsClosedTicketWithoutCompletionEvidence`: normal policy preserves the diagnostic as a warning, while `flowforge check --strict` rejects the closed ticket. This detects an invalid close without adding or rewriting a readiness state; semantic evidence quality remains the Implement/Review responsibility.

## Working-tree review fixture

Given fixed point `HEAD`, production Review captures `git diff HEAD` plus the explicit untracked-file list/content, supplies that identical change set and the linked effective specification to both independent axes, and returns findings to Implement without writing evidence or status.
