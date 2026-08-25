---
flowforge:
  schema: 1
  role: requirement
  id: documentation-contract-refinement
  revision: 2
---

<a id="documentation-contract-refinement"></a>
# FlowForge Documentation Contract Refinement

**Version:** 0.2 approved design baseline  
**Date:** 2026-08-25  
**Status:** Implemented by tickets 01–10; verified by [end-to-end evidence](evidence/end-to-end.md)

This document is the current requirements authority and navigation entry point. Detailed interface documents own their respective design contracts; historical pre-validation baselines live under `history/`.

## Problem Statement

FlowForge v5 has a valuable overall workflow: skills help people clarify a problem and design a solution, Markdown files preserve the resulting artifacts, and the CLI deterministically calculates an executable path through ticket blocking edges with `check` and `frontier`. Its use of explicit `MUST` and `MUST NOT` constraints is also valuable.

However, v5 does not currently guarantee the complete engineering memory previously represented across v1-v4: what the requirement is, why it matters, what solution was selected, and exactly what each task must accomplish. A ticket can be DAG-ready while still requiring an execution model to rediscover code, recover design intent, or make architecture decisions.

Restoring the old document formats unchanged is not acceptable. v1-v4 artifacts often contained repeated restatements, template-filling prose, vague actions, implementation commonplaces, and content that was difficult for humans to scan. The desired outcome is therefore not “more documentation.” It is complete coverage with high information value.

Two additional failure modes must be avoided:

1. Readiness must not recreate the v2/v3 card-state failure mode, where stale state or a flawed transition model prevents useful work from continuing.
2. Machine traceability identifiers such as `REQ-1` must not become the primary human reading interface. Machine-oriented identity should support documents without making prose harder to understand.

## Solution

Refine the v5 documentation contract without abandoning its local-first Markdown, skill collaboration, tracer-bullet tickets, or deterministic DAG engine.

The target is an end-to-end, progressively disclosed engineering memory:

```text
problem and evidence
  -> requirements and observable acceptance
  -> design decisions and constraints
  -> executable ticket increments
  -> DAG execution order
  -> implementation and verification evidence
```

Each layer owns information that is unique to it. Downstream documents provide a short semantic summary and link to upstream authority instead of copying whole sections. Human-readable prose remains authoritative; stable machine identifiers exist as supporting metadata.

Readiness is calculated from current files and DAG facts. It is explainable and overridable, not a manually maintained workflow state. The default workflow remains usable when quality diagnostics warn about incomplete content.

## User Stories

1. As a feature author, I want to record why a change is needed and what observable outcome is required, so that implementation does not substitute a different problem.
2. As a designer, I want to record the selected modules, interfaces, seams, responsibilities, constraints, and rejected alternatives, so that execution does not reopen settled decisions.
3. As a planner, I want to turn the selected design into tracer-bullet tickets with explicit blocking edges, so that FlowForge can calculate the execution frontier.
4. As an execution agent with limited context or reasoning capacity, I want a ticket to identify the approved seam, affected area, ordered change, constraints, completion conditions, and verification, so that I do not need to redesign the feature.
5. As a human reviewer, I want to understand a ticket without resolving opaque requirement or decision numbers, so that review remains fast and natural.
6. As a CLI consumer, I want stable machine metadata linking requirements, decisions, tickets, and evidence, so that traceability can be checked deterministically.
7. As a document reader, I want each fact or decision to have one authoritative home, so that repeated copies do not drift.
8. As a document reader, I want short semantic context beside a link, so that deduplication does not turn a ticket into an unreadable list of references.
9. As a maintainer, I want project-wide constraints defined once and referenced by features and tickets, so that every ticket does not repeat permanent rules.
10. As a maintainer, I want feature-specific and ticket-specific `MUST` and `MUST NOT` constraints to have explicit scope, so that they are actionable rather than rhetorical.
11. As a reviewer, I want vague or non-decidable language called out, so that phrases such as “properly,” “as needed,” and “follow best practices” do not masquerade as requirements.
12. As a reviewer, I want content to pass a deletion test, so that removing any retained paragraph would lose a fact, decision, constraint, action, verification method, unknown, or evidence.
13. As a project owner, I want completeness and information value assessed separately, so that solving missing information does not recreate verbose v1-v4 artifacts.
14. As a project owner, I want readiness diagnostics to explain every pass, warning, and failure, so that I can judge the actual execution risk.
15. As a project owner, I want quality concerns to be warnings by default, so that an imperfect template or heuristic cannot stop all work.
16. As a project owner, I want deterministic DAG errors distinguished from content-quality warnings, so that only genuine execution-graph failures block the normal frontier.
17. As a project owner, I want an explicit way to include unready work, so that I can proceed knowingly when a check is inapplicable or the situation is exceptional.
18. As an auditor, I want a persistent waiver to name the waived check and its reason, so that bypassing a diagnostic does not silently claim the risk is gone.
19. As an implementer, I want closing a ticket to preserve verification evidence and deviations without bloating the original execution instructions, so that input and outcome remain separately readable.
20. As a FlowForge user, I want specification and design files excluded from the executable frontier, so that only actual tickets are returned as work.
21. As a team adopting the refinement gradually, I want ordinary tickets to remain lightweight and higher-risk work to opt into a richer execution-ready profile, so that documentation cost is proportional to risk.
22. As a FlowForge maintainer, I want the CLI to validate structure, links, DAG facts, and evidence presence without accepting long document content as command arguments, so that the local-first boundary remains intact.

