---
id: DES-CR26062801-djkmcrj23oo9
title: 需求分析阶段：扩展 index 模式 + 卡片任务协同演进
type: design
status: draft
importance: should
tags:
    - card-evolution
    - design-skill
    - requirement-analysis
links:
    - target: CONV-djdotphklxrt
      relation: constrains
    - target: PROP-CR26062801
      relation: belongs_to
    - target: REQ-CR26062801-djkhctmg6kke
      relation: satisfies
    - target: TASK-CR26062801-a-djkhn6v3gg2k
      relation: references
created: 2026-06-28T10:43:44.156894373Z
updated: 2026-06-28T10:43:44.157484172Z
source: CR26062801
domain: skill-design
---

# 需求分析阶段：扩展 index 模式 + 卡片任务协同演进

## Goal

让 design skill 在拿到用户模糊输入后，有系统化的需求分析方法论，同时卡片与分析任务可以协同演进而不是串行。

## Decision

**扩展 index 模式**，不新增独立模式。index 模式增加前置需求分析步骤：

1. 需求分析（识别核心目标、边界、隐含假设、风险点）
2. 拆解为主题条目
3. 写入 STR / 创建 REQ 卡

**卡片与任务协同演进**：分析任务不是关卡。创建分析任务后，每完成一个子分析立即 `card update` 更新对应卡片。卡片从 `draft` 转 `stable` 的条件是：所有关联分析任务 done + 所有 Open Questions 闭合。

## Rationale

- index 模式本质就是新需求入口，分析是其自然前置动作，不应新增模式破坏 7 模式结构
- 需求分析检查清单（目标、边界、假设、风险）提供思考框架但不强绑 Agent
- 协同演进符合真实工作流：分析是渐进的，结论逐步形成

## Constraints

- 不新增第八个模式，保持 mode selection 表结构
- 分析检查清单是启发式指引，不是硬规则
- `draft → stable` 转换判据写入 workflow-rules.md

## Impact

- `workflow-rules.md`：index 模式描述、新增协同演进口袋规则
- `card-templates.md`：可能不需要改动（现有模板已够用）

## Verification

- Agent 用一条模糊需求验证：能否产出结构化分析结论 → STR 条目 → REQ 卡
- 验证协同演进：创建一个含 3 个子问题的 analysis task，分步更新 REQ 卡，验证流程

## Follow-up Tasks

- 修改 workflow-rules.md index 模式描述
- 新增协同演进口袋规则
- 写测试验证新流程

## Links

### Outgoing

- [PROP-CR26062801](../../../../03-proposal/CR26062801_优化-design-skill-探索能力需求分析信息.md) [proposal] - 优化 design skill 探索能力：需求分析、信息探索来源与外部知识集成
- [CONV-djdotphklxrt](../../../../02-library/60-conventions/CONV-djdotphklxrt_design-skill-不写长设计文档.md) [convention] - Design SKILL 不写长设计文档
- [TASK-CR26062801-a-djkhn6v3gg2k](TASK-CR26062801-a-djkhn6v3gg2k_分析需求分析阶段的模式定位与卡片演进机制.md) [task] - 分析需求分析阶段的模式定位与卡片演进机制
- [REQ-CR26062801-djkhctmg6kke](REQ-CR26062801-djkhctmg6kke_需求分析方法从模糊输入到结构化需求.md) [requirement] - 需求分析方法：从模糊输入到结构化需求

### Incoming

- [TASK-CR26062801-i-djkmdw0v1uf4](TASK-CR26062801-i-djkmdw0v1uf4_更新-workflow-rulesmd需求分析步骤-协同演进规则-探索判据.md) [task] - 更新 workflow-rules.md：需求分析步骤 + 协同演进规则 + 探索判据
- [TASK-CR26062801-i-djkmdw0v1uf4](TASK-CR26062801-i-djkmdw0v1uf4_更新-workflow-rulesmd需求分析步骤-协同演进规则-探索判据.md) [task] - 更新 workflow-rules.md：需求分析步骤 + 协同演进规则 + 探索判据

