---
id: DES-djdoxatu6nko
title: 共享聚类→审查→批次执行流程
type: design
status: draft
importance: should
links:
    - target: FIND-djdoa1ftfow2
      relation: references
    - target: STR-djdowcqvi5bx
      relation: indexes
created: 2026-06-20T15:15:55.974590164+08:00
updated: 2026-06-20T15:15:55.975667064+08:00
---

两种导入模式汇合后的共享流程：阶段3聚类（按概念聚类，每个聚类对应一个STR索引卡含3-15个知识单元，不按原文结构）、阶段4生成审查计划（拟议STR索引卡+原子知识卡列表+标注重复/合并/警告，用户审查确认）、阶段5生成执行计划并写入计划卡（转化为分批条目，搜索library去重标记create/merge/skip）、阶段6分批执行（每轮激活处理一批，更新计划卡进度）、全部完成后index rebuild。

## Links

### Outgoing

- [STR-djdowcqvi5bx]() [structure] - 知识策展与 Library 导入
- [FIND-djdoa1ftfow2]() [finding] - Curation Plan: docs/ 知识导入

