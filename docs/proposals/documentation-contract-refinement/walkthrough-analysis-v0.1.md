# Workflow Walkthrough Analysis v0.1

**Date:** 2026-08-25  
**Status:** Analysis of the throwaway prototype; recommendations are not yet approved

## Question

Can FlowForge advance a simple requirement and a complex requirement from the user's first message through analysis, design, DAG execution, and verified completion while preserving complete engineering memory without recreating verbose templates or a blocking workflow state machine?

## Executive finding

The current v5 Skills contain most of the necessary reasoning disciplines, but their present hand-offs do not form a complete artifact lifecycle.

The workflow should not be one fixed Skill sequence. It needs:

- one small-feature path in which requirement, design increment, ticket, and evidence are logically complete but physically compact;
- one complex-feature path in which the same logical roles progressively expand into separate artifacts;
- separate bug and fog-of-war on-ramps rather than forcing them through the small-feature path;
- vocabulary disciplines such as domain modeling and codebase design running underneath analysis instead of appearing as mandatory phases;
- dynamic, explainable diagnostics rather than persisted readiness states.

The first prototype is useful precisely because it exposes where its initial flow is wrong.

## Walkthrough method

Each step is evaluated against five questions:

1. What unresolved question triggers this Skill?
2. What unique information does it add?
3. Which artifact owns that information?
4. What allows the workflow to advance without a persisted phase state?
5. Can the step be skipped without losing information?

A step fails the information-value test if its output only restates upstream content. A transition fails the robustness test if work cannot continue solely because a manually maintained state value is stale.

## Walkthrough A: the initial “simple requirement” prototype

The prototype currently uses “`frontier` must not return `spec.md`.” This is a bug with a tight deterministic reproduction, not a representative small feature.

### Step-by-step analysis

| Prototype step | Actual Skill contract | Finding |
|---|---|---|
| `route` identifies a local change | `route` sends broken behavior to the diagnosis on-ramp | Correct routing, but this evaluates a bug path rather than a simple requirement path. |
| `diagnose` produces a fact | `diagnose` owns the complete feedback-loop, hypothesis, regression-test, fix, and cleanup cycle | The prototype understates its responsibility. Diagnosis is not merely research before planning. |
| `plan` creates one ticket | `route` says a one-session change can move directly to `implement`; diagnosis may already perform the fix | A one-ticket plan adds little unless a durable work item or cross-session hand-off is actually needed. |
| `tdd` drives the fix | `diagnose` already requires a regression test; `implement` also invokes TDD | Three Skills claim parts of the same implementation loop. The caller needs one owner, not three serial phases. |
| `review` records evidence | Current `review` reports Standards and Spec findings but does not own durable completion evidence | The proposed evidence artifact has no current Skill owner or schema. |

### Verdict

This walkthrough should be retained as a **bug on-ramp test**, not called the simple-requirement scenario.

Recommended bug flow:

```text
raw bug report
  -> route
  -> diagnose owns reproduce / explain / regress / fix
  -> review when the change warrants independent review
  -> completion evidence
```

Planning is conditional. It is introduced only when diagnosis reveals multiple deliverables, a design gap, multiple sessions, or blocking edges. TDD is a discipline used inside diagnosis rather than a second workflow phase. This reduces duplicated instructions and artifact churn.

## Walkthrough B: a corrected simple requirement

Use a small, well-bounded feature as the representative scenario:

> In human-readable `frontier` output, show the number of warnings associated with each otherwise-ready ticket; JSON output remains stable.

This is intentionally one observable behavior with a likely existing seam and no expected domain or architecture decision.

### Proposed interaction

#### 1. Intake and routing

User input supplies the desired behavior. The agent inspects existing output and tests to establish facts.

The router asks no ceremonial questions. It classifies the work using observable properties:

- scope appears local;
- one behavior is requested;
- no unresolved product choice is visible;
- an existing output seam is likely;
- the work fits one context window.

