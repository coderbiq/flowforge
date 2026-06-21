---
id: CONV-djdoqp5js3bq
title: Requirement Index 只索引 REQ 和 STR
type: convention
status: draft
importance: should
links:
    - target: FIND-djdoa1ftfow2
      relation: references
    - target: STR-djdonhsylv9v
      relation: indexes
created: 2026-06-20T15:07:18.608510822+08:00
updated: 2026-06-20T15:07:18.609984522+08:00
---

顶层需求索引卡STR-{proposalId}-REQ的indexes目标只允许requirement和structure卡片。禁止进入requirement index的类型：design、task、log、finding、decision、convention、module。索引卡直接条目建议上限为15条，超过时应裂变为子STR形成索引树。单个提案有且只有一个顶层需求索引入口。

## Links

### Outgoing

- [STR-djdonhsylv9v]() [structure] - 卡片架构不变量与约束
- [FIND-djdoa1ftfow2]() [finding] - Curation Plan: docs/ 知识导入

