---
id: DES-djdopaa5b1p6
title: 卡片 ID 规范
type: design
status: draft
importance: should
links:
    - target: FIND-djdoa1ftfow2
      relation: references
    - target: STR-djdocjmfl2g0
      relation: indexes
created: 2026-06-20T15:05:27.870635382+08:00
updated: 2026-06-20T15:05:27.871498383+08:00
---

卡片ID格式：{TYPE}-{proposalTs}-{cardTs}。proposalTs为proposal创建时间的Base36编码（如2x9k3m00），cardTs为卡片创建时间的Base36编码。任务ID特殊格式：TASK-{proposalTs}-{type}-{taskTs}，其中任务类型字母编码：a(analysis)、i(implementation)、t(test)、d(docs)、f(fix)、r(refactor)、c(config)。子任务使用父ID加字母后缀（-a、-b、-c）。全局卡片（CONV/MOD/STR）无proposal归属，使用序号如CONV-001。

## Links

### Outgoing

- [STR-djdocjmfl2g0]() [structure] - 卡片系统核心模型
- [FIND-djdoa1ftfow2]() [finding] - Curation Plan: docs/ 知识导入

