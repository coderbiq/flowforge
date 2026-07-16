---
id: TASK-CR26062801-i-djkmdw0v1uf4
title: 更新 workflow-rules.md：需求分析步骤 + 协同演进规则 + 探索判据
type: task
status: done
importance: should
links:
    - target: DES-CR26062801-djkmcrj23oo9
      relation: related
    - target: DES-CR26062801-djkmcrj23oo9
      relation: implements
    - target: DES-CR26062801-djkmcrjk40gc
      relation: related
    - target: DES-CR26062801-djkmcrjk40gc
      relation: implements
    - target: PROP-CR26062801
      relation: belongs_to
    - target: REQ-CR26062801-djkhctmg6kke
      relation: satisfies
    - target: REQ-CR26062801-djkhctmvwq80
      relation: satisfies
created: 2026-06-28T10:45:12.305323292Z
updated: 2026-06-28T18:50:15.093317105+08:00
source: CR26062801
---

# 更新 workflow-rules.md：需求分析步骤 + 协同演进规则 + 探索判据

## Goal
将设计结论写入 design skill 的 workflow-rules.md 参考文件中。

## Inputs
- DES-CR26062801-djkmcrj23oo9（需求分析阶段设计）
- DES-CR26062801-djkmcrjk40gc（探索深度判据设计）
- 现有 workflow-rules.md

## Deliverables
- workflow-rules.md：index 模式增加需求分析步骤
- workflow-rules.md：新增卡片与任务协同演进口袋规则
- workflow-rules.md：新增探索深度判据区间（三层来源、硬规则+启发式）

## Acceptance
- index 模式描述包含前置需求分析步骤（目标、边界、假设、风险）
- 协同演进规则覆盖：任务创建不阻塞卡片更新、draft→stable 条件
- 探索判据包含三来源优先级、硬规则、启发式指导

## Out of Scope
- 修改 library-discovery.md
- 修改 card-templates.md

## Read Before Work
- DES-CR26062801-djkmcrj23oo9
- DES-CR26062801-djkmcrjk40gc
- CONV-djdovkpot49d

## Links

### Outgoing

- [PROP-CR26062801](../../../../03-proposal/CR26062801_优化-design-skill-探索能力需求分析信息.md) [proposal] - 优化 design skill 探索能力：需求分析、信息探索来源与外部知识集成
#### implements
- [DES-CR26062801-djkmcrj23oo9](DES-CR26062801-djkmcrj23oo9_需求分析阶段扩展-index-模式-卡片任务协同演进.md) [design] - 需求分析阶段：扩展 index 模式 + 卡片任务协同演进
- [DES-CR26062801-djkmcrjk40gc](DES-CR26062801-djkmcrjk40gc_探索深度判据三层优先级-硬规则与启发式.md) [design] - 探索深度判据：三层优先级 + 硬规则与启发式
#### related
- [DES-CR26062801-djkmcrj23oo9](DES-CR26062801-djkmcrj23oo9_需求分析阶段扩展-index-模式-卡片任务协同演进.md) [design] - 需求分析阶段：扩展 index 模式 + 卡片任务协同演进
- [DES-CR26062801-djkmcrjk40gc](DES-CR26062801-djkmcrjk40gc_探索深度判据三层优先级-硬规则与启发式.md) [design] - 探索深度判据：三层优先级 + 硬规则与启发式
#### satisfies
- [REQ-CR26062801-djkhctmg6kke](REQ-CR26062801-djkhctmg6kke_需求分析方法从模糊输入到结构化需求.md) [requirement] - 需求分析方法：从模糊输入到结构化需求
- [REQ-CR26062801-djkhctmvwq80](REQ-CR26062801-djkhctmvwq80_设计思维逻辑探索分析决策的推理链路.md) [requirement] - 设计思维逻辑：探索→分析→决策的推理链路

## Summary

已更新 workflow-rules.md：index 模式增加需求分析步骤、新增卡片任务协同演进规则、新增探索深度判据区间

