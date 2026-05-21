# 任务拆分标准与大型提案分阶段执行

- Status: archived
- Question: 任务拆分应该有哪些标准，大型提案如何分阶段执行任务和跟踪？
- Owner: Codex
- Created: 2026-05-21T00:00:00+08:00
- Updated: 2026-05-21T00:00:00+08:00

## Context

当前 `FlowForge` 的提案已经能表达任务列表和完成条件，但还没有一套稳定、可复用的任务拆分标准。对于小改动，简单任务列表足够；对于大型提案，如果没有阶段边界、依赖规则和跟踪口径，任务图会变成“按文件清单”而不是“按交付路径”。

这个探索要回答两个问题：

- 什么样的拆分才算是好的任务拆分
- 大型提案应该如何按阶段推进，并在 proposal 生命周期里持续跟踪

## Current understanding

- 任务拆分应围绕可交付成果，而不是围绕代码文件或实现手段。
- 大型提案通常需要至少三层视角：阶段、任务、检查点。
- `task-map` 已经有 `priority`、`depends_on` 和 `completion_definition`，说明阶段化管理可以先从规范层收敛，而不一定先扩 schema。
- 任务跟踪需要明确“卡在哪一层”：需求未冻结、设计未冻结、实现未完成，还是验收未完成。

## Findings

- [F-001](./findings/F-001-task-map-already-has-basic-structure-for-staged-work.md) 现有 task map 已经具备表达阶段化工作的基础字段。
- [F-002](./findings/F-002-archive-gates-imply-task-completion-must-be-explicit.md) 归档门槛要求任务完成状态明确，因此阶段跟踪必须可验证。

## Candidate decisions

- [ADR-001](../../decisions/ADR-001-task-splitting-deliverable-first.md) 已归档为稳定决策记录。

## Open questions

- 任务粒度是否应该有建议范围，例如单任务目标时长或复杂度上限？
- 大型提案是否需要显式阶段字段，还是只靠任务分组和依赖关系表达？
- 阶段状态应由 proposal 记录，还是由外部任务后端推导？

## Archived result

- `workflow/guides/task-splitting.md` 已成为 canonical 任务拆分规范。
- `workflow/templates/docs/proposals/task-map.md` 已更新为 milestone-first 模板。
- 相关稳定决策已归档到 [ADR-001](../../decisions/ADR-001-task-splitting-deliverable-first.md)。
