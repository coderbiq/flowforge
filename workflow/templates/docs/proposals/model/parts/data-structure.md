---
doc_type: note
title: Data Structure
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

# Data Structure

Use this part for the model's field table and any table-level notes.

## Ownership summary

- Primary module: <type:target or none>
- System / architecture targets: <type:target or none>
- Convention targets: <type:target or none>
- Canonical reading path: this data-structure part

Default columns:

| Field | Type | Physical/JSON | Nullable | Description | Convention ref |
| ----- | ---- | ------------- | -------- | ----------- | -------------- |
| id | varchar(50) | Physical | no | Primary key | conventions/id-fields.md |

Column notes:

- `Physical` means the field becomes a real column.
- `JSON` means the field is stored inside a JSON details column on the same record.
- `Convention ref` links to the convention that governs this field type or shape, when one exists.

## When to customize

Customize this part when the project needs extra columns such as:

- `Master table`
- `Length`
- `Default`
- `Sensitive level`
- other project-specific field annotations

If the table shape changes materially, copy the whole part into the workspace-local template area and edit the copied version directly.
