---
id: CONV-djdoq3w0onh0
title: 子卡主动指向上游
type: convention
status: draft
importance: should
links:
    - target: FIND-djdoa1ftfow2
      relation: references
    - target: STR-djdonhsylv9v
      relation: indexes
created: 2026-06-20T15:06:32.319816164+08:00
updated: 2026-06-20T15:06:32.320538565+08:00
---

链接方向遵循子卡主动指向上游原则：log主动records->TASK/REQ/DES/ROOT，finding主动discovers->TASK/REQ/DES，task主动implements/designs/satisfies/requires/constrains->上游卡，普通卡主动belongs_to->ROOT。任务卡、root card、需求索引卡不因每个新证据反复回写。反向关系由sqlite索引生成：查看任务详情时查询所有records->task的log卡，查看需求详情时查询所有satisfies/analyzes/designs->requirement的卡。

## Links

### Outgoing

- [STR-djdonhsylv9v]() [structure] - 卡片架构不变量与约束
- [FIND-djdoa1ftfow2]() [finding] - Curation Plan: docs/ 知识导入

