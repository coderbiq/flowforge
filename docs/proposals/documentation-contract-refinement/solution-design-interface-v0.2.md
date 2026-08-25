# `flowforge-solution-design` Interface Design v0.2

**Date:** 2026-08-25  
**Status:** Approved and scenario-validated interface baseline; production Skill designed but not implemented  
**Authority:** Defines the responsibility and external interface of the proposed `flowforge-solution-design` Skill

## Purpose

Advance approved feature requirements into an authoritative solution design that `flowforge-plan` can decompose without reopening architecture decisions.

The Skill is a deep module: callers provide requirement authority and repository context; the Skill hides fact discovery, decision-frontier management, supporting-discipline orchestration, incremental document maintenance, coverage analysis, and information-value compression behind one interface.

It answers:

> How will the approved requirement be realized, and which design facts must every planner, implementer, and reviewer share?

It does not answer why the feature is needed, how work is divided into tickets, or whether a particular implementation diff is complete.

## Invocation interface

### Trigger

Invoke `flowforge-solution-design` when at least one of these conditions is present:

- module responsibility changes;
- an interface or seam is introduced, removed, or moved;
- the behavior introduces or changes responsibility, information flow, or ordering across two or more modules;
- two or more credible solution alternatives exist;
- migration or compatibility order affects whether the system remains usable;
- ticket blocking edges depend on an unsettled design choice;
- the verification seam is unsettled;
- implementation has returned a design gap that changes another caller's required knowledge.

Do not invoke it solely because a workflow stage says “design.” A local change at an existing seam can carry its design increment directly in a compact ticket.

### Required inputs

- The authoritative feature requirement content, whether it is a dedicated requirement artifact or the requirement role inside a compact feature artifact.
- Repository context sufficient to locate project vocabulary, applicable ADRs, documented constraints, and existing implementation facts.
- Any known design question, returned implementation gap, or user-supplied constraint that triggered invocation.

### Optional inputs

- Cited research notes.
- Prototype verdicts.
- Existing design artifacts being refined.
- Related tickets whose interfaces or blocking edges may change.
- Prior rejected alternatives that remain credible.

### Caller-visible result

The Skill returns:

1. semantic links to the authoritative design hub and any lazily created child design artifacts;
2. a concise list of decisions added or changed;
3. a coverage report, scoped by design area, explaining which planning-relevant facts are settled;
4. explicit open design questions and their affected scopes;
5. a recommendation to proceed to planning, continue design investigation, return a requirement gap, or proceed only by visible override;
6. an information-value review summary.

The result is not a persisted `design-ready` Boolean.

## Responsibility

`flowforge-solution-design` owns:

- current implementation facts that materially constrain the solution;
- target modules and their responsibilities;
- interfaces and seams;
- owner, provider, consumer, and composition owner relationships;
- relevant data, control, and lifecycle flows;
- selected design decisions;
- credible rejected alternatives and their concise rejection reasons;
- migration and compatibility strategy;
- feature and module design constraints;
- verification seams and strategy;
- open design questions that affect planning or implementation.

It MUST NOT own:

- the user problem, outcome, scope, or acceptance authority;
- project glossary content or qualifying ADR authority;
- ticket decomposition and blocking-edge publication;
- production implementation;
- completion evidence;
- a duplicate summary of every upstream requirement;
- manually maintained workflow readiness state.

## Artifact interface

### Design hub

Complex work has one authoritative `design.md` hub. The hub keeps the design navigable and owns cross-cutting decisions:

- concise target design summary;
- module responsibility map;
- cross-module interfaces and seams;
- key ownership and flow relationships;
- key decisions and credible rejected alternatives;
- migration and verification navigation;
- open design questions.

The hub is not required to use those labels as fixed headings. It must communicate the applicable semantics without template-filling prose.

### Child design artifacts

Create `design/<module-or-concern>.md` only when the content:

- can be reviewed independently;
- will be referenced independently by multiple tickets;
- has a distinct owner or seam;
- or makes the hub materially harder to understand.

The hub links to child authority using semantic text. It MUST NOT copy the child content.

### Adaptive packaging

For a small feature promoted into solution design, requirement and design roles may initially live in the same artifact. Split the design role into a hub only when it earns independent authority.

Physical packaging is an implementation choice; the logical design responsibility remains the same.

## Design-content contract

The Skill resolves the following semantics when they apply:

- Current facts that constrain the solution.
- Target modules and responsibility allocation.
- Interfaces and seams, including invariants, ordering, error modes, and required configuration when relevant.
- Ownership relationships.
- Relevant data, control, and lifecycle flows.
- Selected decisions and credible rejected alternatives.
- Migration and compatibility behavior.
- Normative constraints that exclude concrete failure paths.
- Verification seams and strategy.
- Open design questions and affected work.

