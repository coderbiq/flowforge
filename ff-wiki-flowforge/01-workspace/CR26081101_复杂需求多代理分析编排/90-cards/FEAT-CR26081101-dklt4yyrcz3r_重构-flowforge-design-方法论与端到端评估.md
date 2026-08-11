---
id: FEAT-CR26081101-dklt4yyrcz3r
title: 重构 flowforge-design 方法论与端到端评估
type: feature
status: done
importance: should
links:
    - target: CONV-djkddmqleudi
      relation: constrains
    - target: DEC-CR26081101-dkltaxm2rgcf
      relation: references
    - target: DES-djdocmn5evg2
      relation: references
    - target: FEAT-CR26081101-dklt4yui88kf
      relation: requires
    - target: FEAT-CR26081101-dklt4yvz37gy
      relation: requires
    - target: FEAT-CR26081101-dklt4yxctvsf
      relation: requires
    - target: PROP-CR26081101
      relation: belongs_to
created: 2026-08-11T03:51:19.597145236Z
updated: 2026-08-11T14:42:47.660377845+08:00
source: CR26081101
---

# 重构 flowforge-design 方法论与端到端评估

## Summary

把 `flowforge-design` 从一次性的四阶段流程重构为可选择复杂度、生成调查计划、多轮综合并最终进入 stage gate 的薄适配器，同时建立跨宿主的质量与成本评估。

## Motivation

现有 SKILL 没有说明如何验证用户论断、恢复现状、分析根因、规划调查或判断停止，导致昂贵 Analyst 承担无目标扫描。没有端到端评估也无法确认新编排是否降低成本而不损失设计质量。

## Design

### Key Decisions

- **SKILL.md 保持薄适配器**：继续低于 200 tokens，完整方法放 references，因为 description 和短流程负责命中而非承载全部知识。
- **先分级后编排**：简单需求直接设计；仅有多个独立证据缺口、超长上下文或跨领域事实时创建调查计划。
- **分析模式可组合**：支持 Claim Review、Current-State、Root Cause、Boundary Analysis、New Design、Implementation Design，不强制全部执行。
- **证据和解释分离**：区分 Fact、Interpretation、Assumption、Unknown、Risk，并标记关键论断支持程度。
- **以语义就绪度推进 stage**：检查目标、现状、决策、实施边界、验证和剩余未知项，而非只检查格式。
- **同时评估结果和路径**：度量质量、角色调用、重复调查、恢复能力和用户可见性。

### Architecture

设计循环扩为 `intake → frame → seed → provisional-design → assess → investigate/synthesize* → clarify → design gate → plan gate`。Analyst 首轮必须产出目标/非目标、FEATURE 拆分、半成品方案、证据缺口、2–6 个独立 work 和 re-entry 条件；如果无需调查则直接进入 gate。

综合轮逐条将 FIND 判定为 accepted/rejected/conflicting，更新 FEATURE/DEC，并显式发布新 revision 或完成设计。只有无 required work、无用户问题、设计门禁通过后才能 planned。

references 至少包含 `analysis-workflow.md`、`delegation-brief.md`、`evidence-rules.md`、`readiness-gates.md` 和 `card-templates.md`。所有正文由 Agent 直接编辑；CLI 仅创建骨架、seal 事件、校验、stage 和索引。

评估矩阵覆盖简单单 FEATURE、复杂跨模块、外部资料、证据冲突、调查阻塞、用户决策、计划修订、session 丢失、SQLite 重建、实施反馈和 Codex/OpenCode 差异。

### Alternatives Considered

- 完整方法塞入 SKILL.md：降低命中准确性并违反 token 限制。
- 固定专家流水线：不同需求需要的分析模式不同。
- 只评最终答案：无法发现隐藏重复调用和恢复失败。
- 使用 `--body` 写复杂正文：格式化和局部编辑问题已被实践证伪。

## Constraints

- description 必须与 implement、feedback、curate、review 互斥。
- design SKILL 不执行产品代码。
- 评估使用 model profile 并记录实际映射，不绑定单一模型 ID。
- 文件变更后必须 validate；实现阶段必须运行 `go test ./internal/...`。

## Implementation Plan

### Step 1: 重写 flowforge-design 薄入口

<!-- step-status: done -->

- **Goal**: 正确触发复杂分析循环并保持 SKILL.md 小于 200 tokens。
- **Files**: `assets/skills/flowforge-design/SKILL.md`, `.agents/skills/flowforge-design/SKILL.md`
- **Approach**: 压缩为 intake→frame→iterate→gate，明确直接编辑正文与 CLI 结构化操作边界。
- **Edge Cases**: 单卡创建、feedback、implementation、archive 不应误触发。
- **Dependencies**: 前三张 FEATURE 契约稳定。
- **Parallel**: no。
- **Verification**: token 计数、description 反例和 init/sync 测试通过。

### Step 2: 编写方法与交接 references

<!-- step-status: done -->

