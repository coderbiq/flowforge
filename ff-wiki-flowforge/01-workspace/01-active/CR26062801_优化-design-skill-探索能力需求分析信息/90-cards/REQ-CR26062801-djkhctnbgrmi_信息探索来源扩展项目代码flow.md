---
id: REQ-CR26062801-djkhctnbgrmi
title: 信息探索来源扩展：项目代码、FlowForge知识库与外部文档
type: requirement
status: draft
importance: should
tags:
    - design-skill
    - exploration
    - external-sources
links:
    - target: PROP-CR26062801
      relation: belongs_to
created: 2026-06-28T06:48:43.218376583Z
updated: 2026-06-28T06:48:43.219104247Z
source: CR26062801
domain: skill-design
---

# 信息探索来源扩展：项目代码、FlowForge知识库与外部知识库

## Summary
Design skill 的信息探索来源应不仅限于当前项目代码和 FlowForge 知识库（library），还应支持探索通过配置指定的外部知识库。外部知识库可能是团队历史文档库、非卡片形式的团队知识积累等场景。FlowForge 应保持开放设计，在分析与设计阶段能探索不止内部的知识库。

核心机制：通过配置文件指定外部知识库的位置，在 proposal 中可引用这些外部知识源。

## Source
用户反馈：探索信息来源需要扩展，不止项目代码和 FlowForge 知识库。

## Acceptance
- 提供配置外部知识库位置的机制（如配置文件中的 knowledge_sources 字段）
- Design skill 的探索流程（workflow-rules 或 library-discovery）中增加外部知识库探索指引
- 外部知识库的发现、查询和引用流程与 FlowForge library 不同但有互操作路径

## Scope
- 扩展信息探索的工作流程
- 可能涉及 CLI 新增配置命令或探索命令
- 定义外部知识库的发现、评估、引用机制

## Links

### Outgoing

- [PROP-CR26062801](../../../../03-proposal/CR26062801_优化-design-skill-探索能力需求分析信息.md) [proposal] - 优化 design skill 探索能力：需求分析、信息探索来源与外部知识集成

### Incoming

- [TASK-CR26062801-a-djkhn6w5spyt](TASK-CR26062801-a-djkhn6w5spyt_分析外部知识库的查询接口可信度与优先级策略.md) [task] - 分析外部知识库的查询接口、可信度与优先级策略

## Open Questions
- 外部知识库的"查询"机制：Agent 直接读取文件系统？还是通过 CLI 提供统一查询接口？
- 外部知识库的可信度如何评估？配置文件是否包含可信度标注？
- 外部知识库与 FlowForge library 的优先级关系（先查 library 还是先查外部库）？

