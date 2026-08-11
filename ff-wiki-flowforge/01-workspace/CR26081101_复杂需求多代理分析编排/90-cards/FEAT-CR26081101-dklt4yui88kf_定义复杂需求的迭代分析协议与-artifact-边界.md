---
id: FEAT-CR26081101-dklt4yui88kf
title: 定义复杂需求的迭代分析协议与 Artifact 边界
type: feature
status: done
importance: should
links:
    - target: DEC-CR26081101-dkltaxm2rgcf
      relation: references
    - target: DES-djdotwt01934
      relation: references
    - target: PROP-CR26081101
      relation: belongs_to
created: 2026-08-11T03:51:19.33994262Z
updated: 2026-08-11T05:54:51.876336097Z
source: CR26081101
---

# 定义复杂需求的迭代分析协议与 Artifact 边界

## Summary

建立复杂需求的多轮分析协议，使需求理解、证据调查、半成品方案和后续问题可以跨 Agent、跨会话持续演进。协议明确 Design Analyst、Coordinator、Investigator 与正式 Artifact 的责任边界。

## Motivation

当前设计把调查、架构分析和计划全部交给昂贵 Design Analyst，既增加成本，也缺少可恢复的中间状态。若没有明确循环和 Artifact 边界，Subagent 只能依赖临时返回文本交流，轻量 Coordinator 还需要自行解释历史和做高难度判断。

## Design

### Key Decisions

- **强规划、弱调度**：Design Analyst 制定和修订分析计划，因为复杂度判断和调查方向选择需要高能力模型；Coordinator 只执行已登记计划，避免轻量模型承担架构推理。
- **单层委托**：所有 Worker 由 Coordinator 调用，因为 Codex/OpenCode 对嵌套委托支持不一致，且用户需要看到后台调度过程。
- **角色不按信息来源拆分**：使用一个通用 Investigator，因为代码、Library、网络和日志只是工具来源，按来源建角色会导致角色爆炸。
- **分析是循环而非流水线**：采用 frame → investigate → synthesize → gate 循环，因为每轮证据都可能修改假设、方案和下一轮问题。
- **半成品必须持久化**：FEATURE 保存当前理解和 Working Design，因为不同 Agent 不能只依赖会话摘要恢复复杂分析。

### Architecture

分析状态依次为 `SEED → PLAN → INVESTIGATE → SYNTHESIZE → GATE`，实施或评审发现设计缺口时进入 `REOPEN`。Analyst 在首轮拆分 1–5 张 FEATURE，写入当前假设和调查计划；Coordinator 派发 ready 项；Investigator 收集证据；所有必需结果返回后 Analyst 重新进入，接受、拒绝或标记冲突，再决定下一 revision 或 stage gate。

FEATURE 至少维护 `Objective`、`Current Understanding`、`Evidence`、`Working Design`、`Rejected or Revised Assumptions`、`Open Questions` 与 `Next Investigation`。跨卡决策使用 DEC，调查证据写入 work item 指定的 FIND。Investigator 只能编辑该 FIND 的 Evidence、Source、Impact 和 Open Questions，并通过 CLI 校验后返回；不得修改 FEATURE、DEC 或产品代码。

全局不变量：Coordinator 不创造调查方向；Investigator 不做最终设计判断；未经 Analyst 接受的 follow-up 不得自动派发；设计事实以 Artifact 为准；调度事实以 Journal 为准；每轮都有预算、完成条件和 Analyst re-entry 条件；产品、兼容、迁移、安全和冲突目标由用户决定。

### Alternatives Considered

- Design Analyst 递归委派：宿主兼容和用户可见性差。
- Code/Library/Web 多 Scout：角色会随数据源增长。
- 完整共享聊天：上下文成本持续增长且难以重建。
- 一次性收集后设计：无法处理证据驱动的方案修正。

## Constraints

- 不修改现有 FEATURE stage schema；analysis revision 是协作概念。
- 简单需求可以由 Analyst 直接设计，不强制创建调查计划。
- 正文由 Agent 直接编辑，CLI 只负责结构、校验和门禁。
- 调查重复执行必须安全，不得依赖存活的宿主 session。

## Implementation Plan

### Step 1: 定义分析循环和 Artifact 内容契约

