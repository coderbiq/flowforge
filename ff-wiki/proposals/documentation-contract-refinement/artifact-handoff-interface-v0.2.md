# Artifact Hand-off Interface v0.2

**Date:** 2026-08-25  
**Status:** Implemented by tickets 01–10; end-to-end verified
**Scope:** Requirement → Solution design → Ticket → Evidence

## Purpose

Define the minimum information passed between the authoritative artifact roles in the FlowForge requirement-to-delivery flow.

The interface preserves complete engineering memory without requiring every feature to create the same files. Logical roles are stable; physical packaging adapts to whether information has independent authority.

No hand-off depends on a manually advanced readiness state. Each downstream owner reads current authority and derives gaps, warnings, blockers, and upstream changes.

## Core invariants

- Human-readable semantic content is authoritative.
- Machine metadata supports identity, traceability, revision comparison, and deterministic diagnostics.
- A durable fact has one authoritative owner.
- Downstream artifacts summarize only the local implication and link to upstream authority.
- Downstream artifacts MUST NOT override upstream authority.
- Missing optional structure MUST NOT block the default v5 frontier.
- Artifact roles MUST NOT be inferred solely from a file being under a proposal directory.
- Requirement, design, overview/spec, research, prototype, and evidence artifacts MUST NOT enter the executable frontier.

## Requirement interface

### Responsibility

Requirement authority defines why work exists and what observable result is required. It does not decide implementation modules, seams, migration steps, or ticket decomposition.

### Minimum semantics

Requirement authority resolves, when applicable:

1. **Problem and evidence** — the current externally meaningful problem and the facts proving it exists.
2. **Observable outcomes** — what becomes possible or different after delivery.
3. **Scope and out of scope** — the change envelope and explicit exclusions.
4. **Scenarios and acceptance** — representative behavior and how a human or test observes success.
5. **Requirement constraints** — product, policy, compatibility, or operational rules that shape valid solutions.
6. **Unknowns** — unresolved requirement questions and the work they affect.

These are semantic obligations, not mandatory headings. Empty sections, prose quotas, and exhaustive user-story lists are not required.

### Physical packaging

- A small feature may keep requirement authority in the `Why` or equivalent role of one compact ticket.
- A complex or independently reviewed feature promotes requirement authority to `requirements.md` or an equivalent dedicated file.
- Promotion preserves semantic identity and links; it is not a workflow state transition.

### Hand-off to solution design

Solution design can proceed for an area when its relevant outcome, scope, acceptance, and requirement constraints are known.

- A missing implementation choice is not a requirement gap.
- An ambiguous outcome, conflicting scope, or unobservable acceptance condition is a requirement gap.
- A requirement gap returns through alignment ownership and affects only dependent design areas.

## Solution-design interface

The current owner contract is [`flowforge-solution-design` Interface v0.2](solution-design-interface-v0.2.md).

### Requirement coverage

Design records a human-readable semantic link and a concise response for each implementation-affecting requirement:

```markdown
## Requirement coverage

- [Application does not depend on implementation adapters]
  Realized through the application composition interface.
```

The full requirement is not copied. Supporting metadata associates stable semantic identities.

### Hand-off to planning

For each independently affected design area, planning needs:

- a requirement response;
- settled module responsibilities;
- selected interfaces and seams;
- ownership relationships;
- relevant information, control, ordering, or lifecycle flow;
- migration facts sufficient to derive blocking edges;
- a feasible verification seam and strategy;
- resolved or scoped open design questions.

Planning MAY preserve an obvious design increment at an existing seam. If it must select, introduce, or move a seam—or change responsibility, information flow, or ordering across modules—it returns the question to solution design.

## Ticket interface

### Responsibility

A ticket defines one independently verifiable execution increment and its DAG edges. It does not restate the feature, choose architecture, or store raw execution history.

### Minimum semantics

#### Delivery

The one observable increment this ticket makes available. Do not duplicate it as Goal, Objective, What to build, Summary, and Acceptance.

#### Design context

One locally sufficient semantic summary plus human-readable links to the relevant design authority.

Example:

```markdown
Configured-provider selection belongs to application composition; concrete
construction remains inside the web adapter.

See [Configured web adapter construction](../design.md#configured-web-adapter-construction).
```

