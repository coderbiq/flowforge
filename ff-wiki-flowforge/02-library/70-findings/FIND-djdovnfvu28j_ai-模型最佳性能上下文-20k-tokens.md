---
id: FIND-djdovnfvu28j
title: AI 模型最佳性能上下文 ≤ 20K tokens
type: finding
status: draft
importance: should
links:
    - target: FIND-djdoa1ftfow2
      relation: references
    - target: STR-djdouollezyy
      relation: indexes
created: 2026-06-20T15:13:46.700677483+08:00
updated: 2026-06-20T15:13:46.701869583+08:00
---

行业研究表明模型最佳性能上下文区间为≤20K tokens。超过50K tokens后性能显著退化。工具输出可占总上下文81%（需严格控制）。Lost-in-the-Middle现象导致中段信息准确率下降30%+。结论：问题不是如何塞更多token，而是如何只保留正确的token。

## Links

### Outgoing

- [STR-djdouollezyy]() [structure] - 上下文预算与聚合策略
- [FIND-djdoa1ftfow2]() [finding] - Curation Plan: docs/ 知识导入

