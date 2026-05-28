---
doc_type: "finding"
title: "F-005 三类归档目标共享的是追踪信息层，而不是同一套正文结构"
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
finding_id: "F-005-target-types-share-traceability-not-body-structure"
evidence_sources: []
---

# F-005 三类归档目标共享的是追踪信息层，而不是同一套正文结构

## Ownership summary

- Primary module: modules/workflow-core
- System / architecture targets: architecture/proposal-archive-document-structure.md
- Convention targets: none
- Canonical reading path: proposal-archive-document-structure/findings/F-005-target-types-share-traceability-not-body-structure.md

## Statement

模块、architecture 和 decision 三类归档目标都保留了状态、来源和关联提案等追踪信息，但它们的正文结构明显不同，不能被压成同一套共享模板。

## Why it matters

共享模板如果扩展到正文层，会把三类目标的语义差异抹平，最后得到一个“什么都能写一点，但什么都写不深”的归档结果。更合理的抽象边界是把元信息头部尽量统一，而正文章节按目标类型分开定义。

## References

- [modules/README.md](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/templates/docs/modules/README.md)
- [modules/design.md](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/templates/docs/modules/design.md)
- [modules/api.md](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/templates/docs/modules/api.md)
- [modules/history.md](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/templates/docs/modules/history.md)
- [architecture/system.md](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/templates/docs/architecture/system.md)
- [ADR-template.md](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/templates/docs/decisions/ADR-template.md)