## Implementation Decisions

### Agreed requirements

- Preserve the v5 skill collaboration model and overall discuss-design-plan-execute workflow.
- Preserve local-first Markdown artifacts under the configured documentation directory.
- Preserve one-file-per-ticket blocking edges and deterministic `check` / `frontier` graph calculation.
- Preserve normative `MUST`, `MUST NOT`, `SHOULD`, and `MAY` language where it excludes a concrete failure path.
- Restore complete coverage from requirement through design, ticket execution, and verification evidence.
- Optimize for information value rather than document length or mandatory paragraph count.
- Make documents directly useful to humans as well as executable by limited-context models.
- Do not restore the old heavy card CLI, journal store, or manually maintained multi-stage workflow state machine.
- Do not introduce a CLI interface that carries long-form document content.

### Approved artifact model

Use four logical artifact roles:

1. A requirement artifact owns the problem, outcomes, scope, scenarios, acceptance criteria, terms, constraints, and unknowns.
2. A design artifact owns the current and target models, modules and interfaces, seams, responsibility mapping, flows, decisions, alternatives, migration strategy, constraints, and verification strategy.
3. Ticket artifacts own only the executable increment: delivery, relevant design context, blockers, affected area, ordered changes, ticket-specific constraints, completion conditions, and exact verification.
4. Evidence artifacts own results, commands and outputs, deviations, follow-ups, and implementation references after execution.

Large refactors may split the design by module while retaining one navigation and cross-module decision hub. Small features should not be forced to create empty files or sections merely to satisfy the model.

### Information-value contract

- Each retained statement must contribute at least one fact, requirement, decision, constraint, action, verification method, unknown, or item of evidence.
- A statement that can be removed without losing any of those meanings should be deleted, merged, or replaced by a reference.
- A fact has one authoritative location. Downstream artifacts use a concise semantic summary plus a link when local understanding requires context.
- Templates must require semantic questions to be resolved, not require prose to fill every possible heading.
- Empty or non-applicable sections should normally be omitted. `None` is useful only when proving that an explicitly required check was performed and found nothing.
- Avoid overlapping headings that restate the same goal, such as separate Summary, Goal, Objective, and What to build sections.
- Avoid non-decidable terms unless immediately defined by observable conditions.
- Do not copy project-wide constraints into every feature or ticket.
- Use complexity budgets and review signals rather than minimum word counts or acceptance-criterion quotas.

### Normative constraint contract

- Every normative constraint identifies its subject, required or forbidden behavior, and scope or condition.
- Important constraints have an observable review, test, or static-check method when practical.
- Constraints are scoped as project, feature/design, module, or ticket constraints.
- Tickets reference upstream invariants and expand only the constraints that are especially relevant or easy to violate in that increment.
- Normative language is not used for attitudes such as “MUST follow best practices”; it is used to exclude a specific incorrect outcome.

