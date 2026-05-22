---
doc_type: "note"
title: "项目文档"
status: "draft"
workspace: "default"
module_scope: []
system_scope: []
convention_scope: []
ownership: []
information_class: "note"
topics: []
related_docs: []
archive_target: "default:README.md"
created: "2026-05-22T08:17:52.067Z"
updated: "2026-05-22T08:17:52.067Z"
---

# 项目文档

## Ownership summary

- Primary module: none
- System / architecture targets: none
- Convention targets: none
- Canonical reading path: README.md

## 结构

- `explorations/`：提案创建前的探索工件
- `proposals/`：进行中和已归档的变更提案工作集
- `modules/`：按模块组织的最终文档
- `architecture/`：跨模块和系统级设计文档
- `conventions/`：可复用的共识规范，例如"这一类问题统一用某种方案解决"
- `decisions/`：稳定架构决策和 ADR

## 分类与归档

- 每个 exploration 和 proposal 需要声明 `size_class` 和 `ownership`，参见 [`workflow/guides/sizing.md`](../workflow/guides/sizing.md) 和 [`workflow/guides/ownership.md`](../workflow/guides/ownership.md)。
- 归档目标根据 ownership 决定落在 modules、architecture、conventions、decisions 中的哪一个。

## 当前重点

- [Monorepo Document Workspaces](./architecture/monorepo-document-workspaces.md)
- [workflow-core](./modules/workflow-core/README.md)
- [Monorepo Document Workspaces ADR](./decisions/ADR-003-monorepo-document-workspaces.md)
- [CR26052001 Monorepo 文档工作区支持](./proposals/CR26052001-monorepo-document-workspace-support/proposal.md)
- [Task Splitting canonical guide](../workflow/guides/task-splitting.md)
- [Monorepo 文档工作区支持](./explorations/monorepo-document-workspaces/index.md)
- [提案归档时生成模块和架构文档的结构](./explorations/proposal-archive-document-structure/index.md)
