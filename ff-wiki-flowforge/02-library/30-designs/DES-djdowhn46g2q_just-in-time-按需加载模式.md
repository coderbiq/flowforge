---
id: DES-djdowhn46g2q
title: Just-in-Time 按需加载模式
type: design
status: draft
importance: should
links:
    - target: FIND-djdoa1ftfow2
      relation: references
    - target: STR-djdouollezyy
      relation: indexes
created: 2026-06-20T15:14:52.44142353+08:00
updated: 2026-06-20T15:14:52.443140931+08:00
---

Agent不预先加载大量内容，而是持有轻量标识符，在需要时通过工具调用按需获取。flowforge context输出卡片ID+摘要（轻量引用），Agent通过card read按需获取完整内容。这模仿人类工作方式：不记忆全文，使用文件系统和索引。渐进式探索：每次交互获取的上下文指导下一次决策。总消耗约3000-5000 tokens按需增长。

## Links

### Outgoing

- [STR-djdouollezyy]() [structure] - 上下文预算与聚合策略
- [FIND-djdoa1ftfow2]() [finding] - Curation Plan: docs/ 知识导入

