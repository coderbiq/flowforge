---
doc_type: "decision"
title: "D-001 升级应区分受管 payload 和用户态数据"
status: "draft"
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
decision_id: "D-001-safe-upgrade-should-preserve-user-owned-data"
decision_status: "candidate"
---

# D-001 升级应区分受管 payload 和用户态数据

## Ownership summary

- Primary module: modules/workflow-core
- System / architecture targets: architecture/installed-flowforge-upgrade-policy.md
- Convention targets: none
- Canonical reading path: installed-flowforge-safe-upgrade/decisions/D-001-safe-upgrade-should-preserve-user-owned-data.md

## Decision

升级流程默认覆盖 FlowForge 的受管安装产物，但必须保留：

- `.flowforge/config.json`
- `.flowforge/state/`

如果 `.flowforge` 根目录下存在额外的本地文件，则必须先定义它们是受管文件还是用户自定义文件，再决定是否覆盖。

## Alternatives considered

- 直接全量替换 `.flowforge/`，只保留 `config.json`
- 只更新单个脚本或单个模板文件，避免覆盖整个受管目录

## Risks

- 如果 root-level 本地文件被误当成受管文件，升级可能覆盖项目自定义内容。
- 如果受管清单没有明确维护，新版本可能漏掉某些安装产物。

## Validation needed

- 确认安装脚本的受管文件集合是否足够稳定。
- 确认 `state/` 和任何平台包装文件的保留规则。
