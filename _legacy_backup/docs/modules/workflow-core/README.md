---
doc_type: "module"
title: "workflow-core"
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
related_docs:
  - "default:proposals/CR26052001"
archive_target: "default:modules/workflow-core"
created: "2026-05-22T08:17:52.067Z"
updated: "2026-05-22T08:17:52.067Z"
module_name: "workflow-core"
module_status: "active"
---

# workflow-core

## Ownership summary

- Primary module: modules/workflow-core
- System / architecture targets: none
- Convention targets: none
- Canonical reading path: workflow-core/README.md

## Purpose

`workflow-core` is the shared runtime layer that powers document workspace resolution, proposal creation and validation, and archive target updates.

## Key behaviors

- load project configuration from `.flowforge/config.json` when present
- fall back to a default single workspace for simple projects
- resolve workspaces from `cwd`, metadata, and explicit command flags
- create proposal skeletons with canonical corpus metadata
- validate proposal schema and workspace rules
- update module, architecture, and decision targets during archive
- keep installed payload upgrades aligned with `config.json`, `state/`, and platform command wrappers

## Important links

- [Design](./design.md)
- [API](./api.md)
- [History](./history.md)

<!-- flowforge:proposal:CR26052101 -->
## Archived proposals

- CR26052101: 已安装 FlowForge 安全升级策略
