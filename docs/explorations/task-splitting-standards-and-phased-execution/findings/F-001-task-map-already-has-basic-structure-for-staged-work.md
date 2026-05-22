---
doc_type: "finding"
title: "F-001 现有 task map 已经具备阶段化工作的基础结构"
status: "validated"
workspace: "default"
module_scope: []
system_scope: []
convention_scope: []
ownership:
  - type: "convention"
    target: "conventions/task-splitting.md"
    role: "primary"
  - type: "module"
    target: "modules/workflow-core"
    role: "secondary"
information_class: "exploration"
topics: []
related_docs:
  - "default:explorations/task-splitting-standards-and-phased-execution/index.md"
archive_target: "default:conventions/task-splitting.md"
created: "2026-05-22T08:17:52.067Z"
updated: "2026-05-22T08:17:52.067Z"
exploration_slug: "task-splitting-standards-and-phased-execution"
finding_id: "F-001-task-map-already-has-basic-structure-for-staged-work"
evidence_sources: []
---

# F-001 现有 task map 已经具备阶段化工作的基础结构

## Ownership summary

- Primary module: modules/workflow-core
- System / architecture targets: none
- Convention targets: conventions/task-splitting.md
- Canonical reading path: task-splitting-standards-and-phased-execution/findings/F-001-task-map-already-has-basic-structure-for-staged-work.md

## Statement

当前 task map 已经包含 `priority`、`depends_on` 和 `completion_definition`，足以表达任务优先级、依赖关系和验收条件，因此阶段化执行可以先通过规范和约定建立，不必立即引入全新的任务模型。

## Why it matters

这意味着“大型提案如何分阶段推进”首先是一个拆分方法问题，其次才是一个 schema 扩展问题。如果阶段、检查点和依赖关系能被稳定描述，后续就能在现有结构上形成一致的执行节奏。

## References

- [task-map.schema.yaml](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/schema/task-map.schema.yaml)
- [lifecycle.md](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/guides/lifecycle.md)
