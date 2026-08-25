# FlowForge artifact contract

Use this reference when a Skill authors or consumes proposal artifacts.

## Roles and authority

- Requirement owns the problem, observable outcomes, scope, scenarios, constraints, terms, and requirement-changing unknowns.
- Solution design owns modules, responsibilities, interfaces, seams, flows, migration, alternatives, and verification strategy.
- Ticket owns one independently verifiable execution increment and its genuine DAG edges.
- Evidence owns delivered behavior, actual verification results, review dispositions, deviations, and implementation references.

Human semantic prose is authoritative. Machine IDs and revisions support deterministic traceability; links shown to people state their meaning.

## Packaging

Keep compact work in one ticket. Promote a role to its own file when it has independent consumers, review, lifecycle, cross-ticket scope, or would obscure the local increment. Omit empty sections.

Only `issues/*.md` can be executable. New physical artifacts use schema v1 under the `flowforge` YAML envelope. Existing v5 tickets without the envelope remain compatible.

## Hand-offs

Downstream artifacts summarize only the locally needed meaning and link to upstream authority. They do not copy rationale or settled decisions. A semantic consumer records the authority revision it reviewed.

A ticket states, in this order near its title so legacy parsers can read it:

```markdown
**Blocked by:** 01, 02
**Status:** open
```

Then it provides Delivery, Design context, Touch points, ordered Changes, scoped Constraints, and paired Done and verify conditions. An action that still selects a responsibility, interface, seam, information flow, or ordering is a design gap, not implementation work.

## Diagnostics

- `warning`: work can proceed with a named risk; strict policy may filter it.
- `gap`: a required decision or feasible observation method is missing; normal frontier excludes affected work, while an explicit override may include it without hiding the fact.
- `blocker`: a DAG edge or external fact prevents work; override cannot make it executable.

Readiness is derived from current files and DAG facts. Never persist `design-ready`, `execution-ready`, or a comparable workflow phase.

An `open_items` entry names one unresolved fact, explanation anchor, severity, and affected tickets or authority areas. A persistent waiver matches one diagnostic and exact target, records a reason, and retains the original diagnostic. Blanket skips are invalid.

## Information value

Every retained statement contributes a fact, requirement, decision, constraint, action, verification method, unknown, or evidence. Delete template filler, implementation commonplaces, exploration that changed no decision, repeated upstream prose, and synonymous headings. Keep `MUST NOT` only for a concrete failure path and place the positive target behavior beside it.

