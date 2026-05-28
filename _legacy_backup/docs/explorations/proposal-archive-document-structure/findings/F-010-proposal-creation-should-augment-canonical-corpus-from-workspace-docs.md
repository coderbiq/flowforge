---
doc_type: "finding"
title: "F-010 proposal 创建时应从 workspace 现有最终文档补充 canonical corpus"
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
finding_id: "F-010-proposal-creation-should-augment-canonical-corpus-from-workspace-docs"
evidence_sources: []
---

# F-010 proposal 创建时应从 workspace 现有最终文档补充 canonical corpus

## Ownership summary

- Primary module: modules/workflow-core
- System / architecture targets: architecture/proposal-archive-document-structure.md
- Convention targets: none
- Canonical reading path: proposal-archive-document-structure/findings/F-010-proposal-creation-should-augment-canonical-corpus-from-workspace-docs.md

## Statement

proposal 创建不应只依赖显式声明的 archive targets 来构造 canonical corpus，还应该从当前 workspace 的现有 `docs/modules/`、`docs/architecture/` 和 `docs/decisions/` 中补充已经存在的最终文档，形成更完整的 baseline 候选集合。

## Why it matters

仅按 archive targets 生成 canonical corpus，容易漏掉 workspace 内已经存在但没有被当前 proposal 直接声明的正式知识。补充 workspace 现有最终文档后，proposal 更容易把自己放到真实知识背景里，而不是只围绕当前变更目标自循环。

## References

- [flowforge.js](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/scripts/lib/flowforge.js)
