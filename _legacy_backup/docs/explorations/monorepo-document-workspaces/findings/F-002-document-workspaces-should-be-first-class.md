---
doc_type: "finding"
title: "F-002 文档工作区应当是一等概念"
status: "validated"
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
related_docs:
  - "default:explorations/monorepo-document-workspaces/index.md"
archive_target: "default:architecture/monorepo-document-workspaces.md"
created: "2026-05-22T08:17:52.067Z"
updated: "2026-05-22T08:17:52.067Z"
exploration_slug: "monorepo-document-workspaces"
finding_id: "F-002-document-workspaces-should-be-first-class"
evidence_sources: []
---

# F-002 文档工作区应当是一等概念

## Ownership summary

- Primary module: modules/workflow-core
- System / architecture targets: architecture/monorepo-document-workspaces.md
- Convention targets: none
- Canonical reading path: monorepo-document-workspaces/findings/F-002-document-workspaces-should-be-first-class.md

## 结论

Monorepo 支持应当通过“命名的文档工作区”来表达，每个工作区同时包含 docs root 和对应的 code scope，而不是依赖临时拼接的相对路径。

## 为什么重要

工作区抽象能让工作流稳定回答三个问题：

- 长期文档应该存放在哪里
- 这些文档对应代码库中的哪一部分
- 命令在不同目录下执行时应如何解析默认值

没有这个抽象，一旦存在多棵文档树，proposal 的放置位置、archive 解析和任务归属都会变得模糊。

## 参考

- [Monorepo Document Workspace Support](../index.md)
