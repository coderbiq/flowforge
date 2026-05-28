---
doc_type: "finding"
title: "F-007 decision 目标已经有固定的 ADR 正文骨架"
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
finding_id: "F-007-decision-target-has-a-fixed-adr-body"
evidence_sources: []
---

# F-007 decision 目标已经有固定的 ADR 正文骨架

## Ownership summary

- Primary module: modules/workflow-core
- System / architecture targets: architecture/proposal-archive-document-structure.md
- Convention targets: none
- Canonical reading path: proposal-archive-document-structure/findings/F-007-decision-target-has-a-fixed-adr-body.md

## Statement

decision 目标不是自由散文式说明，而是遵循 ADR 结构的固定正文骨架，至少包括 context、decision、alternatives 和 consequences。

## Why it matters

这使 decision 归档天然偏向“稳定结论记录”，而不是实现细节堆积。归档生成时可以追加新的更新块，但 ADR 的正文应该维持足够稳定，确保后续能直接追溯为什么做出这个决策。

## References

- [ADR-template.md](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/templates/docs/decisions/ADR-template.md)
