---
doc_type: "finding"
title: "F-011 canonical corpus 补充应按 archive target 类型过滤同类最终文档"
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
finding_id: "F-011-canonical-corpus-augmentation-should-be-type-filtered"
evidence_sources: []
---

# F-011 canonical corpus 补充应按 archive target 类型过滤同类最终文档

## Ownership summary

- Primary module: modules/workflow-core
- System / architecture targets: architecture/proposal-archive-document-structure.md
- Convention targets: none
- Canonical reading path: proposal-archive-document-structure/findings/F-011-canonical-corpus-augmentation-should-be-type-filtered.md

## Statement

在自动补充 canonical corpus 时，应以 proposal 的 archive target 类型为过滤轴，只增补同类型的 workspace 最终文档，而不是把整个 workspace 的所有最终文档都作为 baseline。

## Why it matters

这种筛选方式能让 canonical corpus 保持相关性和可读性。proposal 的 baseline 应该尽量贴近本次变更目标，否则阅读者会被大量无关文档干扰，baseline 也会失去“本次变化参照系”的意义。

## References

- [flowforge.js](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/scripts/lib/flowforge.js)
