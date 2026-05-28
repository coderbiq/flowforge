---
doc_type: "finding"
title: "F-004 归档更新采用追加式且带标记的幂等写入"
status: "validated"
workspace: "default"
module_scope: []
system_scope: []
convention_scope: []
ownership:
  - type: "system"
    target: "architecture/proposal-archive-document-structure.md"
    role: "primary"
  - type: "module"
    target: "modules/workflow-core"
    role: "secondary"
information_class: "exploration"
topics: []
related_docs:
  - "default:explorations/proposal-archive-document-structure/index.md"
archive_target: "default:architecture/proposal-archive-document-structure.md"
created: "2026-05-22T08:17:52.067Z"
updated: "2026-05-22T08:17:52.067Z"
exploration_slug: "proposal-archive-document-structure"
finding_id: "F-004-archive-updates-are-append-only-and-marker-based"
evidence_sources: []
---

# F-004 归档更新采用追加式且带标记的幂等写入

## Ownership summary

- Primary module: modules/workflow-core
- System / architecture targets: architecture/proposal-archive-document-structure.md
- Convention targets: none
- Canonical reading path: proposal-archive-document-structure/findings/F-004-archive-updates-are-append-only-and-marker-based.md

## Statement

当前归档实现不会无条件覆盖目标文档，而是先检查是否已存在对应 proposal 的历史标记，再决定是否追加新的归档块。

## Why it matters

这使归档生成更接近“长期文档维护”而不是“导出一次性报告”。因此目标文档结构必须允许多次追加，同时避免重复写入和覆盖人工编辑内容。

## References

- [flowforge.js](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/scripts/lib/flowforge.js)