<!-- step-status: done -->

- **Goal**: 形成角色 Prompt、SKILL 和 CLI 共用的状态机与权威边界。
- **Files**: `docs/analysis-orchestration.md`, `docs/architecture.md`
- **Approach**: 文档化状态、责任、re-entry、停止条件和一轮/多轮 walkthrough。
- **Edge Cases**: 简单需求无需调查；调查推翻目标；planned FEATURE 被实施反馈重开。
- **Dependencies**: None。
- **Parallel**: no。
- **Verification**: walkthrough 能唯一判断当前状态、下一角色和权威数据位置。

### Step 2: 扩展复杂分析 Artifact 校验

<!-- step-status: done -->

- **Goal**: 保证当前理解、证据、半成品方案和开放问题可恢复。
- **Files**: `internal/core/validate.go`, `internal/core/validate_test.go`, `assets/skills/flowforge-design/references/card-templates.md`
- **Approach**: 对启用复杂分析的 FEATURE/FIND 检查约定章节、证据来源和 work ID；正文保持自由格式。
- **Edge Cases**: 历史卡没有新章节；章节只有 TBD；简单卡不应被强制升级。
- **Dependencies**: Step 1。
- **Parallel**: no。
- **Verification**: 新旧卡兼容、缺失章节告警和引用校验测试通过。

## Verification

- 简单请求不会产生无意义调查轮次。
- 两轮以上分析可仅凭 Artifact 与 Journal 恢复。
- Coordinator 无需解释证据即可判断派发或重新调用 Analyst。
- 实施反馈可重开分析并保留旧决策和证据链。

## History

- 2026-08-11：固定强 Analyst、弱 Coordinator、单层委托和 Artifact-first 原则。
- 2026-08-11：实施时发现 Proposal health 依赖尚未实现的 Journal 状态查询，与 F2 Step 4 重复；合并到 F2 以消除循环依赖。
- 2026-08-11T13:43:52+08:00 | progress | 完成复杂分析状态机、角色/权威边界、Artifact 契约及单轮/多轮/reopen walkthrough 文档
- 2026-08-11T13:46:50+08:00 | progress | 增加复杂分析 opt-in Artifact 校验：FEATURE/FIND 必需章节、非 TBD 内容、稳定 work ID、证据 FIND 引用与模板说明；core 测试通过
- 2026-08-11T13:48:50+08:00 | progress | 移除原 Step 3；Proposal health 分析一致性检查合并到 FEAT-CR26081101-dklt4yvz37gy Step 4，以消除对 Journal 状态查询的循环依赖
- 2026-08-11T13:54:51+08:00 | progress | 完成复杂分析协议、架构文档、FEATURE/FIND opt-in 校验和模板；Proposal health 合并到 F2 Step 4 以消除循环依赖

## Links

### Outgoing

- [PROP-CR26081101](../../../03-proposal/CR26081101_复杂需求多代理分析编排.md) [proposal] - 复杂需求多代理分析编排
#### references
- [DEC-CR26081101-dkltaxm2rgcf](DEC-CR26081101-dkltaxm2rgcf_正文直接编辑与-cli-结构化操作边界.md) [decision] - 正文直接编辑与 CLI 结构化操作边界
- [DES-djdotwt01934](../../../02-library/30-designs/DES-djdotwt01934_design-skill-主流程.md) [design] - Design SKILL 主流程

### Incoming

#### requires
- [FEAT-CR26081101-dklt4yvz37gy](FEAT-CR26081101-dklt4yvz37gy_扩展-journal-调度事件sqlite-派生视图与-cli.md) [feature] - 扩展 Journal 调度事件、SQLite 派生视图与 CLI
- [FEAT-CR26081101-dklt4yxctvsf](FEAT-CR26081101-dklt4yxctvsf_重构-subagent-角色拓扑prompt-与宿主适配.md) [feature] - 重构 Subagent 角色拓扑、Prompt 与宿主适配
- [FEAT-CR26081101-dklt4yyrcz3r](FEAT-CR26081101-dklt4yyrcz3r_重构-flowforge-design-方法论与端到端评估.md) [feature] - 重构 flowforge-design 方法论与端到端评估

## Open Questions

None。

## Dependencies

None。
