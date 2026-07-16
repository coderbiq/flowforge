---
id: DES-djdoptgehibt
title: 任务状态流转
type: design
status: draft
importance: should
links:
    - target: FIND-djdoa1ftfow2
      relation: references
    - target: STR-djdocjmfl2g0
      relation: indexes
created: 2026-06-20T15:06:09.607713853+08:00
updated: 2026-06-20T15:06:09.608797353+08:00
---

任务状态流转：backlog -> not_ready -> ready -> in_progress -> done。backlog表示已识别但尚未准备好的任务。not_ready表示已拆出但依赖未确认假设、open question或分析结论。ready表示可执行任务。blocked状态可从in_progress进入，解除阻塞后回到ready。cancelled可从任意非done状态进入。task ready命令只返回ready状态的任务。

## Links

### Outgoing

- [STR-djdocjmfl2g0]() [structure] - 卡片系统核心模型
- [FIND-djdoa1ftfow2]() [finding] - Curation Plan: docs/ 知识导入

