---
id: CONV-djdotd0xwzhj
title: task vs card 命令边界
type: convention
status: draft
importance: should
links:
    - target: FIND-djdoa1ftfow2
      relation: references
    - target: STR-djdoqf3qx922
      relation: indexes
created: 2026-06-20T15:10:47.301028827+08:00
updated: 2026-06-20T15:10:47.302324328+08:00
---

task命令是card create --type task的快捷入口。任务是一张卡片（type: task），但操作频率高、流程固定，因此提供独立task命令组。task create底层调用card create --type task。未显式传--proposal时，task create默认使用当前项目的current proposal；如果没有current proposal，则创建为全局任务卡。

## Links

### Outgoing

- [STR-djdoqf3qx922]() [structure] - CLI 命令体系设计
- [FIND-djdoa1ftfow2]() [finding] - Curation Plan: docs/ 知识导入