Output: no artifact yet. Routing has no durable information of its own.

Advance condition: enough evidence exists to distinguish a small feature from a bug, design problem, or multi-session effort.

#### 2. Compact alignment

The agent confirms only decisions that materially affect behavior:

- whether warnings change text output only;
- what count is shown;
- behavior when the count is zero;
- JSON compatibility.

This is not a full mandatory grilling sequence. If the user's request and current behavior already settle these questions, the agent states the interpretation and proceeds.

Output: an approved compact requirement statement held in the conversation until publication. No `requirements.md` is created solely to satisfy a stage.

Advance condition: there is one observable outcome and no implementation-changing open question.

#### 3. Publish one compact ticket

The one ticket physically collapses the logical requirement, design increment, execution contract, and later evidence while keeping the roles visually distinct:

```text
Why / observable behavior
Design context at the existing output seam
Changes
Constraints
Done when
Verify
Completion evidence (written only after execution)
```

The ticket uses a semantic link when upstream context exists; it does not invent a requirement document or opaque ID.

Advance condition: blocker edges are valid. Content diagnostics may warn but do not create a persisted `execution-ready` state.

#### 4. Execute through one owner

`implement` owns the ticket. It uses TDD internally at the already agreed output seam and invokes review as appropriate. The user should not need to orchestrate `implement -> tdd -> implement -> review` as separate state transitions.

Output: implementation and behavior tests.

Advance condition: verification runs and review finds no unresolved specification or standards failure.

#### 5. Record completion evidence

For a small ticket, evidence can be a compact final section rather than a separate file:

- commands run;
- observed result;
- deviations or “none” where an explicit deviation check is required;
- implementation reference.

The original instruction remains readable because raw logs are not appended.

### Simple-flow verdict

Logical completeness does not require physical document multiplication. One well-shaped ticket can carry the complete lifecycle of a small feature if it keeps requirement, design increment, execution, and result distinguishable and avoids repeating them.

## Walkthrough C: complex cross-module architecture refinement

The complex scenario is representative, but the initial prototype incorrectly presents vocabulary Skills as serial phases and leaves several artifact ownership gaps.

### Step 1: route by uncertainty and size

Trigger: the request crosses multiple modules, changes responsibility or interfaces, and cannot be planned without closing design questions.

Unique output: none. Routing selects a flow; it does not own feature content.

Recommended route:

- use `align` when the decision tree fits a working session;
- use `wayfinder` when the decision tree itself spans sessions;
- feed primary-source research into either flow when facts are missing.

Finding: “complex” and “foggy” are different. A large but already understood refactor needs alignment and design; a genuinely unknowable path needs wayfinding decision tickets.

### Step 2: align owns feature requirement truth

Current behavior: `align` calls grilling and domain modeling. Domain modeling writes only glossary terms and qualifying ADRs. Nothing currently owns a feature-level requirement artifact during the interview.

Risk: requirements remain only in the conversation until `to-spec`. This contradicts the goal of surviving compaction and makes the main flow depend on one unbroken context window.

Recommendation: alignment should maintain a concise feature working document as decisions settle. It should contain only:

- problem and evidence;
- outcomes and scope;
- observable scenarios and acceptance;
- explicit unknowns;
- settled product constraints.

It must not copy glossary entries or implementation design into that file.

Advance condition: the current decision frontier contains no unresolved requirement question that would change the solution space. This is calculated from explicit unknowns, not a stored `aligned` state.

### Step 3: research is a fact-producing detour, not a phase

Current `route` correctly says research feeds alignment rather than replacing it.

Recommendation:

- trigger research when a decision depends on a fact not recoverable cheaply in the active context;
- cite its primary-source note from the feature document;
- copy only the decision-relevant conclusion, not the research narrative;
- return to the same unresolved question after the fact arrives.

Research can occur before or during alignment and design. It should not be a mandatory box in the workflow.

### Step 4: domain modeling runs underneath the flow

