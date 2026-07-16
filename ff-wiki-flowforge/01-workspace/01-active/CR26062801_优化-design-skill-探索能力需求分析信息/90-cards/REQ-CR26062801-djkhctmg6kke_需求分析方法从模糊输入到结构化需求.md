---
id: REQ-CR26062801-djkhctmg6kke
title: 需求分析方法：从模糊输入到结构化需求
type: requirement
status: draft
importance: should
tags:
    - design-skill
    - requirement-analysis
links:
    - target: PROP-CR26062801
      relation: belongs_to
created: 2026-06-28T06:48:43.165986Z
updated: 2026-06-28T06:48:43.166559Z
source: CR26062801
domain: skill-design
---

# 需求分析方法：从模糊输入到结构化需求

## Summary
Design skill 需要具备将用户模糊、自然语言输入转化为结构化需求（requirement card）的分析能力。对于复杂需求，应创建分析任务（analysis task）来管理分析过程，但分析结果应增量更新到需求卡中，而非等分析任务完成后才创建需求卡。卡片与任务协同演进。

## Source
用户反馈：design skill 探索能力需要优化，需求分析逻辑需要明确。

## Acceptance
- Design skill 的 workflow-rules 中增加需求分析阶段指引（分析方法、检查清单）
- 复杂需求可通过 analysis task 管理分析过程，分析中间结果增量更新到对应卡片
- 用户输入分析产出至少包含：识别核心目标、边界范围、隐含假设、风险点

## Scope
- 改进 design skill 的需求分析流程
- 不涉及其他 skill（implement、feedback、curate）

## Links

### Outgoing

- [PROP-CR26062801](../../../../03-proposal/CR26062801_优化-design-skill-探索能力需求分析信息.md) [proposal] - 优化 design skill 探索能力：需求分析、信息探索来源与外部知识集成

### Incoming

- [TASK-CR26062801-a-djkhn6v3gg2k](TASK-CR26062801-a-djkhn6v3gg2k_分析需求分析阶段的模式定位与卡片演进机制.md) [task] - 分析需求分析阶段的模式定位与卡片演进机制
#### satisfies
- [DES-CR26062801-djkmcrj23oo9](DES-CR26062801-djkmcrj23oo9_需求分析阶段扩展-index-模式-卡片任务协同演进.md) [design] - 需求分析阶段：扩展 index 模式 + 卡片任务协同演进
- [TASK-CR26062801-i-djkmdw0v1uf4](TASK-CR26062801-i-djkmdw0v1uf4_更新-workflow-rulesmd需求分析步骤-协同演进规则-探索判据.md) [task] - 更新 workflow-rules.md：需求分析步骤 + 协同演进规则 + 探索判据

## Open Questions
- 需求分析应发生在哪个模式？是 index 模式的第一步，还是单独的 clarify 模式？
- 分析任务与需求卡片的协同演进机制如何设计（任务"完成"的判据、卡片何时从 draft 转为 stable）？

