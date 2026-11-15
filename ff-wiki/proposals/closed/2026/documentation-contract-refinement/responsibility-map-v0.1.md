# Skill and Artifact Responsibility Map v0.1

**Date:** 2026-08-25  
**Status:** Design draft following walkthrough approval  
**Scope:** From a user requirement entering FlowForge through verified delivery

## Purpose

Define one clear owner for every durable kind of information and every workflow decision before designing Markdown schemas or CLI behavior.

The map follows four rules:

1. A Skill is a module with one small interface and one primary responsibility.
2. A durable fact has one authoritative artifact.
3. A hand-off occurs because an unresolved question changes, not because a persisted workflow state advances.
4. A Skill MUST NOT create an artifact when it has no new authoritative information for that artifact.

## Skill classes

### Main-flow owners

These Skills own a durable transformation in the normal requirement-to-delivery flow.

| Skill | Trigger / input interface | Primary responsibility | Durable output | MUST NOT |
|---|---|---|---|---|
| `flowforge-route` | A raw request plus repository context | Classify the next unresolved question and select a flow | None | Restate the request, create feature documents, or persist a workflow phase |
| `flowforge-align` | A feature whose required outcome, scope, scenarios, or constraints remain unsettled | Establish and maintain feature requirement truth | Requirement role in the feature artifact set | Decide implementation seams, decompose tickets, or fill generic template prose |
| `flowforge-solution-design` **(new responsibility)** | Approved requirements plus code facts and project authorities | Decide the modules, interfaces, seams, ownership, flows, migration, constraints, and verification strategy | Design role in the feature artifact set | Reopen product requirements silently, publish execution tickets, or treat vocabulary guidance as the design itself |
| `flowforge-to-spec` **(re-scoped candidate)** | Authoritative requirement and design material that needs a portable review baseline | Synthesize navigation and a compact agreed baseline | Optional feature overview/spec role | Invent missing decisions, duplicate upstream authority, require exhaustive user-story padding, or act as the first persistence point |
| `flowforge-plan` | Settled requirements and solution design | Convert design into independently verifiable execution increments and blocking edges | Ticket artifacts | Make new architecture decisions, hide unknowns in “What to build,” or publish a ticket solely to represent a phase |
| `flowforge-implement` | One frontier ticket plus its effective linked specification | Own the implementation turn, invoking TDD and review disciplines and returning design gaps | Code, tests, and completion evidence | Silently settle design gaps, expand scope without a ticket/design update, or require the user to manually orchestrate nested Skills |
| `flowforge-review` | A change set, fixed comparison point, and effective specification | Independently assess Standards and Specification conformance | Review findings consumed by completion evidence | Become the implementation owner, merge the two axes, or declare completion without verification evidence |

### On-ramps

These Skills enter or re-enter the main flow from a special situation.

| Skill | Trigger | Responsibility | Merge point | MUST NOT |
|---|---|---|---|---|
| `flowforge-triage` | Raw external requests or bugs accumulating in an inbox | Verify, classify, and turn raw external input into a trustworthy brief | Align, diagnose, or direct implementation when already agent-ready | Re-triage tickets produced by planning |
| `flowforge-diagnose` | Broken, failing, throwing, flaky, or slow behavior | Own the feedback loop, reproduction, cause, regression proof, fix, and cleanup | Review/evidence; plan or solution design only if the diagnosis expands scope | Stop after speculation, duplicate a separate TDD phase, or force planning for a one-session fix |
| `flowforge-wayfinder` | A destination whose decision path cannot fit in one session | Maintain a DAG of decision questions until the route becomes visible | Requirements and solution design, then normal planning | Produce delivery tickets before decisions are resolved or implement from the decision map directly |
| `flowforge-improve-architecture` | A codebase-health scan identifies a deepening opportunity | Produce a candidate architectural problem | Align | Implement the candidate or treat survey findings as an approved solution |

### Supporting disciplines and fact-producing detours

These are invoked inside another owner. They do not represent mandatory serial phases.

