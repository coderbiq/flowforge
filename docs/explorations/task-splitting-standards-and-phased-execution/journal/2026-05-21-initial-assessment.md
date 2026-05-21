# Journal Entry

- Timestamp: 2026-05-21T00:00:00+08:00
- Actor: Codex

## What changed

开始梳理任务拆分标准和大型提案的阶段化执行方式，重点确认现有 task-map 和 lifecycle 已经能支持哪些表达，哪些部分还需要通过规范补齐。

## Evidence

- `workflow/schema/task-map.schema.yaml`
- `workflow/guides/lifecycle.md`
- `workflow/guides/archive-rules.md`

## New questions

- 阶段边界应该如何定义才不会和任务边界重复？
- 跟踪状态是应该落在 proposal 文档里，还是依赖任务后端状态？
