---
name: flowforge-plan
description: Convert settled requirement and solution-design authority into independently verifiable tracer tickets with genuine DAG edges. Use when implementation increments and execution order need to be published.
disable-model-invocation: true
---

# Plan tracer tickets

Use the shared contract's [packaging](../_shared/ARTIFACT-CONTRACT.md#packaging), [hand-offs](../_shared/ARTIFACT-CONTRACT.md#hand-offs), and [information-value test](../_shared/ARTIFACT-CONTRACT.md#information-value). Read [schema v1](../_shared/SCHEMA-V1.md) when publishing ticket metadata.

## Process

### 1. Resolve effective authority

Read the current requirement and solution-design authorities, including their scoped open items and consumed revisions. Plan only areas with settled requirement behavior, responsibility, interfaces, seams, flow/order, migration, and a feasible verification strategy.

For schema authority being consumed, run `flowforge check --dir <feature-dir> --strict` before drafting tickets. Resolve diagnostics caused by the authority publication with its owning Skill; this validates current document relationships and does not create a readiness state. After the check, Plan still presents title, Delivery, and genuine DAG edges for user approval before creating issue files.

Inspect the codebase to find exact file paths and symbols for Touch points and Changes. The implementer may be a lightweight model that cannot search the codebase; the ticket must provide coordinates. A planning action that selects or moves a responsibility, interface, seam, information flow, or ordering returns to `flowforge-solution-design` instead of becoming a ticket step.

Read the project's standards extraction guide (configured at `standards.guide`, default `agents/standards.md` under `docs_dir`). Follow the logic described in that guide to determine which project standards apply to this ticket, based on its Touch points and Write set. For each applicable standard, extract a `must` or `must not` statement that captures what the implementer must do or avoid, with a semantic link to its source. Do not copy whole sections—restate only the locally needed meaning. Plan does not judge standard correctness; it only follows the extraction guide.

### 2. Draft tracer increments

Each ticket delivers one narrow, end-to-end behavior that can be verified independently in one fresh context. Use only blocking edges that genuinely prevent the increment from starting.

For a wide mechanical refactor that cannot land vertically, use expand–migrate–contract: add the compatible form, migrate green batches, then remove the old form after every batch. Keep independent migrations parallel in the DAG.

Present each proposed title, Delivery, and Blocked by relationship for approval. Revise the graph until those three facts are accepted.

### 3. Publish execution contracts

Write one schema v1 `role: ticket` file per increment under `<docs_dir>/proposals/<feature>/issues/`. Record reviewed authority revisions in `consumes`; make human-readable semantic links the reading interface. Immediately after the title, retain legacy-compatible fields in this order:

```markdown
**Blocked by:** None
**Status:** open
```

Then write three information tiers. Include file paths and symbols as coordinates for the implementer; do not include code snippets—the implementer writes the code from the mechanical description.

**Tier 1 — human-priority** (no file paths):

- **Delivery:** the single observable increment, stated once.
- **Design context:** a locally sufficient design summary plus semantic authority links.

**Tier 2 — shared execution contract** (human and agent both read):

- **Touch points:** specific file paths and symbols (e.g. `internal/tracker/catalog.go` — `Catalog` struct, `IdentityIndex` map).
- **Changes:** ordered, mechanical actions, each as `- [ ] N. <action naming the target file and symbol>`.
- **Constraints:** ticket-specific invariants, plus extracted `must`/`must not` standards clauses for hard invariants (violation means failure). Include a `Write set:` line listing the only directories or files the implementer may modify.
- **Done and verify:** pair each observable condition with an exact command and the expected result (e.g. "all pass, 0 failures" or named test cases that must pass).

**Tier 3 — agent execution detail** (after a `---` separator):

- **Execution detail:** subsections for `### Settled decisions`, `### Expected tests`, and `### Conventions` that the implementer needs but a human reviewer can skip. Extracted `must`/`must not` standards clauses for conventions (softer, not a direct completion gate) go here alongside non-obvious code conventions.
- **Implementation note:** left empty by Plan; written by the implementer after execution.
- **Review rounds:** left empty by Plan; accumulated by the review agent after each review round.

After writing the ticket, perform an extraction self-check. For tickets with a `Write set` (code changes), the Constraints section must carry one of: `must`/`must not` standards clauses (extraction done), `standards: none found per guide` (extraction attempted, no applicable standards), or `standards: pending` (extraction not yet done). Tickets without a `Write set` (pure documentation) do not require this marker.

<ticket-template>

```markdown
---
flowforge:
  schema: 1
  role: ticket
---

# <NN>: <Ticket title>

**Blocked by:** None
**Status:** open

## Delivery

<one observable increment>

## Design context

<locally sufficient summary>

See the design authority at `../design.md#<anchor>`.

## Touch points

- `<file path>` — <struct/function/module>
- `<file path>` — <struct/function/module>

## Changes

- [ ] 1. <mechanical action at file/symbol>
- [ ] 2. <mechanical action at file/symbol>

## Constraints

- <ticket-specific invariant>
- Write set: <allowed directories/files only>

## Done and verify

- <observable condition>: `<exact command>` — <expected result>
- <observable condition>: `<exact command>` — <expected result>

---

## Execution detail

### Settled decisions

- <design fact the implementer must know>

### Expected tests

- `<test name>` — <what it verifies>

### Conventions

- <non-obvious code convention in the touch area>
```

</ticket-template>

Keep `Blocked by` human-visible even when metadata carries consumption. Omit empty roles; do not impose word counts or repeat upstream rationale. A small existing-seam change stays one compact ticket when splitting would add no independent delivery or real edge.

### 4. Validate the published graph

Run `flowforge check --dir <docs_dir>/proposals` and `flowforge frontier --dir <docs_dir>/proposals`. Correct cycles, dangling edges, blockers, and gaps. Explain retained warnings and any exact waiver or explicit gap override; do not persist a readiness phase or rewrite status merely to make work appear executable.

## Completion

Return created ticket links, the approved DAG, consumed authority revisions, diagnostics and dispositions, and the current frontier. Every ticket is independently verifiable, locally sufficient, free of hidden design choices, and small enough for one fresh context.
