---
name: flowforge-plan
description: Convert settled requirement and solution-design authority into new independently verifiable tracer tickets with genuine DAG edges. NOT for fix changes — flowforge-review appends `Fix:` Changes to existing tickets; plan only creates new tickets for settled authority. Use when implementation increments and execution order need to be published.
---

# Plan tracer tickets

Use the shared contract's [packaging](../_shared/ARTIFACT-CONTRACT.md#packaging), [hand-offs](../_shared/ARTIFACT-CONTRACT.md#hand-offs), and [information-value test](../_shared/ARTIFACT-CONTRACT.md#information-value). Read [schema v1](../_shared/SCHEMA-V1.md) when publishing ticket metadata.

## Process

### 1. Resolve effective authority

Read the current requirement and solution-design authorities, including their scoped open items and consumed revisions. Plan only areas with settled requirement behavior, responsibility, interfaces, seams, flow/order, migration, and a feasible verification strategy.

For schema authority being consumed, run `flowforge check --dir <feature-dir> --strict` before drafting tickets. Resolve diagnostics caused by the authority publication with its owning Skill; this validates current document relationships and does not create a readiness state. After the check, Plan still presents title, Delivery, and genuine DAG edges for user approval before creating issue files.

Inspect the codebase to find exact file paths and symbols for Touch points and Changes. The implementer may be a lightweight model that cannot search the codebase; the ticket must provide coordinates. A planning action that selects or moves a responsibility, interface, seam, information flow, or ordering returns to `flowforge-solution-design` instead of becoming a ticket step.

Read the design authority's `## Standards clauses` section. Transcribe each `must`/`must not` clause into the ticket at the tier indicated by its tag (`[Constraints]` → `## Constraints`, `[Conventions]` → `### Conventions`). Do not modify clause content, re-judge tier placement, or read the standards guide. If the design authority has no Standards clauses section and the ticket has a Write set, mark `standards: pending` and return to Solution Design.

### 2. Draft tracer increments

Select the complete behavior one ticket must deliver before deciding whether to split: weigh interface, implementation, caller, failure paths, and generation relationships, and treat that whole behavior as the increment. By default one ticket carries its related Changes as one continuous group whose checks run once at the end, because every extra ticket costs a Refine, a fresh start, an independent review, and a closeout. A split is earned by an independently usable delivery, a real start dependency, risk isolation, or context capacity — and it states that reason in one line instead of inheriting a default. Use only blocking edges that genuinely prevent the increment from starting.

