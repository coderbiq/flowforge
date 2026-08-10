---
id: FEAT-CR26081001-dkkymgqx8c63
title: 重构 Subagent 路由与职责模型
type: feature
status: done
importance: should
links:
    - target: FEAT-CR26081001-dkkymba7v03t
      relation: requires
    - target: PROP-CR26081001
      relation: belongs_to
created: 2026-08-10T03:56:36.089015227Z
updated: 2026-08-10T13:05:33.283570085Z
source: CR26081001
---

# 重构 Subagent 路由与职责模型

## Summary

根据 FlowForge Proposal 生命周期重新设计 subagent 协作，不继承 AOK 固定四角色。首版采用 Coordinator、Design Analyst、Executor、按风险调用的 Reviewer，并统一角色、Skill、模型 profile、权限与 FlowForge 门禁。

## Motivation

AOK 的职责隔离有价值，但 Task Assessor 的部分职责可由 FlowForge CLI 确定性完成，完整 handoff 又与最小上下文原则冲突。不重构会增加不必要的模型调用、重复状态和角色边界。

## Design

### Key Decisions

- Coordinator 是唯一交互与委托角色，首版不直接修改产品代码；因为这样可去除自我授权问题和常驻 Task Assessor。
- Design Analyst 可编辑 Proposal/FEATURE/DEC/FIND 但不可修改产品代码；因为设计应直接落入正式 artifact。
- Executor 只实施 planned FEATURE 的指定 Step；因为实施者不能静默重设计。
- Reviewer 逻辑独立、按风险调用，首版可复用 Analyst 模型配置；因为不是所有任务都值得语义复核成本。
- deterministic preflight 优先于模型 assessment；只有范围或风险有语义歧义时才临时调用只读 assessment 能力。
- Skill 由角色默认绑定，任务可叠加领域 Skill；实际 provider/model 由宿主解析 profile。

### Architecture

Coordinator 读取 Proposal Journal 最近记录和 artifact 状态后路由。Design Analyst 默认使用 `flowforge-design`，把设计写入 artifact 并通过 designed/planned 门禁。Executor 默认使用 `flowforge-implement`，从 `context feature --step` 获取执行切片；遇到新设计决策、范围扩大、计划过期或无法验证时停止并记录 Journal。Reviewer 对照 FEATURE、最终 diff 和验证证据做一致性复核，必要时触发 feedback 或阶段回退。

路由映射：Proposal 设计交给 Design Analyst/high-capability/`flowforge-design`；planned Step 交给 Executor/tool-capable/`flowforge-implement`；设计缺口回到 Analyst 并结合 `flowforge-feedback`；高风险实现交给 Reviewer/high-capability-read-only；归档使用 `flowforge-curate`。Agent 返回状态不等于持久化 stage，最终仍由 CLI 门禁决定。

### Alternatives Considered

- 原样保留 AOK 四角色：拒绝，因为所有 mutation 经过 Assessor 成本过高。
- 每个 Skill 一个 subagent：拒绝，因为 Skill 是工作方法，不等于权限职责。
- Coordinator 处理简单修改：首版拒绝以换取清晰边界，未来按数据评估。
- 所有实施强制 Reviewer：拒绝，因为机械小改动收益不足。

## Constraints

- worker 不递归委托、不直接向用户提问，统一返回 Coordinator。
- Analyst 不修改产品代码；Executor 不做新架构或范围决策；Reviewer 只读产品代码。
- Executor 仅在 FEATURE planned 且用户有明确实施意图时启动。
- planned 表示文档完备，不隐含用户授权。
- 角色默认 Skill 不得与现有 Skill description 冲突，模型 profile 不包含实际 model ID。
- Prompt 门禁不能替代 CLI stage gate 或宿主权限。

## Implementation Plan

### Step 1: 定义角色与路由中立 schema

<!-- step-status: done -->

- **Goal**: 表达角色职责、委托拓扑、模型 profile、默认 Skill 和进入/停止条件。
- **Files**: `internal/` orchestration 配置包、中立角色源、校验测试
- **Approach**: 数据驱动定义并校验 Coordinator 唯一、worker 无递归委托、Skill/profile 存在和能力一致。
- **Edge Cases**: 禁用 Reviewer、自定义角色、宿主缺少模型 profile。
- **Dependencies**: 总体架构边界。
- **Parallel**: no
- **Verification**: fixture 覆盖合法拓扑和非法关系。

