---
doc_type: model
title: <ModelName>
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
proposal_id: <CRYYMMDDNN id>
model_name: <ModelName>
model_role: <core | lifecycle | view-facing | shared>
data_scope: <single-record | master-table | event | derived>
model_status_in_proposal: <new | modified | retained>
---

# <ModelName>

## Ownership summary

- Primary module: <type:target or none>
- System / architecture targets: <type:target or none>
- Convention targets: <type:target or none>
- Canonical reading path: this model document

This is the default business-model template.

Use this file when the model can stay in the standard shape. If a project needs extra data-structure columns, like `Master table`, or wants to reorganize several sections at once, copy this template or the relevant part files into the workspace-local template area and edit the copies directly.

## Reading order

This model template is split into parts so an agent can understand and reuse each section separately.

- [Header](./parts/header.md)
- [Purpose](./parts/purpose.md)
- [Data structure](./parts/data-structure.md)
- [Responsibilities](./parts/responsibilities.md)
- [Lifecycle](./parts/lifecycle.md)
- [Validation rules](./parts/validation.md)
- [References](./parts/references.md)
- [Open questions](./parts/open-questions.md)

## Default authoring rule

When a project uses this template as-is, preserve the section order above. When a project copies and edits the template, keep the README or the section headers updated so the agent can still understand the customization.
