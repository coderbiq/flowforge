---
name: flowforge-route
description: Select the FlowForge Skill that owns the next unresolved question. Use when the appropriate engineering flow is unclear; return one next owner and the reason, without creating feature content.
disable-model-invocation: true
---

# Route by unresolved owner

Use the shared contract's [authority role boundaries](../_shared/ARTIFACT-CONTRACT.md#roles-and-authority) and [ticket hand-off checks](../_shared/ARTIFACT-CONTRACT.md#hand-offs). Route chooses one next owner; the target Skill owns its process and artifacts.

## Main routes

- Broken, failing, or regressed behavior needing a cause → `flowforge-diagnose`.
- Raw external bug/request needing classification and a crisp brief → `flowforge-triage`.
- Local PRD, old proposal, brief, or notes are the request's starting material → `flowforge-import`; its classified hand-off returns to Align or Solution Design.
- Requirement outcome, scope, scenario, constraint, or term remains unsettled → `flowforge-align`.
- Requirement is settled but responsibility, interface, seam, cross-module flow/order, migration, verification strategy, or credible alternatives remain → `flowforge-solution-design`.
- Small settled work clearly reuses a seam identified by current design authority → `flowforge-plan` for a compact ticket, or `flowforge-implement` when an existing ticket satisfies the contract's delivery, context, changes, constraints, and verification hand-off.
- Work is too foggy to expose one-session decision frontiers → `flowforge-wayfinder`; its decisions return to Align/Solution Design before planning.
- A selected architecture candidate needs a deeper module boundary → Align if its outcome is unclear, otherwise Solution Design with `flowforge-codebase-design` underneath.

Supporting Skills are detours, not user-managed phases: Research resolves a named primary-source fact; Prototype answers one runnable design question; Domain Modeling owns durable vocabulary/ADRs; Codebase Design supplies module/interface/seam vocabulary.

## Artifact adaptation

Do not classify a feature as a simple or complex workflow state. Keep roles compact until independent consumption, review, lifecycle, or readability earns promotion into requirement/design/evidence files.

## Other routes

- Architecture health scan → `flowforge-improve-architecture`.
- Review a fixed change set → `flowforge-review`.
- Continue a specified ticket → `flowforge-implement`.
- Resolve an active merge/rebase conflict → `flowforge-resolving-conflicts`.
- Teach a concept → `flowforge-teach`.
- Human-only external procedure → `flowforge-wizard`.
- No repository, only an idea to stress-test → `flowforge-grill-me`.

Read [phase boundaries](PHASE-BOUNDARIES.md) only when choosing continue, compact, handoff, or a separate execution context.

## Completion

Return the selected Skill and one sentence naming the unresolved fact it owns. Route creates no requirement, design, ticket, state, or readiness artifact.
