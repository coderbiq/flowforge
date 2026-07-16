---
id: DES-djdori1l319f
title: proposal create 必须创建 ROOT 和 STR-REQ
type: design
status: draft
importance: should
links:
    - target: FIND-djdoa1ftfow2
      relation: references
    - target: STR-djdonhsylv9v
      relation: indexes
created: 2026-06-20T15:08:21.495525619+08:00
updated: 2026-06-20T15:08:21.496693019+08:00
---

flowforge proposal create 必须创建两张卡片：ROOT-{proposalId}.md（type为proposal）和STR-{proposalId}-REQ.md（type为structure）。ROOT卡indexes->STR-{proposalId}-REQ，正文Entries由FlowForge生成使用Markdown链接指向STR卡。REQ index卡belongs_to->ROOT-{proposalId}，初始Entries为- None。同时创建90-cards/目录存放内容卡片。

## Links

### Outgoing

- [STR-djdonhsylv9v]() [structure] - 卡片架构不变量与约束
- [FIND-djdoa1ftfow2]() [finding] - Curation Plan: docs/ 知识导入

