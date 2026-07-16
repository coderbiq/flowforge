---
id: CONV-djdov6wl1nf5
title: Library 查询只通过 CLI
type: convention
status: draft
importance: should
links:
    - target: FIND-djdoa1ftfow2
      relation: references
    - target: STR-djdosg2sqonp
      relation: indexes
created: 2026-06-20T15:13:10.705180933+08:00
updated: 2026-06-20T15:13:10.706330034+08:00
---

library查询遵循三层模型：候选发现（CLI根据当前需求/任务/关键词返回候选卡片摘要）、结构化筛选（Agent根据类型/标签/领域/关系缩小范围）、定点读取（只对少量高相关卡片调用card read读取全文）。禁止行为：Agent直接读取02-library/文件、用shell grep遍历卡片库、没有筛选就批量读全文、把不确定相关的候选全部链接进来。library的查找、筛选、摘要和读取都必须通过CLI。

## Links

### Outgoing

- [STR-djdosg2sqonp]() [structure] - Design SKILL 工作流
- [FIND-djdoa1ftfow2]() [finding] - Curation Plan: docs/ 知识导入

