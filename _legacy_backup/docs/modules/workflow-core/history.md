---
doc_type: "module"
title: "workflow-core History"
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
archive_target: "default:modules/workflow-core/history.md"
created: "2026-05-22T08:17:52.067Z"
updated: "2026-05-22T08:17:52.067Z"
module_name: "workflow-core"
module_status: "active"
---

# workflow-core History

## Ownership summary

- Primary module: modules/workflow-core
- System / architecture targets: none
- Convention targets: none
- Canonical reading path: modules/workflow-core/history.md

## 2026-05-21

- Proposal: CR26052001
- Summary: Added workspace-aware proposal routing, canonical corpus tracking, and archive maintenance rules.
- Result: `workflow-core` now models the canonical corpus as an explicit baseline for later proposals.

## 2026-05-21

- Proposal: CR26052101
- Summary: Defined safe upgrade behavior for installed FlowForge payloads, platform command wrappers, and preserved project-owned config/state.
- Result: the workflow-core entry point now documents upgrade-safe installation as part of the canonical runtime contract.

<!-- flowforge:proposal:CR26052101 -->
## 2026-05-21

- Proposal: CR26052101
- Summary: 已安装 FlowForge 安全升级策略
- Source: ../../proposals/CR26052101-installed-flowforge-safe-upgrade
