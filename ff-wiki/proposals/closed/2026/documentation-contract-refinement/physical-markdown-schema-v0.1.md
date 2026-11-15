# Minimal Physical Markdown Schema v0.1

**Date:** 2026-08-25  
**Status:** Implemented by tickets 01–04; end-to-end verified
**Scope:** Physical representation of FlowForge artifact roles and machine traceability

## Purpose

Add the minimum machine-readable structure needed to distinguish artifact roles, resolve semantic dependencies, compare authority revisions, and keep non-ticket documents out of the executable frontier—without moving human workflow meaning out of Markdown prose or breaking existing v5 tickets.

The schema is intentionally smaller than the full information contract. Semantic completeness remains the responsibility of Skills and review; this schema exposes only deterministic facts the CLI can safely process.

## Design principles

- Directory placement provides the executable safety boundary.
- A small YAML envelope provides artifact identity and traceability.
- Human-visible status and blocker fields remain authoritative for execution.
- Metadata MUST NOT duplicate prose authority without a deterministic need.
- Missing new metadata degrades to legacy behavior with diagnostics.
- Schema version is not a workflow state or CLI version.
- Compact tickets remain one physical artifact even when they contain several logical roles.

## Executable discovery

### Safety rule

Only Markdown files under a feature's `issues/` directory are candidates for executable discovery.

The discovery layer MUST NOT add a file merely because it is named:

- `spec.md`;
- `requirements.md`;
- `design.md`;
- `map.md`;
- or any other conventional artifact name.

Requirement, design, overview/spec, research, prototype, map, and evidence artifacts never enter the frontier unless they are incorrectly placed under `issues/`; incorrect placement is then handled safely by role validation.

### New ticket

A new ticket under `issues/` declares:

```yaml
---
flowforge:
  schema: 1
  role: ticket
---
```

### Legacy ticket

An existing `issues/*.md` file with no FlowForge frontmatter is inferred as:

```text
role: ticket
schema: legacy
```

It continues to participate in DAG calculation and receives a compatibility diagnostic rather than a default-frontier failure.

### Role and location conflict

- A non-`issues/` file declaring `role: ticket` does not become executable and receives a placement warning.
- An `issues/*.md` file explicitly declaring a non-ticket role is excluded from executable discovery and receives a placement warning.
- Strict policy may make either warning produce a non-zero result; its diagnostic severity remains `warning`.
- The CLI MUST NOT move or rewrite the file automatically.

Safety wins over permissive role declarations.

## YAML envelope

All new FlowForge machine metadata is nested under one top-level key:

```yaml
---
flowforge:
  schema: 1
  role: design
---
```

### Required fields for a new artifact

#### `schema`

The metadata structure version.

- Version `1` is this schema.
- It does not represent the FlowForge CLI version.
- It does not represent feature, design, or ticket readiness.
- A missing schema uses legacy interpretation.
- A higher unsupported schema produces a warning and safe degradation by default; strict policy makes that warning produce a non-zero result.

#### `role`

The physical artifact's primary role.

Initial values:

- `requirement`
- `design`
- `spec`
- `ticket`
- `evidence`
- `research`
- `map`

Prototype and handoff artifacts remain outside normal proposal parsing unless a later design gives them a deterministic repository role.

### Optional common fields

#### `id`

A stable semantic identity for an independently consumed authority artifact.

Do not assign an ID merely because a file exists. Compact tickets use their filename identity and do not repeat it here.

#### `revision`

A positive integer representing the semantic revision of the identified authority.

- Increase only when meaning changes.
- Do not increase for spelling, formatting, or link repair.
- Revision is not a readiness or lifecycle state.

#### `areas`

Independently consumed authority areas within one artifact. Use only when different consumers need different revisions.

#### `open_items`

An authority-owning Skill records an explicitly known unresolved fact when it affects downstream work:

```yaml
open_items:
  - id: route-verification-seam
    diagnostic: design-verification-unselected
    severity: gap
    affects: [rpc-route-execution]
    anchor: route-verification-seam
```

