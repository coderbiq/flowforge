---
id: DES-djdot12440a8
title: flowforge library 命令设计
type: design
status: draft
importance: should
links:
    - target: FIND-djdoa1ftfow2
      relation: references
    - target: STR-djdoqf3qx922
      relation: indexes
created: 2026-06-20T15:10:21.250509092+08:00
updated: 2026-06-20T15:10:21.252161893+08:00
---

flowforge library提供四类操作：facets（从library卡片中统计可用facet key/value和常见组合）、classify（从任务/需求/设计卡提取候选facet，只输出建议不写入）、suggest（基于facet组合过滤+关键词排序返回候选卡片摘要，Agent必须读取确认后才链接）、import（外部资料预处理后的安全写入，CLI校验类型/关系/外链）、promote（把proposal内稳定的finding/design/convention复制为library卡，保留原卡追溯）。MVP不依赖embedding，使用本地卡片扫描和关键词打分。

## Links

### Outgoing

- [STR-djdoqf3qx922]() [structure] - CLI 命令体系设计
- [FIND-djdoa1ftfow2]() [finding] - Curation Plan: docs/ 知识导入

