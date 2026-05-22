---
doc_type: note
title: Model Header
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

# Model Header

Use this part for the model's metadata block.

## Ownership summary

- Primary module: <type:target or none>
- System / architecture targets: <type:target or none>
- Convention targets: <type:target or none>
- Canonical reading path: this header part

Typical fields:

- Role: core | lifecycle | view-facing
- Owning modules: <module target(s)>
- Related system / architecture targets: <system or cross-module target(s)>
- Related conventions: links to `docs/conventions/<topic>.md`
- Status in proposal: new | modified | retained

## When to customize

Customize this part when the project needs to add or remove metadata fields that are part of the model's identity, not its data table.
