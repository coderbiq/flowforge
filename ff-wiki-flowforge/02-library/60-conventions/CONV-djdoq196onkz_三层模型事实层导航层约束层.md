---
id: CONV-djdoq196onkz
title: 三层模型：事实层/导航层/约束层
type: convention
status: draft
importance: should
links:
    - target: FIND-djdoa1ftfow2
      relation: references
    - target: STR-djdonhsylv9v
      relation: indexes
created: 2026-06-20T15:06:26.585606583+08:00
updated: 2026-06-20T15:06:26.586400683+08:00
---

卡片架构采用三层模型：事实层（YAML frontmatter的links字段，是查询和索引重建来源，遵循子卡主动指向上游原则）、导航层（正文中面向人类阅读的内部导航，只能由CLI根据frontmatter关系格式化生成标准Markdown链接）、约束层（由CLI和validate命令实现：写入前检查目标存在、wikilink报错、requirement index只收录requirement/structure、所有非ROOT卡至少有一个outbound link）。

## Links

### Outgoing

- [STR-djdonhsylv9v]() [structure] - 卡片架构不变量与约束
- [FIND-djdoa1ftfow2]() [finding] - Curation Plan: docs/ 知识导入

