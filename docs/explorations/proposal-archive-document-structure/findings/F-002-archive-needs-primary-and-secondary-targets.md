---
doc_type: "finding"
title: "F-002 归档需要同时更新主目标和次目标"
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
finding_id: "F-002-archive-needs-primary-and-secondary-targets"
evidence_sources: []
---

# F-002 归档需要同时更新主目标和次目标

## Ownership summary

- Primary module: modules/workflow-core
- System / architecture targets: architecture/proposal-archive-document-structure.md
- Convention targets: none
- Canonical reading path: proposal-archive-document-structure/findings/F-002-archive-needs-primary-and-secondary-targets.md

## Statement

归档流程要求先更新 primary archive target，再更新 secondary targets，因此文档结构需要允许“主文档承载核心结论，次文档承载配套沉淀”的分工。

## Why it matters

如果没有主次结构，归档时就难以判断哪些内容必须进入主文档，哪些内容应该同步到模块文档或 ADR 中作为补充说明。

## References

- [archive-rules.md](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/guides/archive-rules.md)
