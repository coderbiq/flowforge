---
doc_type: "exploration"
title: "已安装 FlowForge 安全升级"
status: "archived"
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
related_docs: []
archive_target: "default:architecture/installed-flowforge-upgrade-policy.md"
created: "2026-05-21T12:22:04Z"
updated: "2026-05-21T12:40:09Z"
exploration_slug: "installed-flowforge-safe-upgrade"
question: "已安装到项目中的 FlowForge 应该如何安全升级到最新版本，同时保留项目级配置和运行态数据？"
reusable_rules: []
expected_size_class: medium
---

# 已安装 FlowForge 安全升级

## Ownership summary

- Primary module: modules/workflow-core
- System / architecture targets: architecture/installed-flowforge-upgrade-policy.md
- Convention targets: none
- Canonical reading path: installed-flowforge-safe-upgrade/index.md

## Context

当前安装脚本会把 FlowForge 的 workflow core、agents 和 scripts 安装到项目内的 `.flowforge/` 目录，并且只在 `config.json` 不存在时才创建默认配置。

这意味着安装后的大部分文件是 FlowForge 的受管产物，而不是项目业务资产。升级策略应该围绕这个边界设计，而不是把整棵 `.flowforge/` 目录视为不可替换的用户目录。

## Canonical corpus consulted

- [FlowForge 安装入门](../../GETTING-STARTED.md)
- [Workflow Guide](../../PROPOSAL-WORKFLOW.md)
- [Architecture](../../ARCHITECTURE.md)
- [Configuration](../../../workflow/guides/configuration.md)
- [Adapter Contract](../../../workflow/guides/adapter-contract.md)
- [Lifecycle](../../../workflow/guides/lifecycle.md)
- [安装脚本](../../../scripts/install.sh)
- [共享库](../../../scripts/lib/flowforge.js)

## Current understanding

- 安装后的 FlowForge 主要作为模板化工具层存在，项目真实知识仍落在 `docs/` 等业务文档目录中。
- `config.json` 是项目级配置边界，必须默认保留。
- `state/` 属于运行态和恢复态，升级时也应保留。
- `.flowforge/` 下是否存在额外的本地包装文件，需要明确是否属于受管 payload。
- Claude Code 和 OpenCode 的 `flowforge` commands 只是脚本入口的包装层，也应纳入升级同步范围。

## Findings

- [F-001](./findings/F-001-managed-payload-is-regenerable.md) 安装产物的核心内容来自受管模板与脚本目录，可重新生成。
- [F-002](./findings/F-002-config-json-is-user-owned.md) `config.json` 已经被实现为仅在缺失时创建的项目级配置。
- [F-003](./findings/F-003-state-is-runtime-owned.md) `.flowforge/state/` 承载运行态和恢复态，不应被升级覆盖。
- [F-004](./findings/F-004-platform-commands-track-the-script-surface.md) 平台 commands 只是 FlowForge 脚本 surface 的适配层，升级时必须保持同步。

## Candidate decisions

- [D-001](./decisions/D-001-safe-upgrade-should-preserve-user-owned-data.md) 升级应明确区分受管 payload 和用户态数据，默认只覆盖受管部分。

## Open questions

- `upgrade` 是否应该成为独立命令，还是沿用 `install.sh` 的扩展模式？
- `.flowforge` 根目录下的额外本地文件是否要纳入受管清单，还是一律视为用户自定义文件？
- 升级时是否需要记录安装版本，以便后续诊断或回滚？

## Proposed next step

- 将探索结论转化为新的变更提案，明确升级边界、保留规则和命令形态。
