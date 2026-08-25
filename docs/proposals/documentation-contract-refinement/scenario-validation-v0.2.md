# `flowforge-solution-design` Scenario Validation v0.2

**Date:** 2026-08-25  
**Status:** Prototype verdict; validates the interface direction, not production Skill instructions

## Artifacts under test

- [Simple compact-ticket fixture](scenario-fixtures/simple-frontier-warning-count.example.md)
- [Tangram-style complex design fixture](scenario-fixtures/tangram-backend-design.example.md)
- [`flowforge-solution-design` pre-validation interface baseline](history/solution-design-interface-v0.1.md)

Both fixtures are examples. They are not executable tickets and do not modify Tangram's architecture authority.

## Question A: does the simple feature skip solution design?

Scenario:

> In human-readable `frontier` output, show the number of warnings associated with each otherwise-ready ticket; JSON remains stable.

### Trigger evaluation

| Invocation signal | Present? | Reason |
|---|---:|---|
| Module responsibility changes | No | Existing diagnostics and rendering retain their owners. |
| Interface or seam introduced/moved | No | The existing text-output seam is reused. |
| Cross-module design coordination | No | The renderer consumes an existing result; no cross-module interface changes. Touching code in more than one package would not alone make this a design problem. |
| Credible solution alternatives with a meaningful trade-off | No | Placement and zero-warning rendering are compact behavior decisions. |
| Migration/compatibility order | No | JSON stability is a ticket constraint, not a migration program. |
| Blocking edges depend on design | No | One independently verifiable increment is expected. |
| Verification seam unsettled | No | Existing frontier rendering tests are the seam. |
| Returned implementation design gap | No | This is initial planning. |

### Verdict

`flowforge-solution-design` SHOULD be skipped.

The compact ticket preserves logical completeness without a separate design artifact:

- `Why` owns the observable requirement increment.
- `Design increment` identifies the existing seam and non-decision.
- `Change`, `Constraints`, and `Done and verify` form the execution contract.
- `Completion evidence` reserves the close-out authority without adding logs before execution.

### Information-value review

- No project background is repeated.
- No requirement or decision IDs interrupt human reading.
- Each normative constraint excludes a concrete regression.
- The three change steps are distinct from the observable done conditions.
- A separate requirements file, design file, spec, and evidence file would add navigation cost without independent authority.
- The sentence explaining that diagnostics are not recomputed is the only design-specific fact a lightweight model would otherwise need to rediscover.

### Risk exposed

The routing rule must distinguish **touching multiple files/packages** from **changing cross-module responsibility or interface**. A mechanical “two modules means complex” heuristic would invoke solution design too often.

## Question B: can the complex design support partial planning without a state gate?

Scenario: Tangram-style refinement of contracts, Agent domain, task lifecycle, RPC, configured web adapters, Spring/container integration, and application composition.

### Trigger evaluation

Every high-value signal is present:

- module responsibilities change;
- interfaces and seams move;
- multiple modules coordinate;
- credible alternatives exist;
- migration order affects buildability;
- blocking edges depend on interface choices;
- verification crosses several public seams.

`flowforge-solution-design` MUST be invoked.

### Planning coverage

| Coverage fact | Result | Evidence in fixture |
|---|---|---|
| Requirements have design responses | Pass for sampled requirements | Requirement context plus four seam decisions address dependency direction, contract purity, configuration, RPC, and lifecycle behavior. |
| Module responsibilities are settled | Pass | Responsibility table names ownership and exclusions for six module roles. |
| Interfaces and seams are settled | Partial | RPC, configured web, and trigger seams are selected; Agent composition type shape remains open. |
| Ownership relationships are settled | Pass with one detail pending | Owner/provider/consumer/composition owner are explicit; Agent interface shape, not owner, is pending. |
| Relevant flows are settled | Pass | Construction, invocation, lifecycle, and migration flows are stated. |
| Migration can produce blocking edges | Pass | Expand, prove, migrate, contract, integrate order is explicit. |
| Verification seam is selected | Pass with warning | Observable seams are selected; exact Gradle commands remain unconfirmed. |
| Credible alternatives are handled | Pass | Three recurring incorrect approaches are rejected with concrete reasons. |
| No planning-changing question remains | Partial by scope | Agent-related planning is affected; RPC, trigger, and web planning are not. |

