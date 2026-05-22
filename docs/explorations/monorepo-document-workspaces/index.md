---
doc_type: "exploration"
title: "Monorepo 文档工作区支持"
status: "active"
workspace: "default"
module_scope: []
system_scope: []
convention_scope: []
ownership:
  - type: "system"
    target: "architecture/monorepo-document-workspaces.md"
    role: "primary"
  - type: "module"
    target: "modules/workflow-core"
    role: "secondary"
information_class: "exploration"
topics: []
related_docs: []
archive_target: "default:architecture/monorepo-document-workspaces.md"
created: "2026-05-20T00:00:00Z"
updated: "2026-05-20T00:00:00Z"
exploration_slug: "monorepo-document-workspaces"
question: "`tg-workflow` 应该如何同时支持单一文档根目录和 monorepo 中多个文档根目录的场景？"
reusable_rules: []
expected_size_class: medium
---

# Monorepo 文档工作区支持

## Ownership summary

- Primary module: modules/workflow-core
- System / architecture targets: architecture/monorepo-document-workspaces.md
- Convention targets: none
- Canonical reading path: monorepo-document-workspaces/index.md

## 背景

`tg-workflow` 当前默认一个项目只有一个顶层文档根目录。这个模型在简单仓库中够用，但在 monorepo 中会失效，因为根目录级架构文档和子项目本地文档都可能需要长期存在，并且都应当是一等工件。

工作流需要同时保留两种能力：

- 简单项目应当保持简单，不需要额外配置就能使用
- monorepo 应当显式建模多个文档工作区，而不是依赖路径约定和人工判断

## 当前理解

- 当前配置只暴露单个 `paths.docs_root`，所以 exploration、proposal 和 archive 全都被解析到同一个文档树。
- 要支持 monorepo，需要引入“一等的多文档工作区模型”，而不是只增加几个新的默认路径。
- 一旦工作跨越多个文档树，proposal 元数据、archive 规则、任务映射和记忆标签都需要增加 workspace 维度。
- workspace 选择规则不能依赖猜测，必须有明确的优先级和歧义处理策略。

## 关键发现

- [F-001](./findings/F-001-single-docs-root-is-insufficient.md) 当前工作流在模型层就是单文档根目录。
- [F-002](./findings/F-002-document-workspaces-should-be-first-class.md) Monorepo 支持需要把文档工作区建模为一等概念，并同时表达 docs root 与 code scope。
- [F-003](./findings/F-003-workspace-awareness-must-propagate-through-the-lifecycle.md) Workspace 感知必须贯穿 proposal、task、archive 和 memory 全生命周期。

## 候选决策

- [D-001](./decisions/D-001-introduce-document-workspaces.md) 用文档工作区模型替代单一 `docs_root`，同时保留简单默认场景。

## 未决问题

- monorepo 命令是否必须显式指定 workspace，还是只在自动推断不明确时要求显式指定？
- `v1` 与 `v2` proposal schema 的迁移边界应该如何控制？

## 下一步建议

- 基于当前提案，将配置模型、schema 与脚本解析逻辑一起升级为 workspace-aware 设计。