Assign each proof to exactly one owning ticket. Cross-ticket proof belongs to a normal `role: ticket` stage ticket that depends on the delivered tickets, and a release or use ticket needing combined proof depends on that stage ticket. A local ticket claims only behavior it proves; when existing text claims more, return to the authority that wrote it instead of deleting the command. A pure-verification ticket delivers a non-empty verification record naming the implementation fixed point, reproducible commands, results, and boundaries, together with the implementation scope it covers. See [verification ownership and reuse](../_shared/ARTIFACT-CONTRACT.md#verification-ownership-and-reuse).

A ticket the requirement, design, or ticket text gates on an existing stage's positive proof also writes that proof: the record's location, the version it must cover, and the role that confirms it before dispatch. An upstream local delivery alone puts no ticket under the gate: a stage ticket never names the proof it will produce as its own entry condition, a release or real-use ticket needing the whole cites that stage record, and a plain local dependency runs on its own proof. The gate is then checkable at launch instead of inferred from a DAG edge — see [stage proof before dispatch](../_shared/ARTIFACT-CONTRACT.md#stage-proof-before-dispatch).

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
- **Constraints:** ticket-specific invariants, plus `must`/`must not` standards clauses transcribed from the design authority's Standards clauses tagged `[Constraints]`. Include a `Write set:` line listing the only directories or files the implementer may modify.
- **Done and verify:** pair each observable condition with an exact command and the expected result (e.g. "all pass, 0 failures" or named test cases that must pass).
- **`### Verification schedule`:** every new ticket adopts this schedule and writes it inside Done and verify. By default it is one compact paragraph covering the ticket's related Changes: when it runs, the exact command the ticket already defines with its expected result, and where the result is saved. A second, earlier group needs its own stated reason — an unverified assumption the group depends on, an irreversible use step, or the rework its failure would force. A heavy check adds one line on whether it shares a stage run, and an unmeasured cost is written as unknown. Self-check the schedule you published: a new ticket missing it has silently fallen back to per-Change acceptance, which is not a saving. Only a ticket that predates this default keeps the previous per-Change execution rules. The element list and the three execution paths are shared: [verification schedule](../_shared/ARTIFACT-CONTRACT.md#verification-schedule).

**Tier 3 — agent execution detail** (after a `---` separator):

- **Execution detail:** Plan publishes only the titled skeleton — the `## Execution detail` heading plus empty `### Verified contracts`, `### Execution scenarios`, `### Expected tests`, `### Generated artifacts`, and `### Conventions` subsection headings. Fact-backed content is filled by `flowforge-refine-ticket`, which verifies repository evidence before writing. Plan must not populate these sections with guessed values.
- **Implementation note:** left empty by Plan; written by the implementer after execution.
- **Review rounds:** left empty by Plan; accumulated by the review agent after each review round.

After writing the ticket, perform a transcription self-check. For tickets with a `Write set` (code changes), the Constraints section must carry one of: `must`/`must not` standards clauses (transcription done), `standards: none found` (design authority had no applicable standards), or `standards: pending` (design authority lacks Standards clauses—return to Solution Design). Tickets without a `Write set` (pure documentation) do not require this marker.

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
**Mode:** lightweight

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

### Verified contracts

- <filled by flowforge-refine-ticket>

### Execution scenarios

- <filled by flowforge-refine-ticket>

### Expected tests

- <filled by flowforge-refine-ticket>

### Generated artifacts

- <filled by flowforge-refine-ticket>

### Conventions

- <filled by flowforge-refine-ticket>
```

</ticket-template>

Keep `Blocked by` human-visible even when metadata carries consumption. Omit empty roles; do not impose word counts or repeat upstream rationale. A small existing-seam change stays one compact ticket when splitting would add no independent delivery or real edge.

`**Mode:**` declares who chose the execution mode so the implementer never self-selects: write `**Mode:** lightweight` for mechanical tickets aimed at a lightweight executor (unchecked mechanical Changes + Write set + Execution detail); omit the line or write `**Mode:** full` when the ticket expects a capable agent that owns implement-review-close. The dispatcher prompt may also declare the mode explicitly and overrides the ticket line.

### 4. Validate the published graph

Run `flowforge check --dir <docs_dir>/proposals` and `flowforge frontier --dir <docs_dir>/proposals`. Correct cycles, dangling edges, blockers, and gaps. Explain retained warnings and any exact waiver or explicit gap override; do not persist a readiness phase or rewrite status merely to make work appear executable.

### 5. Record the journal event

After the tickets are published and validated, append this run's event to the proposal journal `<docs_dir>/proposals/<feature>/journal.md` by hand (v1 ships no `journal stamp` subcommand). When the file is absent, create it from the frontmatter skeleton:

```markdown
---
flowforge:
  schema: 1
  role: journal
  id: <feature>-journal
  events:
    - {round: 1, stage: plan, action: done, ts: 2026-10-12T09:00:00Z, by: flowforge-plan}
---
```

One event is one flow mapping in the `events` list; append it at the end and never edit or delete an earlier line.

- `round` — the journal's active round, the round the published increments belong to
- `stage` — `plan`
- `action` — `done`, or `skipped` with a mandatory `reason`
- `ts` — UTC timestamp of the append
- `by` — `flowforge-plan`
- `reason` — required when `action` is `skipped`
- `source` — stage-less round events only; `round-open` requires `user`, `review`, or `verification`
- `artifacts` — optional artifact file names this event produced
- `backfill` — optional `true`; written by `flowforge journal init`, not by this step

## Completion

Return created ticket links, the approved DAG, consumed authority revisions, diagnostics and dispositions, and the current frontier. Every ticket is independently verifiable, locally sufficient, free of hidden design choices, and small enough for one fresh context.
