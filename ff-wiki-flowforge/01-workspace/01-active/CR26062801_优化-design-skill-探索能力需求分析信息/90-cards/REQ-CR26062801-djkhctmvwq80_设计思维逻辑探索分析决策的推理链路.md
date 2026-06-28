---
id: REQ-CR26062801-djkhctmvwq80
title: 设计思维逻辑：探索→分析→决策的推理链路
type: requirement
status: draft
importance: should
tags:
    - design-skill
    - design-thinking
    - workflow
links:
    - target: PROP-CR26062801
      relation: belongs_to
created: 2026-06-28T06:48:43.192281Z
updated: 2026-06-28T06:48:43.192974Z
source: CR26062801
domain: skill-design
---

# 设计思维逻辑：探索→分析→决策的推理链路

## Summary
Design skill 需要有清晰的推理方法指引，描述从信息探索到分析判断再到设计决策的完整逻辑链路。当前 skill 的 workflow-rules 列出了七种模式（index/clarify/analyze/discover library/design/split tasks/refresh navigation），但在模式切换的推理依据、探索深度判断等环节缺乏明确指引。

## Source
用户反馈：设计的思维逻辑需要优化。

## Acceptance
- Workflow-rules 中增加推理逻辑章节：何时结束探索、何时进入分析、何时做出设计决策
- 明确"探索足够"的判据（不再有未解决的开放性关键问题）

## Scope
- 改进 design skill 的 workflow-rules 推理逻辑
- 不涉及具体卡片内容模板

## Links

### Outgoing

- [PROP-CR26062801](../../../../03-proposal/CR26062801_优化-design-skill-探索能力需求分析信息.md) [proposal] - 优化 design skill 探索能力：需求分析、信息探索来源与外部知识集成

### Incoming

- [TASK-CR26062801-a-djkhn6vpft3o](TASK-CR26062801-a-djkhn6vpft3o_分析探索深度判据与探索结论承载机制.md) [task] - 分析探索深度判据与探索结论承载机制
#### satisfies
- [DES-CR26062801-djkmcrjk40gc](DES-CR26062801-djkmcrjk40gc_探索深度判据三层优先级-硬规则与启发式.md) [design] - 探索深度判据：三层优先级 + 硬规则与启发式
- [TASK-CR26062801-i-djkmdw0v1uf4](TASK-CR26062801-i-djkmdw0v1uf4_更新-workflow-rulesmd需求分析步骤-协同演进规则-探索判据.md) [task] - 更新 workflow-rules.md：需求分析步骤 + 协同演进规则 + 探索判据

## Open Questions
- 探索深度的判据是否应该是硬规则还是启发式指引？
- 是否需要新增"探索总结"卡类型来承载探索结论？