#### Blocked by

Only the tickets or explicit external/design facts that genuinely prevent this increment from starting. Natural-language explanation may accompany machine-readable edges but must not be parsed as a ticket identifier.

Blocking information MUST remain visible to a human reader even when machine metadata carries the edges. Explicit `None` has information value here because it proves the ticket can start immediately.

#### Touch points

The existing seams, modules, stable paths, or symbols needed to locate this increment. Avoid exhaustive affected-file inventories.

#### Changes

Ordered, already-decided actions. If an action contains an architecture choice, the ticket has a design gap.

#### Constraints

Only project/design invariants especially easy to violate in this increment and constraints unique to the ticket. Use semantic links for upstream authority.

#### Done and verify

Observable completion conditions paired with exact commands or another feasible observation method. Do not repeat the same sentence as separate acceptance and verification sections.

If the observable seam is settled but the current environment cannot confirm exact invocation details, emit a warning. If no feasible observation method has been selected, emit a gap.

### Compact-ticket shape

A small feature can package requirement, design increment, execution, and completion roles in one readable artifact:

```text
Why / observable behavior
Design increment
Blocked by
Touch points
Changes
Constraints
Done and verify
Completion evidence (written after execution)
```

The design increment MUST reuse an existing seam. Selecting or moving a seam promotes the work to solution design.

### Hand-off to implementation

Implementation derives whether:

- blocker edges are satisfied;
- the delivery is observable;
- design context resolves;
- ordered changes contain no hidden design choice;
- relevant constraints are explicit or linked;
- verification is feasible.

Content deficiencies produce scoped gaps or warnings. DAG and external facts may produce blockers. No `ticket-ready` state is persisted.

## Evidence interface

### Responsibility and owner

`flowforge-implement` owns completion evidence. `flowforge-review` supplies independent Standards and Specification findings; implementation records their disposition with actual verification results.

### Minimum semantics

- Delivered behavior.
- Commands or observation methods actually used.
- Observed results.
- Review findings and their resolution, reasoned waiver, or follow-up reference.
- Deviations from requirement, design, or ticket and how they were handled.
- Implementation reference such as commit, diff, or changed artifact.

Evidence MUST NOT contain an unfiltered terminal transcript, a copied ticket, or a future implementation plan.

### Physical packaging

Evidence remains an inline `Completion evidence` role for a small ticket unless it has independent value.

Promote it to `evidence/<ticket>.md` or a shared integration-evidence artifact when any of these apply:

- verification spans multiple commands, environments, or actors;
- evidence is independently audited or approved;
- several tickets share one integration proof;
- the ticket would become hard to read;
- evidence remains useful after ticket details become historical.

### Completion derivation

A ticket may close when delivered behavior and verification evidence exist and review findings are resolved, waived with reason, or converted into explicit follow-up work.

Checking boxes alone is insufficient evidence. Completion does not retroactively change requirement or design authority.

## Independent-authority test

Promote a logical role into its own file when at least one condition holds:

- two or more downstream artifacts reference it;
- it requires independent review or approval;
- it has a longer or different lifecycle than its current container;
- it crosses multiple tickets;
- leaving it embedded makes the local execution increment difficult to find;
- changing it affects another executor or reviewer.

Line count is a prompt for review, not a hard threshold.

## Unified diagnostic semantics

### Gap

A downstream owner lacks a required decision or feasible observation method.

- Affected work is excluded from the normal executable set.
- Unaffected work continues.
- A visible override may include the work without resolving the gap.

### Warning

Work can continue, but a specific risk, evidence limitation, stale dependency, or content weakness exists.

- Default frontier includes the work with explanation.
- Strict policy may filter it.

### Blocker

A DAG edge or external fact currently prevents work from starting, such as an unfinished dependency, cycle, unavailable required environment, or missing external approval.

- An override does not claim the blocker fact changed.
- Normal execution resumes only when the fact changes or the dependency model is explicitly corrected.

These terms describe current facts, not workflow phases.

## Human links and machine identity

### Human interface

- Link text states the meaning, never `REQ-1`, `DEC-3`, or “click here.”
- A local summary contains only what the current executor needs.
- Complete rationale, alternatives, and cross-module impact remain upstream.

