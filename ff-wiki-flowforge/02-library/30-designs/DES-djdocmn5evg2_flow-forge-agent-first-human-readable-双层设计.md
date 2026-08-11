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
created: 2026-06-20T06:48:56.044128174Z
updated: 2026-06-20T06:48:56.045187275Z
---

FlowForge 采用双重视角设计：Agent 是主要接口（通过 SKILL 触发、CLI 命令执行、卡片消费结构化知识），同时面向人类开发者（卡片内容人类可读、协议透明、知识网络可作为项目文档查阅）。这种设计确保 Agent 执行过程可追溯、可审计。

## Links

### Outgoing

- [STR-djdocjkam32l](../structures/STR-djdocjkam32l_flow-forge-项目定位与架构设计.md) [structure] - FlowForge 项目定位与架构设计
- [FIND-djdoa1ftfow2](../70-findings/FIND-djdoa1ftfow2_curation-plan-docs-知识导入.md) [finding] - Curation Plan: docs/ 知识导入

### Incoming

- [FEAT-CR26081101-dklt4yyrcz3r](../../01-workspace/CR26081101_复杂需求多代理分析编排/90-cards/FEAT-CR26081101-dklt4yyrcz3r_重构-flowforge-design-方法论与端到端评估.md) [feature] - 重构 flowforge-design 方法论与端到端评估