| Skill | Invoked when | Contribution | Durable output | MUST NOT |
|---|---|---|---|---|
| `flowforge-domain-modeling` | Terms conflict, a durable relationship is resolved, or a decision qualifies for an ADR | Canonical vocabulary and durable architectural context | `CONTEXT.md` and qualifying ADRs | Store feature requirements, implementation detail, or routine decisions |
| `flowforge-codebase-design` | A module interface or seam is being selected | Deep-module vocabulary and design tests | Normally none by itself; decisions belong to solution design | Pretend vocabulary guidance is a complete solution-design phase |
| `flowforge-research` | A decision depends on facts requiring primary-source investigation | Cited facts and bounded conclusions | Research note | Replace the decision owner or copy research narrative into every downstream artifact |
| `flowforge-prototype` | A runnable experiment is needed to settle one design question | Primary-source experiment and explicit verdict | Throwaway prototype reference plus verdict in design | Become production implementation or answer several unrelated questions |
| `flowforge-tdd` | Implementation needs behavior proof at a pre-agreed seam | Red-green discipline and behavior tests | Tests and implementation inside the owning turn | Become a separate user-managed workflow state or choose an unapproved seam |
| `flowforge-grilling` | An owner has a decision frontier requiring user choices | Interview discipline over the current decision tree | Decisions written by the invoking owner | Own feature artifacts or ask the user for discoverable facts |

### Context transport

| Skill | Trigger | Responsibility | Output | MUST NOT |
|---|---|---|---|---|
| `flowforge-handoff` | Work must cross a session, harness, directory, or person and compaction is insufficient | Transport only the context not already preserved by authoritative artifacts | Temporary handoff document with links | Duplicate specs, designs, tickets, ADRs, diffs, or evidence; become long-term authority |

The current approved external interface for the new owner is defined in [`solution-design-interface-v0.2.md`](solution-design-interface-v0.2.md). [v0.1](history/solution-design-interface-v0.1.md) is retained as the pre-validation baseline.

The approved information contract between Requirement, Solution design, Ticket, and Evidence is defined in [Artifact Hand-off Interface v0.2](artifact-handoff-interface-v0.2.md). [v0.1](history/artifact-handoff-interface-v0.1.md) is retained as the pre-validation baseline.

The approved physical artifact representation is defined in [Minimal Physical Markdown Schema v0.1](physical-markdown-schema-v0.1.md).

The coordinated production Skill responsibilities and hand-offs are defined in [Production Skill Coordination Design v0.1](skill-coordination-design-v0.1.md).

## Artifact roles

Artifact **role** and physical **file** are deliberately separate. A small feature may place several roles in one ticket; a complex feature may give each role its own file hierarchy.

| Artifact role | Authority | Created when | Required content | MUST NOT contain |
|---|---|---|---|---|
| Project vocabulary | Canonical domain terms and relationships | A term becomes shared and durable | Term, precise meaning, distinctions | Feature scope, implementation plans, session notes |
| ADR | A hard-to-reverse, surprising trade-off decision | All ADR qualification criteria are met | Context, decision, alternatives, consequences | Routine or easily reversible choices |
| Requirement | Why the feature exists and what observable outcome is required | Behavioral, scope, scenario, or constraint information must survive the conversation | Problem/evidence, outcomes, scope, scenarios, acceptance, requirement constraints, unknowns | Implementation sequence, speculative interfaces, duplicated glossary |
| Solution design | How the approved requirement will be realized | Implementation choices affect modules, interfaces, seams, ownership, migration, or verification | Current/target model, decisions, physical facts where valuable, alternatives, flows, constraints, test seams, open design questions | Repeated requirement narrative, execution logs, unresolved choices disguised as instructions |
| Feature overview/spec | Compact navigation and agreed baseline | Separate requirement/design authorities need a portable entry point or external review package | Semantic links, concise decisions and scope summary | A second copy of all requirement/design content or invented completeness |
| Ticket | One independently verifiable increment and its blocking edges | Work spans a context boundary, needs DAG coordination, or benefits from durable assignment | Delivery, semantic design context, blockers, touch points, ordered changes, ticket constraints, done conditions, exact verification | Whole-feature background, project-wide constraint copies, raw logs, new architecture choices |
| Completion evidence | What actually happened and why the ticket may close | Execution has produced verification and review results | Commands/results, delivered behavior, deviations, review disposition, implementation reference | Future plan, copied ticket instructions, unfiltered terminal transcript |
| Research note | Recoverable primary-source facts | Investigation is substantial enough that repeating it would be wasteful | Question, sources, findings, limits | Product/design decision claimed without its owner |
| Prototype evidence | Runnable primary source for one design verdict | A prototype answers a named question | Question, variants, observation, verdict reference | Production guarantees or unrelated experiments |
| Decision map | Fog-of-war decision DAG | The route to a destination exceeds one session | Decision questions, blockers, resolutions, destination criteria | Delivery tasks masquerading as decisions |
| Handoff | Temporary context delta | Work crosses a context boundary and artifacts are insufficient | Current focus, unresolved local context, links, suggested Skills | Duplicated authoritative content |

