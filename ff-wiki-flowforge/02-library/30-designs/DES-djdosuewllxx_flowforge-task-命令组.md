---
id: DES-djdosuewllxx
title: flowforge task 命令组
type: design
status: draft
importance: should
links:
    - target: FIND-djdoa1ftfow2
      relation: references
    - target: STR-djdoqf3qx922
      relation: indexes
created: 2026-06-20T15:10:06.786480521+08:00
updated: 2026-06-20T15:10:06.787610222+08:00
---

flowforge task提供独立快捷命令组：create（创建任务卡片，type编码：a分析/i实现/t测试/d文档/f修复/r重构/c配置）、ready（列出就绪任务，analysis task还需具备Goal/Inputs/Investigation Plan/Expected Outputs/Done When）、claim（认领任务，status从ready变为in_progress）、done（完成任务）、block/unblock（阻塞管理）、status（查看任务详情读取卡片全文）、sub（创建子任务，自动生成子任务ID）、link-add/remove（管理稳定上下文链接，执行过程产生的log/finding不回写任务卡）。

## Links

### Outgoing

- [STR-djdoqf3qx922]() [structure] - CLI 命令体系设计
- [FIND-djdoa1ftfow2]() [finding] - Curation Plan: docs/ 知识导入

