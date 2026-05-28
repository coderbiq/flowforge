---
doc_type: "finding"
title: "F-008 已归档知识库应成为后续探索的默认来源"
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
finding_id: "F-008-archived-knowledge-base-should-be-the-default-exploration-source"
evidence_sources: []
---

# F-008 已归档知识库应成为后续探索的默认来源

## Ownership summary

- Primary module: modules/workflow-core
- System / architecture targets: architecture/proposal-archive-document-structure.md
- Convention targets: none
- Canonical reading path: proposal-archive-document-structure/findings/F-008-archived-knowledge-base-should-be-the-default-exploration-source.md

## Statement

后续探索不应默认从空白开始，而应先以已归档的模块文档、architecture 文档和 ADR 作为默认参考语料，再针对差距、冲突或新增问题发起新的探索。

## Why it matters

如果归档产物不是后续探索的默认起点，知识库就会被不断绕开，最终变成“写完就停”的静态档案，而不是持续累积的系统知识基础。

## References

- [lifecycle.md](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/guides/lifecycle.md)
- [authoring-rules.md](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/guides/authoring-rules.md)
- [AGENTS.md](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/templates/project/AGENTS.md)