Current Skill correctly limits `CONTEXT.md` to glossary content and ADRs to hard-to-reverse, surprising, genuine trade-offs.

Recommendation: do not show domain modeling as a user-visible serial phase. Invoke it only when:

- a term is ambiguous or conflicts with the glossary;
- a durable domain relationship is resolved;
- a decision meets all ADR criteria.

Output ownership remains `CONTEXT.md` or an ADR. The feature requirement/design links to that authority and does not duplicate it.

### Step 5: codebase design also runs underneath solution design

Current Skill is explicitly a vocabulary and design-discipline layer, not a workflow that promises a `design.md` artifact.

Recommendation: introduce an explicit **solution-design responsibility** in the main flow. Whether this is an expanded `align`, a new Skill, or a re-scoped synthesis Skill remains open. It consumes:

- approved feature requirements;
- primary-source code facts;
- project vocabulary and ADRs;
- deep-module principles.

It produces the target modules, interfaces, seams, ownership, flows, migration strategy, constraints, alternatives, and test seams. `codebase-design` supplies the vocabulary and review criteria underneath it.

Advance condition: no open question remains that would change ticket interfaces, blocking edges, or verification strategy.

### Step 6: prototype is conditional and answers one named question

The prototype detour is correct when a runnable answer is cheaper or more reliable than discussion. It must declare:

- the exact question;
- the competing interpretations;
- the observed verdict;
- which design decision the verdict changes.

The verdict is folded into the authoritative design. The throwaway program is evidence, not a permanent layer of the feature specification.

### Step 7: `to-spec` currently conflicts with the desired information contract

Current problems:

- it waits until late synthesis to create the feature artifact;
- it combines requirement and implementation decisions in one template;
- it requires an “extremely extensive” user-story list, which encourages low-value repetition;
- it forbids specific file paths even when a settled refactor design needs physical seams;
- it asks to confirm test seams but has no explicit source artifact that owns the design.

Recommendation: do not treat current `to-spec` as a neutral compiler. Redesign it either to:

1. compact already-authoritative requirement and design artifacts into a navigation baseline without duplicating them; or
2. disappear from flows where alignment and solution design already maintain those artifacts directly.

It MUST NOT invent decisions to make a template look complete.

### Step 8: planning converts design into execution increments

Current strengths:

- tracer-bullet slices;
- explicit blocking edges;
- user review before publication;
- expand-contract exception for wide refactors.

Current gaps:

- exploration is optional even when no physical execution context exists;
- the ticket schema lacks design context, touch points, ordered changes, ticket constraints, and exact verification;
- it forbids paths by default;
- `ready-for-agent` is written as a status rather than derived confidence.

Recommendation: planning may publish a ticket only when the design decisions it consumes are settled. A complex-feature ticket should contain the minimum local execution packet and semantic links to full design. Machine traceability remains metadata.

Advance condition: blocking edges are valid. Missing content produces explicit diagnostics; only an explicit unresolved design blocker excludes the ticket from the normal executable set.

### Step 9: `check` and `frontier` stay deterministic

Recommended output groups:

```text
Ready
Ready with warnings
Blocked by DAG/design fact
Included by override
```

The CLI may deterministically validate links, roles, blocker edges, explicit unknowns, evidence presence, and waiver shape. It must not pretend to deterministically judge prose quality or architecture wisdom.

Specification, design, research, and evidence files must not be parsed as executable tickets.

### Step 10: implementation owns execution and design-gap return

Current `implement` is too small to define the hand-off contract. It says to implement, use TDD, run review, and commit, but not what happens when a supposedly ready ticket still requires a design decision.

Recommendation:

- one Skill owns the implementation turn;
- TDD and codebase-design are disciplines invoked inside that turn;
- a new architecture choice is not made silently;
- the result explicitly distinguishes implementation progress from a returned design gap;
- scope expansion creates a new ticket or reopens design discussion rather than being hidden in the diff.

No additional persisted workflow state is needed. The existence of an explicit design gap changes the derived diagnostics.

