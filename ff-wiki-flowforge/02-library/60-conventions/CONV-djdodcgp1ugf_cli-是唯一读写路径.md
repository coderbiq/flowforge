---
id: CONV-djdodcgp1ugf
title: CLI 是唯一读写路径
type: convention
status: draft
importance: should
links:
    - target: FIND-djdoa1ftfow2
      relation: references
    - target: STR-djdocjkam32l
      relation: indexes
created: 2026-06-20T14:49:52.250189387+08:00
updated: 2026-06-20T14:49:52.251473488+08:00
---

Agent 通过 CLI 命令读写所有卡片，不得直接操作目标项目文件系统。SKILL 的描述、触发、流程、产出自检通过 SKILL.md 定义，但 Agent 对卡片的增删改查都必须使用 flowforge CLI 命令。CLI 是唯一的写入入口，确保所有卡片操作经过校验和索引。

## Links

### Outgoing

- [STR-djdocjkam32l]() [structure] - FlowForge 项目定位与架构设计
- [FIND-djdoa1ftfow2]() [finding] - Curation Plan: docs/ 知识导入

