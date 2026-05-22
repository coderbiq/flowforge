---
doc_type: "module"
title: "workflow-core Design"
status: "active"
workspace: "default"
module_scope: []
system_scope: []
convention_scope: []
ownership:
  - type: "module"
    target: "modules/workflow-core"
    role: "primary"
information_class: "module"
topics: []
related_docs: []
archive_target: "default:modules/workflow-core/design.md"
created: "2026-05-22T08:17:52.067Z"
updated: "2026-05-22T08:17:52.067Z"
module_name: "workflow-core"
module_status: "active"
---

# workflow-core Design

## Ownership summary

- Primary module: modules/workflow-core
- System / architecture targets: none
- Convention targets: none
- Canonical reading path: modules/workflow-core/design.md

## Current shape

`workflow-core` lives in `scripts/lib/flowforge.js` and is shared by the command wrappers under `scripts/`.

The runtime currently handles:

- project root detection
- tool root detection
- configuration loading and defaulting
- workspace enumeration and lookup
- proposal skeleton creation
- proposal validation
- archive target rendering

## Dependencies

- `workflow/guides/`
- `workflow/schema/`
- `workflow/templates/docs/`
- `docs/` as the canonical corpus

## Invariants

- `docs.default_workspace` must always resolve to a declared workspace, or be synthesized as `default`
- proposal metadata uses relative refs, not repo-absolute paths
- canonical corpus entries must refer to real documents in the workspace
- archive updates should preserve historical facts when replacing existing content
- workspace resolution should prefer explicit input over inferred scope