### Step 11: review and evidence close the loop

Current review has a useful two-axis split: Standards and Spec. Gaps remain:

- it relies on finding a spec and may not treat the originating ticket plus linked design as one effective execution contract;
- it requires a fixed Git point, which may be unavailable for an uncommitted per-ticket review;
- it reports findings but does not own durable completion evidence;
- it does not assess information-value drift in changed documentation.

Recommendation: review resolves an **effective spec** consisting of the ticket plus its semantic requirement/design links. Completion evidence is then written by the implementation owner from actual commands and review outcome. Complex work uses a separate evidence file; small work may use a compact ticket section.

## Corrected candidate flows

### Small feature

```text
user request
  -> route using size / uncertainty / durability facts
  -> compact alignment (only unresolved behavioral decisions)
  -> one compact ticket with logical requirement + design increment + execution contract
  -> DAG diagnostics
  -> implement owns TDD and review
  -> compact completion evidence
```

Artifacts: normally one ticket. A requirement or design file appears only if information becomes independently reusable or the work expands.

### Complex feature

```text
user request
  -> route
  -> alignment maintains feature requirements
       <-> research detours for missing facts
       <-> domain modeling for durable vocabulary / ADRs
  -> solution design maintains design artifacts
       <-> codebase-design discipline
       <-> prototype detours for runnable questions
  -> synthesis only if a navigation baseline adds value
  -> planning publishes execution packets + DAG edges
  -> check / frontier derive executable work and warnings
  -> implement each ticket in fresh context
       <-> return design gaps instead of silently deciding
  -> two-axis review
  -> evidence and closure
```

Artifacts: requirement, design, tickets, and evidence are separate when their information is independently valuable. Vocabulary, ADR, research, and prototype artifacts are created only on demand.

### Bug on-ramp

```text
bug report
  -> route
  -> diagnose owns feedback loop, cause, regression, and fix
  -> plan only if diagnosis expands into multiple deliverables
  -> review and evidence
```

### Fog-of-war on-ramp

```text
destination with no visible path
  -> wayfinder decision DAG
  -> resolved decisions feed requirements / solution design
  -> normal complex-feature planning and execution
```

## Decisions suggested by the walkthrough

These are recommendations for discussion, not accepted decisions:

1. Use **logical artifact roles with adaptive physical packaging**. Small work collapses roles; complex work separates them.
2. Route by unresolved uncertainty, expected session count, and durability of information—not by a persisted phase state.
3. Give every flow exactly one implementation owner. TDD, diagnosis techniques, and design vocabulary are nested disciplines, not redundant hand-offs.
4. Make alignment maintain requirement truth incrementally so compaction cannot erase it.
5. Add an explicit owner for solution-design artifacts; current `codebase-design` is a vocabulary layer and current `to-spec` is too late and too template-heavy.
6. Treat synthesis as optional when authoritative artifacts already exist.
7. Resolve an effective spec from the ticket and semantic links during review.
8. Let evidence scale: inline for a small ticket, separate for complex work.

## Questions the next design round must settle

1. Should solution design become a new Skill, an expanded responsibility of `align`, or a rewritten `to-spec`?
2. What exact signals classify work as small, complex, a bug, or fog-of-war without turning routing into another rigid state machine?
3. What is the minimum compact-ticket shape that remains human-readable and sufficient for a lightweight model?
4. When does information deserve a separate requirement, design, or evidence file instead of a section in one ticket?
5. Who writes completion evidence: `implement`, `review`, or a dedicated close-out responsibility inside implementation?
6. How should review identify the fixed point and effective spec for one uncommitted ticket?
7. Which explicit design-gap representation should remove a ticket from the normal frontier while still allowing a reasoned override?

## Prototype disposition

The existing HTML remains unchanged as the primary source for the first experiment. Its simple tab should be revised only after the above candidate simple-feature flow is accepted; changing it now would erase the discrepancy that produced the finding.
