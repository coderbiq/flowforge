---
doc_type: note
title: Workflow Rules
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

# Workflow Rules

This directory seeds the project-local workflow rule bundle that is installed
into `docs/flowforge/_rules/`.

## Ownership summary

- Primary module: none
- System / architecture targets: none
- Convention targets: none
- Canonical reading path: this workflow-rules bundle

Use these files as the initial project-editable defaults for workflow posture,
input-package handling, exploration analysis, proposal writing, and archive
behavior.

This bundle is meant to be loaded by the workflow skill and command surfaces
before they reason about a project task. It is not passive documentation.

## File order

1. `workflow.md`
2. `intake.md`
3. `explore.md`
4. `propose.md`
5. `archive.md`

## Guidance

- Treat the workflow core as the stable mechanism layer.
- Treat these files as the editable project seed rules that sit above the core
  for analysis and writing defaults.
- If a project needs different behavior, edit the copied project-local rules
  bundle rather than patching the core workflow guides first.