Empty sections, minimum paragraph counts, exhaustive user-story lists, and generic `None` values are not required.

### Physical code facts

The design MAY include:

- stable paths and symbols identifying an existing seam;
- selected target interface or type shapes;
- source-to-target migration mappings;
- dependency-direction changes;
- exact existing verification commands when they are part of the design.

It MUST NOT include:

- undigested search results;
- exhaustive affected-file inventories that belong to execution;
- line-by-line mechanical instructions for all future tickets;
- transient implementation details that do not constrain another caller.

The test is whether losing the physical fact would force a planner or implementer to rediscover a settled design choice.

## Interaction model

### 1. Establish the design frontier

Read the requirement authority, project vocabulary, ADRs, existing design, and relevant code facts. Construct a decision tree and expose only questions whose factual prerequisites are known.

Facts are the agent's responsibility. User decisions are requested only when real alternatives or trade-offs remain.

### 2. Resolve facts before asking decisions

Use repository inspection directly for local facts. Invoke a supporting detour when the fact requires a bounded specialist method:

| Need | Supporting Skill | Return contract |
|---|---|---|
| Primary-source facts require substantial investigation | `flowforge-research` | Cited facts and limits return to the blocked design question |
| Vocabulary conflicts or a qualifying durable decision appears | `flowforge-domain-modeling` | Glossary or ADR authority is linked from design |
| An interface or seam needs deep-module analysis | `flowforge-codebase-design` | Vocabulary and design tests shape the decision; design remains authoritative |
| A runnable model or visual comparison is needed | `flowforge-prototype` | One named question, observation, and verdict return to design |
| User decisions form a branching frontier | `flowforge-grilling` | Settled decisions return to incremental design maintenance |

Supporting Skills are not persisted workflow phases.

### 3. Maintain design incrementally

After each settled decision:

- update the authoritative design immediately;
- link to requirement, project, research, prototype, and ADR authority;
- remove superseded statements rather than appending a contradictory history;
- preserve only credible rejected alternatives;
- keep unresolved questions explicit.

Do not write raw exploration, complete conversation history, or speculative possibilities into the authoritative design.

### 4. Handle requirement gaps

When a design decision depends on an ambiguous, conflicting, or unobservable requirement:

1. state the gap precisely;
2. gather any discoverable facts;
3. request the user's requirement decision through alignment rules;
4. update the requirement authority under alignment ownership;
5. resume the affected design frontier.

The Skill MUST NOT fill the gap with an undocumented design assumption.

This return does not require a new session or a persisted phase transition.

### 5. Analyse planning coverage

Before recommending planning, derive and report the following facts for each independently affected design area:

- each implementation-affecting requirement has a design response;
- module responsibilities are settled;
- cross-module interfaces and seams are settled;
- ownership relationships are settled;
- relevant flows are settled;
- migration strategy is sufficient to derive blocking edges;
- verification seams and strategy are selected;
- credible alternatives are selected or rejected;
- no open question remains that would change ticket interface, order, or verification.

The report identifies evidence, gaps, and the work affected by each gap. One incomplete design area MUST NOT collapse the feature into a global not-ready verdict. It does not store a readiness state.

Missing exact verification commands have two meanings:

- If the observable seam and feasible verification behavior are settled but the current environment cannot confirm invocation details, report a warning.
- If no feasible way to observe the required behavior has been selected, report a design gap affecting normal planning.

### 6. Compress for information value

Before returning:

- delete synonymous repetition;
- delete implementation commonplaces;
- delete exploration that did not change a decision;
- replace copied project constraints with semantic links;
- reduce repeated requirements to the minimum design context;
- make every normative constraint observable and scoped;
- apply the deletion test to every paragraph;
- confirm no unresolved choice is disguised as an action.

Report material compression actions without assigning a synthetic quality score.

## Unresolved design facts

Human-facing design questions use semantic titles and explanations. Stable semantic identity and affected scopes live in supporting metadata.

Conceptual example:

```yaml
design_questions:
  provider-construction-owner:
    affects:
      - application-composition
    status: open
```

The approved physical representation is the `open_items` contract in [Minimal Physical Markdown Schema v0.1](physical-markdown-schema-v0.1.md). The example below illustrates its semantic input; the physical schema remains authoritative.

Rules:

- Question status describes the fact of that question, not a feature workflow phase.
- Only affected work is excluded from the normal executable set.
- Unaffected design and tickets may advance.
- An override retains the question and warning.
- Resolving a question records the decision in human-readable design content; metadata supports traceability.

## Handoff interfaces

### From alignment

The Skill consumes requirement authority rather than an `aligned` status. If requirement-changing unknowns remain, it may design unaffected areas while reporting the affected frontier as unavailable.

### To synthesis

