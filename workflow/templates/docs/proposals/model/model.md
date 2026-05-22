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

Use this file when the model can stay in the standard shape. If a project needs
extra data-structure columns, like `Master table`, or wants to reorganize the
sections, copy this whole template into the workspace-local template area and
edit the copy directly.

## Purpose

State what this model represents and why it exists as a separate model.

## Data Structure

Use this section for the model's field table and any table-level notes.

Default columns:

| Field | Type | Physical/JSON | Nullable | Description | Convention ref |
| ----- | ---- | ------------- | -------- | ----------- | -------------- |
| id | varchar(50) | Physical | no | Primary key | conventions/id-fields.md |

Column notes:

- `Physical` means the field becomes a real column.
- `JSON` means the field is stored inside a JSON details column on the same record.
- `Convention ref` links to the convention that governs this field type or shape, when one exists.

## Responsibilities

- What this model is allowed to decide
- What this model must delegate to other models or services

## Lifecycle

- States this model can be in
- Transitions and who can trigger them
- Audit footprint

## Validation Rules

- Required combinations of fields
- Cross-field constraints
- External constraints such as uniqueness, referential integrity, or dictionary membership

## References

### Referenced by

- `<OtherModel>` - relationship description

### References

- `<OtherModel>` - relationship description

## Open Questions

- Question that should be resolved before implementation
