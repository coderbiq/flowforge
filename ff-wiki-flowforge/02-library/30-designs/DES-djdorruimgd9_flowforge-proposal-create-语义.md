---
id: DES-djdorruimgd9
title: flowforge proposal create 语义
type: design
status: draft
importance: should
links:
    - target: FIND-djdoa1ftfow2
      relation: references
    - target: STR-djdoqf3qx922
      relation: indexes
created: 2026-06-20T15:08:42.835950941+08:00
updated: 2026-06-20T15:08:42.837816342+08:00
---

proposal create在当前项目01-workspace/01-active/下创建提案目录，初始化ROOT card、顶层STR-REQ卡和90-cards/目录。创建后默认激活该提案写入当前提案指针。ROOT卡indexes->STR-REQ。提案解析顺序：显式--project > sqlite currentProjectId > 单项目自动选中 > 报错要求先project use。

## Links

### Outgoing

- [STR-djdoqf3qx922]() [structure] - CLI 命令体系设计
- [FIND-djdoa1ftfow2]() [finding] - Curation Plan: docs/ 知识导入

