---
id: DES-djdow5s2jp1y
title: 三层上下文加载模型
type: design
status: draft
importance: should
links:
    - target: FIND-djdoa1ftfow2
      relation: references
    - target: STR-djdouollezyy
      relation: indexes
created: 2026-06-20T15:14:26.61962761+08:00
updated: 2026-06-20T15:14:26.620763611+08:00
---

上下文加载采用三层模型：Level 0永久层（始终加载<500 tokens，包括项目元信息、SKILL触发摘要、活跃proposal概要）、Level 1摘要层（按需加载<3000 tokens，包括相关卡片id+title+summary、按importance排序）、Level 2完整层（Agent主动通过flowforge card read按需获取完整内容，每张卡片约100-300 tokens，受maxTokens预算控制）。初始只加载卡片摘要，完整内容通过CLI按需获取。

## Links

### Outgoing

- [STR-djdouollezyy]() [structure] - 上下文预算与聚合策略
- [FIND-djdoa1ftfow2]() [finding] - Curation Plan: docs/ 知识导入