- `id` is stable within the artifact and supports exact waiver targeting.
- `diagnostic` is a stable diagnostic code; prose at `anchor` explains the missing fact in human terms.
- `severity` is `gap` when informed execution is possible only by explicit override, or `blocker` when execution cannot safely proceed. Ordinary risks remain derived warnings and are not persisted here.
- `affects` names only authority areas or tickets actually affected. An empty or omitted scope is invalid; one open item MUST NOT create a feature-wide gate implicitly.
- Resolving the fact removes the item and updates the human authority. This records a fact, not a readiness transition.

#### `waivers`

A persistent waiver records an intentional policy exception without deleting or downgrading its diagnostic:

```yaml
waivers:
  - diagnostic: upstream-changed
    target: rpc-route-execution
    reason: "Reviewed revision 4; this ticket uses only the unchanged request shape."
```

- `diagnostic`, `target`, and non-empty `reason` are required.
- A waiver matches one diagnostic code and one artifact, area, open-item, or ticket target. Blanket targets and `skip_checks` are invalid.
- The Catalog still returns the original diagnostic plus waiver metadata; command policy decides whether the matched diagnostic affects filtering or exit status.
- A waiver becomes stale when its target no longer resolves or the diagnostic no longer exists. Stale and malformed waivers produce warnings and never suppress another diagnostic.
- Command overrides are ephemeral and do not create or mutate waivers.

## Human-visible execution metadata

Ticket lifecycle and DAG edges remain visible in the Markdown body:

```markdown
**Status:** open

**Blocked by:** None
```

or:

```markdown
**Blocked by:** 01, 03
```

### Authority

- `**Status:**` remains the ticket lifecycle source parsed by existing compatibility logic.
- `**Blocked by:**` remains the DAG-edge source.
- YAML MUST NOT duplicate status or blocking edges in schema v1.
- Explicit `None` is encouraged for a ticket with no blockers because it conveys immediate executability to a human.

### Parsing constraint

The blocker parser reads only identifiers from the blocker field. Additional explanatory prose MUST NOT be interpreted as part of an identifier.

If explanation is valuable, place it in following prose rather than inside the machine-parsed identifier list.

## Ticket identity and work kind

### DAG identity

Ticket identity continues to come from the filename:

```text
issues/03-migrate-rpc-route-execution.md
       └─ 03 is the feature-local DAG identity
```

The title may repeat the ID for human readability, but YAML does not declare another ticket ID.

### Artifact role versus work kind

`flowforge.role` answers:

> What kind of repository artifact is this?

Legacy `**Type:** bug`, `task`, or similar values answer:

> What kind of work does this ticket represent?

The parser model should evolve from one overloaded `Type` field toward:

```text
ArtifactRole
WorkKind
```

Work kind never decides whether a file enters executable discovery.

## Requirement authority

### Whole-artifact revision

Use one authority identity and revision until independent consumption proves a need for finer granularity:

```yaml
---
flowforge:
  schema: 1
  role: requirement
  id: backend-refinement-requirements
  revision: 2
---
```

Do not pre-assign identities to every scenario or acceptance condition.

### Area promotion

Promote independently consumed requirement areas only after unrelated semantic changes create real stale-warning noise or separate lifecycles emerge. Requirement areas then use the same shape as design areas.

## Design authority areas

One design artifact may own several independently consumed areas:

```yaml
---
flowforge:
  schema: 1
  role: design
  id: backend-refinement-design
  areas:
    rpc-route-execution:
      revision: 3
      anchor: rpc-route-execution
    configured-web-construction:
      revision: 2
      anchor: configured-web-construction
---
```

Corresponding human content uses stable invisible anchors:

```markdown
<a id="rpc-route-execution"></a>
## RPC route execution
```

Rules:

- Area identity is feature-local and stable across title changes.
- `anchor` must identify an explicit anchor in the same file.
- Ordinary headings and paragraphs need no identity.
- One area changing increments only that area's revision.
- A whole-file revision is unnecessary when all independently consumed semantics live in areas.
- The Skill that changes an area owns revision judgment and anchor maintenance; CLI checks structure and mismatch, not semantic meaning.

## Consumer dependencies

A ticket declares the authority revisions it has reviewed:

```yaml
---
flowforge:
  schema: 1
  role: ticket
  consumes:
    requirements:
      backend-refinement-requirements: 2
    design:
      rpc-route-execution: 3
---
```

