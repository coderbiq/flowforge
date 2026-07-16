---
id: CONV-djdorf97z66p
title: 内部导航使用 Markdown 链接而非 wikilink
type: convention
status: draft
importance: should
links:
    - target: FIND-djdoa1ftfow2
      relation: references
    - target: STR-djdonhsylv9v
      relation: indexes
created: 2026-06-20T15:08:15.426892869+08:00
updated: 2026-06-20T15:08:15.42811057+08:00
---

内部卡片导航只能由 FlowForge CLI 根据 frontmatter 关系格式化并插入正文，渲染结果使用标准 Markdown 相对路径链接，例如[REQ-xxx](90-cards/REQ-xxx_title.md)。Agent 不手写内部卡片链接。Agent 在卡片正文中手写 Markdown 链接时只能用于外部资料引用。CLI 不生成[[wikilink]]格式。frontmatter links是事实关系来源，正文链接只是人类可读导航层。

## Links

### Outgoing

- [STR-djdonhsylv9v]() [structure] - 卡片架构不变量与约束
- [FIND-djdoa1ftfow2]() [finding] - Curation Plan: docs/ 知识导入

