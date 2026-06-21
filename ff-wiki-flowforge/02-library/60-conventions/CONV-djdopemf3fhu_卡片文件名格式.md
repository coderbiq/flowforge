---
id: CONV-djdopemf3fhu
title: 卡片文件名格式
type: convention
status: draft
importance: should
links:
    - target: FIND-djdoa1ftfow2
      relation: references
    - target: STR-djdocjmfl2g0
      relation: indexes
created: 2026-06-20T15:05:37.31979963+08:00
updated: 2026-06-20T15:05:37.32058593+08:00
---

卡片文件名格式：{ID}_{slug}.md。ID包含完整卡片ID（含类型、proposal、时间戳），slug为标题短横线化（kebab-case），支持中文。slug最大长度50字符，超出截断。依赖关系不编码在文件名中，通过frontmatter的links字段记录，由CLI构建sqlite缓存索引。

## Links

### Outgoing

- [STR-djdocjmfl2g0]() [structure] - 卡片系统核心模型
- [FIND-djdoa1ftfow2]() [finding] - Curation Plan: docs/ 知识导入