### Identity scope

- Requirement, design, and area IDs are unique within a feature.
- Same-feature references use the bare semantic ID.
- Cross-feature references use `feature/id`.
- Project invariants and ADRs retain their existing path or ADR identity and are outside feature namespace rules.

### Revision comparison

When the current authority revision exceeds the consumed revision:

```text
warning: upstream-changed
authority: rpc-route-execution
consumed: 3
current: 4
```

- The ticket remains default-frontier eligible unless another gap or blocker exists.
- Unrelated area revisions do not warn the ticket.
- After rereading the changed authority, plan or implementation updates the consumed revision.
- CLI does not rewrite ticket content.

### Missing or future revision

- Missing authority: gap or dangling semantic dependency, scoped to the consumer.
- Consumed revision below current: `upstream-changed` warning.
- Consumed revision equal to current: satisfied.
- Consumed revision greater than current: invalid dependency warning; strict policy makes the warning produce a non-zero result.

## Human semantic links

Machine dependencies do not replace human links.

A ticket that consumes a design area includes a meaningful local summary and link:

```markdown
Application crosses one framework-neutral route-execution seam; Spring owns
discovery and runtime invocation.

See [RPC route execution](../design.md#rpc-route-execution).
```

Diagnostics:

- Machine dependency without a corresponding human semantic link: warning.
- Human link to an identified authority area without a consumed revision: `untracked-upstream` warning.
- Strict policy may make either warning produce a non-zero result; its diagnostic severity remains `warning`.
- CLI does not infer all semantic dependencies from arbitrary Markdown links.
- CLI does not automatically insert links or metadata.

## Artifact placement

Recommended feature layout:

```text
proposals/<feature>/
  requirements.md
  design.md
  design/*.md
  spec.md
  issues/*.md
  evidence/*.md
  research/*.md
  map.md
```

| Role | Recommended location | Executable |
|---|---|---:|
| requirement | `<feature>/requirements.md` | No |
| design | `<feature>/design.md`, `<feature>/design/*.md` | No |
| spec | `<feature>/spec.md` | No |
| ticket | `<feature>/issues/*.md` | Yes |
| evidence | `<feature>/evidence/*.md` | No |
| research | `<feature>/research/*.md` | No |
| map | `<feature>/map.md` | No |

Non-standard placement is allowed when safe, but receives a placement diagnostic. Only `issues/*.md` can be executable.

### Compact ticket

A compact ticket may contain requirement, design-increment, execution, and evidence sections. Its physical role remains `ticket`; nested logical roles do not receive separate role declarations.

## Status compatibility

The parser continues to recognize existing statuses:

- `open`
- `needs-triage`
- `needs-info`
- `ready-for-agent`
- `ready-for-human`
- `claimed`
- `resolved`
- `closed`
- `wontfix`

New authoring behavior:

- `flowforge-plan` writes `open` by default.
- `flowforge-implement` writes `closed` after completion evidence exists.
- `ready-for-agent` remains an executable legacy alias but new Skills do not emit it.
- Triage-only states remain on external-request on-ramps and are not generated by planning.
- `claimed` remains available for collaboration occupancy pending separate lifecycle evaluation.

This schema does not add a requirement, design, or execution readiness state.

## Parse and compatibility behavior

### Production parsing seam

The approved production direction is:

```text
proposal Markdown files
  -> Artifact Catalog + Diagnostics
  -> Executable Ticket Projection
  -> Ticket DAG
```

The catalog discovers requirement, design, spec, ticket, evidence, research, and map metadata for cross-artifact validation. DAG construction consumes only its executable ticket projection.

Production parsing first separates YAML frontmatter from the Markdown body, then applies legacy human-field parsing to `**Status:**`, `**Type:**`, and `**Blocked by:**` in the body.

The model distinguishes `ArtifactRole`, `WorkKind`, `TicketStatus`, `Diagnostic`, `AuthorityArea`, and `ConsumedAuthority`; the current overloaded issue type is not extended to carry all meanings.

Discovery diagnostics propagate to `check` and `frontier`. Safe degradation MUST NOT become silent degradation, including in automation-oriented output.

### Missing FlowForge frontmatter

