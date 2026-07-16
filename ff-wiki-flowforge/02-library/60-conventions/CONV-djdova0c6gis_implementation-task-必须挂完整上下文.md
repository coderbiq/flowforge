---
id: CONV-djdova0c6gis
title: Implementation task 必须挂完整上下文
type: convention
status: draft
importance: should
links:
    - target: FIND-djdoa1ftfow2
      relation: references
    - target: STR-djdosg2sqonp
      relation: indexes
created: 2026-06-20T15:13:17.46250014+08:00
updated: 2026-06-20T15:13:17.463844541+08:00
---

implementation task进入ready必须同时满足：至少关联一个requirement、至少关联一个design card、Acceptance可验证、Deliverables明确、Out of Scope明确、相关convention/module/decision已通过library确认。满足以下任一条件只能为not_ready：只有requirement没有design、设计依赖用户未确认的业务假设、影响范围未知、验收标准不可验证、涉及跨项目边界但项目职责未确认。not_ready task必须链接阻塞来源（open question/analysis task/finding/design assumption）。

## Links

### Outgoing

- [STR-djdosg2sqonp]() [structure] - Design SKILL 工作流
- [FIND-djdoa1ftfow2]() [finding] - Curation Plan: docs/ 知识导入

