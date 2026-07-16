---
id: CONV-djdoqs9dshib
title: 每张卡片必须至少有一个 outbound link
type: convention
status: draft
importance: should
links:
    - target: FIND-djdoa1ftfow2
      relation: references
    - target: STR-djdonhsylv9v
      relation: indexes
created: 2026-06-20T15:07:25.370664919+08:00
updated: 2026-06-20T15:07:25.37174592+08:00
---

除proposal root（ROOT）和全局入口索引（STR-HOME）外，所有卡片必须至少有一个outbound frontmatter link。proposal作用域内创建的普通卡片必须自动补belongs_to->ROOT-{proposalId}。非proposal作用域创建卡片时如果没有任何--links，CLI必须拒绝写入。这确保没有孤儿卡片存在于知识网络中。

## Links

### Outgoing

- [STR-djdonhsylv9v]() [structure] - 卡片架构不变量与约束
- [FIND-djdoa1ftfow2]() [finding] - Curation Plan: docs/ 知识导入

