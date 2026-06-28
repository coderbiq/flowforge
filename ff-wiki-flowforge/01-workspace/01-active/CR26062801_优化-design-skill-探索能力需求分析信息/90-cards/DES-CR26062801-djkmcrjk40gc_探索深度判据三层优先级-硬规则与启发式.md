---
id: DES-CR26062801-djkmcrjk40gc
title: 探索深度判据：三层优先级 + 硬规则与启发式
type: design
status: draft
importance: should
tags:
    - depth-criterion
    - design-skill
    - exploration
links:
    - target: CONV-djdovkpot49d
      relation: constrains
    - target: PROP-CR26062801
      relation: belongs_to
    - target: REQ-CR26062801-djkhctmvwq80
      relation: satisfies
    - target: TASK-CR26062801-a-djkhn6vpft3o
      relation: references
created: 2026-06-28T10:43:44.186814011Z
updated: 2026-06-28T10:43:44.187733726Z
source: CR26062801
domain: skill-design
---

# 探索深度判据：三层优先级 + 硬规则与启发式

## Goal

明确 design skill 的探索来源、深度判据（何时结束探索）和结论承载方式。

## Decision

**三层探索来源**（由近到远）：

1. FlowForge Library（最高优先级）→ `library suggest` / `card search --scope library`
2. 外部知识库（配置的 knowledge_sources）→ Agent 文件读取
3. 项目源代码（信息参照，非规范性依据）

**探索深度判据**：

| 类型 | 判据 |
|------|------|
| 硬规则 | 需求卡 Open Questions 中仍有阻塞当前设计决策的未决问题时，探索不充分 |
| 硬规则 | 分析任务指定输入源未检查完，任务不可 done |
| 启发式 | 按 Library→外部→源码顺序查完，三者均无新发现 → 可停止 |
| 启发式 | 连续两轮探索无新信息 + 无新搜索词 → 已充分 |

**结论承载**：不需要新卡类型。

- requirement：闭合或保留的 Open Questions
- finding：可复用的探索发现
- log："查了什么、命中/未命中"审计记录
- design：综合探索结论后的最终设计决策

## Rationale

- 三层来源优先级保证探索效率（先查最直接约束层，再补充领域知识）
- 硬规则防 Agent 偷懒，启发式防无限探索
- 不新增卡类型保持系统简洁

## Constraints

- 三层来源中，Library 和外部源是"规范性"依据，源码只提供"事实"信息
- Agent 必须在 log 中记录探索了哪些源、命中/未命中什么

## Impact

- `workflow-rules.md`：新增探索判据区间
- `library-discovery.md`：扩展为三层探索模型

## Verification

- 模拟一个需求：Library 无命中 → 外部源有相关文档 → 源码确认现状。验证 Agent 是否正确按优先级探索
- 验证探索不会无限循环

## Follow-up Tasks

- 修改 workflow-rules.md 新增探索判据区间
- 修改 library-discovery.md 扩展为三层模型

## Links

### Outgoing

- [PROP-CR26062801](../../../../03-proposal/CR26062801_优化-design-skill-探索能力需求分析信息.md) [proposal] - 优化 design skill 探索能力：需求分析、信息探索来源与外部知识集成
- [CONV-djdovkpot49d](../../../../02-library/60-conventions/CONV-djdovkpot49d_每轮-design-必须输出固定格式.md) [convention] - 每轮 design 必须输出固定格式
- [TASK-CR26062801-a-djkhn6vpft3o](TASK-CR26062801-a-djkhn6vpft3o_分析探索深度判据与探索结论承载机制.md) [task] - 分析探索深度判据与探索结论承载机制
- [REQ-CR26062801-djkhctmvwq80](REQ-CR26062801-djkhctmvwq80_设计思维逻辑探索分析决策的推理链路.md) [requirement] - 设计思维逻辑：探索→分析→决策的推理链路

### Incoming

#### 
- [TASK-CR26062801-i-djkmdw0v1uf4](TASK-CR26062801-i-djkmdw0v1uf4_更新-workflow-rulesmd需求分析步骤-协同演进规则-探索判据.md) [task] - 更新 workflow-rules.md：需求分析步骤 + 协同演进规则 + 探索判据
- [TASK-CR26062801-i-djkmdw1gfm4q](TASK-CR26062801-i-djkmdw1gfm4q_更新-library-discoverymd三层探索模型-ab-策略-嵌入格式.md) [task] - 更新 library-discovery.md：三层探索模型 + A/B 策略 + 嵌入格式
#### related
- [TASK-CR26062801-i-djkmdw0v1uf4](TASK-CR26062801-i-djkmdw0v1uf4_更新-workflow-rulesmd需求分析步骤-协同演进规则-探索判据.md) [task] - 更新 workflow-rules.md：需求分析步骤 + 协同演进规则 + 探索判据
- [TASK-CR26062801-i-djkmdw1gfm4q](TASK-CR26062801-i-djkmdw1gfm4q_更新-library-discoverymd三层探索模型-ab-策略-嵌入格式.md) [task] - 更新 library-discovery.md：三层探索模型 + A/B 策略 + 嵌入格式

