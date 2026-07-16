---
id: DES-djdotwt01934
title: Design SKILL 主流程
type: design
status: draft
importance: should
links:
    - target: FIND-djdoa1ftfow2
      relation: references
    - target: STR-djdosg2sqonp
      relation: indexes
created: 2026-06-20T15:11:30.356504636+08:00
updated: 2026-06-20T15:11:30.357629937+08:00
---

Design SKILL的主流程：解析当前project/proposal（project current + proposal current + proposal inspect + context proposal）→更新需求索引树（STR卡保持7-15条，使用structure add/remove）→拆出原子requirement卡（一个用户可感知的行为/约束/验收点）→对不确定点创建analysis task→通过CLI发现library上下文（library suggest / card search --scope library / card read --summary）→结论稳定时创建设计卡→可执行时创建implementation task→每轮记录log并汇报。这个流程不是一次性阶段，而是可回退的循环。

## Links

### Outgoing

- [STR-djdosg2sqonp]() [structure] - Design SKILL 工作流
- [FIND-djdoa1ftfow2]() [finding] - Curation Plan: docs/ 知识导入

