# Production Skill Coordination Design v0.1

**Date:** 2026-08-25  
**Status:** Design baseline; production Skill files are unchanged  
**Scope:** Route, Align, Solution Design, Synthesis, Plan, Implement, TDD, Review, and supporting disciplines

## Objective

Make each Skill a deep module with one primary responsibility and a small hand-off interface, while preserving the v5 Skill collaboration model and removing reliance on one unbroken context window.

The design uses:

- [`flowforge-solution-design` Interface v0.2](solution-design-interface-v0.2.md);
- [Artifact Hand-off Interface v0.2](artifact-handoff-interface-v0.2.md);
- [Physical Markdown Schema v0.1](physical-markdown-schema-v0.1.md);
- [Artifact Catalog Interface v0.1](artifact-catalog-interface-v0.1.md).

## Information hierarchy

Every production Skill keeps:

1. ordered steps and completion criteria in its main `SKILL.md`;
2. shared artifact semantics in one disclosed reference rather than copied into every Skill;
3. branch-specific guidance behind pointers reached only when that branch fires;
4. project facts in repository files and CLI output, not cached in Skill prose.

The Skill system uses the same leading words consistently:

- **requirement authority** — why and observable outcome;
- **solution design** — modules, interfaces, seams, ownership, migration, verification;
- **execution contract** — one ticket's implementable increment;
- **effective specification** — resolved project + requirement + design + ticket view;
- **gap / warning / blocker** — derived diagnostic facts;
- **frontier** — DAG-ready work under the selected policy.

## Shared reference design

Create one shared reference for the artifact hand-off contract used by Align, Solution Design, Plan, Implement, and Review.

Candidate production location:

```text
assets/skills/references/flowforge-artifact-contract.md
```

It contains only stable cross-Skill reference:

- artifact roles and authority order;
- semantic-link and machine-metadata rules;
- compact-ticket versus promoted-artifact rules;
- gap/warning/blocker meanings;
- effective-specification resolution;
- information-value deletion test.

It does not repeat each Skill's ordered steps.

If the runtime skill packaging cannot reliably resolve a shared sibling reference, install the reference into a predictable FlowForge namespace and make each pointer explicit. Do not duplicate the contract as fallback.

## `flowforge-route`

### Responsibility

Select the owner of the next unresolved question. Route does not create feature content.

### Main paths

```text
broken behavior              -> diagnose
raw external request         -> triage
small, settled feature       -> compact alignment / plan or implement
solution-affecting feature   -> align -> solution-design
fog-of-war effort            -> wayfinder -> align/design
architecture candidate       -> align -> solution-design
```

### Required change

- Add `flowforge-solution-design` as the owner of implementation responsibility/seam decisions.
- Replace the fixed implication “working directory means always start with full align flow” with routing by unresolved requirement/design uncertainty and session durability.
- Keep domain modeling and codebase design underneath owners rather than as user-managed phases.
- Describe artifact promotion instead of simple/complex workflow states.
- Point to the shared artifact contract rather than restating it.

### Completion criterion

Route is complete when one next owner and the reason for invoking it are clear. It produces no durable artifact.

## `flowforge-align`

### Responsibility

Establish and incrementally maintain requirement authority.

### Required steps

1. Locate or create the feature's requirement authority only when requirement information must survive the conversation.
2. Inspect repository facts before asking the user.
3. Use grilling for the current requirement decision frontier.
4. Use domain modeling only for durable vocabulary or qualifying ADRs.
5. After each settled requirement decision, update authority immediately.
6. Remove superseded statements instead of appending conversation history.
7. Finish with requirement coverage and explicit unknowns by affected scope.

### Adaptive output

- Small feature: requirement role remains inside a compact ticket/conversation hand-off.
- Complex or independently reviewed feature: `requirements.md` with schema v1 role metadata.

### Completion criterion

Every requirement question that changes the solution space is resolved, explicitly scoped as unknown, or consciously waived. No `requirements-ready` state is written.

### Guardrails

- Align does not choose module interfaces or seams.
- `CONTEXT.md` remains glossary-only.
- ADRs remain sparse and trade-off driven.
- Requirement authority contains no template padding or implementation plan.

## `flowforge-solution-design` (new)

### Responsibility

Advance approved requirements into authoritative solution design that planning can consume without architecture invention.

### Main Skill shape

The main `SKILL.md` contains:

