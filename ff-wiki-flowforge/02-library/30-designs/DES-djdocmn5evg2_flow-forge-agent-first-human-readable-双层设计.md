---
id: DES-djdocmn5evg2
title: FlowForge Agent-First Human-Readable 双层设计
type: design
status: draft
importance: should
links:
    - target: FIND-djdoa1ftfow2
      relation: references
    - target: STR-djdocjkam32l
      relation: indexes
created: 2026-06-20T14:48:56.044128174+08:00
updated: 2026-06-20T14:48:56.045187275+08:00
---

FlowForge 采用双重视角设计：Agent 是主要接口（通过 SKILL 触发、CLI 命令执行、卡片消费结构化知识），同时面向人类开发者（卡片内容人类可读、协议透明、知识网络可作为项目文档查阅）。这种设计确保 Agent 执行过程可追溯、可审计。

## Links

### Outgoing

- [STR-djdocjkam32l]() [structure] - FlowForge 项目定位与架构设计
- [FIND-djdoa1ftfow2]() [finding] - Curation Plan: docs/ 知识导入

