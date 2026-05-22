---
doc_type: "finding"
title: "F-002 `config.json` 已经形成项目级持久配置边界"
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
finding_id: "F-002-config-json-is-user-owned"
evidence_sources: []
---

# F-002 `config.json` 已经形成项目级持久配置边界

## Ownership summary

- Primary module: modules/workflow-core
- System / architecture targets: architecture/installed-flowforge-upgrade-policy.md
- Convention targets: none
- Canonical reading path: installed-flowforge-safe-upgrade/findings/F-002-config-json-is-user-owned.md

## Statement

`install.sh` 只有在 `.flowforge/config.json` 不存在时才创建默认配置；配置文档也把它定义为项目级配置入口。这说明 `config.json` 不是升级时应覆盖的模板文件，而是用户态数据。

## Why it matters

安全升级必须默认保留项目级配置，否则会破坏 workspace 定义、任务后端设置和记忆提供器配置。

## References

- [Configuration guide](../../../workflow/guides/configuration.md)
- [Adapter contract](../../../workflow/guides/adapter-contract.md)
- [Installation script](../../../scripts/install.sh)
