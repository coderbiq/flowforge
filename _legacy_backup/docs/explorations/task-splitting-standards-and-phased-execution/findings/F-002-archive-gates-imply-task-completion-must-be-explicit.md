---
doc_type: "finding"
title: "F-002 归档门槛要求任务完成状态必须可验证"
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
finding_id: "F-002-archive-gates-imply-task-completion-must-be-explicit"
evidence_sources: []
---

# F-002 归档门槛要求任务完成状态必须可验证

## Ownership summary

- Primary module: modules/workflow-core
- System / architecture targets: none
- Convention targets: conventions/task-splitting.md
- Canonical reading path: task-splitting-standards-and-phased-execution/findings/F-002-archive-gates-imply-task-completion-must-be-explicit.md

## Statement

归档要求先确认任务后端没有未关闭任务，然后才能更新 archive targets 并把 proposal 状态改为 `archived`。这意味着任务拆分和跟踪不能只停留在“做完了大概什么”的层面，必须能被明确验证。

## Why it matters

如果任务完成定义不明确，归档就会被迫依赖人工判断，导致大型提案在后期收口时出现遗漏。阶段跟踪需要把“可验收”作为设计前提，而不是事后补充。

## References

- [archive-rules.md](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/guides/archive-rules.md)
- [lifecycle.md](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/workflow/guides/lifecycle.md)
