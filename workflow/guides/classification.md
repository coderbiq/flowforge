# Classification

`FlowForge` classifies newly discovered exploration information in two steps:

1. classify the finding into a canonical bucket
2. derive routing metadata from the bucket and the project module registry

This guide defines the mechanism. Project-local seed rules provide the registry
and project-specific overrides.

## Canonical buckets

Use one of these buckets for each classified item:

- `module`
- `system`
- `cross-module`
- `convention`
- `decision`
- `exploration`

The bucket is the primary classification result. It is not a review state.

## Classification outputs

For each classified exploration item, the workflow may record:

- `classification_bucket`
- `module_name`
- `ownership`
- `module_scope`
- `archive_target`
- `needs_review`
- `review_status`
- `confidence`

The workflow should always emit a bucket. Weak evidence falls back to
`exploration` rather than leaving the item unclassified.

## Project module registry

Project-local rules may define a module registry that maps canonical module
names to archive routing information.

The registry is the source of truth for:

- canonical `module_name` values
- module aliases
- module archive targets
- canonical entry points for module docs

Recommended project-local shape:

```yaml
module_registry:
  data-service:
    path: modules/data-service/design.md
    canonical_entry: modules/data-service/design.md
    aliases:
      - data-service-config
  workflow-core:
    path: modules/workflow-core/README.md
    canonical_entry: modules/workflow-core/README.md
```

## Routing rules

- If the bucket is `module`, resolve `module_name` from the registry and derive
  `ownership.target`, `module_scope`, and `archive_target` from the registry
  entry.
- If the bucket is `system`, route to the corresponding architecture target.
- If the bucket is `cross-module`, route to architecture plus affected module
  history updates.
- If the bucket is `convention`, route to `docs/conventions/<topic>.md`.
- If the bucket is `decision`, route to an ADR target.
- If the evidence is weak, classify as `exploration` and keep the item in the
  exploration corpus until it is strengthened or superseded.

## Review markers

`needs_review` and `review_status` are separate from the bucket.

- `needs_review: true` means the item should be surfaced for later attention.
- `review_status: pending` means the item has not yet been reviewed.
- `review_status: reviewed` means the item has been checked.
- `review_status: waived` means the project has chosen not to review it.

These markers should trigger archive-time warnings when unreviewed items are
being promoted, but they must not block classification itself.

## Boundary

- This guide defines the core classification mechanism.
- Project seed rules define project-specific module registries and overrides.
- Human review policy lives in the workflow usage docs, not in the project
  registry itself.
