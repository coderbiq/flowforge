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

### 2. Deliver at approved seams

Use `flowforge-tdd` internally at the pre-agreed verification seam. Run focused checks through the loop and the full relevant verification at the end. Preserve unrelated worktree changes.

Classify discoveries immediately:

- a repository fact that corrects stale location, symbol, or command detail updates the ticket/evidence without changing authority meaning;
- a local implementation detail inside the approved responsibility and seam is implementation work;
- a responsibility, interface, seam, information-flow, ordering, migration, or verification-strategy change returns to `flowforge-solution-design` with affected areas and work preserved;
- an observable outcome/scope/constraint change returns to `flowforge-align`.

### 3. Review one fixed change set

Invoke `flowforge-review` with a resolvable fixed point, an explicit committed or working-tree diff scope, and the effective specification links/revisions. Keep Standards and Specification findings separate.

Resolve every finding by correction or an authority-owned disposition. A Specification finding that means required delivery is absent keeps the affected ticket open until corrected, or until the requirement/design owner explicitly rescopes the effective specification or grants an exact reasoned waiver. Creating a follow-up alone is not resolution; it can permit close only when its recorded blocking relationship prevents premature close, or the applicable authority explicitly accepts the finding as nonblocking. Review does not close the ticket.

### 4. Record evidence, then close

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