- Legacy `issues/*.md`: parse as a ticket and emit a compatibility warning.
- Other files: do not infer executable role.

### Invalid frontmatter

- Default mode: emit a parse warning and attempt safe legacy ticket parsing only for `issues/*.md`.
- Strict mode: retain the parse warning and produce a non-zero result.
- Never execute a non-`issues/` file due to failed metadata fallback.

### Unsupported schema

- Default mode: warn and use only safely understood fields.
- Strict mode: retain the unsupported-schema warning and produce a non-zero result.
- Unknown optional fields are ignored unless they alter a known safety decision.

### No automatic migration

- Existing repositories continue working.
- Skills add schema v1 incrementally when they materially modify an artifact.
- File path is temporary identity when an old authority lacks semantic identity.
- CLI reports diagnostics but does not rewrite files.

## Deterministic diagnostics in schema scope

The physical schema supports deterministic checks for:

- artifact role and placement mismatch;
- non-ticket exclusion from executable discovery;
- invalid or unsupported frontmatter;
- missing authority identity;
- missing area anchor;
- duplicate feature-local semantic identity;
- missing consumed authority;
- revision mismatch;
- consumed revision newer than authority;
- tracked machine dependency without a corresponding marked semantic link;
- marked semantic link without a consumed revision;
- legacy compatibility.
- scoped open-item projection;
- exact persistent-waiver matching and stale-waiver detection.

It does not deterministically judge:

- whether prose has high information value;
- whether a seam is well designed;
- whether a revision should have incremented;
- whether an artifact deserves promotion into a separate file;
- whether all implicit semantic links have been captured.

## Minimal examples

### Legacy ticket

```markdown
# 01: Existing ticket

**Status:** ready-for-agent
**Blocked by:** None
```

It remains executable with a compatibility warning.

### Schema v1 ticket

```markdown
---
flowforge:
  schema: 1
  role: ticket
---

# 01: New ticket

**Status:** open
**Blocked by:** None
```

### Schema v1 design consumer

```yaml
---
flowforge:
  schema: 1
  role: ticket
  consumes:
    requirements:
      backend-refinement-requirements: 2
    design:
      rpc-route-execution: 3
---
```

Human semantic links remain in the body.

## Verification scenarios

1. `spec.md` never appears in issue count or frontier.
2. A legacy `issues/*.md` ticket remains executable with a compatibility warning.
3. A schema v1 ticket is discovered only under `issues/`.
4. An `issues/` file declaring `role: evidence` is excluded and warned.
5. A non-`issues/` file declaring `role: ticket` remains non-executable and warned.
6. Human `Blocked by: None` and machine discovery agree.
7. Existing status aliases retain current behavior.
8. One design area revision changes and warns only its consumers.
9. A formatting-only edit changes no revision and emits no stale warning.
10. Missing anchor, duplicate ID, and future consumed revision produce deterministic diagnostics.
11. Human semantic links and machine consumption mismatches warn without rewriting files.
12. Invalid schema v1 metadata degrades safely in default mode and fails strict mode.
13. Compact tickets remain one physical ticket role.
14. Requirement, design, spec, evidence, research, and map files never enter the DAG.
15. A scoped gap affects only named work and remains visible under command override.
16. A persistent waiver matches one diagnostic target, retains the original diagnostic, and records a reason.
17. Blanket, malformed, and stale waivers do not suppress diagnostics.

## Out of scope

- Final CLI command and flag names.
- Diagnostic rendering and strict-policy configuration.
- Parser implementation in this design phase.
- Skill instruction changes.
- Automatic Markdown mutation.
- Content-quality scoring.
- Mandatory repository-wide migration.

## Next work

1. Design the Artifact Catalog interface and diagnostic projections into `check` and `frontier`.
2. Validate that interface against real legacy proposal fixtures.
3. Design coordinated production changes for align, solution-design, plan, implement, and review.
4. Only then create tracer-bullet implementation tickets.

The approved Catalog and command-projection design is defined in [Artifact Catalog Interface v0.1](artifact-catalog-interface-v0.1.md).

The parser feasibility run is captured in [Schema v1 Parser Prototype Verdict](parser-prototype-verdict-v0.1.md). Its production-seam recommendations are incorporated above.
