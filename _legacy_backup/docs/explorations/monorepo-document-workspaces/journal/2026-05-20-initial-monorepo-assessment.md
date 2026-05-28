---
doc_type: "journal"
title: "过程记录"
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
related_docs:
  - "default:explorations/monorepo-document-workspaces/index.md"
archive_target: "default:architecture/monorepo-document-workspaces.md"
created: "2026-05-22T08:17:52.067Z"
updated: "2026-05-22T08:17:52.067Z"
journal_date: "2026-05-20"
---

# 过程记录

## Ownership summary

- Primary module: modules/workflow-core
- System / architecture targets: architecture/monorepo-document-workspaces.md
- Convention targets: none
- Canonical reading path: monorepo-document-workspaces/journal/2026-05-20-initial-monorepo-assessment.md

## 本次变化

针对“根文档目录和子项目文档目录同时存在”的 monorepo 场景审视了当前 `tg-workflow` 模型。确认当前实现建立在单一 docs root 之上，而 monorepo 支持需要更强的抽象层。

## 证据

- `workflow/guides/configuration.md`
- `workflow/schema/proposal.schema.yaml`
- `scripts/lib/flowforge.js`
- `docs/PROPOSAL-WORKFLOW.md`

## 新问题

- 文档工作区应该如何在 config 和 metadata 中表达？
- 生命周期中的哪些操作必须显式带上 workspace？
