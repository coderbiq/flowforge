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

This directory contains one document per business model introduced or modified
by the proposal.

## Ownership summary

- Primary module: <type:target or none>
- System / architecture targets: <type:target or none>
- Convention targets: <type:target or none>
- Canonical reading path: this model index

Use the directory when the proposal is `large`, or when a `medium` proposal
introduces two or more business models. Single-model `small` proposals may
describe the model inline in `design.md` instead.

The default model template is the single-file `model.md` document. Copy the
whole file into the workspace-local template area and edit the copy directly
when a project needs a specialized field table or different section ordering.

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
- Cross-model relationships live in `design/model.md`, not duplicated inside
  each model file.
- The default template is a single file; do not rely on `parts/` for the
  standard shape.
