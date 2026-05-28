---
doc_type: "finding"
title: "F-012 当目标类型尚无现有最终文档时，baseline 缺口应提示而不是阻断"
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
finding_id: "F-012-baseline-gap-should-be-a-warning-not-an-error"
evidence_sources: []
---

# F-012 当目标类型尚无现有最终文档时，baseline 缺口应提示而不是阻断

## Ownership summary

- Primary module: modules/workflow-core
- System / architecture targets: architecture/proposal-archive-document-structure.md
- Convention targets: none
- Canonical reading path: proposal-archive-document-structure/findings/F-012-baseline-gap-should-be-a-warning-not-an-error.md

## Statement

如果当前 workspace 里某个 archive target 类型还没有任何现有最终文档，proposal 创建不应该直接失败，而应该发出 baseline gap 警告，提醒这是该类型知识库的初始建立阶段。

## Why it matters

这避免了新 workspace 或首次建立某类文档时被“必须先有 canonical corpus”卡住。知识库维护应当支持从零开始建立基线，而不是只服务于已有成熟体系。

## References

- [flowforge.js](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/scripts/lib/flowforge.js)