## Adaptive physical packaging

### Small feature

One ticket may contain four visually distinct roles:

```text
Requirement — Why / observable behavior
Design increment — existing seam and selected change
Execution — changes / constraints / done / verify
Completion — evidence written after execution
```

Separate requirement, design, or evidence files are created only if the information becomes independently reusable, separately reviewable, or too large to remain locally understandable.

### Complex feature

The normal physical layout is a candidate, not yet a schema commitment:

```text
proposals/<feature>/
  requirements.md
  design.md
  design/<module>.md       # only when progressive disclosure earns its cost
  spec.md                  # optional navigation/review baseline
  issues/<ticket>.md
  evidence/<ticket>.md
  research/<question>.md   # only when needed
```

Project vocabulary and ADRs remain outside the feature directory because they have broader authority.

## Handoff contracts

No hand-off depends on a manually advanced phase state.

| From | To | Required facts at the interface | Derived warning/blocking behavior |
|---|---|---|---|
| Route | Align | There is an unsettled requirement decision | No artifact gate |
| Route | Direct compact ticket | Outcome is observable, scope is local, no design-changing question is visible | Warn when uncertainty or size evidence is weak |
| Align | Solution design | Requirement unknowns that change the solution space are resolved or explicitly waived | Explicit unresolved requirement question blocks normal design finalization, not all work |
| Solution design | Plan | Interfaces, seams, ownership, migration, and verification choices needed for slicing are settled | Explicit open design question prevents affected tickets from normal frontier; override remains visible |
| Plan | Implement | DAG edges resolve and ticket has an effective specification | Content gaps warn by default; graph errors block normal frontier |
| Implement | Solution design | Execution exposes a new architecture choice | Ticket returns a design gap without pretending to complete |
| Implement | Review | A concrete change set and effective specification exist | Missing fixed point or spec source is explained, not hidden |
| Review | Completion | Findings are resolved, waived with reason, or converted to follow-up work; verification is available | Evidence records the disposition |

## Ownership conflicts found in current v5

1. `align` interviews but does not currently own a feature requirement artifact; `to-spec` becomes the first persistence point too late.
2. No current Skill owns solution-design artifacts. `codebase-design` is correctly a vocabulary layer, while `to-spec` mixes synthesis with design-shaped content.
3. `diagnose`, `tdd`, and `implement` overlap unless diagnosis and TDD are treated as disciplines inside one implementation owner.
4. `review` reports findings, but no current responsibility explicitly writes completion evidence.
5. `plan` labels tickets `ready-for-agent` even though readiness should be derived from DAG facts, explicit unknowns, and diagnostics.
6. Current `to-spec` and `plan` templates encourage either repetition or missing physical execution context.

## Decisions recorded from the discussion

- A distinct solution-design responsibility is accepted.
- It is separated from requirement alignment and execution planning.
- Vocabulary Skills remain supporting disciplines rather than serial phases.
- Readiness remains derived, explainable, and overridable.
- Human semantic prose remains authoritative; machine identity is supporting metadata.
- Logical artifact roles adapt to feature complexity instead of forcing the same number of files for every feature.

## Next decisions

1. Minimal physical schemas after the artifact-interface walkthrough.
2. Deterministic routing signals for bug, small feature, complex feature, and fog-of-war paths.
3. Production Skill instructions and the final name or re-scope of optional synthesis.
4. Parser and CLI diagnostics, override, waiver, and compatibility implementation.

These decisions should be resolved before changing templates, parser behavior, or CLI commands.
