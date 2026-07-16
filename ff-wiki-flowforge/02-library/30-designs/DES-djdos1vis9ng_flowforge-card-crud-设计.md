---
id: DES-djdos1vis9ng
title: flowforge card CRUD 设计
type: design
status: draft
importance: should
links:
    - target: FIND-djdoa1ftfow2
      relation: references
    - target: STR-djdoqf3qx922
      relation: indexes
created: 2026-06-20T15:09:04.664513708+08:00
updated: 2026-06-20T15:09:04.665720709+08:00
---

flowforge card是通用卡片CRUD命令：create（自动生成文件名，支持--links）、read（支持--summary只读摘要和--section只读指定段落）、update（支持更新title/body/links/status/importance，文件名保持创建时稳定）、refresh（刷新CLI生成的内部导航）、delete（仅draft状态可直接删除）、list（基于类型目录+frontmatter筛选）、related（图遍历，支持--relation过滤和--depth深度控制）、link/unlink（维护frontmatter链接关系）、search（先按query命中再按type/status/domain/tag缩小范围）。

## Links

### Outgoing

- [STR-djdoqf3qx922]() [structure] - CLI 命令体系设计
- [FIND-djdoa1ftfow2]() [finding] - Curation Plan: docs/ 知识导入