1. invocation signals;
2. required inputs;
3. establish design frontier;
4. resolve facts before decisions;
5. orchestrate supporting Skills by need;
6. update design authority after each settled decision;
7. return requirement gaps through Align ownership;
8. derive scoped planning coverage;
9. run information-value compression;
10. report authority links, decisions, gaps, warnings, and planning recommendation.

Branch-specific reference may cover design hub/child-file packaging and area-level revision mechanics.

### Completion criterion

Each affected design area has an explicit coverage result. Planning may proceed only for areas whose required module, interface, seam, ownership, migration, and verification decisions are settled or visibly overridden.

### Guardrails

Return an authoritative solution design, scoped open items, and an explicit planning handoff. Keep readiness derived, leave ticket publication to Plan, leave production changes to Implement, return requirement ambiguity to Align, and retain only exploration that changed a decision.

Hard boundaries:

- MUST NOT persist a `design-ready` state.
- MUST NOT publish tickets or production implementation.
- MUST NOT silently resolve a requirement ambiguity.

## `flowforge-to-spec` (optional synthesis)

### Responsibility

Create a concise navigation/review baseline when separately authoritative requirements and design need one entry point.

### Required change

- Remove the mandatory “extremely extensive” user-story list.
- Remove responsibility for inventing implementation decisions.
- Stop acting as the first persistence point.
- Replace the monolithic template with a compact overview containing semantic links, scope, key decisions, verification navigation, and remaining gaps.
- Emit schema v1 `role: spec`.
- Explicitly state that the result is not executable and not a second authority.

### Invocation

Use for multi-session features, external review packages, or feature navigation. Skip when one compact ticket is the clearest artifact.

### Completion criterion

A new reader can find requirement, design, tickets, and current gaps from the overview without duplicated source content.

## `flowforge-plan`

### Responsibility

Convert settled requirement/design authority into independently verifiable execution contracts and DAG edges.

### Required steps

1. Resolve the effective requirement and solution-design authorities.
2. Read scoped design coverage; plan unaffected areas only.
3. Return any seam/responsibility choice to solution design.
4. Draft tracer increments and real blocking edges.
5. Show the user title, delivery, and blockers for approval.
6. Publish schema v1 tickets under `issues/`.
7. Run Catalog/DAG checks and report gaps/warnings without writing readiness state.

### Ticket output

Use the seven semantic roles:

- Delivery
- Design context
- Blocked by
- Touch points
- Changes
- Constraints
- Done and verify

Write `**Status:** open`. Preserve `**Blocked by:** None` visibly when unblocked. Include machine consumption metadata only for independently versioned authorities.

### Completion criterion

Every published ticket is one verifiable increment with valid DAG edges and no hidden design choice. Diagnostics and overrides are explicit.

### Guardrails

- Codebase exploration is not optional when touch points or existing seams are unknown.
- Planning may preserve an obvious increment at an existing seam; it does not choose or move seams.
- No file-path ban: include stable physical facts whose absence would force rediscovery.
- No `ready-for-agent` emission; DAG and diagnostics derive executability.

## `flowforge-implement`

### Responsibility

Own one implementation turn from effective specification through evidence-backed close-out.

### Required steps

1. Resolve the effective specification and current consumed revisions.
2. Run Catalog/frontier diagnostics for the ticket.
3. Stop only affected work on a gap/blocker; make overrides visible.
4. Use TDD internally at pre-agreed seams.
5. Run focused checks regularly and full relevant verification at the end.
6. Classify discoveries as factual correction, local design detail, or design change.
7. Return design changes to solution design; do not silently decide them.
8. Invoke two-axis review with a fixed point and effective specification.
9. Resolve findings or create explicit follow-up work.
10. Write completion evidence and close the ticket.

### Evidence output

- Inline for compact work unless independent authority is justified.
- Separate evidence artifact for multi-environment, shared integration, or independently audited proof.
- Actual commands/results, review disposition, deviations, and implementation reference only.

### Completion criterion

The observable delivery is verified, findings are resolved/waived/followed up, deviations are recorded, evidence exists, and only then is status written `closed`.

## `flowforge-tdd`

### Responsibility

Remain the red-green behavior discipline inside Diagnose or Implement.

### Required change

- Resolve the pre-agreed seam from the effective specification instead of always asking the user again.
- Ask only when the seam is absent or conflicting; that absence becomes a design gap.
- Do not become a separate persisted phase.