### Step 2: 重写角色协议与 FlowForge Skill

<!-- step-status: done -->

- **Goal**: 四个角色均使用 Journal 与正式 artifact 交接。
- **Files**: 中立角色 Markdown、`assets/skills/flowforge-*`、review 协议
- **Approach**: 压缩 AOK 规则并收紧 Executor 修改 FEATURE 的权限，只允许进度、证据和无设计含义事实。
- **Edge Cases**: 设计缺口、用户改目标、FEATURE 执行前变化、测试无法运行。
- **Dependencies**: Step 1 与 Journal CLI。
- **Parallel**: no
- **Verification**: 协议测试验证角色允许/禁止行为和结束日志。

### Step 3: 实现 preflight 与风险复核策略

<!-- step-status: done -->

- **Goal**: 用 CLI 事实决定能否实施及是否需要 Reviewer。
- **Files**: `internal/command/` context/validate 扩展、风险策略与测试
- **Approach**: 检查 stage、Step、依赖、用户意图、工作区和验证要求；按公共 API、迁移、安全、跨模块、scope drift 等信号选择复核级别。
- **Edge Cases**: Files 为 glob、实际 diff 超范围、用户跳过可选复核。
- **Dependencies**: Step 2 停止状态与 Reviewer 契约。
- **Parallel**: no
- **Verification**: 表驱动测试覆盖可执行、blocked、stale、semantic review 和 human-required。

## Verification

- Coordinator 无产品代码写权限仍能驱动完整流程。
- Analyst 正式结果落入 artifact，Journal 只保存摘要引用。
- Executor 遇到新设计、scope expansion、stale plan 或 verification failure 时停止。
- Reviewer 仅按风险调用，并能识别测试通过但偏离设计的实现。
- 角色、Skill、模型 profile 与门禁映射唯一且可检查。

## History

- 2026-08-10 [decision] 不继承 AOK 固定四角色。
- 2026-08-10 [decision] 首版取消常驻 Task Assessor，Coordinator 不直接修改产品代码。
- 2026-08-10 [decision] Reviewer 逻辑独立、按风险调用。
- 2026-08-10T16:05:20+08:00 | progress | 完成中立 subagent policy：定义 Coordinator、Design Analyst、Executor、可选 Reviewer 的职责、能力、模型 profile、默认 Skill、进入条件、停止条件与委托拓扑；校验唯一 Coordinator、worker 无递归委托、Skill/profile/能力一致性。orchestration 包测试通过；go test ./internal/... 仍仅受既有 Proposal/Structure 测试失败阻断。
- 2026-08-10T20:43:55+08:00 | progress | 完成 FlowForge 原生角色协议：新增 Coordinator、Design Analyst、Executor、Reviewer 协议与 flowforge-review Skill；移除 AOK 长 handoff/session-state 依赖；Executor 禁止改写设计、约束和计划，遇 design_gap/scope_expanded/plan_stale/verification_failed 时转 feedback/design。聚焦 policy/部署/Journal 测试通过；全量测试仍受既有 Proposal/Structure 用例阻断。
- 2026-08-10T21:05:33+08:00 | progress | 完成 execution preflight 与 risk review：context preflight 检查 planned、实施意图、Step 字段、开放问题和 requires 依赖；context risk-review 基于完成 Step 和 Git diff 路由独立复核。

## Links

### Outgoing

- [PROP-CR26081001](../../../03-proposal/CR26081001_subagent-协作机制整合.md) [proposal] - Subagent 协作机制整合
- [FEAT-CR26081001-dkkymba7v03t](FEAT-CR26081001-dkkymba7v03t_引入-proposal-journal-协作日志.md) [feature] - 引入 Proposal Journal 协作日志

### Incoming

#### requires
- [FEAT-CR26081001-dkkylwkup2d3](FEAT-CR26081001-dkkylwkup2d3_整合-proposal-journal-与-subagent-协作层.md) [feature] - 整合 Proposal Journal 与 Subagent 协作层
- [FEAT-CR26081001-dkkymf4v3m1u](FEAT-CR26081001-dkkymf4v3m1u_生成并部署宿主-subagent-adapter.md) [feature] - 生成并部署宿主 Subagent Adapter

## Open Questions

None

## Dependencies

- Proposal Journal
- FEATURE stage 与 `context feature` 门禁