### Machine interface

Each independently consumed authority area may carry a stable semantic slug and lightweight semantic revision. Do not assign identity to every heading or paragraph.

An authority file may contain several identified areas when different tickets consume them independently:

```yaml
design_areas:
  rpc-route-execution:
    revision: 3
  configured-web-construction:
    revision: 2
```

A consumer may record the revision it reviewed:

```yaml
design_dependencies:
  rpc-route-execution: 2
```

This schema is illustrative. The approved semantics are:

- identity is stable across human-title edits;
- identity is created only when an authority area has an independent consumer or lifecycle;
- revision increases only when meaning changes;
- formatting, spelling, and link repairs do not increase revision;
- revision does not represent readiness or workflow state;
- CLI does not automatically rewrite consumers.

When an authority area's current revision exceeds the revision consumed by a dependent artifact, emit `upstream-changed` as a warning and identify that consumer. An unrelated area changing in the same file MUST NOT warn consumers that do not depend on it. After review, the consumer records the revision it now understands.

## Conflict resolution

Conflicts are resolved by information responsibility, not modification time:

1. Project invariants constrain every feature.
2. Requirement authority owns external behavior and scope.
3. Solution design owns implementation responsibility and seams.
4. Ticket owns the current execution increment.
5. Evidence describes results and does not change normative authority.

Therefore:

- Ticket versus design is a design gap.
- Design versus requirement is a requirement gap.
- Feature content versus a project invariant requires changing project authority explicitly or rejecting the feature solution.
- Evidence showing a mismatch creates a correction, design return, or follow-up; it does not silently redefine expected behavior.

Newer timestamps never settle semantic conflicts automatically.

## Effective specification

An implementer or reviewer resolves the effective specification from:

- applicable project invariants;
- relevant requirement authority;
- relevant solution-design authority;
- the originating ticket;
- approved waivers or overrides, with reasons.

The effective specification is a resolved view, not another copied artifact. Review follows semantic links and machine dependencies instead of assuming one monolithic `spec.md` owns everything.

## Backward compatibility

- Existing v5 tickets continue to parse and participate in DAG calculation.
- Missing new semantics or metadata produce diagnostics, not default-frontier failure.
- A file path acts as temporary machine identity when no stable semantic identity exists.
- Relevant Skills upgrade an artifact incrementally the next time they modify it.
- No mandatory repository-wide migration is required.
- Strict mode or project policy may require the new contract for selected work.
- The parser must reliably distinguish executable tickets from requirement, design, overview/spec, research, prototype, and evidence artifacts before strict content enforcement is considered.

## Verification scenarios

1. A small feature remains one compact ticket through completion.
2. A compact ticket promotes requirement or design authority without losing semantic identity.
3. A complex feature links requirements, design areas, tickets, and evidence without copying content.
4. One open design gap excludes only affected tickets.
5. A warning remains executable in default frontier and visible under strict policy.
6. A blocker reflects a real DAG or external fact and is not erased by override.
7. A semantic design revision produces `upstream-changed` only for consumers of that authority.
8. A spelling-only edit does not change semantic revision.
9. Ticket/design conflict returns to solution design rather than being resolved by timestamp.
10. Existing v5 tickets remain executable with diagnostics.
11. Non-ticket proposal artifacts never enter frontier.
12. Review resolves an effective specification from semantic links rather than requiring one monolithic spec.

## Out of scope

- Final Markdown headings or frontmatter schema.
- Parser and CLI implementation.
- Automatic prose-quality judgment.
- Automatic downstream rewriting.
- Repository-wide migration.
- Production Skill instructions.

## Next work

The simple and complex hand-off walkthrough is complete; see [Artifact Hand-off Scenario Validation v0.3](artifact-scenario-validation-v0.3.md). It produced the refinements incorporated in this version.

The approved physical representation of these semantics is defined in [Minimal Physical Markdown Schema v0.1](physical-markdown-schema-v0.1.md).

Next:

1. Prototype parser behavior against representative legacy and schema v1 artifacts.
2. Design production Skill changes for align, solution-design, plan, implement, and review.
3. Create implementation tickets only after both designs survive their prototypes.
