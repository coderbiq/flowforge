---
name: flowforge-implement
description: Deliver one proposal ticket from effective linked specification through verification, two-axis review, completion evidence, and closeout. Use for an executable frontier ticket or equivalent compact contract.
disable-model-invocation: true
---

# Implement with evidence

Use the shared contract's [hand-offs](../_shared/ARTIFACT-CONTRACT.md#hand-offs), [diagnostics](../_shared/ARTIFACT-CONTRACT.md#diagnostics), and [information-value test](../_shared/ARTIFACT-CONTRACT.md#information-value). Read [schema v1](../_shared/SCHEMA-V1.md) only when promoting evidence.

## Process

### 1. Resolve the effective specification

Read project invariants, the ticket, its linked requirement/design authorities and consumed revisions, and applicable exact waivers. Run Catalog/frontier diagnostics for this ticket.

- A blocker stops dependent work.
- A gap stops only affected work unless the caller explicitly includes gaps; preserve the diagnostic in the report.
- A warning remains visible; strict policy may exclude it.
- An upstream revision mismatch is reviewed against current meaning before execution.

Do not advance status to express readiness. Begin only when the ticket's Delivery, design context, decided Changes, Constraints, and feasible Done and verify method form an executable contract.

**Standards pre-flight:** Before execution, check the ticket's Constraints for standards extraction status. The ticket must carry one of: `must`/`must not` standards clauses (extraction done — proceed), `standards: none found per guide` (extraction attempted, no applicable standards — proceed), or `standards: pending` (extraction not done — return the ticket to Plan). If the ticket has Changes and a `Write set` but lacks any of these markers, return it to Plan for standards extraction. Tickets without a `Write set` (pure documentation) skip this check. Do not extract standards yourself; return to Plan so Plan can perform the extraction.

### 2. Determine execution mode

If the ticket has unchecked `- [ ]` Changes with mechanical steps, a `Write set:` in Constraints, and an `Execution detail` section, use **lightweight mode** (step 3). This mode is for implementers that follow explicit instructions but cannot reliably search the codebase or self-review.

Otherwise, use **full mode** (step 4). This mode is for capable agents that own the entire turn: implement, review, resolve findings, write evidence, and close.

### 3. Lightweight mode

Execute unchecked Changes mechanically, self-check, and stop. Do not review, close the ticket, or make design decisions.

#### 3a. Execute unchecked Changes

Work through each `- [ ]` item in order. For each Change:

1. Read the named file and symbol from the Change description and Touch points.
2. Implement the described mechanical action at that location.
3. Run `go build` (or the project's compile command) after each Change.
4. Run focused tests relevant to the Change if test names are given in Expected tests or Done and verify.

If a Change cannot be completed (file not found, symbol moved, build fails after the change), stop, leave the Change unchecked, and note the blocker in the Implementation note.

If a Change requires a design decision (the action is ambiguous, or completing it would change a responsibility, interface, seam, or ordering), stop, leave the Change unchecked, and note it as a **design return** in the Implementation note.

Do not modify files outside the ticket's `Write set:`.

#### 3b. Mechanical self-check

Run all commands listed in Done and verify. Record pass/fail and observed output for each. If any command fails, stop and report the failure in the Implementation note—do not attempt to debug or fix.

#### 3c. Check off completed Changes

Change `- [ ]` to `- [x]` for each completed Change in the ticket. Leave unchecked any Change that could not be completed.

#### 3d. Write Implementation note

Write a `## Implementation note` section in the ticket (after the `---` separator, before Review rounds if present) recording:

- which Changes were completed (numbers) and which were not (with reason);
- commands run and their results (pass/fail, output summary);
- files modified;
- write-set compliance: "All modifications within write set" or list violations.

#### 3e. Stop

Do not invoke `flowforge-review`. Do not write Completion evidence. Do not change `**Status:**`. Commit the implementation work and return. The review agent will pick up from the Implementation note.

### 4. Full mode

Use `flowforge-tdd` internally at the pre-agreed verification seam. Run focused checks through the loop and the full relevant verification at the end. Preserve unrelated worktree changes.

Classify discoveries immediately:

- a repository fact that corrects stale location, symbol, or command detail updates the ticket/evidence without changing authority meaning;
- a local implementation detail inside the approved responsibility and seam is implementation work;
- a responsibility, interface, seam, information-flow, ordering, migration, or verification-strategy change returns to `flowforge-solution-design` with affected areas and work preserved;
- an observable outcome/scope/constraint change returns to `flowforge-align`.

#### 4a. Review one fixed change set

Invoke `flowforge-review` with a resolvable fixed point, an explicit committed or working-tree diff scope, and the effective specification links/revisions. Keep Standards and Specification findings separate.

Resolve every finding by correction or an authority-owned disposition. A Specification finding that means required delivery is absent keeps the affected ticket open until corrected, or until the requirement/design owner explicitly rescopes the effective specification or grants an exact reasoned waiver. Creating a follow-up alone is not resolution; it can permit close only when its recorded blocking relationship prevents premature close, or the applicable authority explicitly accepts the finding as nonblocking. Review does not close the ticket.

#### 4b. Record evidence, then close

Write concise Completion evidence inline by default. Promote schema v1 `role: evidence` only for multi-environment or multi-actor verification, shared integration proof, independent audit/lifecycle, or when inline proof obscures the ticket.

Evidence records:

- delivered behavior;
- commands or observation methods actually run and their observed results;
- both review axes and every finding disposition;
- deviations and how their authority owner handled them;
- implementation reference such as commit, diff, or changed artifact.

Summarize results; do not copy the ticket or dump terminal output. Write `**Status:** closed` only after observable delivery is verified and every review finding is corrected or has an authority-owned disposition that proves the current delivery acceptable. Required behavior deferred to a blocker keeps the ticket open until that blocker resolves. Checking boxes or creating a nonblocking follow-up is not completion evidence.

### 5. Publish the next frontier

Run full relevant tests, `flowforge check`, and `flowforge frontier`. Commit the completed ticket and implementation. Return evidence location, implementation reference, diagnostic/override facts, review dispositions, and the new frontier.

## Completion

The delivery is observable, consumed authority is current or explicitly dispositioned, verification results and both review axes are recorded, deviations are owned, evidence exists, and only then the ticket is closed and committed.