`flowforge-to-spec` is optional. When used, it creates a navigation and review baseline with semantic links and concise decisions. It is not a second authority and MUST NOT invent missing content.

### To planning

`flowforge-plan` consumes the actual requirement and design authorities plus the derived coverage report. It rechecks the facts it needs and does not trust a stale Boolean.

Planning MAY preserve an already-obvious design increment at an existing seam. If planning must select, introduce, or move a seam—or change responsibility, information flow, or ordering across modules—it returns that question to solution design.

User approval means “this design may be used for planning,” not “advance a stored state.”

### From implementation

Classify implementation discoveries by interface impact:

| Discovery | Handling |
|---|---|
| Factual correction that does not alter the solution | Correct authority and note it |
| Local design detail that preserves interfaces, responsibility, DAG, and acceptance | Refine design within ticket scope |
| Change to responsibility, interface/seam, dependency direction, migration order, acceptance, or blocking edges | Return to solution design |

### To review

Design artifacts participate in the effective specification resolved from the originating ticket and its semantic links. Review must not rely on finding one monolithic spec file.

## Escalating a compact feature

If implementation of a compact ticket exposes independent design complexity:

1. pause affected implementation;
2. preserve or extract the requirement authority;
3. promote the design role into an authoritative hub;
4. retain semantic links and machine identity;
5. resolve the new design frontier;
6. allow planning to recalculate ticket scope and blocking edges;
7. continue, replace, or narrow the original ticket explicitly.

This is artifact-role promotion, not a `simple -> complex` state transition.

## Normative invariants

- The Skill MUST NOT be invoked solely to satisfy a workflow phase.
- The Skill MUST NOT persist `design-ready` or another manually advanced phase state.
- The Skill MUST NOT silently modify requirement meaning.
- The Skill MUST NOT publish execution tickets.
- The Skill MUST NOT implement production behavior.
- The Skill MUST NOT treat supporting artifacts as design authority.
- The Skill MUST NOT force fixed headings, empty sections, or prose quotas.
- The Skill MUST record physical code facts when they encode a settled seam or migration decision that downstream work would otherwise have to rediscover.
- The Skill MUST keep open design questions human-readable and machine-traceable.
- The Skill MUST make coverage gaps and override risk visible.
- The Skill MUST perform an information-value review before returning.

## Failure and continuation modes

| Condition | Result |
|---|---|
| Missing discoverable fact | Investigate directly or use research; continue unaffected frontier |
| Requirement gap | Return precise gap through alignment ownership; continue unaffected design |
| Runnable design question | Prototype detour; fold verdict back into design |
| User decision unavailable | Preserve explicit open question and affected scope; continue unaffected design |
| Planning requested with open affected question | Recommend against normal planning for affected work; offer visible override |
| Existing design contradicts code fact | Surface contradiction and resolve authority; do not silently choose one |
| Design artifact becomes too large | Split independent authority lazily and keep semantic navigation in hub |

The Skill does not label itself blocked merely because part of the design frontier is unavailable.

## Verification seam for the Skill

The highest useful seam is an end-to-end design session driven from requirement authority to design artifacts and a caller-visible coverage report.

Representative verification scenarios:

1. A local feature at an existing seam is rejected as unnecessary solution-design ceremony.
2. A cross-module feature produces a design hub with responsibility, seam, migration, and verification decisions.
3. A requirement ambiguity returns to alignment without being silently assumed.
4. A missing primary-source fact triggers research and resumes the same question.
5. An interface uncertainty triggers a prototype and records only the verdict in design authority.
6. An open question blocks only affected planning and remains visible under override.
7. A large hub lazily splits one independently valuable module design without duplication.
8. Information-value review removes template repetition without deleting a required design fact.
9. Implementation returns a design-changing discovery while handling a factual correction locally.
10. Planning consumes design facts without relying on a persisted readiness Boolean.

Evaluation should inspect both artifact semantics and interaction behavior. File existence alone is insufficient.

## Out of scope for this interface

- CLI parser, diagnostic, override, or waiver implementation.
- Exact routing thresholds for all workflow paths.
- Requirement, ticket, and evidence template design.
- Migration of existing v5 proposals.
- Production code and Skill files.

## Implementation handoff

The simple-feature and Tangram-style complex walkthroughs are complete; see [Scenario Validation v0.2](scenario-validation-v0.2.md). They confirmed adaptive packaging and partial planning, and produced the refinements incorporated in this version.

The approved downstream and upstream contract is defined in [Artifact Hand-off Interface v0.2](artifact-handoff-interface-v0.2.md).

The supporting machine representation is defined in [Minimal Physical Markdown Schema v0.1](physical-markdown-schema-v0.1.md).

The scenario, hand-off, physical-schema, Catalog, and Skill-coordination reviews are complete. Implementation planning consumes these authorities and preserves their dependency order.
