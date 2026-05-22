---
doc_type: note
title: Workflow Rules - Intake
status: draft
workspace: default
module_scope: []
system_scope: []
convention_scope: []
ownership: []
information_class: note
topics: []
related_docs: []
archive_target: none
created: <ISO-8601 timestamp>
updated: <ISO-8601 timestamp>
---

# Intake Rules

## Ownership summary

- Primary module: none
- System / architecture targets: none
- Convention targets: none
- Canonical reading path: this intake rule file

## Default intake behavior

- An input package is an optional pre-exploration entry point for durable
  initial materials.
- The input package may be a single file or a multi-file directory.
- The project may extend the required input-package shape with its own files
  and sections.
- Exploration must read and analyze the input package before generating its own
  skeleton.
- Exploration must record which input-package materials informed the analysis
  and skeleton.

## Update behavior

- Input packages may be updated while exploration is in progress.
- When the input package changes, exploration must re-read the updated
  materials before continuing.
- Updated input should be treated as new evidence, not as an automatic
  overwrite of prior conclusions.

