---
doc_type: "finding"
title: "F-001 现有归档规则已经区分目标类型"
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
finding_id: "F-001-archive-rules-already-distinguish-target-types"
evidence_sources: []
---

# F-001 现有归档规则已经区分目标类型

## Ownership summary

- Primary module: modules/workflow-core
- System / architecture targets: architecture/proposal-archive-document-structure.md
- Convention targets: none
- Canonical reading path: proposal-archive-document-structure/findings/F-001-archive-rules-already-distinguish-target-types.md

## Statement

当前归档规则已经明确区分 `module`、`architecture` 和 `decision` 三种目标类型，并分别给出了典型映射关系。

## Why it matters

既然目标类型已经被区分，后续就可以进一步定义每一类目标的文档结构，而不是把所有归档内容都塞进同一种模板里。

## References

- [archive-rules.md](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/guides/archive-rules.md)
- [lifecycle.md](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/guides/lifecycle.md)
