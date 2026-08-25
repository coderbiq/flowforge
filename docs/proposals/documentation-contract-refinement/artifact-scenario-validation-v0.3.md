# Artifact Hand-off Scenario Validation v0.3

**Date:** 2026-08-25  
**Status:** Prototype verdict; fixtures are not production schemas or executable work

## Fixtures

### Simple feature

- [Compact frontier warning-count ticket](scenario-fixtures/v3/simple-ticket.example.md)

### Complex feature

- [Backend refinement requirements](scenario-fixtures/v3/complex-requirements.example.md)
- [Backend refinement design](scenario-fixtures/v3/complex-design.example.md)
- [RPC migration ticket](scenario-fixtures/v3/rpc-ticket.example.md)
- [RPC completion evidence](scenario-fixtures/v3/rpc-evidence.example.md)

## Question A: does every hand-off provide unique downstream value?

### Simple path

The one compact ticket supplies all required roles without external navigation:

| Role | Information passed | Unique value |
|---|---|---|
| Requirement | Warning count is visible in text; JSON and zero-warning text remain stable | Defines external behavior and compatibility |
| Design increment | Existing renderer consumes existing diagnostics | Prevents recomputation and avoids a new seam |
| Execution | Touch points, three ordered changes, three concrete constraints | Makes the increment directly implementable |
| Evidence | Placeholder owned by future implementation | Keeps close-out authority local without inventing results |

Verdict: pass. Separate requirement, design, spec, and evidence files would not have independent authority in this scenario.

### Complex path

| Hand-off | Information passed | Information deliberately not copied |
|---|---|---|
| Requirement → Design | Outcomes, scenarios, constraints, revision 2 | Full problem narrative is linked, not repeated |
| Design → RPC ticket | One seam summary, revision 3, migration actions, contract constraints | Responsibility table, web/trigger decisions, alternatives |
| Ticket → Evidence | Ticket identity and effective-spec revisions | Ticket instructions and raw logs |
| Evidence → Completion | Actual observations, review disposition, deviations, implementation reference | Future plan or normative changes |

Verdict: pass. Each artifact owns information with a distinct lifecycle and consumer.

## Question B: can planning consume the complex design without redesigning it?

### RPC area

The ticket can be derived mechanically from settled design:

- delivery follows the requirement scenario;
- the route-execution seam and responsibility are selected;
- expand/adapt/migrate/remove ordering follows migration design;
- constraints follow requirement and design authority;
- the observable verification seam is known.

The ticket does not need to choose a contract shape category, Spring ownership, migration strategy, or test seam. Exact Gradle invocation remains a warning because the behavior and feasible observation method are already selected.

Verdict: pass.

### Agent area

The design explicitly lacks the application-composition interface shape. Planning cannot write ordered, already-decided changes without selecting that seam.

No Agent execution ticket is published in the fixture. The design remains useful for RPC, configured web, and Trigger planning.

Verdict: pass. The gap is scoped and does not become a feature-global readiness state.

## Question C: is the effective specification understandable?

For the RPC ticket, a human resolves:

