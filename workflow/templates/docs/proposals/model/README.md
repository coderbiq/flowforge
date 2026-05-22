---
doc_type: note
title: Models
status: draft
workspace: default
module_scope: []
system_scope: []
convention_scope: []
ownership: []
information_class: model
topics: []
related_docs: []
archive_target: none
created: <ISO-8601 timestamp>
updated: <ISO-8601 timestamp>
---

# Models

This directory contains one document per business model introduced or modified by the proposal.

## Ownership summary

- Primary module: <type:target or none>
- System / architecture targets: <type:target or none>
- Convention targets: <type:target or none>
- Canonical reading path: this model index

Use the directory when the proposal is `large`, or when a `medium` proposal introduces two or more business models. Single-model `small` proposals may describe the model inline in `design.md` instead.

The model template is intentionally split into parts under `parts/` so a project can reuse only the sections it needs or copy the whole template and adapt it. Projects that need a specialized field table, such as an extra `Master table` column, should edit the copied data-structure part or the copied full template rather than expecting automatic template merging.

See [`parts/README.md`](./parts/README.md) for the section order and customization guidance.

## Model index

### Core configuration models

- [<ModelName>](./<ModelName>.md)

### Lifecycle and governance models

- [<ModelName>](./<ModelName>.md)

### View-facing helper models

- [<ModelName>](./<ModelName>.md)

## Authoring rules

- One file per business model.
- Each file follows `workflow/templates/docs/proposals/model/model.md`.
- Cross-model relationships live in `design/model.md`, not duplicated inside each model file.
- Use the `parts/` directory as the default explanation of model sections and customization points.
