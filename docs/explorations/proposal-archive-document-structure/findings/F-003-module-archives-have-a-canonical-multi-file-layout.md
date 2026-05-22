---
doc_type: "finding"
title: "F-003 模块归档目标已经隐含固定的多文件目录骨架"
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
finding_id: "F-003-module-archives-have-a-canonical-multi-file-layout"
evidence_sources: []
---

# F-003 模块归档目标已经隐含固定的多文件目录骨架

## Ownership summary

- Primary module: modules/workflow-core
- System / architecture targets: architecture/proposal-archive-document-structure.md
- Convention targets: none
- Canonical reading path: proposal-archive-document-structure/findings/F-003-module-archives-have-a-canonical-multi-file-layout.md

## Statement

模块归档目标不是单个 markdown 文件，而是一个固定目录结构，至少包含 `README.md`、`design.md`、`api.md` 和 `history.md`。

## Why it matters

这意味着归档生成逻辑不能只关心“写一个最终文档”，而要关心“生成并维护一个目录级文档包”。主文档负责入口和边界，子文档负责设计、接口和历史沉淀。

## References

- [modules/README.md](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/templates/docs/modules/README.md)
- [modules/design.md](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/templates/docs/modules/design.md)
- [modules/api.md](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/templates/docs/modules/api.md)
- [modules/history.md](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/templates/docs/modules/history.md)
