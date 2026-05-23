# Rule Loading

`FlowForge` adapters and skill prompts must load project rules in a single,
canonical order before they reason about exploration, proposal, archive, or
status work.

The canonical executable entrypoint for assembling that bundle is
`scripts/flowforge-rules-context.js`.

## Loading order

When the project contains `docs/flowforge/_rules/`, load files in this order:

1. `docs/flowforge/_rules/README.md`
2. `docs/flowforge/_rules/workflow.md`
3. `docs/flowforge/_rules/classification.md`
4. `docs/flowforge/_rules/intake.md`
5. `docs/flowforge/_rules/explore.md`
6. `docs/flowforge/_rules/propose.md`
7. `docs/flowforge/_rules/archive.md`

If a file is missing, skip it and continue with the remaining files.

## Precedence

- Core workflow guides remain authoritative for lifecycle, schema, and
  validation mechanics.
- Project rules refine working posture, analysis emphasis, and archive
  emphasis.
- Project rules must not override core lifecycle, schema, or validation
  contracts.
- Project rules provide project-specific configuration; the core guides own
  classification mechanics and routing semantics.

## Intake bridge

Exploration entrypoints should pair the project rules bundle with an intake
package context, using:

- `scripts/flowforge-rules-context.js`
- `scripts/flowforge-intake-context.js`
- `scripts/flowforge-explore-context.js`

## Adapter behavior

- Load the rule bundle before producing guidance for exploration, proposal,
  archive, or status actions.
- Keep the loading order stable so projects can reason about which defaults are
  in force.
- When describing behavior to a user, make clear whether a statement comes
  from core workflow guidance or project-local rules.
