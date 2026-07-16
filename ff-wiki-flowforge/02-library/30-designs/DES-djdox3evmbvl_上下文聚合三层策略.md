---
id: DES-djdox3evmbvl
title: 上下文聚合三层策略
type: design
status: draft
importance: should
links:
    - target: FIND-djdoa1ftfow2
      relation: references
    - target: STR-djdouollezyy
      relation: indexes
created: 2026-06-20T15:15:39.832532725+08:00
updated: 2026-06-20T15:15:39.833557925+08:00
---

上下文聚合按三层策略执行：Level 1精确匹配（始终输出当前proposal直接关联卡片、importance:must约定卡片、活跃任务依赖卡片），Level 2图遍历扩展（按token预算取一阶邻居links+backlinks，按relation优先级排序：constrains>implements/satisfies>records/discovers>references>related），Level 3 Structure Note摘要（如有剩余预算提供相关领域STR概要，不含完整内容）。

## Links

### Outgoing

- [STR-djdouollezyy]() [structure] - 上下文预算与聚合策略
- [FIND-djdoa1ftfow2]() [finding] - Curation Plan: docs/ 知识导入