1. Project invariants from the host repository.
2. [Backend behavior and constraints](scenario-fixtures/v3/complex-requirements.example.md).
3. [RPC responsibility and seam](scenario-fixtures/v3/complex-design.example.md#rpc-route-execution).
4. [The RPC execution increment](scenario-fixtures/v3/rpc-ticket.example.md).
5. Any visible override or waiver; none exists in the fixture.

The ticket contains enough local meaning to begin work; links provide rationale and cross-module impact. A reader does not need to interpret `REQ-1` or `DEC-3`.

Machine metadata resolves the same chain using semantic identities and consumed revisions. The metadata is secondary; deleting it would reduce diagnostics but would not make the human specification unintelligible.

Verdict: pass.

## Question D: does semantic revision avoid a state machine?

Initial state:

```text
Requirement authority: revision 2
Design authority: revision 3
RPC ticket consumes: requirement 2, design 3
Result: no upstream-changed warning
```

Simulated semantic change:

```text
Design authority: revision 4
Change: route errors now cross the seam as a passive error value
RPC ticket consumes: design 3
Result: warning upstream-changed (consumed 3, current 4)
```

Consequences:

- RPC remains a DAG-ready ticket with a visible warning in default mode.
- Plan or implementation rereads the changed design.
- If the ticket still conforms, its consumed revision becomes 4.
- If the change alters ordered work, plan updates the ticket.
- No feature or ticket phase state changes.
- A spelling-only design edit does not increment revision and creates no warning.

Verdict: pass, with one usability concern: authors must know when a change is semantic. This is a Skill/review judgment and should not be presented as a perfectly deterministic CLI decision.

## Question E: did the artifact chain create repeated prose?

### Necessary local repetition

The requirement states the external RPC outcome. The design states the ownership/seam decision. The ticket restates one sentence of that seam because an executor needs it locally. Evidence states the observed result.

These sentences are related but perform different roles. Deleting any one loses requirement, design, execution, or completion meaning.

### Removed repetition

- The ticket does not repeat the six-module responsibility map.
- The ticket links the contract constraint instead of copying the complete requirement section.
- Evidence does not copy `Changes`, `Touch points`, or all constraints.
- Design links requirement outcomes and adds responses instead of rewriting the full scenarios.
- No artifact contains generic clean-architecture prose or the chronological exploration that produced the design.

Verdict: pass.

## Question F: are physical promotion rules sensible?

### Simple feature

None of the roles has independent authority. One ticket remains best.

### Complex feature

- Requirements affect several design areas and tickets, so a separate requirement artifact is justified.
- Design affects several tickets and outlives one execution increment, so a separate hub is justified.
- RPC evidence spans focused behavior, joint compilation, static contract inspection, and independent review, so a separate evidence artifact is plausible.
- No child design file is yet justified: the hub remains readable and the fixture has only one published execution ticket.

Verdict: pass. Splitting by authority is more reliable than splitting by line count.

## Defects and refinements found

### 1. `Blocked by` must remain visibly human-readable

The simple fixture stores an empty edge list only in illustrative metadata. A human reader cannot see that the ticket is unblocked without inspecting frontmatter.

Recommendation: ticket prose should expose blockers or explicitly state “None” where proving immediate executability has value. Unlike empty template sections, this `None` carries DAG meaning.

### 2. Stable identity should attach to authority, not every heading

One file-level design identity is insufficient for precise revision warnings when only RPC changes, but assigning IDs and revisions to every heading would burden authors and readers.

Recommendation: create semantic identity only for independently consumed design areas. The design hub may own several identified areas; ordinary paragraphs do not need IDs.

### 3. Revision granularity must match dependency granularity

The fixture uses one design revision, so a web-only semantic change would warn the RPC ticket unnecessarily.

Recommendation: consumers depend on identified design areas such as `rpc-route-execution`, each with its own semantic revision, rather than the whole design file revision when practical.

### 4. Evidence promotion is a judgment, not a deterministic threshold

The RPC evidence could reasonably remain inline during a small implementation even though it has several verification modes. The independent-authority rules guide the owner but do not produce one objectively correct file choice.

Recommendation: keep evidence packaging a Skill decision and avoid CLI enforcement beyond valid role/link parsing.

## Overall verdict

The Artifact Hand-off Interface v0.1 passes both scenarios with three design refinements:

1. `Blocked by` remains human-visible even when machine edges exist.
2. Semantic identity is assigned to independently consumed authority areas, not every section.
3. Revision and dependency granularity follow those authority areas to avoid unrelated stale warnings.

The experiment does not justify final frontmatter syntax. It does justify moving next to a minimal physical schema prototype with explicit artifact roles, human-visible blockers, area-level semantic identities, and backward-compatible parsing.

