---
id: DEC-djdopoyggch5
title: 任务卡片是一等公民
type: decision
status: draft
importance: should
links:
    - target: FIND-djdoa1ftfow2
      relation: references
    - target: STR-djdocjmfl2g0
      relation: indexes
created: 2026-06-20T15:05:59.815498333+08:00
updated: 2026-06-20T15:05:59.816595033+08:00
---

任务是卡片的一种类型（type: task），存在于知识网络中。任务通过typed links关联需求（satisfies）、设计（implements）、约束（constrains）和依赖（blocks），形成完整追溯链。提供独立的flowforge task命令组用于高频操作（create/ready/claim/done/block等），task create底层调用card create --type task。

## Links

### Outgoing

- [STR-djdocjmfl2g0]() [structure] - 卡片系统核心模型
- [FIND-djdoa1ftfow2]() [finding] - Curation Plan: docs/ 知识导入