### Readiness and workflow contract

- Persist only the minimum ticket execution state. The current recommendation is `open` and `closed`, with `in-progress` considered only if it provides demonstrated value.
- `design-ready` and `execution-ready` are not manually advanced workflow states. They are derived diagnostics calculated from current content and graph facts.
- DAG readiness and content quality are separate dimensions.
- Deterministic graph failures may block the normal frontier: cycles, missing blockers, self-blocking edges, unparseable ticket structure, or a ticket explicitly blocked by an unresolved design question.
- Content concerns normally produce warnings: missing exact paths, broad touch points, vague acceptance, likely repetition, missing verification commands, or weak design context.
- Readiness output explains resolved, missing, and risky information rather than returning only a Boolean.
- The normal frontier remains useful and displays ready work, ready work with warnings, and blocked work.
- A command-level override can include unready work without modifying document state or hiding warnings.
- A persistent waiver names the specific diagnostic and records a reason. A blanket `skip_checks: true` escape hatch is not acceptable.
- An override permits informed execution; it does not mark a design question resolved or claim high execution confidence.

### Human and machine traceability contract

- Human-readable prose and headings are the authoritative reading interface.
- Human-facing prose uses semantic link text instead of requiring readers to resolve identifiers such as `REQ-1` or `DEC-2`.
- Stable machine identity is stored as supporting metadata or a non-disruptive anchor.
- Prefer semantic identifiers such as `req-application-extension-decoupling` over positional identifiers such as `REQ-1`.
- Human-facing links include a concise meaning, for example “Application does not depend on extension implementations.”
- Machine metadata may associate a ticket with requirement, decision, invariant, and evidence identifiers for deterministic traceability checks.
- Renaming a human heading should not silently destroy machine traceability.

### Skill responsibility direction

- Alignment owns the requirement truth: problem, scope, scenarios, acceptance, and unknowns.
- Domain modeling owns durable vocabulary and architectural decisions where applicable.
- A distinct solution-design responsibility owns modules, interfaces, seams, responsibilities, flows, migration, verification strategy, and design alternatives.
- Codebase design supplies the deep-module vocabulary and design discipline underneath solution design; it is not itself the artifact owner.
- Specification synthesis compacts agreed material without inventing missing decisions.
- Planning creates executable increments from an approved design; it does not hide unresolved architecture work inside tickets.
- Implementation consumes frontier tickets and returns design gaps rather than silently deciding them.
- Review checks both implementation conformance and document information value.
- Handoff preserves temporary session state but is not an authoritative long-term design source.

## Testing Decisions

Test at the highest stable seam available and keep semantic quality review distinct from deterministic CLI checks.

Deterministic CLI tests should cover:

- DAG cycles, missing edges, self-blocking, and frontier calculation.
- Exclusion of requirement, specification, design, and evidence artifacts from executable frontier results.
- Resolution of human-readable links and stable machine metadata.
- Derived readiness explanations without persisted readiness transitions.
- Default warning behavior, strict filtering if introduced, command-level overrides, and scoped persistent waivers.
- The rule that closing a ticket can require evidence without rewriting its execution instructions.
- Backward compatibility or graceful degradation for existing v5 ticket files.

Skill and document evaluation uses representative artifacts as well as parser fixtures. Tangram's backend architecture refinement is the first tracer case. Compare original and refined artifacts using:

- Human time to identify the goal, selected design, next change, constraints, and verification.
- Additional codebase exploration required by an execution model.
- Number of design decisions made during implementation that should have been settled earlier.
- Implementation deviations and rework discovered after tickets are marked closed.
- Repeated or removable prose.
- Token consumption as a secondary measure, not the quality target.

## Out of Scope

- Rejecting or replacing the v5 skill collaboration system.
- Replacing Markdown with a database-backed content store.
- Restoring the v1-v4 card, journal, or working-memory CLI wholesale.
- Making all features use large, multi-file specifications.
- Requiring every warning to block execution.
- Treating machine identifiers as the primary document language.
- Automatically judging architecture quality or prose value as a fully deterministic CLI operation.
- Automatically rewriting existing proposal files during discovery.
- Treating design approval as permission to change a public CLI signature without its own compatibility review.

