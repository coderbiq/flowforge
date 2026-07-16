---
id: DES-djdosxygvfsx
title: flowforge context 上下文命令
type: design
status: draft
importance: should
links:
    - target: FIND-djdoa1ftfow2
      relation: references
    - target: STR-djdoqf3qx922
      relation: indexes
created: 2026-06-20T15:10:14.499736497+08:00
updated: 2026-06-20T15:10:14.500717497+08:00
---

flowforge context按场景裁剪输出：proposal模式输出proposal根卡+需求索引树入口+活跃任务摘要+焦点卡+反链证据+深读建议；task模式输出任务摘要+直接链接卡+反链证据。上下文聚合分三层：Level 1精确匹配（始终输出当前任务直接链接卡片和importance:must约束）、Level 2图遍历扩展（优先constrains>implements/satisfies>records/discovers>references>related）、Level 3 Structure Note摘要（如有剩余预算）。

## Links

### Outgoing

- [STR-djdoqf3qx922]() [structure] - CLI 命令体系设计
- [FIND-djdoa1ftfow2]() [finding] - Curation Plan: docs/ 知识导入

