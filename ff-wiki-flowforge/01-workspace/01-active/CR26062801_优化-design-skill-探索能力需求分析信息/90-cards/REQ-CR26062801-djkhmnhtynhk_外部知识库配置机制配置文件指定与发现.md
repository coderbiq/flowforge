---
id: REQ-CR26062801-djkhmnhtynhk
title: 外部知识库配置机制：配置文件指定与发现
type: requirement
status: draft
importance: should
links:
    - target: PROP-CR26062801
      relation: belongs_to
created: 2026-06-28T07:01:33.468317Z
updated: 2026-06-28T07:01:33.468321Z
source: CR26062801
---

# 外部知识库配置机制：配置文件指定与发现

## Summary
用户可通过 FlowForge 配置文件指定外部知识库的位置（如团队历史文档路径、非卡片形式的知识库目录），在 proposal 分析设计时可通过这些配置引用外部知识源。外部知识库与 FlowForge library 并存，互补探索。

配置应支持多种格式的知识库（Markdown 目录树、纯文本、结构化文档等），并提供发现/验证机制。

## Source
用户反馈：FlowForge 应该是开放的，能探索不止内部的知识库。在配置中指定知识库位置，在 proposal 中可用。

## Acceptance
- 配置文件支持知识库来源定义（路径、类型、可信度标注等）
- CLI 提供知识库来源的添加/移除/列表命令
- Design skill 的探索流程关联配置的外部知识库

## Scope
- 配置文件扩展（新增 knowledge_sources 配置段）
- 可能涉及 CLI 新增  子命令
- 外部知识库与 library 的互操作（导入路径）

## Links

### Outgoing

- [PROP-CR26062801](../../../../03-proposal/CR26062801_优化-design-skill-探索能力需求分析信息.md) [proposal] - 优化 design skill 探索能力：需求分析、信息探索来源与外部知识集成

### Incoming

- [TASK-CR26062801-a-djkhn6w5spyt](TASK-CR26062801-a-djkhn6w5spyt_分析外部知识库的查询接口可信度与优先级策略.md) [task] - 分析外部知识库的查询接口、可信度与优先级策略
#### satisfies
- [DES-CR26062801-djkmcrk47ecc](DES-CR26062801-djkmcrk47ecc_外部知识库配置knowledge-sources-配置段与混合查询机制.md) [design] - 外部知识库配置：knowledge_sources 配置段与混合查询机制
- [TASK-CR26062801-i-djkmdw1z53bz](TASK-CR26062801-i-djkmdw1z53bz_config-扩展config-结构体新增-knowledge-sources-字段.md) [task] - Config 扩展：Config 结构体新增 KnowledgeSources 字段

## Open Questions
- 配置格式设计：单个路径还是多个？是否需要标注知识库类型和优先级？
- 外部知识库的验证机制：CLI 检查路径是否存在？是否检查文件格式合法性？
- 是否需要在 proposal 级别引用特定的外部知识库（选择性子集）？