## Further Notes

### Workflow prototype

The throwaway interaction model at [workflow-prototype.throwaway.html](workflow-prototype.throwaway.html) explores two guided cases plus a free-play mode:

- A simple parser/frontier behavior change tests whether the workflow can route directly through diagnosis, one ticket, TDD, and review without creating empty requirement and design artifacts.
- A complex cross-module architecture refinement tests progressive expansion through alignment, domain modeling, research and codebase design, optional prototyping, synthesis, planning, DAG calculation, implementation, review, and evidence.
- Free-play actions inject a content warning and apply a reasoned override, making the full derived state visible after every action.

The prototype was an exploratory question. Its validated conclusions are incorporated into the approved interfaces linked below; it is not itself workflow authority.

The first guided walkthrough is evaluated in [walkthrough-analysis-v0.1.md](walkthrough-analysis-v0.1.md). That analysis preserves discrepancies between the prototype and the current Skill contracts instead of silently rewriting the experiment.

The accepted responsibility split and the next design questions are captured in [responsibility-map-v0.1.md](responsibility-map-v0.1.md).

The approved external interface for the new design owner is captured in [`flowforge-solution-design` Interface Design v0.2](solution-design-interface-v0.2.md). Its pre-validation baseline is retained under `history/`.

Its second scenario experiment is captured in [Solution-design Scenario Validation v0.2](scenario-validation-v0.2.md) and the [throwaway v2 interaction prototype](workflow-prototype-v2.throwaway.html). The fixtures are examples, not executable tickets or Tangram architecture authority.

The approved information contract between Requirement, Solution design, Ticket, and Evidence is captured in [Artifact Hand-off Interface v0.2](artifact-handoff-interface-v0.2.md). Its pre-validation baseline is retained under `history/`.

Its simple and complex hand-off scenarios are evaluated in [Artifact Hand-off Scenario Validation v0.3](artifact-scenario-validation-v0.3.md) and the [throwaway v3 interaction prototype](workflow-prototype-v3.throwaway.html).

The approved physical representation for artifact roles and semantic dependencies is captured in [Minimal Physical Markdown Schema v0.1](physical-markdown-schema-v0.1.md).

Its Go/YAML feasibility experiment is captured in [Schema v1 Parser Prototype Verdict](parser-prototype-verdict-v0.1.md).

The approved proposal discovery and diagnostic projection seam is captured in [Artifact Catalog Interface v0.1](artifact-catalog-interface-v0.1.md).

Its real FlowForge and Tangram filesystem validation is captured in [Artifact Catalog Filesystem Prototype Verdict](catalog-filesystem-verdict-v0.1.md).

The approved coordinated changes to the requirement-to-delivery Skill system are captured in [Production Skill Coordination Design v0.1](skill-coordination-design-v0.1.md).

### Current design principles

1. Completeness and information density are independent quality dimensions and must be reviewed separately.
2. Readiness is derived, explainable, and overridable; it is not a manually maintained workflow state.
3. Human-readable prose is authoritative; machine identifiers are supporting metadata.
4. DAG calculation answers which work can start; an execution contract answers how that work can be completed correctly.
5. Local semantic summaries are allowed even when the full authority is linked. Complete deduplication must not make a ticket unreadable.
6. Logical artifact roles are stable, while physical packaging adapts to feature complexity.
7. Requirement alignment, solution design, and execution planning have distinct owners.

### Open design questions

- The threshold for splitting logical artifact roles into separate files for complex features.
- The final filenames and directory layout.
- The minimum ticket execution state vocabulary.
- Which deterministic defects block the default frontier and which only warn.
- The exact override and waiver interface.
- Whether a strict readiness profile is a CLI mode, project policy, ticket profile, or combination.
- The metadata representation for stable semantic identity and traceability.
- How existing v5 proposals migrate or degrade gracefully.
- How information-value review is divided between skills, human review, and optional heuristics.
- The exact highest testing seam for the refinement.

These questions are intentionally preserved as unresolved. This baseline must not be treated as approval of a concrete implementation design.
