---
id: FEAT-CR26081001-dkkylwkup2d3
title: 整合 Proposal Journal 与 Subagent 协作层
type: feature
status: done
importance: should
links:
    - target: FEAT-CR26081001-dkkymba7v03t
      relation: requires
    - target: FEAT-CR26081001-dkkymf4v3m1u
      relation: requires
    - target: FEAT-CR26081001-dkkymgqx8c63
      relation: requires
    - target: PROP-CR26081001
      relation: belongs_to
created: 2026-08-10T03:55:52.186315276Z
updated: 2026-08-10T13:06:35.173752127Z
source: CR26081001
---

# 整合 Proposal Journal 与 Subagent 协作层

## Summary

将 AOK 中经过验证的职责隔离、停止语义和宿主适配思想重写进 FlowForge。FlowForge 负责正式文档、确定性门禁和部署生命周期，宿主负责模型调用、session、权限与 sandbox。

## Motivation

FlowForge 已能管理 FEATURE 生命周期，却未统一定义任务使用哪个角色、模型和 Skill；AOK 定义了角色协作，却缺少持久化事实、CLI 门禁和安装升级。继续分离会形成两套重叠机制。

## Design

### Key Decisions

- AOK 作为研究输入重写，不原样迁移固定四角色与完整 handoff；因为这些机制不适合 FlowForge 的 artifact 模型。
- FlowForge 是唯一持久化事实与状态来源；因为现有 FEATURE、Step 和 Verification 已有 CLI 门禁。
- subagent 是可选执行层而非自建 runtime；因为 OpenCode/Codex 已负责模型、session、权限和 sandbox。
- 中立策略与宿主 adapter 分层；因为角色、Skill 和门禁应跨宿主稳定，权限字段则属于宿主。
- 首版不实现路径级硬授权；必须区分 Prompt 规则、artifact gate 与 runtime enforcement。

### Architecture

架构分为正式 artifact、Proposal Journal、中立协作策略、宿主 adapter 四层。Coordinator 根据用户意图和 FlowForge 状态选择 worker；worker 将结果写入正式 artifact，并向 Journal 追加摘要、引用和下一步；CLI 维护卡片、阶段、进度、校验及 managed assets。

AOK 保留中心化委托、worker 不递归、分析/实施隔离、Executor 停止条件、独立复核、模型能力分层、中立源生成和 runtime smoke。删除固定四角色、强制 Assessor、完整输出复制、九段 session state、长授权文档和独立 installer。

### Alternatives Considered

- 保持两个独立项目：拒绝，因为模型、Skill、门禁与安装属于同一工作流决策。
- AOK 原样放入 `assets/`：拒绝，因为会冻结尚未验证的角色和宿主边界。
- FlowForge 自建 orchestration runtime：拒绝，因为重复宿主能力并扩大范围。

## Constraints

- 不修改 FEATURE schema 来承载宿主 session 或模型输出。
- 不把宿主权限字段写入中立策略，不把 Prompt 授权描述成硬 enforcement。
- 正式设计、进度和验证以 artifact 与 CLI 状态为准，Journal 无业务状态效力。
- 未安装 adapter 时现有单 Agent 工作流必须保持可用。
- 首版 OpenCode 优先、项目级显式安装，AOK 仅作为迁移参考。

## Implementation Plan

### Step 1: 定义整合后的中立协作模型

<!-- step-status: done -->

- **Goal**: 形成角色、模型 profile、Skill、门禁和 adapter 的单一中立定义。
- **Files**: `internal/` orchestration 配置与校验、`docs/` 架构说明
- **Approach**: 从 AOK 提取通用字段，去除固定角色数量和宿主配置。
- **Edge Cases**: 未安装 adapter 时现有工作流保持可用。
- **Dependencies**: Journal 与角色模型设计稳定。
- **Parallel**: no
- **Verification**: 测试未知角色、递归委托、未知 Skill/profile 和宿主字段泄漏。

### Step 2: 迁移 AOK 精华并移除重复机制

<!-- step-status: done -->

- **Goal**: 将有效协议改写为 FlowForge 原生职责和停止语义。
- **Files**: `assets/skills/`、中立角色源、协议测试
- **Approach**: 映射 stale、scope expansion、verification failed 和 review，删除完整 handoff 与独立授权文档。
- **Edge Cases**: Prompt 保证的能力必须明确标为软约束。
- **Dependencies**: Step 1 中立模型。
- **Parallel**: no
- **Verification**: 每个停止状态都有明确 FlowForge 后续动作。

### Step 3: 集成 Journal、路由和 adapter 生命周期

<!-- step-status: done -->

- **Goal**: Proposal 可从 Journal 恢复，并由 CLI 部署宿主角色配置。
- **Files**: `internal/command/`、`internal/core/`、`assets/`、测试
- **Approach**: 串联三个子 FEATURE，保持无 adapter 兼容路径。
- **Edge Cases**: 已有同名 Agent 或配置时报告冲突，不静默覆盖。
- **Dependencies**: Journal、角色模型、adapter 子 FEATURE。
- **Parallel**: no
- **Verification**: 端到端覆盖设计、恢复、实施、升级、复核和 adapter 冲突。

## Verification

- 中立策略不含宿主专有字段，未安装 adapter 时现有工作流不回归。
- artifact、Journal 和宿主 session 职责保持单向依赖。
- AOK 保留和删除的机制均有明确迁移去向与理由。
- `go test ./internal/...` 与 adapter smoke 通过。

## History

- 2026-08-10 [decision] AOK 采用骨架重写，不原样迁移。
- 2026-08-10 [decision] FlowForge 是 artifact、门禁和安装生命周期的唯一事实源。
- 2026-08-10T21:06:35+08:00 | progress | 完成 AOK 精华整合：FlowForge 成为 Journal、artifact、阶段门禁和资产生命周期事实源；中立 policy、角色协议、preflight、风险复核和显式 OpenCode adapter 已实现。

## Links

### Outgoing

- [PROP-CR26081001](../../../03-proposal/CR26081001_subagent-协作机制整合.md) [proposal] - Subagent 协作机制整合
#### requires
- [FEAT-CR26081001-dkkymba7v03t](FEAT-CR26081001-dkkymba7v03t_引入-proposal-journal-协作日志.md) [feature] - 引入 Proposal Journal 协作日志
- [FEAT-CR26081001-dkkymf4v3m1u](FEAT-CR26081001-dkkymf4v3m1u_生成并部署宿主-subagent-adapter.md) [feature] - 生成并部署宿主 Subagent Adapter
- [FEAT-CR26081001-dkkymgqx8c63](FEAT-CR26081001-dkkymgqx8c63_重构-subagent-路由与职责模型.md) [feature] - 重构 Subagent 路由与职责模型

## Open Questions

None

## Dependencies

- Proposal Journal
- Subagent 路由与职责模型
- 宿主 Subagent Adapter