- **Goal**: 提供可执行分析方法而非描述性原则。
- **Files**: `assets/skills/flowforge-design/references/*.md`
- **Approach**: 编写模式选择、计划事件、调查 brief、FIND 写回、冲突、用户决策、readiness 和恢复 walkthrough。
- **Edge Cases**: 无证据命中、计划过度拆分、follow-up 超范围、调查无新增信息。
- **Dependencies**: Step 1。
- **Parallel**: yes。
- **Verification**: 每条规则映射到 CLI 或 Artifact 动作，walkthrough 不依赖隐含会话。

### Step 3: 更新阶段门禁和上下文输出

<!-- step-status: done -->

- **Goal**: designed/planned 反映语义就绪度，并为角色提供最小上下文。
- **Files**: `internal/command/card_evolve.go`, `internal/command/context.go`, 相关测试
- **Approach**: 检查复杂分析章节、证据接受状态、Verification 映射和 active plan；按角色裁剪 context。
- **Edge Cases**: 旧 FEATURE、无需调查、active work 未结束、用户接受风险。
- **Dependencies**: Step 2 和 analysis status。
- **Parallel**: no。
- **Verification**: gate 正反例与上下文 golden tests 通过。

### Step 4: 建立端到端编排评估

<!-- step-status: done -->

- **Goal**: 证明新方案在质量、成本代理指标、恢复和宿主兼容上优于单 Analyst。
- **Files**: `internal/orchestration/*_test.go`, `internal/command/*_test.go`, `tests/orchestration/`
- **Approach**: 场景 fixture 记录角色调用、revision、事件流、Artifact 和验证；比较简单直达与复杂多轮。
- **Edge Cases**: Investigator 失败、重复 return、Coordinator 重启、SQLite 删除、无网络。
- **Dependencies**: 所有前置步骤。
- **Parallel**: no。
- **Verification**: 场景可重复，失败可定位到路由、事件、权限或门禁。

## Verification

- SKILL.md 小于 200 tokens 且命中边界无重叠。
- 简单设计不会产生多余 Subagent；复杂设计形成可恢复的多轮分析。
- 删除 SQLite 或丢失 Agent 会话后仍能继续。
- Codex/OpenCode 均能完成 plan→investigate→synthesize→planned。
- 旗舰调用集中在规划和综合检查点，资料调查使用较低成本 profile。

## History

- 2026-08-11：吸收 deep-analyst、Ledger 和专家团调研，形成 FlowForge 原生方法论与评估设计。
- 2026-08-11T14:41:36+08:00 | progress | 压缩 flowforge-design 为 112 words，并通过 sync --adopt 部署 SKILL 与四个分析 reference
- 2026-08-11T14:41:36+08:00 | progress | 完成可恢复的复杂分析循环、调查 brief、证据规则、readiness 与恢复 walkthrough
- 2026-08-11T14:41:36+08:00 | progress | 增加复杂 FEATURE designed/planned 语义门禁与 Analyst/Coordinator/Investigator/Executor 角色上下文
- 2026-08-11T14:41:36+08:00 | progress | 增加简单直达、复杂多轮、失败、用户决策、恢复及 Codex/OpenCode profile/enforcement 场景评估
- 2026-08-11T14:42:47+08:00 | progress | 实现与定向验证完成；单卡校验通过，risk-review 四步均无需独立审查；全库保留 25 个既有历史链接错误

## Links

### Outgoing

- [PROP-CR26081101](../../../03-proposal/CR26081101_复杂需求多代理分析编排.md) [proposal] - 复杂需求多代理分析编排
- [CONV-djkddmqleudi](../../../02-library/60-conventions/CONV-djkddmqleudi_agentsmd-flowforge-区块包裹部署规范.md) [convention] - AGENTS.md FLOWFORGE 区块包裹部署规范
#### references
- [DEC-CR26081101-dkltaxm2rgcf](DEC-CR26081101-dkltaxm2rgcf_正文直接编辑与-cli-结构化操作边界.md) [decision] - 正文直接编辑与 CLI 结构化操作边界
- [DES-djdocmn5evg2](../../../02-library/30-designs/DES-djdocmn5evg2_flow-forge-agent-first-human-readable-双层设计.md) [design] - FlowForge Agent-First Human-Readable 双层设计
#### requires
- [FEAT-CR26081101-dklt4yui88kf](FEAT-CR26081101-dklt4yui88kf_定义复杂需求的迭代分析协议与-artifact-边界.md) [feature] - 定义复杂需求的迭代分析协议与 Artifact 边界
- [FEAT-CR26081101-dklt4yvz37gy](FEAT-CR26081101-dklt4yvz37gy_扩展-journal-调度事件sqlite-派生视图与-cli.md) [feature] - 扩展 Journal 调度事件、SQLite 派生视图与 CLI
- [FEAT-CR26081101-dklt4yxctvsf](FEAT-CR26081101-dklt4yxctvsf_重构-subagent-角色拓扑prompt-与宿主适配.md) [feature] - 重构 Subagent 角色拓扑、Prompt 与宿主适配

## Open Questions

None。

## Dependencies

- FEAT-CR26081101-dklt4yui88kf：分析协议。
- FEAT-CR26081101-dklt4yvz37gy：Journal/CLI 调度能力。
- FEAT-CR26081101-dklt4yxctvsf：角色与宿主协议。
