---
id: CONV-djdov2ndj2vm
title: 分析任务驱动不确定点
type: convention
status: draft
importance: should
links:
    - target: FIND-djdoa1ftfow2
      relation: references
    - target: STR-djdosg2sqonp
      relation: indexes
created: 2026-06-20T15:13:01.441232362+08:00
updated: 2026-06-20T15:13:01.442395863+08:00
---

Design SKILL的关键不是Agent自己臆测所有答案，而是把不确定点变成analysis task。创建analysis task的场景：需求边界不清楚、需要看代码判断影响范围、需要查library中的规范或历史设计、需要比较多个方案、涉及跨项目/前后端/数据模型/兼容性风险。analysis task必须包含Goal/Inputs/Investigation Plan/Expected Outputs/Done When，不能只有标题。analysis task主动analyzes->REQ或STR，过程中产生的log/finding主动链接analysis task。

## Links

### Outgoing

- [STR-djdosg2sqonp]() [structure] - Design SKILL 工作流
- [FIND-djdoa1ftfow2]() [finding] - Curation Plan: docs/ 知识导入