No artifact authority changes are needed.

## `flowforge-review`

### Responsibility

Independently evaluate Standards and Effective Specification conformance.

### Required steps

1. Resolve the fixed comparison point supplied by the implementation owner or user.
2. Resolve effective specification from project invariants, requirement/design links, ticket, and approved waivers.
3. Run Standards and Specification axes independently.
4. Report findings with authority citations.
5. Return findings to implementation; do not write completion evidence or close tickets.

### Required change

- Do not assume one monolithic spec file.
- Accept one ticket plus linked authority as the specification source.
- Support uncommitted implementation review when the caller supplies a valid fixed point and working-tree diff.
- Add document information-value findings only when documentation is part of the diff; keep them distinct from implementation conformance.

### Completion criterion

Both axes report independently against the same fixed change set; every finding cites its standard or effective-specification authority.

## Supporting Skills

### `flowforge-domain-modeling`

No ownership expansion. Continue to own glossary and qualifying ADRs only.

### `flowforge-codebase-design`

No artifact ownership. Continue as module/interface/seam vocabulary underneath solution design, TDD, and architecture improvement.

### `flowforge-research`

Continue producing cited primary-source notes. Add a completion requirement that its conclusion returns to the named requirement/design question and authority link.

### `flowforge-prototype`

Continue answering one runnable question. Missing packaged references (`LOGIC.md`, `UI.md`) must be fixed in Skill distribution before relying on it as a production branch.

### `flowforge-handoff`

Continue transporting only context delta. Add pointers to requirement/design/ticket/evidence authority and never duplicate them.

## Catalog/CLI responsibility

Skills author and update artifacts. The CLI validates deterministic structure:

- role/location;
- DAG edges;
- identity/anchor uniqueness;
- consumed revisions;
- explicit semantic links;
- legacy/strict compatibility.

Skills and review judge:

- information value;
- requirement/design completeness;
- seam quality;
- whether a semantic revision should increment;
- whether artifact promotion earns its cost.

Neither side pretends to own the other's judgment.

## Router context-load budget

The router should name each main path and one trigger sentence. Detailed artifact rules remain behind pointers. Always-loaded AGENTS guidance needs only:

- route by intent;
- files own content;
- CLI owns graph/catalog checks;
- use solution design for responsibility/interface/seam decisions;
- use planning only after those decisions settle.

Do not inline the full artifact schema or Skill flow in AGENTS.md.

## Production packaging risks

The current installed `flowforge-writing-for-agents` references a missing `SKILL-MECHANICS.md`; `flowforge-prototype` references missing `LOGIC.md` and `UI.md`. Production Skill implementation must audit referenced resources and deployment tests so pointers cannot lead to absent files.

This is a distribution correctness requirement, not a reason to duplicate the missing contents into every Skill.

## Verification scenarios

1. Small feature routes to one compact ticket without solution-design ceremony.
2. Complex feature persists requirement decisions before context compaction.
3. Solution design returns scoped coverage and partial planning.
4. To-spec creates navigation without copying authorities.
5. Plan publishes open tickets with human blockers and machine dependencies.
6. Implement detects stale design revision and records visible disposition.
7. TDD consumes an existing seam without re-asking.
8. Implementation returns a design change but handles a factual correction locally.
9. Review resolves an effective specification from links.
10. Evidence exists before close.
11. Legacy tickets continue through the flow with aggregated compatibility diagnostics.
12. Every Skill pointer resolves to a deployed resource.

## Implementation order recommendation

1. Artifact Catalog and schema parsing seam.
2. `check/frontier` diagnostic projection.
3. Shared artifact-contract reference and deployment correctness.
4. `flowforge-solution-design` new Skill.
5. Align and route hand-off changes.
6. Plan ticket authoring changes.
7. Implement/TDD/review/evidence changes.
8. Optional to-spec synthesis change.
9. End-to-end scenario fixtures and migration documentation.

This order establishes deterministic substrate before Skills emit new metadata, then moves authority producers before consumers.

## Out of scope

- Writing production Skill files in this design document.
- Implementing CLI/parser changes.
- Deciding marketplace/plugin packaging unrelated to deployed Skill resources.
- Rewriting standalone Skills outside the requirement-to-delivery flow.

## Next work

The interfaces and sequence are sufficiently defined to create tracer-bullet implementation tickets after one final dual-axis design review of the proposal artifacts.
