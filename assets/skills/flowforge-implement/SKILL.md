---
name: flowforge-implement
description: Deliver one proposal ticket from effective linked specification through TDD implementation, two-axis review, completion evidence, and closeout — implement owns the build-and-close loop, NOT for design decisions — return those to flowforge-solution-design. Use for an executable frontier ticket or equivalent compact contract when the user says "deliver ticket"/"执行 ticket"/"实现".
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

**Standards pre-flight:** Before execution, check the ticket's Constraints for standards transcription status. The ticket must carry one of: `must`/`must not` standards clauses (transcription done — proceed), `standards: none found` (design authority had no applicable standards — proceed), or `standards: pending` (transcription not done — return the ticket to Plan). If the ticket has Changes and a `Write set` but lacks any of these markers, return it to Plan for standards transcription. Tickets without a `Write set` (pure documentation) skip this check. Do not transcribe standards yourself; return to Plan so Plan can transcribe from the design authority.

**Execution-contract pre-flight:** Before any implementation mode starts, verify the ticket has no `execution-contract-incomplete` gap. If the catalog reports this gap, stop and return the ticket to `flowforge-refine-ticket` — the contract must be filled with verified repository evidence before implementation. This preflight is non-negotiable: `--include-gaps` on the frontier or a caller request to include gaps does not bypass it. A ticket visible only through `frontier --include-gaps` because of an incomplete contract is explicitly not eligible for implementation.

### 2. Determine execution mode

If the ticket body declares `**Mode:** lightweight`, or the dispatch prompt explicitly declares lightweight mode, use **lightweight mode** (step 3). Otherwise, if the ticket has unchecked `- [ ]` Changes with mechanical steps, a `Write set:` in Constraints, and an `Execution detail` section, also use **lightweight mode** (step 3). This mode is for implementers that follow explicit instructions but cannot reliably search the codebase or self-review.

The implementer never self-selects or switches modes; mode is declared by the ticket (Plan/refine-ticket writes `**Mode:**`) or by the dispatcher.

Otherwise, use **full mode** (step 4). This mode is for capable agents that own the entire turn: implement, review, resolve findings, write evidence, and close.

### 3. Lightweight mode

You are an editor, not a designer: every decision has already been made in the ticket. Execute mechanically and report honestly. When an instruction is ambiguous, a repository fact contradicts the ticket, or completing a Change would require a design judgement, the only legal exit is `STATUS: BLOCKED` with the reason and the current state preserved — never improvise.

#### Phase 0: Restate before executing

Before touching any file, write the execution restatement at the top of the Implementation note: for every unchecked Change, one line restating the action plus the acceptance command that will verify it (from Done and verify or Expected tests). If any restatement does not match the ticket, stop with `STATUS: BLOCKED` instead of guessing.

#### Phase 0b: Preflight traversal

List every Constraint and every Execution scenario (both Success and Failure) from the ticket as a checklist, and mark each with its implementation landing point (file and symbol) or a blocker note. A Change that cannot honor a Constraint it must satisfy is a `STATUS: BLOCKED` return, not a silent deviation. Do not modify files outside the ticket's `Write set:`.

Execute unchecked Changes mechanically, self-check, and stop. Do not review, close the ticket, or make design decisions.

#### 3a. Execute unchecked Changes

Work through each `- [ ]` item in order. For each Change:

1. Read the named file and symbol from the Change description and Touch points.
2. Implement the described mechanical action at that location.
3. Run `go build` (or the project's compile command) after each Change.
4. Run focused tests relevant to the Change if test names are given in Expected tests or Done and verify.

If a Change cannot be completed (file not found, symbol moved, build fails after the change), stop, leave the Change unchecked, and note the blocker in the Implementation note.

Fail fast: retry a failed Change's verification command at most 2 times; once the ticket accumulates 5 failed repair rounds, stop and return `STATUS: BLOCKED` with the scene preserved. Do not attempt unbounded self-healing — agents succeed quickly and fail slowly.

If a Change requires a design decision (the action is ambiguous, or completing it would change a responsibility, interface, seam, or ordering), stop, leave the Change unchecked, and note it as a **design return** in the Implementation note.

Do not modify files outside the ticket's `Write set:`.

#### 3b. Mechanical self-check

Run all commands listed in Done and verify. Record pass/fail and observed output for each. If any command fails, stop, and report the failure in the Implementation note—do not attempt to debug or fix.

When a command fails, paste the failing command, its exit code, and the relevant error output verbatim into the Implementation note. Never paraphrase or summarize errors — the verbatim output is the feedback signal for the next agent.

#### 3c. Check off completed Changes

Change `- [ ]` to `- [x]` for each completed Change in the ticket, in the same edit as its evidence quadruple: an indented block under the checked item carrying `- cmd:` (the exact command), `- exit: 0`, `- output:` (1–3 key result lines quoted verbatim, e.g. `N tests completed, 0 failed`), and `- artifact:` (repository-relative path of the delivered file). If the exit code is not `0` — including environment-related failures — the Change stays unchecked; record the disposition, and return `STATUS: BLOCKED` when it blocks completion. Leave unchecked any Change that could not be completed.

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
