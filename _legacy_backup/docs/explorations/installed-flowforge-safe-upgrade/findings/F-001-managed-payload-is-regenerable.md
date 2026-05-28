---
doc_type: "finding"
title: "F-001 安装产物的核心内容可重新生成"
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
finding_id: "F-001-managed-payload-is-regenerable"
evidence_sources: []
---

# F-001 安装产物的核心内容可重新生成

## Ownership summary

- Primary module: modules/workflow-core
- System / architecture targets: architecture/installed-flowforge-upgrade-policy.md
- Convention targets: none
- Canonical reading path: installed-flowforge-safe-upgrade/findings/F-001-managed-payload-is-regenerable.md

## Statement

安装后的 FlowForge 核心内容来自受管模板和脚本目录。`install.sh` 会把 workflow core、agents 和 scripts 安装到项目内的 `.flowforge/`，因此这些内容天然适合在升级时重新覆盖。

## Why it matters

如果安装产物本身就是模板化、可再生成的，那么升级策略就可以把它们视为工具自身的发布物，而不是项目业务数据。

## References

- [Installation guide](../../../GETTING-STARTED.md)
- [Installation script](../../../scripts/install.sh)
