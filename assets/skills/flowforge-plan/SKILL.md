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

Inspect the codebase when stable touch points or existing seams are not already known. A planning action that selects or moves a responsibility, interface, seam, information flow, or ordering returns to `flowforge-solution-design` instead of becoming a ticket step.

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

Then write only the semantic roles that carry information:

- **Delivery:** the single observable increment, stated once.
- **Design context:** a locally sufficient design summary plus semantic authority links.
- **Touch points:** stable seams, modules, paths, or symbols needed to locate the work.
- **Changes:** ordered, already-decided actions.
- **Constraints:** ticket-specific or easy-to-violate invariants, linked upstream when shared.
- **Done and verify:** pair each observable completion condition with an exact command or feasible observation method.

Keep `Blocked by` human-visible even when metadata carries consumption. Omit empty roles; do not impose word counts or repeat upstream rationale. A small existing-seam change stays one compact ticket when splitting would add no independent delivery or real edge.

### 4. Validate the published graph

Run `flowforge check --dir <docs_dir>/proposals` and `flowforge frontier --dir <docs_dir>/proposals`. Correct cycles, dangling edges, blockers, and gaps. Explain retained warnings and any exact waiver or explicit gap override; do not persist a readiness phase or rewrite status merely to make work appear executable.

## Completion

Return created ticket links, the approved DAG, consumed authority revisions, diagnostics and dispositions, and the current frontier. Every ticket is independently verifiable, locally sufficient, free of hidden design choices, and small enough for one fresh context.
