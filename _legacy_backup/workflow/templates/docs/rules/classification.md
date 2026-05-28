---
doc_type: "note"
title: "Workflow Rules - Classification"
status: "draft"
workspace: "default"
module_scope: []
system_scope: []
convention_scope: []
ownership: []
information_class: "note"
topics: []
related_docs: []
archive_target: "none"
created: <ISO-8601 timestamp>
updated: <ISO-8601 timestamp>
---

# Classification Configuration

## Ownership summary

- Primary module: none
- System / architecture targets: none
- Convention targets: none
- Canonical reading path: this classification rule file

## Purpose

Project seed rules only provide project-level configuration. Classification
mechanics live in the FlowForge core guides.

## Module registry

Define the canonical module registry for the workspace here. The registry is
the source of truth for module names and archive routing hints.

```yaml
module_registry:
  data-service:
    canonical_entry: modules/data-service/design.md
  workflow-core:
    canonical_entry: modules/workflow-core/README.md
```

## Optional project defaults

Projects may also record local preferences here, such as:

- module aliases
- preferred canonical entry points
- review policy notes
- any module-specific routing exceptions that the project wants to keep close
  to the registry

Keep this file configuration-oriented. Do not duplicate the core classification
mechanics here.
