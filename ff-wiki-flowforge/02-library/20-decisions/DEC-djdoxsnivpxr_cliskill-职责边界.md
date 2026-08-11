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
created: 2026-06-20T07:16:34.774888927Z
updated: 2026-06-20T07:16:34.776023428Z
---

CLI和SKILL有严格职责边界：CLI负责卡片CRUD、链接管理、索引重建、查询检索、校验；不负责内容理解、语义拆分、知识重组、分类判断。SKILL负责理解长文内容、拆分为原子知识、组织卡片结构、判定知识类型、写入卡片；不负责直接操作文件、自行构建索引。任何需要理解内容的步骤都属于SKILL不属于CLI。CLI只提供卡片粒度的原子操作，SKILL组合这些原子操作完成导入/归档流程。

## Links

### Outgoing

- [STR-djdowcqvi5bx](../structures/STR-djdowcqvi5bx_知识策展与-library-导入.md) [structure] - 知识策展与 Library 导入
- [FIND-djdoa1ftfow2](../70-findings/FIND-djdoa1ftfow2_curation-plan-docs-知识导入.md) [finding] - Curation Plan: docs/ 知识导入

### Incoming

- [DEC-CR26081101-dkltaxm2rgcf](../../01-workspace/CR26081101_复杂需求多代理分析编排/90-cards/DEC-CR26081101-dkltaxm2rgcf_正文直接编辑与-cli-结构化操作边界.md) [decision] - 正文直接编辑与 CLI 结构化操作边界

