# Analysis Workflow

Use this loop for Proposal design. Every transition must leave recoverable Artifact or Journal state.

## 1. Intake and mode selection

1. Run `flowforge proposal inspect <id>`, `flowforge journal recent --proposal <id>`, and `flowforge analysis status --proposal <id>`.
2. Choose the simple path when the objective, current behavior, constraints, and decisive evidence already fit one coherent design.
3. Choose complex mode when two or more independent evidence gaps exist, context crosses modules or domains, a user claim needs verification, or likely findings may revise the design.
4. Add `<!-- analysis-mode: complex -->` only for complex FEATUREs. Do not create FIND cards or plan events for a simple path.

## 2. Frame and provisional design

Create 1--5 FEATURE skeletons with `card init --type feature`, then edit their bodies directly. Record `Objective`, `Current Understanding`, `Working Design`, assumptions, open questions, and `Next Investigation`. Select only the analysis modes that answer a real gap:

| Mode | Question answered | Artifact action |
|------|-------------------|-----------------|
| Claim Review | Is a user or design claim supported? | Record claim and support state in FEATURE Evidence. |
| Current-State | What does the system do now? | Cite code, tests, docs, or a FIND. |
| Root Cause | Why does an observed failure occur? | Separate observed chain from interpretation. |
| Boundary Analysis | Which ownership, compatibility, or security boundary applies? | Update Constraints or create/link a DEC. |
| New Design | What behavior and architecture should exist? | Revise Working Design and Key Decisions. |
| Implementation Design | Can an Executor act without new design choices? | Write self-contained Steps and Verification mapping. |

Run `library suggest --for <feature-id>` before accepting reusable constraints. A no-hit is a recorded gap, not permission to invent a rule.

## 3. Publish a bounded plan

For complex mode, create 2--6 independent work items unless fewer genuinely cover all gaps. Each item needs a stable `workId`, question, scope, role, inputs, sources, designated FIND `evidenceTarget`, dependencies, budget, required flag, `doneWhen`, and `parallelGroup`. Write the unsealed managed Journal event body directly, then seal it with the Journal CLI. Set a revision and re-entry condition before dispatch.

Do not split one question merely to increase parallelism. A work item is independent only when another Investigator can complete it from the brief and persisted inputs without a live session.

## 4. Investigate and synthesize

Coordinator dispatches only work returned by `flowforge analysis ready`. On re-entry, run `flowforge analysis reentry` and inspect each designated FIND. Classify every result as accepted, rejected, conflicting, or inconclusive; update FEATURE Evidence, Working Design, revised assumptions, questions, and DEC cards accordingly.

A follow-up outside the published scope is not dispatchable. Add it to `Next Investigation`, publish a superseding revision with a fresh budget, and then return control. If returned work adds no new information, mark it inconclusive or accepted with no design change; do not rerun it without a changed question or source.

## 5. Stop and gate

Stop dispatch when required work returned, budget expired, evidence conflicts, permissions are missing, or a finding invalidates the objective or plan. Use `readiness-gates.md`. Product behavior, compatibility, migration, security, and conflicting goals require a user decision event and a DEC before proceeding.

## Recovery walkthrough

After a lost session or deleted SQLite index: read Proposal and FEATURE artifacts, run `journal recent`, then `analysis rebuild`, `analysis status`, and `analysis ready`. Resume the action reported by derived state. Never reconstruct facts from memory or edit SQLite. When implementation reports a design gap, preserve Steps, History, and Verification, record re-entry in Journal, and start a new revision.