### Derived caller-visible result

```text
Design authority
- scenario-fixtures/tangram-backend-design.example.md

Decisions added
- Application owns configured-provider selection and binding.
- Web implementation owns configured-adapter construction behind a factory seam.
- RPC contracts carry framework-neutral route values, not reflection handles.
- Task domain owns Trigger lifecycle; scheduler is an adapter.

Can plan normally
- RPC route execution
- Trigger lifecycle
- Configured web construction

Cannot plan normally yet
- Agent composition
- Application Agent migration
  Reason: Agent application composition interface shape is unresolved.

Warnings
- Exact focused Gradle verification commands are not confirmed.

Recommendation
- Continue solution design for the Agent interface.
- Planning may proceed for unaffected tracer paths.
```

No `design-ready` state is needed. The same design is simultaneously usable for unaffected work and incomplete for affected work.

### Override validation

If a user explicitly includes Agent work before resolving its interface:

- the open question remains present;
- affected tickets are labelled as included by override;
- planning and implementation receive the reason;
- no metadata is changed to `resolved`;
- the override does not raise design confidence.

This satisfies the requirement that a flawed diagnostic or exceptional circumstance cannot deadlock all execution.

### Information-value review

Content retained because it changes planning or execution:

- current construction and lifecycle constraints;
- module responsibility exclusions;
- four named seams and ownership relationships;
- expand/migrate/contract shape;
- behavior-level verification strategy;
- three credible rejected alternatives;
- two explicit open questions with different effects.

Content deliberately omitted:

- the original feature's complete background narrative;
- an exhaustive file inventory;
- every class made `internal`;
- chronological exploration and terminal transcripts;
- generic statements about clean architecture;
- ticket decomposition and detailed mechanical edits.

One potential compression remains: if each seam grows enough to be independently referenced by several tickets, it should be promoted to a child design artifact. At the current fixture size, one hub remains readable and splitting would be premature.

## Interface defects found by validation

### 1. “Two or more modules” is too broad without a responsibility qualifier

The interface trigger should mean that behavior requires a new or changed collaboration across modules, not that implementation happens to edit more than one package.

Recommended wording:

> the behavior introduces or changes responsibility, information, or ordering across two or more modules

### 2. Coverage should be scoped, not feature-global

The interface implies scoped open questions, but its coverage list can still be read as one global verdict. The production Skill must report coverage per affected design area and recommend partial planning explicitly.

### 3. Exact command absence needs two classifications

An exact verification command can be:

- a warning when the verification seam and observable behavior are settled but the current environment cannot confirm invocation details;
- a design gap when no feasible way to observe the required behavior has been selected.

The production Skill must distinguish these cases.

### 4. Small-ticket design increment needs an owner

Skipping solution design works only if compact alignment/planning is explicitly allowed to record a small design increment at an existing seam. The future ticket interface must define this without turning `plan` into an architecture owner.

Recommended rule:

> Planning may preserve an already-obvious existing-seam design increment; if it must choose or move a seam, it returns to solution design.

## Overall verdict

The `flowforge-solution-design` interface direction is valid with three refinements:

1. narrow the cross-module trigger to changed collaboration or responsibility;
2. make planning coverage explicitly scoped by design area;
3. classify missing verification commands according to whether the observable seam is settled.

The adaptive packaging decision also survives validation:

- simple feature: one compact ticket is more readable and complete;
- complex feature: one design hub is justified;
- child design files remain lazy;
- open questions can delay only affected work without a persistent readiness state.

## Next experiment

After folding the three interface refinements into the baseline, design the compact-ticket and complex-requirement hand-off interfaces. Those interfaces are now the largest remaining uncertainty before production Skill instructions can be written.
