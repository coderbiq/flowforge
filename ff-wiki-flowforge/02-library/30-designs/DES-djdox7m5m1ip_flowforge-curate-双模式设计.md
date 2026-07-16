---
id: DES-djdox7m5m1ip
title: flowforge-curate 双模式设计
type: design
status: draft
importance: should
links:
    - target: FIND-djdoa1ftfow2
      relation: references
    - target: STR-djdowcqvi5bx
      relation: indexes
created: 2026-06-20T15:15:48.979707637+08:00
updated: 2026-06-20T15:15:48.980742138+08:00
---

flowforge-curate统一的知识策展SKILL处理两种来源：Mode A外部资料导入（读长文文件→识别知识单元→用自己的话重述→标注出处）和Mode B proposal归档（扫描proposal卡片→筛选可复用候选→评估知识类型）。两种模式在提取知识单元后汇入完全相同共享流程：聚类→审查计划→用户确认→计划卡→分批CLI写入。差异点：提取阶段的输入不同（文本vs卡片网络），Mode B多一个提案收尾操作。

## Links

### Outgoing

- [STR-djdowcqvi5bx]() [structure] - 知识策展与 Library 导入
- [FIND-djdoa1ftfow2]() [finding] - Curation Plan: docs/ 知识导入

