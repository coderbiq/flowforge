---
id: DEC-djdoxsnivpxr
title: CLI/SKILL 职责边界
type: decision
status: draft
importance: should
links:
    - target: FIND-djdoa1ftfow2
      relation: references
    - target: STR-djdowcqvi5bx
      relation: indexes
created: 2026-06-20T15:16:34.774888927+08:00
updated: 2026-06-20T15:16:34.776023428+08:00
---

CLI和SKILL有严格职责边界：CLI负责卡片CRUD、链接管理、索引重建、查询检索、校验；不负责内容理解、语义拆分、知识重组、分类判断。SKILL负责理解长文内容、拆分为原子知识、组织卡片结构、判定知识类型、写入卡片；不负责直接操作文件、自行构建索引。任何需要理解内容的步骤都属于SKILL不属于CLI。CLI只提供卡片粒度的原子操作，SKILL组合这些原子操作完成导入/归档流程。

## Links

### Outgoing

- [STR-djdowcqvi5bx]() [structure] - 知识策展与 Library 导入
- [FIND-djdoa1ftfow2]() [finding] - Curation Plan: docs/ 知识导入

