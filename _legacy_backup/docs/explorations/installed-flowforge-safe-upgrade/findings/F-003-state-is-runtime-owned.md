---
doc_type: "finding"
title: "F-003 `.flowforge/state/` 是运行态和恢复态"
status: "validated"
workspace: "default"
module_scope: []
system_scope: []
convention_scope: []
ownership:
  - type: "system"
    target: "architecture/installed-flowforge-upgrade-policy.md"
    role: "primary"
  - type: "module"
    target: "modules/workflow-core"
    role: "secondary"
information_class: "exploration"
topics: []
related_docs:
  - "default:explorations/installed-flowforge-safe-upgrade/index.md"
archive_target: "default:architecture/installed-flowforge-upgrade-policy.md"
created: "2026-05-22T08:17:52.067Z"
updated: "2026-05-22T08:17:52.067Z"
exploration_slug: "installed-flowforge-safe-upgrade"
finding_id: "F-003-state-is-runtime-owned"
evidence_sources: []
---

# F-003 `.flowforge/state/` 是运行态和恢复态

## Ownership summary

- Primary module: modules/workflow-core
- System / architecture targets: architecture/installed-flowforge-upgrade-policy.md
- Convention targets: none
- Canonical reading path: installed-flowforge-safe-upgrade/findings/F-003-state-is-runtime-owned.md

## Statement

FlowForge 的本地状态和恢复态存放在 `.flowforge/state/`，其职责是恢复活跃会话和工作流状态，而不是承载可再安装的工具本体。

## Why it matters

升级不应覆盖运行态数据，否则会打断工作恢复和会话连续性。

## References

- [Architecture](../../../docs/ARCHITECTURE.md)
- [Getting Started](../../../docs/GETTING-STARTED.md)
